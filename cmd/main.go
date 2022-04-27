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
	Logger = *log.NewLogger()
	defer Logger.Close()
	_, _ = Logger.NewCh(log.ChConfig{Type: log.ChSyslog})
	_, _ = Logger.NewCh(log.ChConfig{Encoder: &spewEncoder})
	// lfc, _ := Logger.NewCh()
	// _ = lfc.Out(*log.ChDefaults.Mark)                               // write to identified channel
	// _ = lfc.Out(log.LOG_EMERG, "entry", "with", "severity")         // write to identified channel with severity
	// _ = log.Out(lfc, log.LOG_CRIT, "entry", "with", "severity")     // write to identified channel with severity
	// _ = Logger.Ch[0].Out(*log.ChDefaults.Mark, "bar", 1, 1.1, true) // write directly to the first channel
	// _ = Logger.Out(*log.ChDefaults.Mark)                            // write to all channels
	// _ = Logger.Out(log.LOG_ALERT, *log.ChDefaults.Mark)             // write to all channels with severity
	// _ = log.Out(&Logger, log.LOG_EMERG, "foobar")                   // write to all logger channels with severity
	// _ = log.Out(nil, log.LOG_EMERG, "quux")                         // write to nowhere
	// SomethinRemote()

	// endregion: logger
	// region: db

	dummyConfig := db.Config{
		User:   "testuser",
		Passwd: "testpasswd",
	}
	dummyConfig.SetDefaults() // or db.SetDefaults(&dummyConfig) is also available
	dummyConfig.FormatDSN()   // or db.FormatDSN(&dummyConfig)

	// TODO #1
	// dummyConfig, err = db.ParseDSN("user:pass@tcp(host)/dbname?allowNativePasswords=true&checkConnLiveness=true&collation=utf8&loc=UTC&maxAllowedPacket=4&foo=bar")
	// err = dummyConfig.ParseDSN("user:pass@tcp(host)/dbname?allowNativePasswords=true&checkConnLiveness=true&collation=utf8&loc=UTC&maxAllowedPacket=4&foo=bar")
	// Logger.Out(log.LOG_DEBUG, spew.Sdump(err))
	Logger.Out(log.LOG_DEBUG, dummyConfig)

	// TODO #2 normalize

	// TODO #3 Open
	// _, err = db.Open(db.Config{})
	// ll := log.LOG_DEBUG
	// _, err = db.Open(db.Config{Type: db.DbMariaDB, Logger: Logger, Loglevel: &ll})
	// Logger.Out(log.LOG_DEBUG, err)
	// _, err = db.Open(db.Config{Type: db.DbMySQL, Logger: Logger, Loglevel: &ll})
	// Logger.Out(log.LOG_DEBUG, err)
	// _, err = db.Open(db.Config{Type: db.DbPostgres, Logger: Logger, Loglevel: &ll})
	// Logger.Out(log.LOG_DEBUG, err)
	// _, err = db.Open(db.Config{Type: db.DbSQLite3, Logger: Logger, Loglevel: &ll})
	// Logger.Out(log.LOG_DEBUG, err)

	// endregion: db

}

func SomethinRemote() {
	Logger.Out(log.LOG_DEBUG, "Something remote...")
}
