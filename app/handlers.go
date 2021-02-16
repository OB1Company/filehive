package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OB1Company/filehive/fil"
	"github.com/OB1Company/filehive/repo/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/filecoin-project/go-address"
	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrBadPassword        = errors.New("password is too short")
	ErrWeakPassword       = errors.New("password is too weak")
	ErrIncorrectPassword  = errors.New("password is incorrect")
	ErrInvalidEmail       = errors.New("email address is invalid")
	ErrInvalidJSON        = errors.New("invalid JSON input")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotResetting   = errors.New("password token and email combo not valid")
	ErrDatasetNotFound    = errors.New("dataset not found")
	ErrNotLoggedIn        = errors.New("not logged in")
	ErrInvalidImage       = errors.New("invalid base64 image")
	ErrInvalidAddress     = errors.New("invalid address")
	ErrImageNotFound      = errors.New("image not found")
	ErrInvalidOption      = errors.New("invalid option")
	ErrMissingForm        = errors.New("missing form")
	ErrInsuffientFunds    = errors.New("insufficient funds")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

const jwtExpirationHours = 24 * 7

type claims struct {
	Email string `json:"Email"`
	jwt.StandardClaims
}

func wrapError(err error) string {
	return fmt.Sprintf(`{"error": "%s"}`, err.Error())
}

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !strings.Contains(e, "@") || !strings.Contains(e, ".") {
		return false
	}
	return emailRegex.MatchString(e)
}

func (s *FileHiveServer) loginUser(w http.ResponseWriter, email string) {
	expirationTime := time.Now().Add(jwtExpirationHours * time.Hour)

	claims := &claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	httpOnly := true
	secureToken := false

	if s.useSSL {
		httpOnly = false
		secureToken = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Domain:   s.domain,
		MaxAge:   0,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: httpOnly,
		Secure:   secureToken,
	})
}

func (s *FileHiveServer) handlePOSTLogin(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var creds credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}
	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(creds.Email)).First(&user).Error

	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrIncorrectPassword), http.StatusUnauthorized)
			return
		}
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	hashedPW := hashPassword([]byte(creds.Password), user.Salt)
	if !bytes.Equal(hashedPW, user.HashedPassword) {
		http.Error(w, wrapError(ErrIncorrectPassword), http.StatusUnauthorized)
		return
	}

	s.loginUser(w, creds.Email)
}

func (s *FileHiveServer) handlePOSTLogout(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	_, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	httpOnly := true
	secureToken := false

	if s.useSSL {
		httpOnly = false
		secureToken = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "expired",
		Expires:  time.Time{},
		Domain:   s.domain,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: httpOnly,
		Secure:   secureToken,
	})
}

func (s *FileHiveServer) handlePOSTTokenExtend(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}
	s.loginUser(w, email)
}

func (s *FileHiveServer) handlePOSTUser(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Country  string `json:"country"`
	}
	var d data
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	if !isEmailValid(d.Email) {
		http.Error(w, wrapError(ErrInvalidEmail), http.StatusBadRequest)
		return
	}

	err := s.db.View(func(db *gorm.DB) error {
		var user models.User
		return db.Where("LOWER(email) = ?", strings.ToLower(d.Email)).First(&user).Error

	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, wrapError(ErrUserExists), http.StatusConflict)
		return
	}

	if passwordScore(d.Password) < 3 {
		http.Error(w, wrapError(ErrWeakPassword), http.StatusBadRequest)
		return
	}

	userId, token, err := s.filecoinBackend.CreateUser()
	if err != nil {
		http.Error(w, wrapError(ErrWeakPassword), http.StatusBadRequest)
	}

	newAddress, err := s.walletBackend.NewAddress(token)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	salt := makeSalt()
	hashedPW := hashPassword([]byte(d.Password), salt)

	id, err := makeID()
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	// Generate activation code for creating datasets
	otp, err := GenerateOTP(6)
	if err != nil {
		log.Errorf("error generating OTP code: %v", err)
	}

	user := models.User{
		ID:              id,
		Email:           d.Email,
		Name:            d.Name,
		Country:         d.Country,
		Salt:            salt,
		HashedPassword:  hashedPW,
		FilecoinAddress: newAddress,
		PowergateToken:  token,
		PowergateID:     userId,
		ActivationCode:  otp,
	}

	err = s.db.Update(func(db *gorm.DB) error {
		return db.Save(&user).Error
	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	// Send email notification
	mg := mailgun.NewMailgun(s.mailDomain, s.mailgunKey)

	sender := "administrator@" + s.mailDomain
	subject := "Welcome to Filehive! 🐝"
	body := ""
	recipient := d.Email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	// If you want to use a pre-made template in Mailgun set it here
	//message.SetTemplate("welcome-email")

	pwd, _ := os.Getwd()
	template, err := ioutil.ReadFile(filepath.Join(pwd, "email_templates/welcome-email.tpl"))
	if err != nil {
		log.Debug(err)
	}

	templateString := strings.ReplaceAll(string(template), "%recipient_name%", d.Name)
	templateString = strings.ReplaceAll(templateString, "%domain_name%", s.mailDomain)
	templateString = strings.ReplaceAll(templateString, "%code%", otp)
	templateString = strings.ReplaceAll(templateString, "%email%", url.QueryEscape(d.Email))

	message.SetHtml(templateString)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)
	log.Debugf("Mailgun Response: %v, %v", resp, id)

	if err != nil {
		log.Error(err)
	}

	s.loginUser(w, user.Email)
}

