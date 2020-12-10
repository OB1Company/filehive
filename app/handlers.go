package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OB1Company/filehive/fil"
	"github.com/OB1Company/filehive/repo/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/filecoin-project/go-address"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrBadPassword        = errors.New("password is too short")
	ErrIncorrectPassword  = errors.New("password is incorrect")
	ErrInvalidEmail       = errors.New("email address is invalid")
	ErrInvalidJSON        = errors.New("invalid JSON input")
	ErrUserNotFound       = errors.New("user not found")
	ErrDatasetNotFound    = errors.New("dataset not found")
	ErrNotLoggedIn        = errors.New("not logged in")
	ErrInvalidImage       = errors.New("invalid base64 image")
	ErrInvalidAddress     = errors.New("invalid address")
	ErrImageNotFound      = errors.New("image not found")
	ErrInvalidOption      = errors.New("invalid option")
	ErrMissingForm        = errors.New("missing form")

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

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Domain:   s.domain,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
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
		return db.Where("email = ?", creds.Email).First(&user).Error

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

func (s *FileHiveServer) handlePOSTTokenExtend(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	email, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusInternalServerError)
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
		return db.Where("email = ?", d.Email).First(&user).Error

	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, wrapError(ErrUserExists), http.StatusConflict)
		return
	}

	if len(d.Password) == 0 {
		http.Error(w, wrapError(ErrBadPassword), http.StatusBadRequest)
		return
	}

	newAddress, err := s.walletBackend.NewAddress()
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

	user := models.User{
		ID:              id,
		Email:           d.Email,
		Name:            d.Name,
		Country:         d.Country,
		Salt:            salt,
		HashedPassword:  hashedPW,
		FilecoinAddress: newAddress.String(),
	}

	err = s.db.Update(func(db *gorm.DB) error {
		return db.Save(&user).Error
	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	s.loginUser(w, user.Email)
}

func (s *FileHiveServer) handleGETUser(w http.ResponseWriter, r *http.Request) {
	var email string
	emailFromPath := mux.Vars(r)["email"]
	if emailFromPath != "" {
		email = emailFromPath
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
		return db.Where("email = ?", email).First(&user).Error

	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, wrapError(ErrUserNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusInternalServerError)
		return
	}

	sanitizedJSONResponse(w, struct {
		Email   string
		Name    string
		Country string
		Avatar  string
	}{
		Email:   email,
		Name:    user.Name,
		Country: user.Country,
		Avatar:  user.AvatarFilename,
	})
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
		return db.Where("email = ?", currentEmail).First(&user).Error

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
		if d.Email != "" && d.Email != currentEmail {
			if !isEmailValid(d.Email) {
				return ErrInvalidEmail
			}

			var checkUser models.User
			if err := db.Where("email = ?", d.Email).First(&checkUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserExists
			}

			user.Email = d.Email
			emailChanged = true
		}
		if d.Name != "" {
			user.Name = d.Name
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
		return db.Where("email = ?", email).First(&user).Error

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
		return db.Where("email = ?", email).First(&user).Error

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

	balance, err := s.walletBackend.Balance(addr)
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
		return db.Where("email = ?", email).First(&user).Error

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

	from, err := address.NewFromString(user.FilecoinAddress)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	to, err := address.NewFromString(d.Address)
	if err != nil {
		http.Error(w, wrapError(ErrInvalidAddress), http.StatusBadRequest)
		return
	}

	txid, err := s.walletBackend.Send(from, to, fil.FILtoAttoFIL(d.Amount))
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
		Txid: txid.String(),
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
		return db.Where("email = ?", email).First(&user).Error

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

	addr, err := address.NewFromString(user.FilecoinAddress)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	txs, err := s.walletBackend.Transactions(addr, limit, offset)
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
	addr, err := address.NewFromString(d.Address)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	s.walletBackend.(*fil.MockWalletBackend).GenerateToAddress(addr, fil.FILtoAttoFIL(d.Amount))
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
		return db.Where("email = ?", email).First(&user).Error

	})
	if err != nil {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
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

		// FIXME: we are just saving this to file for now. In the future we will
		// need to ingest this into powergate and also check to make sure the user
		// has enough coins in his account in order to pay for the filecoin storage.
		if part.FormName() == "file" {
			outfile, err := os.Create(path.Join(s.staticFileDir, "files", id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer outfile.Close()

			_, err = io.Copy(outfile, part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
				ImageFilename:    filename,
			}
			containsMetadata = true
		}
	}

	if !containsFile || !containsMetadata {
		http.Error(w, wrapError(ErrMissingForm), http.StatusInternalServerError)
		return
	}
	err = s.db.Update(func(db *gorm.DB) error {
		return db.Save(&dataset).Error

	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}
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
		return db.Where("email = ?", email).First(&user).Error

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
