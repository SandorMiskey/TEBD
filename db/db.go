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
	Config *Config
	Conn   *sql.DB
}

// endregion: types
// region: (pseudo-)constants

const (
	DbUndefined DbType = iota
	DbMariaDB
	DbMySQL
	DbPostgres
	DbSQLite3
)

var DbDrivers = map[DbType]string{
	DbMariaDB:  "mysql",
	DbMySQL:    "mysql",
	DbPostgres: "postgres",
	DbSQLite3:  "sqlite3",
}

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
	dbDefaultAllowedPacket int           = 4 << 20
	dbDefaultReadTimeout   time.Duration = time.Second * 60
	dbDefaultWriteTimeout  time.Duration = time.Second * 60
	dbDefaultTimeout       time.Duration = time.Second * 60

	dbDefaultMaxLifetime  time.Duration = time.Minute * 3
	dbDefaultMaxIdleConns int           = 10
	dbDefaultMaxOpenConns int           = 10

	dbDefaultLoglevel syslog.Priority = log.LOG_DEBUG
)

var DbDefaults = DbDefaultsMySQL

var DbDefaultsMySQL = Config{
	Addr:   "localhost",
	DBName: "tex",
	DSN:    "",
	Net:    "tcp",
	Passwd: "",
	Type:   DbMySQL,
	User:   "",
	Params: map[string]string{
		"allowNativePasswords": "true",                               // Allows the native password authentication method
		"checkConnLiveness":    "true",                               // Check connections for liveness before using them
		"collation":            "utf8_general_ci",                    // Connection collation
		"loc":                  time.UTC.String(),                    // Location for time.Time values
		"maxAllowedPacket":     strconv.Itoa(dbDefaultAllowedPacket), // Max packet size allowed
		// "allowAllFiles":           "false",                                      // Allow all files to be used with LOAD DATA LOCAL INFILE
		// "allowCleartextPasswords": "false",                                      // Allows the cleartext client side plugin
		// "allowOldPasswords":       "false",                                      // Allows the old insecure password method
		// "clientFoundRows":         "false",                                      // Return number of matching rows instead of rows changed
		// "columnsWithAlias":        "false",                                      // Prepend table alias to column names
		// "interpolateParams":       "false",                                      // Interpolate placeholders into query string
		// "multiStatements":         "false",                                      // Allow multiple statements in one query
		// "parseTime":               "false",                                      // Parse time values to time.Time
		// "readTimeout":             dbDefaultMySQLReadTimeout.String(),           // I/O read timeout
		// "rejectReadOnly":          "false",                                      // Reject read-only connections
		// "serverPubKey":            "",                                           // Server public key name
		// "timeout":                 dbDefaultMySQLTimeout.String(),               // Dial timeout
		// "tls":                     "",                                           // TLS configuration name
		// "writeTimeout":            dbDefaultMySQLWriteTimeout.String(),          // I/O read timeout
	},

	MaxLifetime:  &dbDefaultMaxLifetime,
	MaxIdleConns: &dbDefaultMaxIdleConns,
	MaxOpenConns: &dbDefaultMaxOpenConns,

	Logger:   nil,
	Loglevel: &dbDefaultLoglevel,
}

var DbDefaultsPostgres = Config{
	Addr:   "localhost", // The host to connect to. Values that start with / are for unix domain sockets
	DBName: "tex",       // The name of the database to connect to
	DSN:    "",
	Passwd: "",
	Type:   DbPostgres,
	User:   "",
	Params: map[string]string{
		"sslmode":                   "disable", // Whether or not to use SSL (default is require, this is not the default for libpq)
		"application_name":          "foo",     // Specifies a value for the application_name configuration parameter
		"fallback_application_name": "foo",     // An application_name to fall back to if one isn't provided.
		"connect_timeout":           "20",      // Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.
		"sslcert":                   "",        // Cert file location. The file must contain PEM encoded data.
		"sslkey":                    "",        // Key file location. The file must contain PEM encoded data
		"sslrootcert":               "",        // The location of the root certificate file. The file must contain PEM encoded data
		// Valid values for sslmode are:
		// * disable - No SSL
		// * require - Always SSL (skip verification)
		// * verify-ca - Always SSL (verify that the certificate presented by the
		//   server was signed by a trusted CA)
		// * verify-full - Always SSL (verify that the certification presented by
		//   the server was signed by a trusted CA and the server host name
		//   matches the one in the certificate)
		//
		// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
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

func Open(cs ...*Config) (*Db, error) {

	// region: prepare input

	// not sure if this is idiomatic, but this way you can call db.Open() instead of db.Open(db.DbDefaults)

	defaults := DbDefaults
	if len(cs) == 0 {
		cs = append(cs, &defaults)
	}
	if len(cs) > 1 {
		return nil, ErrTooManyParameters
	}
	c := cs[0]

	if c.DSN == "" {
		c.SetDefaults()
		c.FormatDSN()
	} else {
		c.ParseDSN()
		c.SetDefaults()
	}
	log.Out(c.Logger, *c.Loglevel, c)

	// endregion: prepare
	// region: connect

	if DbDrivers[c.Type] == "" {
		return nil, fmt.Errorf("%s: %v", ErrInvalidDbType, c.Type)
	}

	db := Db{Config: c}
	conn, e := sql.Open(DbDrivers[c.Type], c.DSN)
	if e != nil {
		log.Out(c.Logger, *c.Loglevel, e)
		return nil, e
	}
	db.Conn = conn
	log.Out(c.Logger, *c.Loglevel, fmt.Sprintf("%s is connected with %s", c.DSN, DbDrivers[c.Type]))

	db.Conn.SetMaxOpenConns(*c.MaxOpenConns)
	db.Conn.SetMaxIdleConns(*c.MaxIdleConns)
	db.Conn.SetConnMaxLifetime(*c.MaxLifetime)

	// endregion: connect
	// region: ping

	e = db.Conn.Ping()
	if e != nil {
		log.Out(c.Logger, *c.Loglevel, e)
		db.Conn.Close()
		return nil, e
	}

	// endregion: ping

	return &db, nil
}

func (c *Config) Open() (*Db, error) {
	return Open(c)
}

// endregion: open/close