func (s *FileHiveServer) handleGETUsers(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if !user.Admin {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var (
		users []models.User
		count int64
	)
	err = s.db.View(func(db *gorm.DB) error {
		if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
			return err
		}
		return db.Order("created_at DESC").Find(&users).Error

	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Users []models.User `json:"users"`
	}{
		Users: users,
	})

}

func (s *FileHiveServer) handleGETUser(w http.ResponseWriter, r *http.Request) {
	var email, userID string
	emailOrIDFromPath := mux.Vars(r)["emailOrID"]
	if emailOrIDFromPath != "" {
		if isEmailValid(emailOrIDFromPath) {
			email = emailOrIDFromPath
		} else {
			userID = emailOrIDFromPath
		}
	} else {
		emailIface := r.Context().Value("email")

		emailFromToken, ok := emailIface.(string)
		if !ok {
			http.Error(w, wrapError(ErrInvalidCredentials), http.StatusInternalServerError)
			return
		}
		email = emailFromToken
	}

	var user models.User

	err := s.db.View(func(db *gorm.DB) error {
		if email != "" {
			return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error
		} else {
			return db.Where("id = ?", userID).First(&user).Error
		}

	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrUserNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, user)
}

func (s *FileHiveServer) handlePATCHUser(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	currentEmail, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusInternalServerError)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(currentEmail)).First(&user).Error

	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrUserNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	type data struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Country  string `json:"country"`
		Avatar   string `json:"avatar"`
	}
	var d data
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}
	var newPW []byte
	if d.Password != "" {
		newPW = hashPassword([]byte(d.Password), user.Salt)
	}

	var emailChanged bool
	err = s.db.Update(func(db *gorm.DB) error {
		if d.Email != "" && strings.ToLower(d.Email) != strings.ToLower(currentEmail) {
			if !isEmailValid(d.Email) {
				return ErrInvalidEmail
			}

			var checkUser models.User
			if err := db.Where("LOWER(email) = ?", strings.ToLower(d.Email)).First(&checkUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserExists
			}

			user.Email = d.Email
			emailChanged = true
		}
		if d.Name != "" {
			user.Name = d.Name

			if err := db.Model(&models.Dataset{}).Where("user_id = ?", user.ID).Update("username", d.Name).Error; err != nil {
				return err
			}
		}
		if d.Country != "" {
			user.Country = d.Country
		}
		if newPW != nil {
			user.HashedPassword = newPW
		}

		if d.Avatar != "" {
			filename := fmt.Sprintf("avatar-%s.jpg", user.ID)
			if err := saveAvatar(path.Join(s.staticFileDir, "images", filename), d.Avatar); err != nil {
				return err
			}
			user.AvatarFilename = filename
		}

		if err := db.Save(&user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrInvalidImage) {
			http.Error(w, wrapError(err), http.StatusBadRequest)
			return
		} else if errors.Is(err, ErrUserExists) {
			http.Error(w, wrapError(ErrInvalidJSON), http.StatusConflict)
			return
		}
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusInternalServerError)
		return
	}
	if emailChanged {
		s.loginUser(w, d.Email)
	}
}

func (s *FileHiveServer) handleGETImage(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]

	f, err := os.Open(path.Join(s.staticFileDir, "images", filename))
	if err != nil {
		http.Error(w, wrapError(ErrImageNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeContent(w, r, filename, time.Now(), f)
}

func (s *FileHiveServer) handleGETWalletAddress(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	sanitizedJSONResponse(w, struct {
		Address string
	}{
		Address: user.FilecoinAddress,
	})
}

func (s *FileHiveServer) handleGETWalletBalance(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	balance, err := s.walletBackend.Balance(user.FilecoinAddress, user.PowergateToken)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Balance float64
	}{
		Balance: fil.AttoFILToFIL(balance),
	})
}

func (s *FileHiveServer) handlePOSTWalletSend(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	type data struct {
		Address string  `json:"address"`
		Amount  float64 `json:"amount"`
	}
	var d data
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	txid, err := s.walletBackend.Send(user.FilecoinAddress, d.Address, fil.FILtoAttoFIL(d.Amount), user.PowergateToken)
	if err != nil {
		if errors.Is(err, fil.ErrInsuffientFunds) {
			http.Error(w, wrapError(fil.ErrInsuffientFunds), http.StatusBadRequest)
			return
		}
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Txid string
	}{
		Txid: txid,
	})
}

