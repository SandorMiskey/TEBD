// region: packages

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/syslog"
	"strconv"
	"time"

	"github.com/SandorMiskey/TEx-kit/log"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// endregion: packages
// region: types

type DbType int

type Config struct {

	// region: dsn

	Addr   string            // Network address (requires Net)
	DBName string            // Database name
	DSN    string            //
	Net    string            // Network type
	Passwd string            // Password (requires User)
	Type   DbType            //
	User   string            // Username
	Params map[string]string // Connection parameters

	// endregion: dsn
	// region: connection

	MaxLifetime  *time.Duration
	MaxIdleConns *int
	MaxOpenConns *int

	// endregion: connection
	// region: logger

	Logger   interface{}
	Loglevel *syslog.Priority

	// endregion: logger

}

type Db struct {
	Config Config
	Conn   *sql.DB
	DSN    string
}

// endregion: types
// region: constants

const (
	DbUndefined DbType = iota
	DbMariaDB
	DbMySQL
	DbPostgres
	DbSQLite3
)

// endregion: constants
// region: messages

var (
	ErrInvalidDbType       = errors.New("invalid db type")
	errInvalidDSNAddr      = errors.New("invalid DSN: network address not terminated (missing closing brace)")
	errInvalidDSNNoSlash   = errors.New("invalid DSN: missing the slash separating the database name")
	errInvalidDSNUnescaped = errors.New("invalid DSN: did you forget to escape a param value?")
	ErrNotImplementedYet   = errors.New("not implemented yet")
	ErrTooManyParameters   = errors.New("too many parameters")

	// errInvalidDSNUnsafeCollation = errors.New("invalid DSN: interpolateParams can not be used with unsafe collations")
)

// endregion: messages
// region: defaults

var (
	dbDefaultsMaxAllowedPacket int           = 4 << 20
	dbDefaultsReadTimeout      time.Duration = time.Second * 60
	dbDefaultsWriteTimeout     time.Duration = time.Second * 60
	dbDefaultsTimeout          time.Duration = time.Second * 60

	dbDefaultMaxLifetime  time.Duration = time.Minute * 3
	dbDefaultMaxIdleConns int           = 10
	dbDefaultMaxOpenConns int           = 10

	dbDefaultLoglevel syslog.Priority = log.LOG_DEBUG
)

var DbDefaults = Config{
	Addr:   "localhost",
	DBName: "tex",
	DSN:    "",
	Net:    "tcp",
	Passwd: "",
	Type:   DbMySQL,
	User:   "",
	Params: map[string]string{
		"allowNativePasswords": "true",                                   // Allows the native password authentication method
		"checkConnLiveness":    "true",                                   // Check connections for liveness before using them
		"collation":            "utf8_general_ci",                        // Connection collation
		"loc":                  time.UTC.String(),                        // Location for time.Time values
		"maxAllowedPacket":     strconv.Itoa(dbDefaultsMaxAllowedPacket), // Max packet size allowed
		// "allowAllFiles":           "false",                         // Allow all files to be used with LOAD DATA LOCAL INFILE
		// "allowCleartextPasswords": "false",                         // Allows the cleartext client side plugin
		// "allowOldPasswords":       "false",                         // Allows the old insecure password method
		// "clientFoundRows":         "false",                         // Return number of matching rows instead of rows changed
		// "columnsWithAlias":        "false",                         // Prepend table alias to column names
		// "interpolateParams":       "false",                         // Interpolate placeholders into query string
		// "multiStatements":         "false",                         // Allow multiple statements in one query
		// "parseTime":               "false",                         // Parse time values to time.Time
		// "readTimeout":             dbDefaultsReadTimeout.String(),  // I/O read timeout
		// "rejectReadOnly":          "false",                         // Reject read-only connections
		// "serverPubKey":            "",                              // Server public key name
		// "timeout":                 dbDefaultsTimeout.String(),      // Dial timeout
		// "tls":                     "",                              // TLS configuration name
		// "writeTimeout":            dbDefaultsWriteTimeout.String(), // I/O read timeout
	},

	MaxLifetime:  &dbDefaultMaxLifetime,
	MaxIdleConns: &dbDefaultMaxIdleConns,
	MaxOpenConns: &dbDefaultMaxOpenConns,

	Logger:   nil,
	Loglevel: &dbDefaultLoglevel,
}

func SetDefaults(c *Config) {
	if c.Addr == "" {
		c.Addr = DbDefaults.Addr
	}
	if c.DBName == "" {
		c.DBName = DbDefaults.DBName
	}
	if c.DSN == "" {
		c.DSN = DbDefaults.DSN
	}
	if c.Net == "" {
		c.Net = DbDefaults.Net
	}
	if c.Params == nil {
		c.Params = DbDefaults.Params
	}
	if c.Passwd == "" {
		c.Passwd = DbDefaults.Passwd
	}
	if c.Type == DbUndefined {
		c.Type = DbDefaults.Type
	}
	if c.User == "" {
		c.User = DbDefaults.User
	}

	if c.MaxLifetime == nil {
		c.MaxLifetime = DbDefaults.MaxLifetime
	}
	if c.MaxIdleConns == nil {
		c.MaxIdleConns = DbDefaults.MaxIdleConns
	}
	if c.MaxOpenConns == nil {
		c.MaxOpenConns = DbDefaults.MaxOpenConns
	}

	if c.Logger == nil {
		c.Logger = DbDefaults.Logger
	}
	if c.Loglevel == nil {
		c.Loglevel = DbDefaults.Loglevel
	}
}

func (c *Config) SetDefaults() {
	SetDefaults(c)
}

// endregion: defaults
// region: open/close

func Open(cs ...Config) (*Db, error) {

	// region: prepare input

	// not sure if this is idiomatic, but this way you can call db.Open() instead of db.Open(db.DbDefaults)

	if len(cs) == 0 {
		cs = append(cs, DbDefaults)
	}
	if len(cs) > 1 {
		return nil, ErrTooManyParameters
	}
	c := cs[0]

	// if c.DSN == nil {
	// 	// TODO: SetDefaults and FormatDSN
	// } else {
	// 	// TODO: SetDefaults and ParseDSN
	// }
	// log.Out(c.Logger, *c.Loglevel, c)

	// endregion: prepare
	// region: connect

	db := Db{
		Config: c,
	}

	switch c.Type {
	case DbPostgres:
		return nil, ErrNotImplementedYet
	case DbMariaDB, DbMySQL:
		// TODO: implement w/ logging
		return &db, nil
	case DbSQLite3:
		return nil, ErrNotImplementedYet
	default:
		return nil, fmt.Errorf("%s: %v", ErrInvalidDbType, c.Type)
	}

	// TODO: func (c *Config) Open() (*Db, error)
	// TODO: func close

	// endregion: connect

}

// endregion: open/close
