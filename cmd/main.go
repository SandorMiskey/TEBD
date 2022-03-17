// region: packages

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/SandorMiskey/TEx-kit/cfg"
	"github.com/davecgh/go-spew/spew"
)

// endregion: packages
// region: global variables

var (
	Config cfg.Config
)

// endregion: globals

func main() {

	// region: config and cli flags

	Config = *cfg.New(os.Args[0])
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

	something()
}

func something() {
	// cj, err := Config.JSON("", "	")
	// fmt.Printf("--> JSON representation (err: %v)\n", err)
	// fmt.Println(string(cj))
	// fmt.Println("--> Spew dump")
	// fmt.Println(Config.Sdump())
	fmt.Println(spew.Sdump(Config.Entries))
}