func (s *FileHiveServer) handleGETWalletTransactions(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var limit, offset int
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, wrapError(ErrInvalidOption), http.StatusBadRequest)
			return
		}
	} else {
		limit = -1
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, wrapError(ErrInvalidOption), http.StatusBadRequest)
			return
		}
	}

	txs, err := s.walletBackend.Transactions(user.FilecoinAddress, limit, offset)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, txs)
}

func (s *FileHiveServer) handlePOSTGenerateCoins(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Address string  `json:"address"`
		Amount  float64 `json:"amount"`
	}
	var d data
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	s.walletBackend.(*fil.MockWalletBackend).GenerateToAddress(d.Address, fil.FILtoAttoFIL(d.Amount))
}

func (s *FileHiveServer) handleGETDelist(w http.ResponseWriter, r *http.Request) {
	sp := strings.Split(r.URL.Path, "/")
	id := sp[len(sp)-1]

	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if user.Admin {
		err = s.db.Update(func(db *gorm.DB) error {
			if err := db.Model(&models.Dataset{}).Where("id = ?", id).Update("delisted", true).Error; err != nil {
				return err
			}
			return nil
		})
	} else {
		err = s.db.Update(func(db *gorm.DB) error {
			if err := db.Model(&models.Dataset{}).Where("id = ? and email = ?", id, email).Update("delisted", true).Error; err != nil {
				return err
			}
			return nil
		})
	}

}

func (s *FileHiveServer) handleGETRelist(w http.ResponseWriter, r *http.Request) {
	sp := strings.Split(r.URL.Path, "/")
	id := sp[len(sp)-1]

	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if user.Admin {
		err = s.db.Update(func(db *gorm.DB) error {
			if err := db.Model(&models.Dataset{}).Where("id = ?", id).Update("delisted", false).Error; err != nil {
				return err
			}
			return nil
		})
	} else {
		err = s.db.Update(func(db *gorm.DB) error {
			if err := db.Model(&models.Dataset{}).Where("id = ? and email = ?", id, email).Update("delisted", false).Error; err != nil {
				return err
			}
			return nil
		})
	}
}

func (s *FileHiveServer) handlePOSTDataset(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	addr, err := address.NewFromString(user.FilecoinAddress)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	id, err := makeID()
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	var (
		containsFile, containsMetadata bool
		dataset                        models.Dataset
		jobID                          string
		size                           int64
		cid                            string
	)
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if part.FormName() == "file" {
			fileBytes, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "failed to read content of the part", http.StatusInternalServerError)
				return
			}

			jobID, cid, _, err = s.filecoinBackend.Store(bytes.NewReader(fileBytes), addr, user.PowergateToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			size = int64(len(fileBytes))

			containsFile = true
		}

		if part.FormName() == "metadata" {
			type data struct {
				Title            string  `json:"title"`
				ShortDescription string  `json:"shortDescription"`
				FullDescription  string  `json:"fullDescription"`
				Image            string  `json:"image"`
				FileType         string  `json:"fileType"`
				Price            float64 `json:"price"`
				Filename         string  `json:"filename"`
			}
			var d data
			if err := json.NewDecoder(part).Decode(&d); err != nil {
				http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
				return
			}

			filename := fmt.Sprintf("%s.jpg", id)
			if err := saveDatasetImage(path.Join(s.staticFileDir, "images", filename), d.Image); err != nil {
				http.Error(w, wrapError(ErrInvalidImage), http.StatusBadRequest)
				return
			}

			dataset = models.Dataset{
				Title:            d.Title,
				ShortDescription: d.ShortDescription,
				FullDescription:  d.FullDescription,
				FileType:         d.FileType,
				Price:            d.Price,
				UserID:           user.ID,
				ID:               id,
				Username:         user.Name,
				ImageFilename:    filename,
				DatasetFilename:  d.Filename,
			}
			containsMetadata = true
		}
	}
	dataset.FileSize = size
	dataset.ContentID = cid

	if !containsFile || !containsMetadata {
		http.Error(w, wrapError(ErrMissingForm), http.StatusInternalServerError)
		return
	}
	dataset.JobID = jobID
	err = s.db.Update(func(db *gorm.DB) error {
		return db.Save(&dataset).Error
	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}
	sanitizedJSONResponse(w, struct {
		DatasetID string `json:"datasetID"`
	}{
		DatasetID: dataset.ID,
	})
}

