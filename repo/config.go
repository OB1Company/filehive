package repo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/natefinch/lumberjack"
	"github.com/op/go-logging"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	defaultConfigFilename = "filehive.conf"
	defaultLogFilename    = "filehive.log"
)

var log = logging.MustGetLogger("REPO")

var (
	DefaultHomeDir    = AppDataDir("filehive", false)
	defaultConfigFile = filepath.Join(DefaultHomeDir, defaultConfigFilename)

	fileLogFormat   = logging.MustStringFormatter(`%{time:2006-01-02 T15:04:05.000} [%{level}] [%{module}] %{message}`)
	stdoutLogFormat = logging.MustStringFormatter(`%{color:reset}%{color}%{time:15:04:05} [%{level}] [%{module}] %{message}`)
	LogLevelMap     = map[string]logging.Level{
		"debug":    logging.DEBUG,
		"info":     logging.INFO,
		"notice":   logging.NOTICE,
		"warning":  logging.WARNING,
		"error":    logging.ERROR,
		"critical": logging.CRITICAL,
	}
)

// Config defines the configuration options for the crawler.
//
// See LoadConfig for details on the configuration load process.
type Config struct {
	ShowVersion bool   `short:"v" long:"version" description:"Display version information and exit"`
	ConfigFile  string `short:"C" long:"configfile" description:"Path to configuration file"`
	DataDir     string `short:"d" long:"datadir" description:"Directory to store data"`
	LogDir      string `long:"logdir" description:"Directory to log output."`
	LogLevel    string `long:"loglevel" description:"Set the logging level [debug, info, notice, warning, error, critical]." default:"info"`

	Listen  string `short:"l" long:"listen" description:"The interface:port for the app server to bind to." default:"0.0.0.0:8080"`
	UseSSL  bool   `long:"usessl" description:"Set to true if you want to use SSL with the server."`
	SSLCert string `long:"sslcert" description:"Path to the SSL certificate."`
	SSLKey  string `long:"sslkey" description:"Path to the SSL key."`

	DBDialect string `long:"dbdialect" description:"The type of database to use [sqlite3, mysql, postgress]" default:"sqlite3"`
	DBHost    string `long:"dbhost" description:"The host:post location of the database."`
	DBUser    string `long:"dbuser" description:"The database username"`
	DBPass    string `long:"dbpass" description:"The database password"`
}

// LoadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
//
// The above results in OpenBazaar functioning properly without any config settings
// while still allowing the user to override settings with config files and
// command line options.  Command line options always take precedence.
func LoadConfig() (*Config, error) {
	// Default config.
	cfg := Config{
		DataDir:    DefaultHomeDir,
		ConfigFile: defaultConfigFile,
	}

	// Pre-parse the command line options to see if an alternative config
	// file or the version flag was specified.  Any errors aside from the
	// help message error can be ignored here since they will be caught by
	// the final parse below.
	preCfg := cfg
	preParser := flags.NewParser(&cfg, flags.HelpFlag)
	_, err := preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return nil, err
		}
	}
	if cfg.DataDir != "" {
		preCfg.ConfigFile = filepath.Join(cfg.DataDir, defaultConfigFilename)
	}

	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("Use %s -h to show usage", appName)
	if preCfg.ShowVersion {
		fmt.Println(appName, "version", VersionString())
		os.Exit(0)
	}

	// Load additional config from file.
	var configFileError error
	parser := flags.NewParser(&cfg, flags.Default)
	if _, err := os.Stat(preCfg.ConfigFile); os.IsNotExist(err) {
		err := createDefaultConfigFile(preCfg.ConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating a "+
				"default config file: %v\n", err)
		}
	}

	err = flags.NewIniParser(parser).ParseFile(preCfg.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			fmt.Fprintf(os.Stderr, "Error parsing config "+
				"file: %v\n", err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, err
		}
		configFileError = err
	}

	_, ok := LogLevelMap[strings.ToLower(cfg.LogLevel)]
	if !ok {
		return nil, errors.New("invalid log level")
	}

	// Warn about missing config file only after all other configuration is
	// done.  This prevents the warning on help messages and invalid
	// options.  Note this should go directly before the return.
	if configFileError != nil {
		log.Errorf("%v", configFileError)
	}
	setupLogging(cfg.LogDir, cfg.LogLevel)

	cfg.DataDir = cleanAndExpandPath(cfg.DataDir)
	if cfg.LogDir == "" {
		cfg.LogDir = cleanAndExpandPath(path.Join(cfg.DataDir, "logs"))
	}

	return &cfg, nil
}

// createDefaultConfig copies the sample-filehive.conf content to the given destination path,
// and populates it with some randomly generated RPC username and password.
func createDefaultConfigFile(destinationPath string) error {
	// Create the destination directory if it does not exists
	err := os.MkdirAll(filepath.Dir(destinationPath), 0700)
	if err != nil {
		return err
	}

	sampleBytes, err := Asset("sample-filehive.conf")
	if err != nil {
		return err
	}
	src := bytes.NewReader(sampleBytes)

	dest, err := os.OpenFile(destinationPath,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer dest.Close()

	// We copy every line from the sample config file to the destination.
	reader := bufio.NewReader(src)
	for err != io.EOF {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if _, err := dest.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}

// cleanAndExpandPath expands environment variables and leading ~ in the
// passed path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// Expand initial ~ to OS specific home directory.
	if strings.HasPrefix(path, "~") {
		homeDir := filepath.Dir(DefaultHomeDir)
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%,
	// but they variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}

func setupLogging(logDir, logLevel string) {
	backendStdout := logging.NewLogBackend(os.Stdout, "", 0)
	backendStdoutFormatter := logging.NewBackendFormatter(backendStdout, stdoutLogFormat)
	if logDir != "" {
		rotator := &lumberjack.Logger{
			Filename:   path.Join(logDir, defaultLogFilename),
			MaxSize:    10, // Megabytes
			MaxBackups: 3,
			MaxAge:     30, // Days
		}

		backendFile := logging.NewLogBackend(rotator, "", 0)
		backendFileFormatter := logging.NewBackendFormatter(backendFile, fileLogFormat)
		logging.SetBackend(backendStdoutFormatter, backendFileFormatter)
	} else {
		logging.SetBackend(backendStdoutFormatter)
	}
	logging.SetLevel(LogLevelMap[strings.ToLower(logLevel)], "")
}
