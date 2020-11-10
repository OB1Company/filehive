package app

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_Handlers(t *testing.T) {

	// User tests
	runAPITests(t, apiTests{
		{
			name: "Post user success",
			path: "/v1/user",
			method: http.MethodPost,
			statusCode: http.StatusOK,
			body: []byte(`{"email": "brian@ob1.io", "password":"asdf", "name": "Brian", "country": "United_States"}`),
			expectedResponse: func()([]byte, error) {
				return nil, nil
			},
		},
		{
			name: "Post user invalid JSON",
			path: "/v1/user",
			method: http.MethodPost,
			statusCode: http.StatusBadRequest,
			body: []byte(`{"email": "brian@ob1.io "password":"asdf", "name": "Brian", "country": "United_States"}`),
			expectedResponse: func()([]byte, error) {
				return []byte(fmt.Sprintf(`{"error": "%s"}%s`, ErrInvalidJSON.Error(), "\n")), nil
			},
		},
		{
			name: "Post user nil password",
			path: "/v1/user",
			method: http.MethodPost,
			statusCode: http.StatusBadRequest,
			body: []byte(`{"email": "brian2@ob1.io", "password":"", "name": "Brian", "country": "United_States"}`),
			expectedResponse: func()([]byte, error) {
				return []byte(fmt.Sprintf(`{"error": "%s"}%s`, ErrBadPassword.Error(), "\n")), nil
			},
		},
		{
			name: "Post user invalid email",
			path: "/v1/user",
			method: http.MethodPost,
			statusCode: http.StatusBadRequest,
			body: []byte(`{"email": "brian2ob1", "password":"adsf", "name": "Brian", "country": "United_States"}`),
			expectedResponse: func()([]byte, error) {
				return []byte(fmt.Sprintf(`{"error": "%s"}%s`, ErrInvalidEmail.Error(), "\n")), nil
			},
		},
		{
			name: "Post user already exists",
			path: "/v1/user",
			method: http.MethodPost,
			statusCode: http.StatusConflict,
			body: []byte(`{"email": "brian@ob1.io", "password":"", "name": "Brian", "country": "United_States"}`),
			expectedResponse: func()([]byte, error) {
				return []byte(fmt.Sprintf(`{"error": "%s"}%s`, ErrUserExists.Error(), "\n")), nil
			},
		},
	})
}
