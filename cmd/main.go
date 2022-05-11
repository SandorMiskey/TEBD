// region: packages

package main

import (
	"fmt"
	"log/syslog"
	"os"
	"time"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/db"
	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"
)

// endregion: packages
// region: global variables

var (
	Config cfg.Config
	Logger log.Logger
	Db     *db.Db
)

// endregion: globals

func main() {

	// region: config and cli flags

	Config = *cfg.NewConfig(os.Args[0])
	fs := Config.NewFlagSet(os.Args[0])
	fs.Entries = map[string]cfg.Entry{
		"bool":     {Desc: "bool description", Type: "bool", Def: true},
		"duration": {Desc: "duration description", Type: "time.Duration", Def: time.Duration(66000)},
		"float64":  {Desc: "float64 desc", Type: "float64", Def: 77.7},
		"int":      {Desc: "int description", Type: "int", Def: 99},

		"dbUser":        {Desc: "Database user", Type: "string", Def: ""},
		"dbPasswd":      {Desc: "Database password", Type: "string", Def: ""},
		"dbPasswd_file": {Desc: "Database password file", Type: "string", Def: ""},
		"dbName":        {Desc: "Database name", Type: "string", Def: "tex"},
	}

	err := fs.ParseCopy()
	if err != nil {
		panic(err)
	}

	// endregion: cli flags
	// region: logger

	var spewEncoder log.Encoder = func(c *log.Ch, n ...interface{}) (s string, e error) {

		// prefix with severity label, if needed
		if severity, ok := n[0].(syslog.Priority); ok {
			labels := *c.Config.SeverityLabels
			label := labels[severity]
			s = label + s
			_, n = n[0], n[1:]
		}

		// encode
		s = "\n"
		for k, v := range n {
			s = fmt.Sprintf("%s%s%d%s%s", s, *c.Config.Delimiter, k, *c.Config.Delimiter, spew.Sdump(v))
		}
		// s = strings.Replace(s, *c.Config.Delimiter, "", 1)
		// s = strings.TrimSuffix(s, "\n")

		// done
		return s, nil
	}
	severity := log.LOG_INFO
	Logger = *log.NewLogger()
	defer Logger.Close()
	_, _ = Logger.NewCh(log.ChConfig{Type: log.ChSyslog})
	_, _ = Logger.NewCh(log.ChConfig{Encoder: &spewEncoder, Severity: &severity})
	// lfc, _ := Logger.NewCh()
	// _ = lfc.Out(*log.ChDefaults.Mark)                               // write to identified channel
	// _ = lfc.Out(log.LOG_EMERG, "entry", "with", "severity")         // write to identified channel with severity
	// _ = log.Out(lfc, log.LOG_CRIT, "entry", "with", "severity")     // write to identified channel with severity
	// _ = Logger.Ch[0].Out(*log.ChDefaults.Mark, "bar", 1, 1.1, true) // write directly to the first channel
	// _ = Logger.Out(*log.ChDefaults.Mark)                            // write to all channels
	// _ = Logger.Out(log.LOG_ALERT, *log.ChDefaults.Mark)             // write to all channels with severity
	// _ = log.Out(&Logger, log.LOG_EMERG, "foobar")                   // write to all logger channels with severity
	// _ = log.Out(nil, log.LOG_EMERG, "quux")                         // write to nowhere

	// endregion: logger
	// region: db

	// region: config, defaults, dsn

	dbConfig := db.Config{
		User:   Config.Entries["dbUser"].Value.(string),
		Passwd: Config.Entries["dbPasswd"].Value.(string),
		DBName: Config.Entries["dbName"].Value.(string),
		Logger: Logger,
	}
	// dbConfig.SetDefaults() // or db.SetDefaults(&dbConfig) is also available
	// dbConfig.FormatDSN()   // or db.FormatDSN(&dbConfig)
	// _ = dbConfig.ParseDSN("user:pass@tcp(host)/dbname?allowNativePasswords=true&checkConnLiveness=true&collation=utf8_general_ci&loc=UTC&maxAllowedPacket=4&foo=bar")

	// endregion: config, defaults, dsn
	// region: MySQL

	dbConfig.Type = db.MySQL
	dbConfig.Addr = "localhost:23306"
	Db, _ = dbConfig.Open() // or db.Open(dbConfig)
	dbDrill()
	defer Db.Close() // or tdb.Close(db)

	// endregion: MySQL
	// region: MariaDB

	dbConfig.Type = db.MariaDB
	dbConfig.Addr = "localhost:13306"
	dbConfig.DSN = ""
	Db, _ = dbConfig.Open()
	dbDrill()
	defer Db.Close()

	// endregion: MariaDB
	// region: PostgreSQL

	db.Defaults = db.DefaultsPostgres
	dbConfig.Type = db.Postgres
	dbConfig.Addr = "localhost:15432"
	dbConfig.DSN = ""
	dbConfig.Params = nil
	Db, _ = dbConfig.Open()
	dbDrill()
	defer Db.Close()

	// endregion: Postgres
	// region: SQLite3

	dbConfig.Type = db.SQLite3
	dbConfig.Addr = "tex.db"
	dbConfig.DSN = ""
	dbConfig.Params = nil
	Db, err = dbConfig.Open()
	dbDrill()
	defer Db.Close()

	// endregion: SQLite3
	// region: history

	Logger.Out("History len", len(Db.History))
	for k, v := range Db.History {
		if v != nil {
			Logger.Out(k, v.SQL)

		}
	}
	// endregion: history

	// endregion: db

}

