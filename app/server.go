package app

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/OB1Company/filehive/fil"
	"github.com/OB1Company/filehive/repo"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/multiformats/go-multihash"
	"github.com/op/go-logging"
	"golang.org/x/crypto/pbkdf2"
	"net"
	"net/http"
	"os"
	"path"
)

var log = logging.MustGetLogger("APP")

// FileHiveServer is the web server used to serve the FileHive application.
type FileHiveServer struct {
	db              *repo.Database
	walletBackend   fil.WalletBackend
	filecoinBackend fil.FilecoinBackend
	staticFileDir   string
	listener        net.Listener
	handler         http.Handler
	jwtKey          []byte
	domain          string
	shutdown        chan struct{}

	testMode bool
	useSSL   bool
	sslCert  string
	sslKey   string
}

// NewServer instantiates a new FileHiveServer with the provided options.
func NewServer(listener net.Listener, db *repo.Database, staticFileDir string, walletBackend fil.WalletBackend, filecoinBackend fil.FilecoinBackend, opts ...Option) (*FileHiveServer, error) {
	var options Options
	if err := options.Apply(opts...); err != nil {
		return nil, err
	}
	if listener == nil {
		return nil, errors.New("listener is nil")
	}
	if db == nil {
		return nil, errors.New("database is nil")
	}
	if staticFileDir == "" {
		return nil, errors.New("static file dir is empty string")
	}

	if options.TestMode {
		if _, ok := walletBackend.(*fil.MockWalletBackend); !ok {
			return nil, errors.New("MockWalletBackend must be used in testmode")
		}
	}

	if options.JWTKey == nil {
		jwtKey := make([]byte, 32)
		rand.Read(jwtKey)
		options.JWTKey = jwtKey
	}

	if err := os.MkdirAll(path.Join(staticFileDir, "images"), os.ModePerm); err != nil {
		return nil, err
	}

	var (
		s = &FileHiveServer{
			db:            db,
			walletBackend: walletBackend,
			listener:      listener,
			staticFileDir: staticFileDir,
			useSSL:        options.UseSSL,
			sslCert:       options.SSLCert,
			sslKey:        options.SSLKey,
			jwtKey:        options.JWTKey,
			domain:        options.Domain,
			shutdown:      make(chan struct{}),
		}
		topMux = http.NewServeMux()
	)

	r := s.newV1Router()

	csrfKey := make([]byte, 32)
	rand.Read(csrfKey)

	csrfOpts := []csrf.Option{
		csrf.SameSite(csrf.SameSiteLaxMode),
	}
	if options.Domain != "" {
		csrfOpts = append(csrfOpts, csrf.Domain(options.Domain))
	}
	csrfMiddleware := csrf.Protect(csrfKey, csrfOpts...)
	r.Use(
		csrfMiddleware,
		s.setCSRFHeaderMiddleware,
		mux.CORSMethodMiddleware(r),
	)

	topMux.Handle("/api/v1/", r)

	s.handler = topMux
	return s, nil
}

// Close shutsdown the Gateway listener.
func (s *FileHiveServer) Close() error {
	close(s.shutdown)
	return s.listener.Close()
}

// Serve begins listening on the configured address.
func (s *FileHiveServer) Serve() error {
	log.Infof("FileHive server listening on %s\n", s.listener.Addr().String())
	var err error
	if s.useSSL {
		err = http.ListenAndServeTLS(s.listener.Addr().String(), s.sslCert, s.sslKey, s.handler)
	} else {
		err = http.Serve(s.listener, s.handler)
	}
	return err
}

