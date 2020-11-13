package app

import (
	"fmt"
	"github.com/OB1Company/filehive/fil"
	"github.com/OB1Company/filehive/repo"
	"github.com/OB1Company/filehive/repo/models"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
	"gorm.io/gorm"
	"math/big"
	"net/http"
	"testing"
)

func Test_Handlers(t *testing.T) {
	t.Run("User Tests", func(t *testing.T) {
		runAPITests(t, apiTests{
			{
				name:             "Post user success",
				path:             "/api/v1/user",
				method:           http.MethodPost,
				statusCode:       http.StatusOK,
				body:             []byte(`{"email": "brian@ob1.io", "password":"asdf", "name": "Brian", "country": "United_States"}`),
				expectedResponse: nil,
			},
			{
				name:             "Post user invalid JSON",
				path:             "/api/v1/user",
				method:           http.MethodPost,
				statusCode:       http.StatusBadRequest,
				body:             []byte(`{"email": "brian@ob1.io "password":"asdf", "name": "Brian", "country": "United_States"}`),
				expectedResponse: errorReturn(ErrInvalidJSON),
			},
			{
				name:             "Post user nil password",
				path:             "/api/v1/user",
				method:           http.MethodPost,
				statusCode:       http.StatusBadRequest,
				body:             []byte(`{"email": "brian2@ob1.io", "password":"", "name": "Brian", "country": "United_States"}`),
				expectedResponse: errorReturn(ErrBadPassword),
			},
			{
				name:             "Post user invalid email",
				path:             "/api/v1/user",
				method:           http.MethodPost,
				statusCode:       http.StatusBadRequest,
				body:             []byte(`{"email": "brian2ob1", "password":"adsf", "name": "Brian", "country": "United_States"}`),
				expectedResponse: errorReturn(ErrInvalidEmail),
			},
			{
				name:             "Post user already exists",
				path:             "/api/v1/user",
				method:           http.MethodPost,
				statusCode:       http.StatusConflict,
				body:             []byte(`{"email": "brian@ob1.io", "password":"", "name": "Brian", "country": "United_States"}`),
				expectedResponse: errorReturn(ErrUserExists),
			},
			{
				name:       "Get user while logged in",
				path:       "/api/v1/user",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Email   string
					Name    string
					Country string
					Avatar  string
				}{
					Email:   "brian@ob1.io",
					Name:    "Brian",
					Country: "United_States",
				}),
			},
			{
				name:       "Get user from path",
				path:       "/api/v1/user/brian@ob1.io",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Email   string
					Name    string
					Country string
					Avatar  string
				}{
					Email:   "brian@ob1.io",
					Name:    "Brian",
					Country: "United_States",
				}),
			},
			{
				name:             "Get user from path not found",
				path:             "/api/v1/user/chris@ob1.io",
				method:           http.MethodGet,
				statusCode:       http.StatusNotFound,
				expectedResponse: errorReturn(ErrUserNotFound),
			},
			{
				name:             "Patch user success",
				path:             "/api/v1/user",
				method:           http.MethodPatch,
				statusCode:       http.StatusOK,
				body:             []byte(`{"email": "brian@ob1.io", "password":"ffff", "name": "Brian2", "country": "Botswana"}`),
				expectedResponse: nil,
			},
			{
				name:       "Check user patched correctly",
				path:       "/api/v1/user/brian@ob1.io",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Email   string
					Name    string
					Country string
					Avatar  string
				}{
					Email:   "brian@ob1.io",
					Name:    "Brian2",
					Country: "Botswana",
				}),
			},
			{
				name:             "Patch user change email",
				path:             "/api/v1/user",
				method:           http.MethodPatch,
				statusCode:       http.StatusOK,
				body:             []byte(`{"email": "brian2@ob1.io"}`),
				expectedResponse: nil,
			},
			{
				name:       "Check user patched correctly",
				path:       "/api/v1/user/brian2@ob1.io",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Email   string
					Name    string
					Country string
					Avatar  string
				}{
					Email:   "brian2@ob1.io",
					Name:    "Brian2",
					Country: "Botswana",
				}),
			},
			{
				name:             "Check previous email deleted correctly",
				path:             "/api/v1/user/brian@ob1.io",
				method:           http.MethodGet,
				statusCode:       http.StatusNotFound,
				expectedResponse: errorReturn(ErrUserNotFound),
			},
		})
	})

	t.Run("Login Tests", func(t *testing.T) {
		runAPITests(t, apiTests{
			{
				name:             "Post extend token not logged in",
				path:             "/api/v1/token/extend",
				method:           http.MethodPost,
				statusCode:       http.StatusUnauthorized,
				body:             nil,
				expectedResponse: errorReturn(ErrNotLoggedIn),
			},
			{
				name:             "Post login invalid email",
				path:             "/api/v1/login",
				method:           http.MethodPost,
				statusCode:       http.StatusUnauthorized,
				body:             []byte(`{"email": "brian@ob1.io", "password":"asdf"}`),
				expectedResponse: errorReturn(ErrIncorrectPassword),
			},
			{
				name:             "Post login invalid JSON",
				path:             "/api/v1/login",
				method:           http.MethodPost,
				statusCode:       http.StatusBadRequest,
				body:             []byte(`{"email": "brian@ob1.io", "password":"asdf"`),
				expectedResponse: errorReturn(ErrInvalidJSON),
			},
			{
				name:       "Post login incorrect password",
				path:       "/api/v1/login",
				method:     http.MethodPost,
				statusCode: http.StatusUnauthorized,
				setup: func(db *repo.Database, wbe fil.WalletBackend) error {
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
				body:             []byte(`{"email": "brian@ob1.io", "password":"aaaaa"}`),
				expectedResponse: errorReturn(ErrIncorrectPassword),
			},
			{
				name:             "Post login valid",
				path:             "/api/v1/login",
				method:           http.MethodPost,
				statusCode:       http.StatusOK,
				body:             []byte(`{"email": "brian@ob1.io", "password":"asdf"}`),
				expectedResponse: nil,
			},
			{
				name:       "Get user while logged in",
				path:       "/api/v1/user",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Email   string
					Name    string
					Country string
					Avatar  string
				}{
					Email:   "brian@ob1.io",
					Name:    "Brian",
					Country: "United_States",
				}),
			},
			{
				name:             "Post extend token",
				path:             "/api/v1/token/extend",
				method:           http.MethodPost,
				statusCode:       http.StatusOK,
				body:             nil,
				expectedResponse: nil,
			},
		})
	})

	t.Run("Image Tests", func(t *testing.T) {
		runAPITests(t, apiTests{
			{
				name:             "Post user success",
				path:             "/api/v1/user",
				method:           http.MethodPost,
				statusCode:       http.StatusOK,
				body:             []byte(`{"email": "brian@ob1.io", "password":"asdf", "name": "Brian", "country": "United_States"}`),
				expectedResponse: nil,
			},
			{
				name:             "Patch user with avatar",
				path:             "/api/v1/user",
				method:           http.MethodPatch,
				statusCode:       http.StatusOK,
				body:             []byte(fmt.Sprintf(`{"avatar": "%s"}`, jpgTestImage)),
				expectedResponse: nil,
			},
			{
				name:             "Get avatar",
				path:             "/api/v1/image/avatar-1.jpg",
				method:           http.MethodGet,
				statusCode:       http.StatusOK,
				expectedResponse: jpgImageBytes,
			},
		})
	})

	t.Run("Wallet Tests", func(t *testing.T) {
		runAPITests(t, apiTests{
			{
				name:       "Post user success",
				path:       "/api/v1/user",
				method:     http.MethodPost,
				statusCode: http.StatusOK,
				setup: func(db *repo.Database, wbe fil.WalletBackend) error {
					addr, err := address.NewFromString("f1cu3c2dqsbyt7nq63x2yubyy6ofuini2nfvnnahi")
					if err != nil {
						return err
					}
					wbe.(*fil.MockWalletBackend).SetNextAddress(addr)
					return nil
				},
				body:             []byte(`{"email": "brian@ob1.io", "password":"asdf", "name": "Brian", "country": "United_States"}`),
				expectedResponse: nil,
			},
			{
				name:       "Get wallet address",
				path:       "/api/v1/wallet/address",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Address string
				}{
					Address: "f1cu3c2dqsbyt7nq63x2yubyy6ofuini2nfvnnahi",
				}),
			},
			{
				name:       "Get wallet balance",
				path:       "/api/v1/wallet/balance",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				setup: func(db *repo.Database, wbe fil.WalletBackend) error {
					txid, err := cid.Decode("bafkreiewgqfti56ls5zt2kko2utajoliipl3te7cl5lvtiowgny6qb2pde")
					if err != nil {
						return err
					}
					wbe.(*fil.MockWalletBackend).SetNextTxid(txid)
					addr, err := address.NewFromString("f1cu3c2dqsbyt7nq63x2yubyy6ofuini2nfvnnahi")
					if err != nil {
						return err
					}
					amt, _ := new(big.Int).SetString("15500000000000000000", 10)
					wbe.(*fil.MockWalletBackend).GenerateToAddress(addr, amt)
					return nil
				},
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Balance float64
				}{
					Balance: 15.5,
				}),
			},
			{
				name:       "Post wallet send",
				path:       "/api/v1/wallet/send",
				method:     http.MethodPost,
				statusCode: http.StatusOK,
				body:       []byte(`{"address": "f1gyvikksfdmokwhg5jhcrkvfqkyd2sjdy46klgbq", "amount": 1}`),
				setup: func(db *repo.Database, wbe fil.WalletBackend) error {
					txid, err := cid.Decode("bafkreif2mzhq6663465bcb2s3xgqefysbmr3a2bxloobw7s4vrxooj6kva")
					if err != nil {
						return err
					}
					wbe.(*fil.MockWalletBackend).SetNextTxid(txid)
					return nil
				},
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Txid string
				}{
					Txid: "bafkreif2mzhq6663465bcb2s3xgqefysbmr3a2bxloobw7s4vrxooj6kva",
				}),
			},
			{
				name:             "Post wallet send insufficient funds",
				path:             "/api/v1/wallet/send",
				method:           http.MethodPost,
				statusCode:       http.StatusBadRequest,
				body:             []byte(`{"address": "f1gyvikksfdmokwhg5jhcrkvfqkyd2sjdy46klgbq", "amount": 20}`),
				expectedResponse: errorReturn(fil.ErrInsuffientFunds),
			},
			{
				name:       "Get wallet balance",
				path:       "/api/v1/wallet/balance",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON(struct {
					Balance float64
				}{
					Balance: 14.5,
				}),
			},
			{
				name:       "Get wallet transactions",
				path:       "/api/v1/wallet/transactions",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON([]struct {
					To            string  `json:"to"`
					From          string  `json:"from"`
					TransactionID string  `json:"transactionID"`
					Amount        float64 `json:"amount"`
					Timestamp     string  `json:"timestamp"`
				}{
					{
						Timestamp:     "0001-01-01T00:00:00Z",
						Amount:        15.5,
						To:            "f1cu3c2dqsbyt7nq63x2yubyy6ofuini2nfvnnahi",
						From:          "",
						TransactionID: "bafkreiewgqfti56ls5zt2kko2utajoliipl3te7cl5lvtiowgny6qb2pde",
					},
					{
						Timestamp:     "0001-01-01T00:00:00Z",
						Amount:        1,
						To:            "f1gyvikksfdmokwhg5jhcrkvfqkyd2sjdy46klgbq",
						From:          "f1cu3c2dqsbyt7nq63x2yubyy6ofuini2nfvnnahi",
						TransactionID: "bafkreif2mzhq6663465bcb2s3xgqefysbmr3a2bxloobw7s4vrxooj6kva",
					},
				}),
			},
			{
				name:       "Get wallet transactions with limit",
				path:       "/api/v1/wallet/transactions?limit=1",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON([]struct {
					To            string  `json:"to"`
					From          string  `json:"from"`
					TransactionID string  `json:"transactionID"`
					Amount        float64 `json:"amount"`
					Timestamp     string  `json:"timestamp"`
				}{
					{
						Timestamp:     "0001-01-01T00:00:00Z",
						Amount:        15.5,
						To:            "f1cu3c2dqsbyt7nq63x2yubyy6ofuini2nfvnnahi",
						From:          "",
						TransactionID: "bafkreiewgqfti56ls5zt2kko2utajoliipl3te7cl5lvtiowgny6qb2pde",
					},
				}),
			},
			{
				name:       "Get wallet transactions with offset",
				path:       "/api/v1/wallet/transactions?offset=1",
				method:     http.MethodGet,
				statusCode: http.StatusOK,
				expectedResponse: mustMarshalAndSanitizeJSON([]struct {
					To            string  `json:"to"`
					From          string  `json:"from"`
					TransactionID string  `json:"transactionID"`
					Amount        float64 `json:"amount"`
					Timestamp     string  `json:"timestamp"`
				}{
					{
						Timestamp:     "0001-01-01T00:00:00Z",
						Amount:        1,
						To:            "f1gyvikksfdmokwhg5jhcrkvfqkyd2sjdy46klgbq",
						From:          "f1cu3c2dqsbyt7nq63x2yubyy6ofuini2nfvnnahi",
						TransactionID: "bafkreif2mzhq6663465bcb2s3xgqefysbmr3a2bxloobw7s4vrxooj6kva",
					},
				}),
			},
			{
				name:             "Get wallet transactions invalid limit",
				path:             "/api/v1/wallet/transactions?limit=zzz",
				method:           http.MethodGet,
				statusCode:       http.StatusBadRequest,
				expectedResponse: errorReturn(ErrInvalidOption),
			},
			{
				name:             "Get wallet transactions invalid offset",
				path:             "/api/v1/wallet/transactions?offset=zzz",
				method:           http.MethodGet,
				statusCode:       http.StatusBadRequest,
				expectedResponse: errorReturn(ErrInvalidOption),
			},
		})
	})
}