func (s *FileHiveServer) handlePATCHDataset(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	type data struct {
		ID               string  `json:"id"`
		Title            string  `json:"title"`
		ShortDescription string  `json:"shortDescription"`
		FullDescription  string  `json:"fullDescription"`
		Image            string  `json:"image"`
		FileType         string  `json:"fileType"`
		Price            float64 `json:"price"`
	}
	var d data
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	var dataset models.Dataset
	err = s.db.View(func(db *gorm.DB) error {
		return db.Where("id = ?", d.ID).First(&dataset).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrDatasetNotFound), http.StatusBadRequest)
		return
	}

	if dataset.UserID != user.ID {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if d.Image != "" {
		filename := fmt.Sprintf("%s.jpg", d.ID)
		if err := saveDatasetImage(path.Join(s.staticFileDir, "images", filename), d.Image); err != nil {
			http.Error(w, wrapError(ErrInvalidImage), http.StatusBadRequest)
			return
		}
	}

	if d.Title != "" {
		dataset.Title = d.Title
	}
	if d.ShortDescription != "" {
		dataset.ShortDescription = d.ShortDescription
	}
	if d.FullDescription != "" {
		dataset.FullDescription = d.FullDescription
	}
	if d.Price != dataset.Price {
		dataset.Price = d.Price
	}
	if d.FileType != "" {
		dataset.FileType = d.FileType
	}

	err = s.db.Update(func(db *gorm.DB) error {
		return db.Save(&dataset).Error

	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}
}

func (s *FileHiveServer) handleGETDatasetFile(w http.ResponseWriter, r *http.Request) {
	sp := strings.Split(r.URL.Path, "/")
	id := sp[len(sp)-1]

	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	// Get dataset
	var dataset models.Dataset
	err = s.db.View(func(db *gorm.DB) error {
		return db.Where("id = ?", id).First(&dataset).Error
	})
	if err != nil {
		http.Error(w, wrapError(ErrDatasetNotFound), http.StatusBadRequest)
		return
	}

	// Get datset uploader account token
	var uploader models.User
	err = s.db.View(func(db *gorm.DB) error {
		return db.Where("id = ?", dataset.UserID).First(&uploader).Error
	})
	if err != nil {
		http.Error(w, wrapError(ErrUserNotFound), http.StatusUnauthorized)
		return
	}

	fileStream, err := s.filecoinBackend.Get(dataset.ContentID, uploader.PowergateToken)
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+dataset.DatasetFilename)
	w.Header().Set("Content-Type", dataset.FileType)

	io.Copy(w, fileStream)
}

func (s *FileHiveServer) handleGETPurchased(w http.ResponseWriter, r *http.Request) {
	sp := strings.Split(r.URL.Path, "/")
	id := sp[len(sp)-1]

	// Get email address from session
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	// Get user info
	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	// Retrieve matching purchase if it exists
	var purchase models.Purchase
	err = s.db.View(func(db *gorm.DB) error {
		return db.Where("user_id = ? and dataset_id = ?", user.ID, id).First(&purchase).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrDatasetNotFound), http.StatusBadRequest)
		return
	}

	sanitizedJSONResponse(w, struct {
		Success bool `json:"success"`
	}{
		Success: true,
	})
}

func (s *FileHiveServer) handleGETDatasetDeal(w http.ResponseWriter, r *http.Request) {
	sp := strings.Split(r.URL.Path, "/")
	id := sp[len(sp)-1]

	var dataset models.Dataset
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("content_id = ?", id).First(&dataset).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrDatasetNotFound), http.StatusBadRequest)
		return
	}

	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err = s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	jobStatus, err := s.filecoinBackend.JobStatus(dataset.JobID, user.PowergateToken)
	if err != nil {
		http.Error(w, wrapError(ErrDatasetNotFound), http.StatusBadRequest)
		return
	}

	sanitizedJSONResponse(w, jobStatus)
}

func (s *FileHiveServer) handleGETDataset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var dataset models.Dataset
	err := s.db.Update(func(db *gorm.DB) error {
		if err := db.Where("id = ?", id).First(&dataset).Error; err != nil {
			return err
		}

		if err := db.Save(&models.Click{DatasetID: dataset.ID, Timestamp: time.Now()}).Error; err != nil {
			return err
		}

		dataset.Views++
		return db.Save(&dataset).Error

	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrDatasetNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, dataset)
}

func (s *FileHiveServer) handleGETDatasets(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var page int
	pageStr := mux.Vars(r)["page"]
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, wrapError(ErrInvalidOption), http.StatusBadRequest)
			return
		}
	}

	var (
		datasets []models.Dataset
		count    int64
	)
	err = s.db.View(func(db *gorm.DB) error {
		if err := db.Model(&models.Dataset{}).Where("user_id = ?", user.ID).Count(&count).Error; err != nil {
			return err
		}
		return db.Where("user_id = ?", user.ID).Offset(page * 10).Limit(10).Find(&datasets).Error

	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrDatasetNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Pages    int              `json:"pages"`
		Page     int              `json:"page"`
		Datasets []models.Dataset `json:"datasets"`
	}{
		Pages:    (int(count) / 10) + 1,
		Page:     page,
		Datasets: datasets,
	})
}

