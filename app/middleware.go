package app

import (
	"context"
	"github.com/OB1Company/filehive/repo/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/csrf"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func (s *FileHiveServer) setCSRFHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("X-CSRF-Token", csrf.Token(r))
		}
		next.ServeHTTP(w, r)
	})
}

func (s *FileHiveServer) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, wrapError(ErrNotLoggedIn), http.StatusUnauthorized)
				return
			}
			http.Error(w, wrapError(err), http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &claims{}

		if c.Value == "expired" {
			http.Error(w, wrapError(ErrNotLoggedIn), http.StatusUnauthorized)
			return
		}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return s.jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, wrapError(err), http.StatusUnauthorized)
				return
			}
			http.Error(w, wrapError(err), http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, wrapError(err), http.StatusUnauthorized)
			return
		}

		// Check database for disabled account
		var user models.User
		err = s.db.View(func(db *gorm.DB) error {
			return db.Where("LOWER(email) = ? and disabled = false", strings.ToLower(claims.Email)).First(&user).Error

		})
		if err != nil {
			http.Error(w, wrapError(ErrInvalidCredentials), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(context.Background(), "email", claims.Email)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
