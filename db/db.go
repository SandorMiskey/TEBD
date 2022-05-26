// region: packages

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/syslog"
	"strconv"
	"text/template"
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
	// region: internaL

	Logger   interface{}
	Loglevel *syslog.Priority

	History *int

	// endregion: internal

}

type Db struct {
	Config  *Config
	Conn    *sql.DB
	History []*Statement
}

type Statement struct {
	Args         []interface{}
	RowsAffected int64
	Err          error
	LastInsertId int64
	Unprotected  bool
	SQL          string
	Result       interface{}
}

// endregion: types
// region: (pseudo-)constants

const (
	Undefined DbType = iota
	MariaDB
	MySQL
	Postgres
	SQLite3
)

var Drivers = map[DbType]string{
	MariaDB:  "mysql",
	MySQL:    "mysql",
	Postgres: "postgres",
	SQLite3:  "sqlite3",
}

// endregion: constants
// region: messages

var (
	ErrExecFailed             = errors.New("db.Exec() failed")
	ErrExecLastIdFailed       = errors.New("db.Exec was successfull but getting LastInsertId failed")
	ErrExecRowsAffectedFailed = errors.New("db.Exec was successfull but getting RowsAffected failed")
	ErrInvalidDbType          = errors.New("invalid db type")
	ErrInvalidDSNAddr         = errors.New("invalid DSN: network address not terminated (missing closing brace)")
	ErrInvalidDSNNoSlash      = errors.New("invalid DSN: missing the slash separating the database name")
	ErrInvalidDSNUnescaped    = errors.New("invalid DSN: did you forget to escape a param value?")
	ErrNotImplementedYet      = errors.New("not implemented yet")
	ErrTooManyParameters      = errors.New("too many parameters")
	// ErrInvalidDSNUnsafeCollation = errors.New("invalid DSN: interpolateParams can not be used with unsafe collations")

	MsgConnClosed           = "connection closed"
	MsgConnEstablished      = "connection established"
	MsgExecStatement        = "db.Exec() statement"
	MsgExecStatementEscaped = "db.Exec() statement escaped"
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

	dbDefaultHistory int = 10
)

var Defaults = DefaultsMySQL

var DefaultsMySQL = Config{
	Addr:   "localhost",
	DBName: "tex",
	DSN:    "",
	Net:    "tcp",
	Passwd: "",
	Type:   MySQL,
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

	History: &dbDefaultHistory,
}

var DefaultsPostgres = Config{
	Addr:   "localhost", // The host to connect to. Values that start with / are for unix domain sockets
	DBName: "tex",       // The name of the database to connect to
	DSN:    "",
	Passwd: "",
	Type:   Postgres,
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

	History: &dbDefaultHistory,
}

func SetDefaults(c *Config) {
	if c.Addr == "" {
		c.Addr = Defaults.Addr
	}
	if c.DBName == "" {
		c.DBName = Defaults.DBName
	}
	if c.DSN == "" {
		c.DSN = Defaults.DSN
	}
	if c.Net == "" {
		c.Net = Defaults.Net
	}
	if c.Params == nil {
		c.Params = Defaults.Params
	}
	if c.Passwd == "" {
		c.Passwd = Defaults.Passwd
	}
	if c.Type == Undefined {
		c.Type = Defaults.Type
	}
	if c.User == "" {
		c.User = Defaults.User
	}

	if c.MaxLifetime == nil {
		c.MaxLifetime = Defaults.MaxLifetime
	}
	if c.MaxIdleConns == nil {
		c.MaxIdleConns = Defaults.MaxIdleConns
	}
	if c.MaxOpenConns == nil {
		c.MaxOpenConns = Defaults.MaxOpenConns
	}

	if c.Logger == nil {
		c.Logger = Defaults.Logger
	}
	if c.Loglevel == nil {
		c.Loglevel = Defaults.Loglevel
	}

	if c.History == nil {
		c.History = Defaults.History
	}
}

func (c *Config) SetDefaults() {
	SetDefaults(c)
}

// endregion: defaults
// region: logger

// func (db *Db) Logger(n ...interface{}) {
// 	log.Out(db.Config.Logger, *db.Config.Loglevel, n...)
// }

// endregion: logger
// region: open/close

