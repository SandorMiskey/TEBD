// region: packages

package main

import (
	"os"
	"time"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/SandorMiskey/TEx-kit/log"
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

	dc, _ := log.NewCh(log.ChConfig{Type: log.ChSyslog})
	defer dc.Close()
	dc.Out(*log.ChDefaults.Mark)

	Logger = *log.NewLogger()
	lc, _ := Logger.NewCh()
	defer Logger.Close()
	_ = Logger.Ch[0].Out(*log.ChDefaults.Mark, "bar", 1, 1.1, true) // write direct to the first channel
	_ = lc.Out(*log.ChDefaults.Mark)                                // write to identified channel
	_ = Logger.Out(*log.ChDefaults.Mark)                            // write to all channels
	_ = log.Out(lc, log.LOG_EMERG, "entry", "with", "severity")     // write to identified channel with severity
	_ = log.Out(&Logger, log.LOG_EMERG, "foobar")                   // write to all logger channels with severity

	// endregion: logger

}
