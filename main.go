package main

import (
	"crypto/rand"
	"github.com/OB1Company/filehive/app"
	"github.com/OB1Company/filehive/fil"
	"github.com/OB1Company/filehive/repo"
	"github.com/jessevdk/go-flags"
	"github.com/op/go-logging"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path"
)

var log = logging.MustGetLogger("MAIN")

func main() {
	parser := flags.NewParser(&repo.Config{}, flags.Default)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	config, err := repo.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", config.Listen)
	if err != nil {
		log.Fatal(err)
	}

	dbOpts := []repo.Option{
		repo.Host(config.DBHost),
		repo.Dialect(config.DBDialect),
		repo.Username(config.DBUser),
		repo.Password(config.DBPass),
	}
	db, err := repo.NewDatabase(config.DataDir, dbOpts...)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: this will need to be set by a config option when powergate gets wired up.
	wbe, err := fil.NewPowergateWalletBackend()
	if err != nil {
		log.Fatalf("Powergate server is not available: %v", err)
	}

	if err := os.MkdirAll(path.Join(config.DataDir, "files"), os.ModePerm); err != nil {
		log.Fatal(err)
	}
	fbe, err := fil.NewPowergateBackend(path.Join(config.DataDir, "files"), config.PowergateToken, config.PowergateHost)
	if err != nil {
		log.Fatal(err)
	}

	key, err := loadJWTKey(config.DataDir)
	if err != nil {
		log.Fatal(err)
	}

	serverOpts := []app.Option{
		app.JWTKey(key),
		app.Domain(config.Domain),
	}
	if config.UseSSL {
		serverOpts = append(serverOpts, []app.Option{
			app.UseSSL(true),
			app.SSLCert(config.SSLCert),
			app.SSLKey(config.SSLKey),
		}...)
	}

	if config.MailgunKey != "" {
		serverOpts = append(serverOpts, []app.Option{
			app.MailgunKey(config.MailgunKey),
		}...)
	}

	server, err := app.NewServer(listener, db, config.StaticFileDir, wbe, fbe, serverOpts...)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	for sig := range c {
		if sig == os.Kill {
			log.Info("filehive killed")
			os.Exit(1)
		}

		if err := server.Close(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		log.Info("filehive stopping...")
		os.Exit(0)
	}
}

func loadJWTKey(dataDir string) ([]byte, error) {
	filename := path.Join(dataDir, "server.key")
	key, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		jwtKey := make([]byte, 32)
		rand.Read(jwtKey)
		if err := ioutil.WriteFile(filename, key, os.ModePerm); err != nil {
			return nil, err
		}
		return jwtKey, nil
	}
	return key, nil
}
