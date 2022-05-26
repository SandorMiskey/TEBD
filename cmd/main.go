// region: packages

package main

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/SandorMiskey/TEx-kit/cfg"
	tedb "github.com/SandorMiskey/TEx-kit/db"
	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"
)

// endregion: packages
// region: global variables

var (
	Config cfg.Config
	Db     *tedb.Db
	Logger log.Logger
)

// endregion: globals

func main() {

	// region: config and cli flags

	Config = *cfg.NewConfig(os.Args[0])
	fs := Config.NewFlagSet(os.Args[0])
	fs.Entries = map[string]cfg.Entry{
		// "bool":     {Desc: "bool description", Type: "bool", Def: true},
		// "duration": {Desc: "duration description", Type: "time.Duration", Def: time.Duration(66000)},
		// "float64":  {Desc: "float64 desc", Type: "float64", Def: 77.7},

		"dbUser":        {Desc: "Database user", Type: "string", Def: ""},
		"dbPasswd":      {Desc: "Database password", Type: "string", Def: ""},
		"dbPasswd_file": {Desc: "Database password file", Type: "string", Def: ""},
		"dbName":        {Desc: "Database name", Type: "string", Def: "tex"},

		"loggerLevel": {Desc: "Logger min severity", Type: "int", Def: 5},
		"logLevel":    {Desc: "Log level everywhere", Type: "int", Def: 6},
	}

	err := fs.ParseCopy()
	if err != nil {
		panic(err)
	}

	// endregion: cli flags
	// region: logger

	// region: sample encoder

	var spewEncoder log.Encoder = func(c *log.Ch, n ...interface{}) (s string, e error) {
		s = "\n"
		for k, v := range n {
			if severity, ok := v.(syslog.Priority); ok {
				labels := *c.Config.SeverityLabels
				v = labels[severity]
			}
			s = fmt.Sprintf("%s%s%d%s%s", s, *c.Config.Delimiter, k, *c.Config.Delimiter, spew.Sdump(v))
		}
		// s = strings.Replace(s, *c.Config.Delimiter, "", 1)
		// s = strings.TrimSuffix(s, "\n")

		return s, nil
	}

	// endregion: sample encoder
	// region: new logger and channels

	logLevel := syslog.Priority(Config.Entries["logLevel"].Value.(int))
	loggerLevel := syslog.Priority(Config.Entries["loggerLevel"].Value.(int))

	Logger = *log.NewLogger()
	defer Logger.Close()
	_, _ = Logger.NewCh(log.ChConfig{Type: log.ChSyslog})
	lfc, _ := Logger.NewCh(log.ChConfig{Encoder: &spewEncoder, Severity: &loggerLevel})

	// endregion: logger and channels
	// region: sample messages

	_ = lfc.Out(logLevel, "entry1", "with", "severity")      // write to identified channel with severity
	_ = log.Out(lfc, logLevel, "entry2", "with", "severity") // write to identified channel with severity
	_ = log.Out(&Logger, logLevel, "foobar")                 // write to all logger channels with severity
	_ = Logger.Out(logLevel, *log.ChDefaults.Mark)           // write to all channels with severity
	// _ = lfc.Out(*log.ChDefaults.Mark)                            // write to identified channel
	// _ = Logger.Ch[0].Out(*log.ChDefaults.Mark, "bar", 1, 1.1, true) // write directly to the first channel
	// _ = Logger.Out(*log.ChDefaults.Mark)                            // write to all channels
	// _ = log.Out(nil, LogLevel, "quux")                         // write to nowhere

	// endregion: sample messages

	// endregion: logger
	// region: db

	// region: defaults, dsn

	dbDefaults := tedb.Config{
		User:   Config.Entries["dbUser"].Value.(string),
		Passwd: Config.Entries["dbPasswd"].Value.(string),
		DBName: Config.Entries["dbName"].Value.(string),
		Logger: Logger,
	}
	// dbConfig.SetDefaults() // or db.SetDefaults(&dbConfig) is also available
	// dbConfig.FormatDSN()   // or db.FormatDSN(&dbConfig)
	// _ = dbConfig.ParseDSN("user:pass@tcp(host)/dbname?allowNativePasswords=true&checkConnLiveness=true&collation=utf8_general_ci&loc=UTC&maxAllowedPacket=4&foo=bar")

	// endregion: defaults, dsn
	// region: configs

	var dbConfigs []tedb.Config = make([]tedb.Config, 4)

	dbConfigs[0] = dbDefaults
	dbConfigs[0].Type = tedb.MariaDB
	dbConfigs[0].Addr = "localhost:13306"
	dbConfigs[1] = dbDefaults
	dbConfigs[1].Type = tedb.MySQL
	dbConfigs[1].Addr = "localhost:23306"
	dbConfigs[2] = dbDefaults
	dbConfigs[2].Type = tedb.Postgres
	dbConfigs[2].Addr = "localhost:15432"
	dbConfigs[3] = dbDefaults
	dbConfigs[3].Type = tedb.SQLite3
	dbConfigs[3].Addr = "tex.db"

	// endregion: configs
	// region: db connections and drills

	for _, conf := range dbConfigs {

		// region: connection

		tedb.Defaults = tedb.DefaultsMySQL
		if conf.Type == tedb.Postgres {
			tedb.Defaults = tedb.DefaultsPostgres
		}

		Db, err = conf.Open() // or db.Open(conf)
		defer Db.Close()      // or db.Close(dbInstance)
		if err != nil {
			Logger.Out(logLevel, "db err", err)
		}

		// endregion: conn
		// region: Exec()

		// region: DROP TABLE

		dropTable := tedb.Statement{
			SQL: `DROP TABLE IF EXISTS dummy;`,
		}
		dropTable.Exec(Db) // or Db.Exec(dropTable) or db.Exec(Db, dropTable)
		if dropTable.Err != nil {
			Logger.Out(tedb.Drivers[Db.Config.Type], "DROP TABLE ERROR", dropTable.Err)
		}

		// endregion: DROP TABLE
		// region: CREATE TABLE

		createTable := tedb.Statement{}
		if Db.Config.Type == tedb.Postgres || Db.Config.Type == tedb.SQLite3 {
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
		if createTable.Err != nil {
			Logger.Out(tedb.Drivers[Db.Config.Type], "CREATE TABLE ERROR", createTable.Err)

		}

		// endregion: CREATE TABLE
		// region: INSERT

		insertRows := tedb.Statement{
			Args: []interface{}{
				1, "foo",
				2, "bar",
				4, "baz",
				5, "xxx",
			},
		}
		if Db.Config.Type == tedb.Postgres {
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
		Logger.Out(log.LOG_INFO, tedb.Drivers[Db.Config.Type], "INSERT", insertRows.LastInsertId, insertRows.RowsAffected)
		if insertRows.Err != nil {
			Logger.Out(tedb.Drivers[Db.Config.Type], "INSERT ERROR", createTable.Err)
		}

		// endregion: INSERT
		// region: UPDATE

		updateRows := tedb.Statement{Args: []interface{}{"quux", "baz"}}
		if Db.Config.Type == tedb.Postgres {
			updateRows.SQL = `UPDATE dummy SET foo = $1 WHERE foo = $2;`
		} else {
			updateRows.SQL = `UPDATE dummy SET foo = ? WHERE foo = ?;`
		}
		updateRows.Exec(Db)
		Logger.Out(log.LOG_INFO, tedb.Drivers[Db.Config.Type], "UPDATE", updateRows.LastInsertId, updateRows.RowsAffected)
		if updateRows.Err != nil {
			Logger.Out(tedb.Drivers[Db.Config.Type], "UPDATE ERROR", updateRows.Err)
		}

		// endregion: UPDATE
		// region: DELETE

		deleteRows := tedb.Statement{
			SQL: `DELETE FROM dummy WHERE id = 4;`,
		}
		deleteRows.Exec(Db)
		Logger.Out(log.LOG_INFO, tedb.Drivers[Db.Config.Type], "DELETE", deleteRows.Err, deleteRows.LastInsertId, deleteRows.RowsAffected)
		if deleteRows.Err != nil {
			Logger.Out(tedb.Drivers[Db.Config.Type], "DELETE ERROR", deleteRows.Err)
		}

		// endregion: DELETE

		// endregion: Exec()
		// region: history

		Logger.Out(logLevel, "HISTORY LENGTH", len(Db.History))
		for k, v := range Db.History {
			if v != nil {
				Logger.Out(logLevel, "HISTORY ENTRY", k, v.SQL)
			}
		}

		// endregion: history
	}

	// endregion: db connections and drills

	// endregion: db

}
