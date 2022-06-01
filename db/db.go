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

// region: flat

type DbType int
type History []Statement

// endregion: flat
// region: interface

type canExecute interface {
	appendHistory(s *Statement)
	Config() *Config
	exec() interface{}
}

type hasHistory interface {
	appendHistory(s *Statement)
	Config() *Config
	History() History
	setHistory(*History)
}

// endregion: interface
// region: struct

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
	config  *Config
	conn    *sql.DB
	history History
}

type Statement struct {
	Args         []interface{}
	Db           *Db
	Err          error
	LastInsertId int64
	Result       sql.Result
	RowsAffected int64
	SQL          string
	Tx           *Tx
	Unprotected  bool
}

type Tx struct {
	db      *Db
	history History
	session *sql.Tx
}

// endregion: struct

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
	ErrInvalidExec            = errors.New("invalid interface.exec()")
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

	dbDefaultHistory int = 5
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
// region: db

// region: open

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

	db := Db{config: c}
	conn, e := sql.Open(Drivers[c.Type], c.DSN)
	if e != nil {
		log.Out(c.Logger, *c.Loglevel, e, c)
		return nil, e
	}
	db.conn = conn
	db.conn.SetMaxOpenConns(*c.MaxOpenConns)
	db.conn.SetMaxIdleConns(*c.MaxIdleConns)
	db.conn.SetConnMaxLifetime(*c.MaxLifetime)

	db.history = make(History, 0, *db.config.History)

	// endregion: connect
	// region: ping

	e = db.conn.Ping()
	if e != nil {
		log.Out(c.Logger, *c.Loglevel, e, c)
		db.conn.Close()
		return nil, e
	}
	log.Out(c.Logger, *c.Loglevel, fmt.Sprintf("connection is established to %s with %s driver", c.DSN, Drivers[c.Type]))

	// endregion: ping

	return &db, nil
}

func (c *Config) Open() (*Db, error) {
	return Open(c)
}

// endregion: open
// region: close

func Close(db *Db) error {
	e := db.Conn().Close()
	if e != nil {
		log.Out(db.Config().Logger, *db.Config().Loglevel, e, db)
		return e
	}
	conf := *db.Config()
	conf.Logger = nil
	log.Out(db.Config().Logger, *db.Config().Loglevel, "database connection closed", conf)
	return nil
}

func (db *Db) Close() error {
	return Close(db)
}

// endregion: close
// region: getters

func (db *Db) Config() *Config {
	return db.config
}

func (db *Db) Conn() *sql.DB {
	return db.conn
}

func (db *Db) exec() interface{} {
	return db.conn
}

func (db *Db) History() History {
	return db.history
}

// endregion: getters

// endregion: db
// region: tx

// region: begin

func Begin(db *Db) (*Tx, error) {
	s := &Statement{Db: db, SQL: "BEGIN"}

	session, e := db.Conn().Begin()
	if e != nil {
		s.Err = e
		db.appendHistory(s)
		log.Out(db.Config().Logger, *db.Config().Loglevel, e)
		return nil, e
	}

	tx := Tx{db: db, session: session}
	tx.history = make(History, 0, *db.Config().History)
	s.Tx = &tx
	tx.appendHistory(s)
	return &tx, nil
}

func (db *Db) Begin() (*Tx, error) {
	return Begin(db)
}

// endregion: begin
// region: commit

func Commit(tx *Tx) error {
	s := &Statement{Db: tx.Db(), Tx: tx, SQL: "COMMIT"}
	e := tx.Session().Commit()
	if e != nil {
		s.Err = e
		tx.appendHistory(s)
		log.Out(tx.Db().Config().Logger, *tx.Db().Config().Loglevel, e)
		return e
	}

	tx.appendHistory(s)
	return nil
}

func (tx *Tx) Commit() error {
	return Commit(tx)
}

// endregion: commit
// region: getters

func (tx *Tx) Config() *Config {
	return tx.Db().Config()
}

func (tx *Tx) Db() *Db {
	return tx.db
}

