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
	"github.com/ipfs/go-cid"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path"
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

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Domain:   s.domain,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   false,
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

func (s *FileHiveServer) handlePOSTLogout(w http.ResponseWriter, r *http.Request) {
	emailIface := r.Context().Value("email")

	_, ok := emailIface.(string)
	if !ok {
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "expired",
		Expires:  time.Time{},
		Domain:   s.domain,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
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
		return db.Where("email = ?", d.Email).First(&user).Error

	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, wrapError(ErrUserExists), http.StatusConflict)
		return
	}

	if passwordScore(d.Password) < 3 {
		http.Error(w, wrapError(ErrWeakPassword), http.StatusBadRequest)
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
			return db.Where("email = ?", email).First(&user).Error
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
		jobID                          cid.Cid
		size                           int64
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
			_, jobID, size, err = s.filecoinBackend.Store(part, addr)
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
				Username:         user.Name,
				ImageFilename:    filename,
			}
			containsMetadata = true
		}
	}
	dataset.FileSize = size

	if !containsFile || !containsMetadata {
		http.Error(w, wrapError(ErrMissingForm), http.StatusInternalServerError)
		return
	}
	dataset.JobID = jobID.String()
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
		return db.Where("email = ?", email).First(&user).Error

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
		return db.Where("email = ?", email).First(&user).Error

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

	userAddr, err := address.NewFromString(user.FilecoinAddress)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}
	datasetUserAddr, err := address.NewFromString(datasetUser.FilecoinAddress)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	balance, err := s.walletBackend.Balance(userAddr)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	amt := fil.FILtoAttoFIL(dataset.Price)
	if balance.Cmp(amt) < 0 {
		http.Error(w, wrapError(ErrInsuffientFunds), http.StatusBadRequest)
		return
	}

	txid, err := s.walletBackend.Send(userAddr, datasetUserAddr, amt)
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	purchase := models.Purchase{
		UserID:           user.ID,
		Username:         datasetUser.Name,
		ImageFilename:    dataset.ImageFilename,
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

	sanitizedJSONResponse(w, struct {
		Txid string `json:"txid"`
	}{
		Txid: txid.String(),
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
		return db.Where("email = ?", email).First(&user).Error

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
		return db.Where("user_id = ?", user.ID).Offset(page * 10).Limit(10).Find(&purchases).Error

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
		Pages:     (int(count) / 10) + 1,
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
		return db.Order("created_at desc").Offset(page * 10).Limit(10).Find(&recent).Error
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

		if page*10 > int(count-1) {
			page = int(count - 1)
		}

		for _, res := range results[page*10:] {
			var ds models.Dataset
			if err := db.Where("id = ?", res.DatasetID).First(&ds).Error; err != nil {
				return err
			}
			trending = append(trending, ds)
			if len(trending) >= 10 {
				return nil
			}
		}

		var recentCount int64
		if err := db.Model(&models.Dataset{}).Count(&recentCount).Error; err != nil {
			return err
		}

		if len(trending) < 10 {
			trendingPages := (count / 10) + 1
			recentPage := 0
			if page-trendingPages > 0 {
				recentPage = page - trendingPages
			}

			if err := db.Order("created_at desc").Limit((recentPage * 10) + (10 - len(trending))).Find(&recent).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			trending = append(trending, recent...)
		}

		count += int(recentCount)

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
		Pages:    (count / 10) + 1,
		Page:     page,
		Datasets: trending,
	})
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
		query := "SELECT * FROM datasets WHERE MATCH(title, short_description, full_description) AGAINST(? IN NATURAL LANGUAGE MODE)"
		if err := db.Raw(query, searchTerm).Count(&count).Error; err != nil {
			return err
		}
		if err := db.Raw(query, searchTerm).Scan(&results).Offset(page * 10).Limit(10).Error; err != nil {
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