func dbDrill() {

	// region: Exec()

	// region: DROP TABLE

	dropTable := db.Statement{
		SQL: `	DROP TABLE IF EXISTS dummy;
		`,
	}
	dropTable.Exec(Db) // or Db.Exec(dropTable) or tdb.Exec(Db, dropTable)
	// Logger.Out(db.Drivers[Db.Config.Type], "DROP", dropTable.Err)

	// endregion: DROP TABLE
	// region: CREATE TABLE

	createTable := db.Statement{}
	if Db.Config.Type == db.Postgres || Db.Config.Type == db.SQLite3 {
		createTable.SQL = `	CREATE TABLE dummy (
								id		SERIAL			NOT NULL PRIMARY KEY,
								foo		VARCHAR(32)		NOT NULL
							);
		`
	} else {
		createTable.SQL = `	CREATE TABLE dummy (
								id		INT				NOT NULL AUTO_INCREMENT PRIMARY KEY,
								foo		VARCHAR(32)		NOT NULL
							);
		`
	}
	createTable.Exec(Db)
	// Logger.Out(db.Drivers[Db.Config.Type], "CREATE TABLE", createTable.Err)

	// endregion: CREATE TABLE
	// region: INSERT

	insertRows := db.Statement{
		Args: []interface{}{
			1, "foo",
			2, "bar",
			4, "baz",
			5, "xxx",
		},
	}
	if Db.Config.Type == db.Postgres {
		insertRows.SQL = `	INSERT INTO dummy
								(id, 	foo)
							VALUES
								($1,	$2),
								($3,	$4),
								($5,	$6),
								($7,	$8)
							RETURNING id;
		`
	} else {
		insertRows.SQL = `	INSERT INTO dummy
								(id, 	foo)
							VALUES
								(?,		?),
								(?,		?),
								(?,		?),
								(?,		?);
		`
	}
	insertRows.Exec(Db)
	// Logger.Out(db.Drivers[Db.Config.Type], "INSERT", insertRows.Err, insertRows.LastInsertId, insertRows.RowsAffected)

	// endregion: INSERT
	// region: UPDATE

	updateRows := db.Statement{Args: []interface{}{"quux", "baz"}}
	if Db.Config.Type == db.Postgres {
		updateRows.SQL = `UPDATE dummy SET foo = $1 WHERE foo = $2;`
	} else {
		updateRows.SQL = `UPDATE dummy SET foo = ? WHERE foo = ?;`
	}
	updateRows.Exec(Db)
	// Logger.Out(db.Drivers[Db.Config.Type], "UPDATE", updateRows.Err, updateRows.LastInsertId, updateRows.RowsAffected)

	// endregion: UPDATE
	// region: DELETE

	deleteRows := db.Statement{
		SQL: `	DELETE FROM dummy
				WHERE id = 4;
		`,
	}
	deleteRows.Exec(Db)
	// Logger.Out(db.Drivers[Db.Config.Type], "DELETE", deleteRows.Err, deleteRows.LastInsertId, deleteRows.RowsAffected)

	// endregion: DELETE

	// endregion: Exec()

}
