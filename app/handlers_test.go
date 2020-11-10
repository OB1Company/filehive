package app

import (
	"fmt"
	"github.com/OB1Company/filehive/repo"
	"github.com/OB1Company/filehive/repo/models"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func errorReturn(err error) []byte {
	return []byte(fmt.Sprintf(`{"error": "%s"}%s`, err.Error(), "\n"))
}

func Test_Handlers(t *testing.T) {
	t.Run("User Tests", func(t *testing.T) {
		runAPITests(t, apiTests{
			{
				name:       "Post user success",
				path:       "/api/v1/user",
				method:     http.MethodPost,
				statusCode: http.StatusOK,
				body:       []byte(`{"email": "brian@ob1.io", "password":"asdf", "name": "Brian", "country": "United_States"}`),
				expectedResponse: func() ([]byte, error) {
					return nil, nil
				},
			},
			{
				name:       "Post user invalid JSON",
				path:       "/api/v1/user",
				method:     http.MethodPost,
				statusCode: http.StatusBadRequest,
				body:       []byte(`{"email": "brian@ob1.io "password":"asdf", "name": "Brian", "country": "United_States"}`),
				expectedResponse: func() ([]byte, error) {
					return errorReturn(ErrInvalidJSON), nil
				},
			},
			{
				name:       "Post user nil password",
				path:       "/api/v1/user",
				method:     http.MethodPost,
				statusCode: http.StatusBadRequest,
				body:       []byte(`{"email": "brian2@ob1.io", "password":"", "name": "Brian", "country": "United_States"}`),
				expectedResponse: func() ([]byte, error) {
					return errorReturn(ErrBadPassword), nil
				},
			},
			{
				name:       "Post user invalid email",
				path:       "/api/v1/user",
				method:     http.MethodPost,
				statusCode: http.StatusBadRequest,
				body:       []byte(`{"email": "brian2ob1", "password":"adsf", "name": "Brian", "country": "United_States"}`),
				expectedResponse: func() ([]byte, error) {
					return errorReturn(ErrInvalidEmail), nil
				},
			},
			{
				name:       "Post user already exists",
				path:       "/api/v1/user",
				method:     http.MethodPost,
				statusCode: http.StatusConflict,
				body:       []byte(`{"email": "brian@ob1.io", "password":"", "name": "Brian", "country": "United_States"}`),
				expectedResponse: func() ([]byte, error) {
					return errorReturn(ErrUserExists), nil
				},
			},
			{
				name:       "Get user while logged in",
				path:       "/api/v1/user",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: func() ([]byte, error) {
					return marshalAndSanitizeJSON(struct {
						Email   string
						Name    string
						Country string
					}{
						Email:   "brian@ob1.io",
						Name:    "Brian",
						Country: "United_States",
					})
				},
			},
			{
				name:       "Get user from path",
				path:       "/api/v1/user/brian@ob1.io",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: func() ([]byte, error) {
					return marshalAndSanitizeJSON(struct {
						Email   string
						Name    string
						Country string
					}{
						Email:   "brian@ob1.io",
						Name:    "Brian",
						Country: "United_States",
					})
				},
			},
			{
				name:       "Get user from path not found",
				path:       "/api/v1/user/chris@ob1.io",
				method:     http.MethodGet,
				statusCode: http.StatusNotFound,
				expectedResponse: func() ([]byte, error) {
					return errorReturn(ErrUserNotFound), nil
				},
			},
		})
	})

	t.Run("Login Tests", func(t *testing.T) {
		runAPITests(t, apiTests{
			{
				name:       "Post login invalid email",
				path:       "/api/v1/login",
				method:     http.MethodPost,
				statusCode: http.StatusUnauthorized,
				body:       []byte(`{"email": "brian@ob1.io", "password":"asdf"}`),
				expectedResponse: func() ([]byte, error) {
					return errorReturn(ErrIncorrectPassword), nil
				},
			},
			{
				name:       "Post login invalid JSON",
				path:       "/api/v1/login",
				method:     http.MethodPost,
				statusCode: http.StatusBadRequest,
				body:       []byte(`{"email": "brian@ob1.io", "password":"asdf"`),
				expectedResponse: func() ([]byte, error) {
					return errorReturn(ErrInvalidJSON), nil
				},
			},
			{
				name:       "Post login incorrect password",
				path:       "/api/v1/login",
				method:     http.MethodPost,
				statusCode: http.StatusUnauthorized,
				setup: func(db *repo.Database) error {
					return db.Update(func(tx *gorm.DB) error {
						salt := []byte("salt")
						pw := hashPassword([]byte("asdf"), salt)
						return tx.Save(&models.User{
							Email:          "brian@ob1.io",
							Country:        "United_States",
							Name:           "Brian",
							Salt:           salt,
							HashedPassword: pw,
						}).Error
					})
				},
				body: []byte(`{"email": "brian@ob1.io", "password":"aaaaa"}`),
				expectedResponse: func() ([]byte, error) {
					return errorReturn(ErrIncorrectPassword), nil
				},
			},
			{
				name:       "Post login valid",
				path:       "/api/v1/login",
				method:     http.MethodPost,
				statusCode: http.StatusOK,
				setup: func(db *repo.Database) error {
					return db.Update(func(tx *gorm.DB) error {
						salt := []byte("salt")
						pw := hashPassword([]byte("asdf"), salt)
						return tx.Save(&models.User{
							Email:          "brian@ob1.io",
							Country:        "United_States",
							Name:           "Brian",
							Salt:           salt,
							HashedPassword: pw,
						}).Error
					})
				},
				body: []byte(`{"email": "brian@ob1.io", "password":"asdf"}`),
				expectedResponse: func() ([]byte, error) {
					return nil, nil
				},
			},
			{
				name:       "Get user while logged in",
				path:       "/api/v1/user",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: func() ([]byte, error) {
					return marshalAndSanitizeJSON(struct {
						Email   string
						Name    string
						Country string
					}{
						Email:   "brian@ob1.io",
						Name:    "Brian",
						Country: "United_States",
					})
				},
			},
		})
	})
}
