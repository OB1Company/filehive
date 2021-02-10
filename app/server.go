package app

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/OB1Company/filehive/fil"
	"github.com/OB1Company/filehive/repo"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/multiformats/go-multihash"
	"github.com/op/go-logging"
	"golang.org/x/crypto/pbkdf2"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
)

var log = logging.MustGetLogger("APP")

// FileHiveServer is the web server used to serve the FileHive application.
type FileHiveServer struct {
	db              *repo.Database
	walletBackend   fil.WalletBackend
	filecoinBackend fil.FilecoinBackend
	filecoinAddress string
	staticFileDir   string
	listener        net.Listener
	handler         http.Handler
	jwtKey          []byte
	domain          string
	mailgunKey      string
	mailDomain      string
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
		if _, ok := walletBackend.(*fil.PowergateWalletBackend); !ok {
			return nil, errors.New("MockWalletBackend must be used in testmode")
		}
		if _, ok := filecoinBackend.(*fil.PowergateBackend); !ok {
			return nil, errors.New("MockFilecoinBackend must be used in testmode")
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
			db:              db,
			walletBackend:   walletBackend,
			filecoinBackend: filecoinBackend,
			filecoinAddress: options.FilecoinAddress,
			listener:        listener,
			staticFileDir:   staticFileDir,
			useSSL:          options.UseSSL,
			sslCert:         options.SSLCert,
			sslKey:          options.SSLKey,
			jwtKey:          options.JWTKey,
			domain:          options.Domain,
			mailgunKey:      options.MailgunKey,
			mailDomain:      options.MailDomain,
			shutdown:        make(chan struct{}),
		}
		topMux = http.NewServeMux()
	)

	r := s.newV1Router()

	csrfKey := make([]byte, 32)
	rand.Read(csrfKey)

	csrfOpts := []csrf.Option{
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.Path("/"),
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
	var err error
	if s.useSSL {
		err = http.ServeTLS(s.listener, s.handler, s.sslCert, s.sslKey)
	} else {
		err = http.Serve(s.listener, s.handler)
	}
	if err == nil {
		log.Infof("FileHive server listening on %s\n", s.listener.Addr().String())
	}
	return err
}

func (s *FileHiveServer) newV1Router() *mux.Router {
	r := mux.NewRouter()
	r.Methods("OPTIONS")

	// Unauthenticated Routes
	r.HandleFunc("/api/v1/user", s.handlePOSTUser).Methods("POST")
	r.HandleFunc("/api/v1/user/{emailOrID}", s.handleGETUser).Methods("GET")
	r.HandleFunc("/api/v1/login", s.handlePOSTLogin).Methods("POST")
	r.HandleFunc("/api/v1/image/{filename}", s.handleGETImage).Methods("GET")
	r.HandleFunc("/api/v1/dataset/{id}", s.handleGETDataset).Methods("GET")
	r.HandleFunc("/api/v1/latest", s.handleGETRecent).Methods("GET")
	r.HandleFunc("/api/v1/trending", s.handleGETTrending).Methods("GET")
	r.HandleFunc("/api/v1/search", s.handleGETSearch).Methods("GET")
	r.HandleFunc("/api/v1/confirm", s.handleGETConfirm).Methods("GET")
	r.HandleFunc("/api/v1/passwordreset", s.handleGETPasswordReset).Methods("GET")
	r.HandleFunc("/api/v1/passwordreset", s.handlePOSTPasswordReset).Methods("POST")
	r.HandleFunc("/api/v1/checkresetcode", s.handleGETCheckResetCode).Methods("GET")

	if s.testMode {
		r.HandleFunc("/api/v1/generatecoins", s.handlePOSTGenerateCoins).Methods("POST")
	}

	// Authenticated Routes
	subRouter := r.PathPrefix("/api/v1").Subrouter()
	subRouter.Use(s.authenticationMiddleware)

	subRouter.HandleFunc("/logout", s.handlePOSTLogout).Methods("POST")
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
	subRouter.HandleFunc("/datasetdeal/{id}", s.handleGETDatasetDeal).Methods("GET")
	subRouter.HandleFunc("/purchase/{id}", s.handlePOSTPurchase).Methods("POST")
	subRouter.HandleFunc("/purchases", s.handleGETPurchases).Methods("GET")
	subRouter.HandleFunc("/purchased/{id}", s.handleGETPurchased).Methods("GET")
	subRouter.HandleFunc("/sales", s.handleGETSales).Methods("GET")
	subRouter.HandleFunc("/download/{cid}", s.handleGETDatasetFile).Methods("GET")

	return r
}

// Options represents the filehive server options.
type Options struct {
	JWTKey          []byte
	Domain          string
	UseSSL          bool
	FilecoinAddress string
	SSLCert         string
	SSLKey          string
	TestMode        bool
	MailgunKey      string
	MailDomain      string
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

func MailgunKey(mailgunKey string) Option {
	return func(o *Options) error {
		o.MailgunKey = mailgunKey
		return nil
	}
}

func MailDomain(mailDomain string) Option {
	return func(o *Options) error {
		o.MailDomain = mailDomain
		return nil
	}
}

func FilecoinAddress(filecoinAddress string) Option {
	return func(o *Options) error {
		o.FilecoinAddress = filecoinAddress
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

func passwordScore(pw string) (score int) {
	matchLower := regexp.MustCompile(`[a-z]`)
	matchUpper := regexp.MustCompile(`[A-Z]`)
	matchNumber := regexp.MustCompile(`[0-9]`)
	matchSpecial := regexp.MustCompile(`[\!\@\#\$\%\^\&\*\(\\\)\-_\=\+\,\.\?\/\:\;\{\}\[\]~]`)
	if len(pw) < 8 {
		return 0
	}
	if matchLower.MatchString(pw) {
		score++
	}
	if matchUpper.MatchString(pw) {
		score++
	}
	if matchNumber.MatchString(pw) {
		score++
	}
	if matchSpecial.MatchString(pw) {
		score++
	}
	return score
}