func (s *FileHiveServer) newV1Router() *mux.Router {
	r := mux.NewRouter()
	// Unauthenticated Routes
	r.HandleFunc("/api/v1/user", s.handlePOSTUser).Methods("POST")
	r.HandleFunc("/api/v1/user/{emailOrID}", s.handleGETUser).Methods("GET")
	r.HandleFunc("/api/v1/login", s.handlePOSTLogin).Methods("POST")
	r.HandleFunc("/api/v1/image/{filename}", s.handleGETImage).Methods("GET")
	r.HandleFunc("/api/v1/dataset/{id}", s.handleGETDataset).Methods("GET")

	if s.testMode {
		r.HandleFunc("/api/v1/generatecoins", s.handlePOSTGenerateCoins).Methods("POST")
	}

	// Authenticated Routes
	subRouter := r.PathPrefix("/api/v1").Subrouter()
	subRouter.Use(s.authenticationMiddleware)

	subRouter.HandleFunc("/token/extend", s.handlePOSTTokenExtend).Methods("POST")
	subRouter.HandleFunc("/user", s.handleGETUser).Methods("GET")
	subRouter.HandleFunc("/user", s.handlePATCHUser).Methods("PATCH")
	subRouter.HandleFunc("/wallet/address", s.handleGETWalletAddress).Methods("GET")
	subRouter.HandleFunc("/wallet/balance", s.handleGETWalletBalance).Methods("GET")
	subRouter.HandleFunc("/wallet/send", s.handlePOSTWalletSend).Methods("POST")
	subRouter.HandleFunc("/wallet/transactions", s.handleGETWalletTransactions).Methods("GET")
	subRouter.HandleFunc("/dataset", s.handlePOSTDataset).Methods("POST")
	subRouter.HandleFunc("/dataset", s.handlePATCHDataset).Methods("PATCH")
	subRouter.HandleFunc("/datasets", s.handleGETDatasets).Methods("GET")
	subRouter.HandleFunc("/purchase/{id}", s.handlePOSTPurchase).Methods("POST")
	subRouter.HandleFunc("/purchases", s.handleGETPurchases).Methods("GET")

	return r
}

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
		ctx := context.WithValue(context.Background(), "email", claims.Email)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

// Options represents the filehive server options.
type Options struct {
	JWTKey   []byte
	Domain   string
	UseSSL   bool
	SSLCert  string
	SSLKey   string
	TestMode bool
}

// Apply sets the provided options in the main options struct.
func (o *Options) Apply(opts ...Option) error {
	for i, opt := range opts {
		if err := opt(o); err != nil {
			return fmt.Errorf("option %d failed: %s", i, err)
		}
	}
	return nil
}

// Option represents a db option.
type Option func(*Options) error

// JWTKey represents a JSON Web Token key for the server.
// Use this if you want to persist the key to disk. If
// This option is nil a random key will be generated.
func JWTKey(key []byte) Option {
	return func(o *Options) error {
		o.JWTKey = key
		return nil
	}
}

// Domain sets the domain the server is running on.  Defaults to the current domain of the request
// only (recommended).
//
// This should be a hostname and not a URL. If set, the domain is treated as
// being prefixed with a '.' - e.g. "example.com" becomes ".example.com" and
// matches "www.example.com" and "secure.example.com".
func Domain(domain string) Option {
	return func(o *Options) error {
		o.Domain = domain
		return nil
	}
}

// TestMode option allows exposes an additional API
// to generate mock coins.
func TestMode(testMode bool) Option {
	return func(o *Options) error {
		o.TestMode = testMode
		return nil
	}
}

// UseSSL option allows you to set SSL on the server.
func UseSSL(useSSL bool) Option {
	return func(o *Options) error {
		o.UseSSL = useSSL
		return nil
	}
}

// SSLCert is required if using the UseSSL option.
func SSLCert(sslCert string) Option {
	return func(o *Options) error {
		o.SSLCert = sslCert
		return nil
	}
}

// SSLKey is required if using the UseSSL option.
func SSLKey(sslKey string) Option {
	return func(o *Options) error {
		o.SSLKey = sslKey
		return nil
	}
}

func hashPassword(pw, salt []byte) []byte {
	return pbkdf2.Key(pw, salt, 100000, 256, sha512.New512_256)
}

func makeSalt() []byte {
	salt := make([]byte, 32)
	rand.Read(salt)
	return salt
}

func makeID() (string, error) {
	idBytes := make([]byte, 32)
	rand.Read(idBytes)
	encoded, err := multihash.Encode(idBytes, multihash.IDENTITY)
	if err != nil {
		return "", err
	}
	id, err := multihash.Cast(encoded)
	if err != nil {
		return "", err
	}
	return id.B58String(), nil
}
