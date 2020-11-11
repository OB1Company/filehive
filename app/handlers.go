package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OB1Company/filehive/repo/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"regexp"
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
	ErrNotLoggedIn        = errors.New("not logged in")

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

	salt := makeSalt()
	hashedPW := hashPassword([]byte(d.Password), salt)

	user := models.User{
		Email:          d.Email,
		Name:           d.Name,
		Country:        d.Country,
		Salt:           salt,
		HashedPassword: hashedPW,
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
	}{
		Email:   email,
		Name:    user.Name,
		Country: user.Country,
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
		http.Error(w, wrapError(ErrInvalidCredentials), http.StatusInternalServerError)
		return
	}

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
	var newPW []byte
	if d.Password != "" {
		newPW = hashPassword([]byte(d.Password), user.Salt)
	}

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

			if err := db.Where("email = ?", currentEmail).Delete(&models.User{}).Error; err != nil {
				return err
			}
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
		return db.Save(&user).Error
	})
	if err != nil {
		if errors.Is(err, ErrInvalidEmail) {
			http.Error(w, wrapError(ErrInvalidJSON), http.StatusBadRequest)
			return
		} else if errors.Is(err, ErrUserExists) {
			http.Error(w, wrapError(ErrInvalidJSON), http.StatusConflict)
			return
		}
		http.Error(w, wrapError(ErrInvalidJSON), http.StatusInternalServerError)
		return
	}
}