func (s *FileHiveServer) handlePOSTDisableUsers(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if !user.Admin {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	type Users struct {
		Users []string `json:"users"`
	}
	var disabledUsers Users
	if err := json.NewDecoder(r.Body).Decode(&disabledUsers); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	for _, userId := range disabledUsers.Users {

		// Delist their listings
		var (
			datasets []models.Dataset
		)
		err = s.db.View(func(db *gorm.DB) error {
			return db.Where("user_id = ?", userId).Find(&datasets).Error
		})

		for _, ds := range datasets {
			err = s.db.Update(func(db *gorm.DB) error {
				if err := db.Model(&models.Dataset{}).Where("id = ?", ds.ID).Update("delisted", true).Error; err != nil {
					return err
				}
				return nil
			})
		}

		// Disable their account
		err = s.db.Update(func(db *gorm.DB) error {
			if err := db.Model(&models.User{}).Where("id = ?", userId).Update("disabled", true).Error; err != nil {
				return err
			}
			return nil
		})
	}

	log.Debug(disabledUsers)
}

func (s *FileHiveServer) handlePOSTMakeAdmin(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if !user.Admin {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	type Users struct {
		Users []string `json:"users"`
	}
	var adminUsers Users
	if err := json.NewDecoder(r.Body).Decode(&adminUsers); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	for _, userId := range adminUsers.Users {
		err = s.db.Update(func(db *gorm.DB) error {
			if err := db.Model(&models.User{}).Where("id = ?", userId).Update("admin", true).Error; err != nil {
				return err
			}
			return nil
		})
	}

	log.Debug(adminUsers)

}

func (s *FileHiveServer) handlePOSTMakeUser(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if !user.Admin {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	type Users struct {
		Users []string `json:"users"`
	}
	var adminUsers Users
	if err := json.NewDecoder(r.Body).Decode(&adminUsers); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	for _, userId := range adminUsers.Users {
		err = s.db.Update(func(db *gorm.DB) error {
			if err := db.Model(&models.User{}).Where("id = ?", userId).Update("admin", false).Error; err != nil {
				return err
			}
			return nil
		})
	}

	log.Debug(adminUsers)

}

func (s *FileHiveServer) handlePOSTEnableUsers(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if !user.Admin {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	type Users struct {
		Users []string `json:"users"`
	}
	var enabledUsers Users
	if err := json.NewDecoder(r.Body).Decode(&enabledUsers); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	for _, userId := range enabledUsers.Users {

		// Relist
		var (
			datasets []models.Dataset
		)
		err = s.db.View(func(db *gorm.DB) error {
			return db.Where("user_id = ?", userId).Find(&datasets).Error
		})

		for _, ds := range datasets {
			err = s.db.Update(func(db *gorm.DB) error {
				if err := db.Model(&models.Dataset{}).Where("id = ?", ds.ID).Update("delisted", false).Error; err != nil {
					return err
				}
				return nil
			})
		}

		// Re-enable user
		err = s.db.Update(func(db *gorm.DB) error {
			if err := db.Model(&models.User{}).Where("id = ?", userId).Update("disabled", false).Error; err != nil {
				return err
			}
			return nil
		})
	}

	log.Debug(enabledUsers)
}

func (s *FileHiveServer) handlePOSTPurchase(w http.ResponseWriter, r *http.Request) {
	sp := strings.Split(r.URL.Path, "/")
	id := sp[len(sp)-1]

	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var (
		dataset     models.Dataset
		datasetUser models.User
	)
	err = s.db.View(func(db *gorm.DB) error {
		if err := db.Where("id = ?", id).First(&dataset).Error; err != nil {
			return ErrDatasetNotFound
		}

		return db.Where("id = ?", dataset.UserID).First(&datasetUser).Error

	})
	if err != nil {
		if errors.Is(err, ErrDatasetNotFound) {
			http.Error(w, wrapError(ErrDatasetNotFound), http.StatusBadRequest)
			return
		} else {
			http.Error(w, wrapError(err), http.StatusInternalServerError)
			return
		}
	}

	purchaseID, err := makeID()
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	balance, err := s.walletBackend.Balance(user.FilecoinAddress, "")
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	amt := fil.FILtoAttoFIL(dataset.Price)
	if balance.Cmp(amt) < 0 {
		http.Error(w, wrapError(ErrInsuffientFunds), http.StatusBadRequest)
		return
	}

	feeAmount := new(big.Int).Set(amt)

	if s.filecoinAddress != "" {
		// Send fee amount to Filehive if address is specified
		feeAmount.Div(amt, new(big.Int).SetInt64(20))
		_, err := s.walletBackend.Send(user.FilecoinAddress, s.filecoinAddress, feeAmount, user.PowergateToken)
		if err != nil {
			http.Error(w, wrapError(err), http.StatusInternalServerError)
			return
		}

		amt.Sub(amt, feeAmount)
	}

	txid, err := s.walletBackend.Send(user.FilecoinAddress, datasetUser.FilecoinAddress, amt, user.PowergateToken)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	purchase := models.Purchase{
		UserID:           user.ID,
		SellerID:         datasetUser.ID,
		Username:         datasetUser.Name,
		ImageFilename:    dataset.ImageFilename,
		Price:            dataset.Price,
		ShortDescription: dataset.ShortDescription,
		FileType:         dataset.FileType,
		DatasetID:        dataset.ID,
		Title:            dataset.Title,
		Timestamp:        time.Now(),
		ID:               purchaseID,
	}

	err = s.db.Update(func(db *gorm.DB) error {
		dataset.Purchases++
		if err := db.Save(&dataset).Error; err != nil {
			return err
		}
		return db.Save(&purchase).Error
	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	var seller models.User

	err = s.db.View(func(db *gorm.DB) error {
		return db.Where("id = ?", purchase.SellerID).First(&seller).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrUserNotFound), http.StatusNotFound)
			return
		}
	}

	// Get image thumb for dataset
	f, err := os.Open(path.Join(s.staticFileDir, "images", purchase.ImageFilename))
	if err != nil {
		http.Error(w, wrapError(ErrImageNotFound), http.StatusNotFound)
		return
	}

	image, err := jpeg.Decode(f)
	thumb := resize.Thumbnail(48, 48, image, resize.NearestNeighbor)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumb, &jpeg.Options{100})
	imageBit := buf.Bytes()
	thumbBase64 := base64.StdEncoding.EncodeToString([]byte(imageBit))

	// Send email to seller
	mg := mailgun.NewMailgun(s.mailDomain, s.mailgunKey)

	sender := "administrator@" + s.mailDomain
	subject := "You've made a sale on Filehive! 🤑"
	body := ""
	recipient := seller.Email

	message := mg.NewMessage(sender, subject, body, recipient)

	pwd, _ := os.Getwd()
	template, err := ioutil.ReadFile(filepath.Join(pwd, "email_templates/sale.tpl"))
	if err != nil {
		http.Error(w, wrapError(err), http.StatusBadRequest)
		return
	}

	templateString := strings.ReplaceAll(string(template), "%recipient_name%", seller.Name)
	templateString = strings.ReplaceAll(templateString, "%domain_name%", s.mailDomain)
	templateString = strings.ReplaceAll(templateString, "%api_domain%", "api."+s.mailDomain)
	templateString = strings.ReplaceAll(templateString, "%customer%", user.Name)
	templateString = strings.ReplaceAll(templateString, "%image%", "data:image/png;base64, "+thumbBase64)
	templateString = strings.ReplaceAll(templateString, "%dataset_name%", purchase.Title)
	templateString = strings.ReplaceAll(templateString, "%dataset_shortdescription%", purchase.ShortDescription)
	templateString = strings.ReplaceAll(templateString, "%dataset_price%", fmt.Sprintf("%f FIL", purchase.Price))
	templateString = strings.ReplaceAll(templateString, "%order_id%", purchase.ID)
	templateString = strings.ReplaceAll(templateString, "%timestamp%", purchase.Timestamp.Format("2006-01-02 15:04:05"))
	templateString = strings.ReplaceAll(templateString, "%email%", url.QueryEscape(user.Email))

	message.SetHtml(templateString)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)
	log.Debugf("Mailgun Response: %v, %v", resp, id)

	if err != nil {
		http.Error(w, wrapError(err), http.StatusBadRequest)
		return
	}

	sanitizedJSONResponse(w, struct {
		Txid string `json:"txid"`
	}{
		Txid: txid,
	})
}

func (s *FileHiveServer) handleGETAdminSales(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if !user.Admin {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var (
		sales []models.Purchase
		count int64
	)
	err = s.db.View(func(db *gorm.DB) error {
		if err := db.Model(&models.Purchase{}).Count(&count).Error; err != nil {
			return err
		}
		return db.Order("created_at DESC").Find(&sales).Error

	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Users []models.Purchase `json:"sales"`
	}{
		Users: sales,
	})
}

func (s *FileHiveServer) handleGETSales(w http.ResponseWriter, r *http.Request) {
	pagesize := 1000

	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var seller models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&seller).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var page int
	pageStr := mux.Vars(r)["page"]
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, wrapError(ErrInvalidOption), http.StatusBadRequest)
			return
		}
	}

	var (
		sales []models.Purchase
		count int64
	)
	err = s.db.View(func(db *gorm.DB) error {
		if err := db.Model(&models.Purchase{}).Where("seller_id = ?", seller.ID).Count(&count).Error; err != nil {
			return err
		}
		return db.Where("seller_id = ?", seller.ID).Offset(page * pagesize).Limit(pagesize).Find(&sales).Error

	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Pages int               `json:"pages"`
		Page  int               `json:"page"`
		Sales []models.Purchase `json:"sales"`
	}{
		Pages: (int(count) / pagesize) + 1,
		Page:  page,
		Sales: sales,
	})
}