func Open(cs ...*Config) (*Db, error) {

	// region: prepare input

	// not sure if this is idiomatic, but this way you can call db.Open() instead of db.Open(db.DbDefaults)

	defaults := Defaults
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

	// endregion: prepare
	// region: connect

	if Drivers[c.Type] == "" {
		return nil, fmt.Errorf("%s: %v", ErrInvalidDbType, c.Type)
	}

	db := Db{Config: c}
	conn, e := sql.Open(Drivers[c.Type], c.DSN)
	if e != nil {
		log.Out(c.Logger, *c.Loglevel, e, c)
		return nil, e
	}
	db.Conn = conn
	db.Conn.SetMaxOpenConns(*c.MaxOpenConns)
	db.Conn.SetMaxIdleConns(*c.MaxIdleConns)
	db.Conn.SetConnMaxLifetime(*c.MaxLifetime)

	db.History = make([]*Statement, 0, *db.Config.History)

	// endregion: connect
	// region: ping

	e = db.Conn.Ping()
	if e != nil {
		log.Out(c.Logger, *c.Loglevel, e, c)
		db.Conn.Close()
		return nil, e
	}
	log.Out(c.Logger, *c.Loglevel, fmt.Sprintf("connection is established to %s with %s driver", c.DSN, Drivers[c.Type]))

	// endregion: ping

	return &db, nil
}

func (c *Config) Open() (*Db, error) {
	return Open(c)
}

func Close(db *Db) error {
	e := db.Conn.Close()
	if e != nil {
		log.Out(db.Config.Logger, *db.Config.Loglevel, e, db)
		return e
	}
	conf := *db.Config
	conf.Logger = nil
	log.Out(db.Config.Logger, *db.Config.Loglevel, "database connection closed", conf)
	return nil
}

func (db *Db) Close() error {
	return Close(db)
}

// endregion: open/close
// region: history

func (db *Db) HistoryAppend(s *Statement) {
	x := *db.Config.History
	if len(db.History) < *db.Config.History {
		x = len(db.History)
	}
	if len(db.History) != 0 {
		db.History = db.History[:x]
	}
	db.History = append([]*Statement{s}, db.History...)
}

// endregion: history
// region: exec

func Exec(db *Db, stmnt Statement) Statement {
	log.Out(db.Config.Logger, *db.Config.Loglevel, MsgExecStatement, stmnt)

	// region: xss protection

	// TODO: check sting and []byte if content is valid JSON?

	if !stmnt.Unprotected {
		for k, v := range stmnt.Args {
			switch v.(type) {
			case string:
				stmnt.Args[k] = template.HTMLEscaper(v)
			case []byte:
				stmnt.Args[k] = template.HTMLEscaper(string(v.([]byte)))
			case sql.NullString:
				stmnt.Args[k] = sql.NullString{
					String: template.HTMLEscaper(v.(sql.NullString).String),
					Valid:  v.(sql.NullString).Valid,
				}
			default:
				stmnt.Args[k] = v
			}
		}
		stmnt.SQL = template.HTMLEscaper(stmnt.SQL)
		log.Out(db.Config.Logger, *db.Config.Loglevel, MsgExecStatementEscaped, stmnt)
	}

	// endregion: injection protection
	// region: actual

	var err error
	stmnt.Result, err = db.Conn.Exec(stmnt.SQL, stmnt.Args...)
	if err != nil {
		stmnt.Err = fmt.Errorf("%s: %w", ErrExecFailed, err)
		log.Out(db.Config.Logger, *db.Config.Loglevel, stmnt.Err)
		return stmnt
	}

	// endregion: actual
	// region: last id and rows affected

	if db.Config.Type != Postgres {
		stmnt.LastInsertId, err = stmnt.Result.(sql.Result).LastInsertId()
		if err != nil {
			stmnt.Err = fmt.Errorf("%s: %w", ErrExecLastIdFailed, err)
			log.Out(db.Config.Logger, *db.Config.Loglevel, stmnt.Err)
			return stmnt
		}
	}

	stmnt.RowsAffected, err = stmnt.Result.(sql.Result).RowsAffected()
	if err != nil {
		stmnt.Err = fmt.Errorf("%s: %w", ErrExecRowsAffectedFailed, err)
		log.Out(db.Config.Logger, *db.Config.Loglevel, stmnt.Err)
		return stmnt
	}

	// endregion: last id and rows affected

	db.HistoryAppend(&stmnt)
	return stmnt
}

func (db *Db) Exec(stmnt Statement) Statement {
	return Exec(db, stmnt)
}

func (stmnt *Statement) Exec(db *Db) {
	statement := Exec(db, *stmnt)
	stmnt.Err = statement.Err
	stmnt.RowsAffected = statement.RowsAffected
	stmnt.LastInsertId = statement.LastInsertId
	stmnt.Result = statement.Result
	db.History[0] = stmnt
}

// endregion: exec
// region: query

func Query(db *Db, st Statement) Statement {

	return st
}

func (db *Db) Query(st Statement) Statement {
	return Query(db, st)
}

func (st *Statement) Query(db *Db) {
	statement := Query(db, *st)
	st.Err = statement.Err
	st.Result = statement.Result
}

// endregion: query
