// region: packages

package cfg

import (
	"encoding/json"
	"time"

	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/davecgh/go-spew/spew"
)

// endregion: packages
// region: types

// region: entrySource to indicate where the value is comming from (not sure if it will be implemented at all)

type entrySource int

const (
	EntrySourceEnv entrySource = iota + 1
	EntrySourceCli
	EntrySourceDef
	EntrySourceDb
)

// endregion: entrySource
// region: Entry

type Entry struct {
	createdAt  time.Time
	createdBy  string
	modifiedAt time.Time
	modifiedBy string

	Desc   string
	Type   string
	Def    interface{}
	Source entrySource
	Value  interface{}
}

// endregion: Entry
// region: Config

type Config struct {
	createdAt  time.Time
	createdBy  string
	modifiedAt time.Time
	modifiedBy string

	Entries map[string]Entry
	FlagSet map[string]*FlagSet
	Name    string
}

// endregion: Config

// endregion: types
// region: constructor

func New(name string) *Config {
	_, _, caller := log.Trace()
	return &Config{
		createdAt:  time.Now().UTC(),
		createdBy:  caller,
		modifiedAt: time.Now().UTC(),
		modifiedBy: caller,

		Entries: make(map[string]Entry),
		FlagSet: make(map[string]*FlagSet),
		Name:    name,
	}
}

// endregion: constructor
// region: dumps

func (c *Config) Sdump() string {
	return spew.Sdump(c)
}

func (c *Config) JSON(prefix string, indent string) ([]byte, error) {
	return json.MarshalIndent(c, prefix, indent)
}

// endregion: dumps