func (s *FileHiveServer) handleGETPurchases(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	var page int
	pageStr := mux.Vars(r)["page"]
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, wrapError(ErrInvalidOption), http.StatusBadRequest)
			return
		}
	}

	var (
		purchases []models.Purchase
		count     int64
	)
	err = s.db.View(func(db *gorm.DB) error {
		if err := db.Model(&models.Purchase{}).Where("user_id = ?", user.ID).Count(&count).Error; err != nil {
			return err
		}
		return db.Where("user_id = ?", user.ID).Order("timestamp DESC").Offset(page * 1000).Limit(1000).Find(&purchases).Error

	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Pages     int               `json:"pages"`
		Page      int               `json:"page"`
		Purchases []models.Purchase `json:"purchases"`
	}{
		Pages:     (int(count) / 1000) + 1,
		Page:      page,
		Purchases: purchases,
	})
}

func (s *FileHiveServer) handleGETRecent(w http.ResponseWriter, r *http.Request) {
	var (
		page int
		err  error
	)
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, wrapError(ErrInvalidOption), http.StatusBadRequest)
			return
		}
	}

	var (
		recent []models.Dataset
		count  int64
	)

	err = s.db.View(func(db *gorm.DB) error {
		if err := db.Model(&models.Dataset{}).Count(&count).Error; err != nil {
			return err
		}
		return db.Order("created_at desc").Where("delisted = false").Offset(page * 10).Limit(10).Find(&recent).Error
	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Pages    int              `json:"pages"`
		Page     int              `json:"page"`
		Datasets []models.Dataset `json:"datasets"`
	}{
		Pages:    (int(count) / 10) + 1,
		Page:     page,
		Datasets: recent,
	})
}

