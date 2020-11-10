package app

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/OB1Company/filehive/repo"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"net"
	"net/http"
)

var log = logging.MustGetLogger("APP")

// FileHiveServer is the web server used to serve the FileHive application.
type FileHiveServer struct {
	db       *repo.Database
	listener net.Listener
	handler  http.Handler
	shutdown chan struct{}

	useSSL  bool
	sslCert string
	sslKey  string
}

// NewServer instantiates a new FileHiveServer with the provided options.
func NewServer(listener net.Listener, db *repo.Database, opts ...Option) (*FileHiveServer, error) {
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

	var (
		s = &FileHiveServer{
			db:       db,
			listener: listener,
			useSSL:   options.UseSSL,
			sslCert:  options.SSLCert,
			sslKey:   options.SSLKey,
			shutdown: make(chan struct{}),
		}
		topMux = http.NewServeMux()
	)

	r := s.newV1Router()

	csrfKey := make([]byte, 32)
	rand.Read(csrfKey)

	csrfMiddleware := csrf.Protect(csrfKey)
	r.Use(csrfMiddleware)
	r.Use(s.setCSRFHeaderMiddleware)

	topMux.Handle("/v1/", r)

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
	r.HandleFunc("/v1/user", s.handlePOSTUser).Methods("POST")
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

// Options represents the filehive server options.
type Options struct {
	UseSSL  bool
	SSLCert string
	SSLKey  string
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
