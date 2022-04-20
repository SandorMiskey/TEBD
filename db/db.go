// region: packages

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/syslog"
	"time"

	// "github.com/go-sql-driver/mysql"
	"github.com/SandorMiskey/TEx-mysql-driver"
	_ "github.com/lib/pq"
)

// endregion: packages
// region: types

type DbType int

type Config struct {
	// mysql.config embeds these fields:
	//
	// User                    string            // Username
	// Passwd                  string            // Password (requires User)
	// Net                     string            // Network type
	// Addr                    string            // Network address (requires Net)
	// DBName                  string            // Database name
	// Params                  map[string]string // Connection parameters
	// Collation               string            // Connection collation
	// Loc                     *time.Location    // Location for time.Time values
	// MaxAllowedPacket        int               // Max packet size allowed
	// ServerPubKey            string            // Server public key name
	//
	// TLSConfig               string            // TLS configuration name
	//
	// Timeout                 time.Duration     // Dial timeout
	// ReadTimeout             time.Duration     // I/O read timeout
	// WriteTimeout            time.Duration     // I/O write timeout
	//
	// AllowAllFiles           bool // Allow all files to be used with LOAD DATA LOCAL INFILE
	// AllowCleartextPasswords bool // Allows the cleartext client side plugin
	// AllowNativePasswords    bool // Allows the native password authentication method
	// AllowOldPasswords       bool // Allows the old insecure password method
	// CheckConnLiveness       bool // Check connections for liveness before using them
	// ClientFoundRows         bool // Return number of matching rows instead of rows changed
	// ColumnsWithAlias        bool // Prepend table alias to column names
	// InterpolateParams       bool // Interpolate placeholders into query string
	// MultiStatements         bool // Allow multiple statements in one query
	// ParseTime               bool // Parse time values to time.Time
	// RejectReadOnly          bool // Reject read-only connections
	//
	Parsed mysql.Config

	Type DbType
	DSN  string

	MaxLifetime  *time.Duration
	MaxIdleConns *int
	MaxOpenConns *int

	Logger   *interface{}
	Loglevel *syslog.Priority
}

type Db struct {
	Config Config
	Conn   *sql.DB
}

// endregion: types
// region: constants

const (
	DbUndefined DbType = iota
	DbMySQL
	DbPQ
)

// endregion: constants
// region: defaults

// region: DbDefaults

var (

// dbDefaultHost         string          = "localhost"
// dbDefaultName         string          = "default"
// dbDefaultPort         int             = 3306
// dbDefaultMaxLifetime  time.Duration   = time.Minute * 3
// dbDefaultMaxIdleConns int             = 10
// dbDefaultMaxOpenConns int             = 10
// dbDefaultLogger       interface{}     = nil
// dbDefaultLoglevel     syslog.Priority = syslog.LOG_DAEMON
)

var DbDefaults = Config{
	// Config: mysql.Config{
	// 	User:      "",
	// 	Passwd:    "",
	// 	Net:       "",
	// 	Addr:      "localhost",
	// 	DBName:    "",
	// 	Params:    nil,
	// 	Collation: "utf8mb4_general_ci",
	// 	Loc:
	// },
	// Host:         &dbDefaultHost,
	// Name:         &dbDefaultName,
	// Password:     nil,
	// Port:         &dbDefaultPort,
	// Type:         DbMySQL,
	// User:         nil,
	// MaxLifetime:  &dbDefaultMaxLifetime,
	// MaxIdleConns: &dbDefaultMaxIdleConns,
	// MaxOpenConns: &dbDefaultMaxOpenConns,
	// Logger:       &dbDefaultLogger,
	// Loglevel:     &dbDefaultLoglevel,
}

// endregion: DbDefaults
// region: messages

var (
	ErrInvalidDbType     = errors.New("invalid db type")
	ErrNotImplementedYet = errors.New("not implemented yet")
	ErrTooManyParameters = errors.New("too many parameters")
)

// endregion: messages

// endregion: defaults
// region: open/close

func Open(cs ...Config) (*Db, error) {

	// region: prepare input

	// not sure if this is idiomatic, but this way you can call db.New() instead of db.New(db.DbDefaults)

	if len(cs) == 0 {
		cs = append(cs, DbDefaults)
	}
	if len(cs) > 1 {
		return nil, ErrTooManyParameters
	}
	c := cs[0]

	// endregion: prepare
	// region: check/set defaults

	// if c.Host == nil {
	// 	c.Host = DbDefaults.Host
	// }
	// if c.Name == nil {
	// 	c.Name = DbDefaults.Name
	// }
	// if c.Password == nil {
	// 	c.Password = DbDefaults.Password
	// }
	// if c.Port == nil {
	// 	c.Port = DbDefaults.Port
	// }
	// if c.Type == DbUndefined {
	// 	c.Type = DbDefaults.Type
	// }
	// if c.User == nil {
	// 	c.User = DbDefaults.User
	// }
	// if c.MaxLifetime == nil {
	// 	c.MaxLifetime = DbDefaults.MaxLifetime
	// }
	// if c.MaxIdleConns == nil {
	// 	c.MaxIdleConns = DbDefaults.MaxIdleConns
	// }
	// if c.MaxOpenConns == nil {
	// 	c.MaxOpenConns = DbDefaults.MaxOpenConns
	// }
	// if c.Logger == nil {
	// 	c.Logger = DbDefaults.Logger
	// }
	// if c.Loglevel == nil {
	// 	c.Loglevel = DbDefaults.Loglevel
	// }

	// endregion: check/set defaults
	// region: logging

	// TODO: implement

	// endregion: logging
	// region: connect

	db := Db{
		Config: c,
	}

	switch c.Type {
	case DbPQ:
		// TODO: logging
		return nil, ErrNotImplementedYet
	case DbMySQL:
		// TODO: implement w/ logging
		return &db, nil
	default:
		// TODO: logging
		return nil, fmt.Errorf("%s: %v", ErrInvalidDbType, c.Type)
	}

	// endregion: connect

}

// endregion: open/close