func (s *FileHiveServer) handleGETTrending(w http.ResponseWriter, r *http.Request) {
	pageSize := 1000

	var (
		page int
		err  error
	)
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, wrapError(ErrInvalidOption), http.StatusBadRequest)
			return
		}
	}

	var (
		trending []models.Dataset
		recent   []models.Dataset
		count    int
	)

	err = s.db.Update(func(db *gorm.DB) error {
		if err := db.Where("timestamp < ?", time.Now().Add(-time.Hour*24)).Delete(&models.Click{}).Error; err != nil {
			return err
		}

		type result struct {
			DatasetID string
			Count     int64
		}

		var clicks []models.Click
		if err := db.Distinct("dataset_ID").Find(&clicks).Error; err != nil {
			return err
		}

		results := make([]result, 0, len(clicks))
		for _, click := range clicks {
			var count int64
			if err := db.Model(&models.Click{}).Where("dataset_id = ?", click.DatasetID).Count(&count).Error; err != nil {
				return err
			}

			results = append(results, result{DatasetID: click.DatasetID, Count: count})
		}

		sort.Slice(results, func(i, j int) bool { return results[i].Count > results[j].Count })
		count = len(results)

		if page > 0 && page*pageSize > int(count-1) {
			page = int(count - 1)
		}

		for _, res := range results[page*pageSize:] {
			var ds models.Dataset
			if err := db.Where("id = ? and delisted = 0", res.DatasetID).First(&ds).Error; err != nil {
				log.Debug("found a delisted dataset")
			} else {
				trending = append(trending, ds)
			}

			if len(trending) >= pageSize {
				return nil
			}
		}

		var recentCount int64
		if err := db.Model(&models.Dataset{}).Count(&recentCount).Error; err != nil {
			return err
		}

		if len(trending) < pageSize {
			trendingPages := (count / pageSize) + 1
			recentPage := 0
			if page-trendingPages > 0 {
				recentPage = page - trendingPages
			}

			tx := db.Order("created_at desc").Limit((recentPage * pageSize) + (pageSize - len(trending)))

			for _, r := range results {
				tx.Where("id != ?", r.DatasetID)
			}
			if err := tx.Find(&recent).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			trending = append(trending, recent...)
		}

		count += int(recentCount) - count

		return nil
	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Pages    int              `json:"pages"`
		Page     int              `json:"page"`
		Datasets []models.Dataset `json:"datasets"`
	}{
		Pages:    (count / pageSize) + 1,
		Page:     page,
		Datasets: trending,
	})
}

func (s *FileHiveServer) handleGETCheckResetCode(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	var user models.User
	success := true

	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ? and reset_token = ? and reset_valid > ?", strings.ToLower(email), code, time.Now().Format(time.RFC3339)).First(&user).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			success = false
		}
	}

	sanitizedJSONResponse(w, struct {
		Success bool `json:"success"`
	}{
		Success: success,
	})

}

