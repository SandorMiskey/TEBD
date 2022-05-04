// region: packages

package main

import (
	"fmt"
	"log/syslog"
	"os"
	"time"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/db"
	tdb "github.com/SandorMiskey/TEx-kit/db"
	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"
)

// endregion: packages
// region: global variables

var (
	Config cfg.Config
	Logger log.Logger
)

// endregion: globals

func main() {

	// region: config and cli flags

	Config = *cfg.NewConfig(os.Args[0])
	fs := Config.NewFlagSet(os.Args[0])
	fs.Entries = map[string]cfg.Entry{
		"bool":        {Desc: "bool description", Type: "bool", Def: true},
		"duration":    {Desc: "duration description", Type: "time.Duration", Def: time.Duration(66000)},
		"float64":     {Desc: "float64 desc", Type: "float64", Def: 77.7},
		"int":         {Desc: "int description", Type: "int", Def: 99},
		"string":      {Desc: "string description", Type: "string", Def: "string"},
		"string_file": {Desc: "string_file description", Type: "string", Def: ""},

		"dbUser":   {Desc: "Database user", Type: "string", Def: ""},
		"dbPasswd": {Desc: "Database password", Type: "string", Def: ""},
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
		for k, v := range n {
			s = fmt.Sprintf("%s%s%d%s%s\n", s, *c.Config.Delimiter, k, *c.Config.Delimiter, spew.Sdump(v))
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

	dbConfig := tdb.Config{
		User:   Config.Entries["dbUser"].Value.(string),
		Passwd: Config.Entries["dbPasswd"].Value.(string),
		DBName: "tex",
		Logger: Logger,
	}
	// dbConfig.SetDefaults() // or db.SetDefaults(&dbConfig) is also available
	// dbConfig.FormatDSN()   // or db.FormatDSN(&dbConfig)
	// _ = dbConfig.ParseDSN("user:pass@tcp(host)/dbname?allowNativePasswords=true&checkConnLiveness=true&collation=utf8_general_ci&loc=UTC&maxAllowedPacket=4&foo=bar")

	// MySQL
	dbConfig.Type = tdb.MySQL
	dbConfig.Addr = "localhost:23306"
	db, _ := dbConfig.Open() // or db.Open(dbConfig)
	dbDrill(db)
	db.Close()

	// MariaDB
	dbConfig.Type = tdb.MariaDB
	dbConfig.Addr = "localhost:13306"
	dbConfig.DSN = ""
	db, _ = dbConfig.Open()
	dbDrill(db)
	db.Close()

	// PostgreSQL
	tdb.Defaults = tdb.DefaultsPostgres
	dbConfig.Type = tdb.Postgres
	dbConfig.Addr = "localhost:15432"
	dbConfig.DSN = ""
	dbConfig.Params = nil
	db, _ = dbConfig.Open()
	dbDrill(db)
	db.Close()

	// SQLite3
	dbConfig.Type = tdb.SQLite3
	dbConfig.Addr = "tex.db"
	dbConfig.DSN = ""
	dbConfig.Params = nil
	db, _ = dbConfig.Open()
	dbDrill(db)
	defer db.Close() // or tdb.Close(db)

	// endregion: db

}

func dbDrill(db *db.Db) {
	Logger.Out(log.LOG_INFO, "DB DRILL...")
}