func (tx *Tx) exec() interface{} {
	return tx.session
}

func (tx *Tx) History() History {
	return tx.history
}

func (tx *Tx) Session() *sql.Tx {
	return tx.session
}

// endregion: getters

// endregion: tx
// region: history

func appendHistory(i hasHistory, s *Statement) {
	var h History = i.History()
	var l int = *i.Config().History - 1

	if len(h) < l {
		l = len(h)
	}
	if len(h) != 0 {
		h = h[:l]
	}

	h = append(History{*s}, h...)
	i.setHistory(&h)
}

func (db *Db) appendHistory(s *Statement) {
	appendHistory(db, s)
}

func (tx *Tx) appendHistory(s *Statement) {
	appendHistory(tx, s)
	appendHistory(tx.Db(), s)
}

func (s *Statement) appendHistory(i hasHistory) {
	i.appendHistory(s)
}

func (db *Db) setHistory(h *History) {
	db.history = *h
}

func (tx *Tx) setHistory(h *History) {
	tx.history = *h
}

// endregion: history
// region: exec

func Exec(i canExecute, s *Statement) error {
	log.Out(i.Config().Logger, *i.Config().Loglevel, MsgExecStatement, s.SQL, s.Args)

	// region: xss protection

	// TODO: check sting and []byte if content is valid JSON?

	if !s.Unprotected {
		for k, v := range s.Args {
			switch v.(type) {
			case string:
				s.Args[k] = template.HTMLEscaper(v)
			case []byte:
				s.Args[k] = template.HTMLEscaper(string(v.([]byte)))
			case sql.NullString:
				s.Args[k] = sql.NullString{
					String: template.HTMLEscaper(v.(sql.NullString).String),
					Valid:  v.(sql.NullString).Valid,
				}
			default:
				s.Args[k] = v
			}
		}
		s.SQL = template.HTMLEscaper(s.SQL)
		log.Out(i.Config().Logger, *i.Config().Loglevel, MsgExecStatementEscaped, *s)
	}

	// endregion: injection protection
	// region: execution

	s.Err = nil

	switch i.exec().(type) {
	case *sql.DB:
		s.Result, s.Err = i.exec().(*sql.DB).Exec(s.SQL, s.Args...)
		s.Db = i.(*Db)
	case *sql.Tx:
		s.Result, s.Err = i.exec().(*sql.Tx).Exec(s.SQL, s.Args...)
		s.Tx = i.(*Tx)
		s.Db = i.(*Tx).Db()
	default:
		s.Err = ErrInvalidExec
		i.appendHistory(s)
		log.Out(i.Config().Logger, *i.Config().Loglevel, s.Err)
		return s.Err
	}
	if s.Err != nil {
		s.Err = fmt.Errorf("%s: %s", ErrExecFailed, s.Err)
		i.appendHistory(s)
		log.Out(i.Config().Logger, *i.Config().Loglevel, s.Err)
		return s.Err
	}

	// endregion: execution
	// region: last id and rows affected

	if i.Config().Type != Postgres {
		s.LastInsertId, s.Err = s.Result.LastInsertId()
		if s.Err != nil {
			s.Err = fmt.Errorf("%s: %w", ErrExecLastIdFailed, s.Err)
			i.appendHistory(s)
			log.Out(i.Config().Logger, *i.Config().Loglevel, s.Err)
			return s.Err
		}
	}

	s.RowsAffected, s.Err = s.Result.RowsAffected()
	if s.Err != nil {
		s.Err = fmt.Errorf("%s: %w", ErrExecRowsAffectedFailed, s.Err)
		i.appendHistory(s)
		log.Out(i.Config().Logger, *i.Config().Loglevel, s.Err)
		return s.Err
	}

	// endregion: last id and rows affected

	i.appendHistory(s)
	return nil
}

func (db *Db) Exec(s *Statement) error {
	return Exec(db, s)
}

func (tx *Tx) Exec(s *Statement) error {
	return Exec(tx, s)
}

func (s *Statement) Exec(i canExecute) error {
	return Exec(i, s)
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