func (s *FileHiveServer) handlePOSTPasswordReset(w http.ResponseWriter, r *http.Request) {
	type passwordReset struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}
	var newPasswordReset passwordReset
	if err := json.NewDecoder(r.Body).Decode(&newPasswordReset); err != nil {
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
		return
	}

	// Get user for the salt
	var user models.User
	err := s.db.View(func(db *gorm.DB) error {
		return db.Where("LOWER(email) = ?", strings.ToLower(newPasswordReset.Email)).First(&user).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrUserNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	if passwordScore(newPasswordReset.Password) < 3 {
		http.Error(w, wrapError(ErrWeakPassword), http.StatusBadRequest)
		return
	}

	// Hash the password
	var newPW []byte
	if newPasswordReset.Password != "" {
		newPW = hashPassword([]byte(newPasswordReset.Password), user.Salt)
	}

	// Update the user password, clear code where email and code match
	err = s.db.Update(func(db *gorm.DB) error {
		if err := db.Model(&models.User{}).Where("LOWER(email) = ? and reset_token = ?", strings.ToLower(newPasswordReset.Email), newPasswordReset.Code).Update("hashed_password", newPW).Update("reset_token", "").Error; err != nil {
			return err
		}
		return nil
	})

	success := true
	if err != nil {
		success = false
	}

	sanitizedJSONResponse(w, struct {
		Success bool `json:"success"`
	}{
		Success: success,
	})
}

func (s *FileHiveServer) handleGETPasswordReset(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	// Fix email if has space, it's supposed to be a +
	email = strings.Replace(email, " ", "+", 1)

	otp, err := GenerateOTP(12)
	if err != nil {
		log.Error(err)
	}

	var user models.User
	// Update user record with reset token and time limit
	err = s.db.Update(func(db *gorm.DB) error {
		if err := db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error; err != nil {
			return err
		}
		if err := db.Model(&models.User{}).Where("LOWER(email) = ?", strings.ToLower(email)).Update("reset_token", otp).Update("reset_valid", time.Now().Add(time.Hour*24).Format(time.RFC3339)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrUserNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, wrapError(err), http.StatusBadRequest)
		return
	}

	// Send email notification
	mg := mailgun.NewMailgun(s.mailDomain, s.mailgunKey)

	sender := "administrator@" + s.mailDomain
	subject := "Password Reset Instructions for Filehive Account"
	body := ""
	recipient := user.Email

	message := mg.NewMessage(sender, subject, body, recipient)

	pwd, _ := os.Getwd()
	template, err := ioutil.ReadFile(filepath.Join(pwd, "email_templates/password-reset.tpl"))
	if err != nil {
		http.Error(w, wrapError(err), http.StatusBadRequest)
		return
	}

	templateString := strings.ReplaceAll(string(template), "%recipient_name%", user.Name)
	templateString = strings.ReplaceAll(templateString, "%domain_name%", s.mailDomain)
	templateString = strings.ReplaceAll(templateString, "%code%", otp)
	templateString = strings.ReplaceAll(templateString, "%email%", url.QueryEscape(user.Email))

	message.SetHtml(templateString)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)
	log.Debugf("Mailgun Response: %v, %v", resp, id)

	if err != nil {
		http.Error(w, wrapError(err), http.StatusBadRequest)
		return
	}
}

func (s *FileHiveServer) handleGETConfirm(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	// Fix email if has space, it's supposed to be a +
	email = strings.Replace(email, " ", "+", 1)

	err := s.db.View(func(db *gorm.DB) error {
		if err := db.Model(&models.User{}).Where("LOWER(email) = ? and activation_code = ?", strings.ToLower(email), code).Update("activated", true).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusBadRequest)
		return
	}

}

func (s *FileHiveServer) handleGETSearch(w http.ResponseWriter, r *http.Request) {
	var (
		page int
		err  error
	)
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, wrapError(ErrInvalidOption), http.StatusBadRequest)
			return
		}
	}

	searchTerm := r.URL.Query().Get("query")

	var (
		results []models.Dataset
		count   int64
	)

	err = s.db.View(func(db *gorm.DB) error {

		searchTerm = fmt.Sprintf("%%%s%%", searchTerm)

		if err := db.Model(&models.Dataset{}).Where("(datasets.delisted = 0 and datasets.title LIKE ? OR datasets.short_description LIKE ? OR datasets.full_description LIKE ?) and users.Disabled = false", searchTerm, searchTerm, searchTerm).Joins("left join users on datasets.user_id=users.id").Count(&count).Error; err != nil {
			return err
		}
		if err := db.Model(&models.Dataset{}).Where("(datasets.delisted = 0 and datasets.title LIKE ? OR datasets.short_description LIKE ? OR datasets.full_description LIKE ?) and users.Disabled = false", searchTerm, searchTerm, searchTerm).Joins("left join users on datasets.user_id=users.id").Scan(&results).Offset(page * 10).Limit(10).Error; err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, wrapError(ErrDatasetNotFound), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Pages    int              `json:"pages"`
		Page     int              `json:"page"`
		Datasets []models.Dataset `json:"datasets"`
	}{
		Pages:    (int(count) / 10) + 1,
		Page:     page,
		Datasets: results[page*10:],
	})
}
