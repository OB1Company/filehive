package main

import (
	"github.com/OB1Company/filehive/app"
	"github.com/OB1Company/filehive/repo"
	"github.com/jessevdk/go-flags"
	"github.com/op/go-logging"
	"net"
	"os"
	"os/signal"
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

	serverOpts := []app.Option{
		app.Domain(config.Domain),
	}
	if config.UseSSL {
		serverOpts = append(serverOpts, []app.Option{
			app.UseSSL(true),
			app.SSLCert(config.SSLCert),
			app.SSLKey(config.SSLKey),
		}...)
	}

	server, err := app.NewServer(listener, db, serverOpts...)
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
