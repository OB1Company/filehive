package app

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OB1Company/filehive/repo/models"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
)

var (
	ErrUserExists = errors.New("user already exists")
	ErrBadPassword = errors.New("password is too short")
	ErrInvalidEmail = errors.New("email address is invalid")
	ErrInvalidJSON = errors.New("invalid JSON input")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)


func wrapError(err error) string {
	return fmt.Sprintf(`{"error": "%s"}`, err.Error())
}

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !strings.Contains(e, "@") || !strings.Contains(e, "."){
		return false
	}
	return emailRegex.MatchString(e)
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

	salt := make([]byte, 32)
	rand.Read(salt)

	hashedPW:= pbkdf2.Key([]byte(d.Password), salt, 100000, 256, sha512.New512_256)

	user := models.User{
		Email: d.Email,
		Name: d.Name,
		Country: d.Country,
		Salt: salt,
		HashedPassword: hashedPW,
	}

	err = s.db.Update(func(db *gorm.DB) error {
		return db.Save(&user).Error
	})
	if err != nil {
		http.Error(w, wrapError(err), http.StatusInternalServerError)
		return
	}

	// TODO: I think we should log them in here.
}
