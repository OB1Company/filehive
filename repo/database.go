package repo

import (
	"errors"
	"fmt"
	"github.com/OB1Company/filehive/repo/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	glog "log"
	"os"
	"path"
	"strings"
	"sync"
)

const (
	dbName = "filehive"
)

var silentLogger = logger.New(
	glog.New(os.Stdout, "\r\n", glog.LstdFlags),
	logger.Config{
		LogLevel: logger.Silent,
	},
)

// Database is a mutex wrapper around a GORM db.
type Database struct {
	db  *gorm.DB
	mtx sync.Mutex
}

// NewDatabase returns a new database with the given options.
// Sqlite3, Mysql, and Postgress is supported.
func NewDatabase(dataDir string, opts ...Option) (*Database, error) {
	options := Options{
		Host:    "localhost",
		Dialect: "sqlite3",
	}
	if err := options.Apply(opts...); err != nil {
		return nil, err
	}

	dbPath := path.Join(dataDir, dbName)
	var dialector gorm.Dialector

	switch strings.ToLower(options.Dialect) {
	case "memory":
		dbPath = ":memory:"
		dialector = sqlite.Open(dbPath)
	case "mysql":
		dbPath = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True", options.User, options.Password, options.Host, options.Port, dbName)
		dialector = mysql.Open(dbPath)
	case "postgress":
		dbPath = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", options.Host, options.Port, options.User, dbName, options.Password)
		dialector = postgres.Open(dbPath)
	case "sqlite3":
		dbPath = dbPath + ".db"
		dialector = sqlite.Open(dbPath)
		break
	default:
		return nil, errors.New("unknown database dialect")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            silentLogger,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Dataset{}, &models.Purchase{}, &models.Click{}); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// View is used for read access to the db. Reads are made
// inside and open transaction.
func (d *Database) View(fn func(db *gorm.DB) error) error {
	if err := fn(d.db); err != nil {
		return err
	}
	return nil
}

// Update is used for write access to the db. Updates are made
// inside an open transaction.
func (d *Database) Update(fn func(db *gorm.DB) error) error {
	tx := d.db.Begin()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// Options represents the database options.
type Options struct {
	Host     string
	Port     uint
	Dialect  string
	User     string
	Password string
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

// Host option allows you to set the host for mysql or postgress dbs.
func Host(host string) Option {
	return func(o *Options) error {
		o.Host = host
		return nil
	}
}

// Port option sets the port for mysql or postgress dbs.
func Port(port uint) Option {
	return func(o *Options) error {
		o.Port = port
		return nil
	}
}

// Dialect sets the database type...sqlite3, mysql, postress.
func Dialect(dialect string) Option {
	return func(o *Options) error {
		o.Dialect = dialect
		return nil
	}
}

// Password is the password for the mysql or postgress dbs.
func Password(pw string) Option {
	return func(o *Options) error {
		o.Password = pw
		return nil
	}
}

// Username is the username for the mysql or postgress dbs.
func Username(user string) Option {
	return func(o *Options) error {
		o.User = user
		return nil
	}
}
