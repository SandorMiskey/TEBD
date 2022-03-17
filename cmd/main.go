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
		"foo":      {Desc: "foo description", Type: "time.Duration", Def: time.Duration(66)},
		"foo_file": {Desc: "foo file desc", Type: "string", Def: ""},
		"bar":      {Desc: "bar description", Type: "int", Def: 55},
		"baz":      {Desc: "baz description", Type: "string", Def: "default baz"},
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
