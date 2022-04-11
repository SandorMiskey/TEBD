// region: packages

package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// endregion: packages
// region: types

type FlagSet struct {
	config     *Config
	createdAt  time.Time
	createdBy  string
	modifiedAt time.Time
	modifiedBy string

	Arguments     []string
	Entries       map[string]Entry
	ErrorHandling flag.ErrorHandling
	FileSuffix    string
	FlagSet       *flag.FlagSet
	Name          string
	Output        io.Writer
	Usage         func() `json:"-"`
}

// endregion: types
// region: defaults

// default list of arguments
var FlagSetArguments = os.Args[1:]

// default -h message
var FlagSetUsage = func(fs *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage of %s:\n\n", fs.Name())
		fs.PrintDefaults()
		fmt.Printf("\n")
		fmt.Printf("  Paramters can also be passed via env. variables, like `HTTPPORT=80` insted of '-httpPort=80'. Order of precedence:\n")
		fmt.Printf("  1. Command line options\n")
		fmt.Printf("  2. Environment variables\n")
		fmt.Printf("  3. Default values\n\n")
	}
}

// default error handling
var FlagSetErrorHandling = flag.ExitOnError

// default output for usage and error messages
var FlagSetOutput = os.Stderr

// default suffix for values in files
var FlagSetFileSuffix = "_file"

// endregion: defaults
// region: flagset constructor

func NewFlagSet(name string) *FlagSet {
	// _, _, caller := log.Trace()
	fs := FlagSet{
		config:    nil,
		createdAt: time.Now().UTC(),
		// createdBy:  caller,
		modifiedAt: time.Now().UTC(),
		// modifiedBy: caller,

		Arguments:     FlagSetArguments,
		Entries:       make(map[string]Entry),
		ErrorHandling: FlagSetErrorHandling,
		FileSuffix:    FlagSetFileSuffix,
		Name:          name,
		Output:        FlagSetOutput,
	}
	fs.FlagSet = flag.NewFlagSet(name, fs.ErrorHandling)
	fs.Usage = FlagSetUsage(fs.FlagSet)
	return &fs
}

func (c *Config) NewFlagSet(name string) *FlagSet {
	// _, _, caller := log.Trace()
	fs := NewFlagSet(name)
	// fs.createdBy = caller
	// fs.modifiedBy = caller
	fs.config = c
	fs.FlagSet = flag.NewFlagSet(name, fs.ErrorHandling)
	fs.Usage = FlagSetUsage(fs.FlagSet)
	c.FlagSet[name] = fs
	return fs
}

// endregion: flagset constructor
// region: flagset parse

//
// 1. os.LookupEnv() environment variables into `env`, also checking whether the value of the variable is a null string or it is unset ("value"|""|nil)
// 2. add flag entries and parse (w/ errorhandling and usage) into `cli`
// 3. merge `env` and `cli`
// 	a. copy `env` entries to `fs.Entries` where `env.Value` is not null
// 	b. copy `env` entries w/ null value into `fs.Entries` w/ default value
// 	c. flag.Visit parsed cli flags and overwrite `fs.Entries[key]`
// 4. resolve flags w/ FlagSetFileSuffix if value != "" or null
// 	a. flag w/o suffix exists
// 	b. flag w/o suffix doesn't exist
// 5. double check if value and type match
//

func (fs *FlagSet) Parse() (err error) {

	// region: 1. os.LookupEnv() environment variables into `env`, also checking whether the value of the variable is a null string or it is unset ("value"|""|nil)

	env := make(map[string]Entry)

	for key := range fs.Entries {
		entry := fs.Entries[key]
		entry.createdAt = time.Now().UTC()
		entry.modifiedAt = time.Now().UTC()
		// _, _, entry.createdBy = log.Trace()
		// _, _, entry.modifiedBy = log.Trace()

		value, set := os.LookupEnv(strings.ToUpper(key))
		if set {
			entry.Value = value
		} else {
			entry.Value = nil
		}

		env[key] = entry
	}

	// endregion: env
	// region: 2. add flag entries and parse (w/ errorhandling and usage) into `cli`

	cli := make(map[string]Entry)

	for key := range fs.Entries {
		entry := fs.Entries[key]
		entry.createdAt = time.Now().UTC()
		entry.modifiedAt = time.Now().UTC()
		// _, _, entry.createdBy = log.Trace()
		// _, _, entry.modifiedBy = log.Trace()

		switch entry.Type {
		case "bool":
			entry.Value = fs.FlagSet.Bool(key, entry.Def.(bool), entry.Desc)
		case "time.Duration":
			entry.Value = fs.FlagSet.Duration(key, entry.Def.(time.Duration), entry.Desc)
		case "float64":
			entry.Value = fs.FlagSet.Float64(key, entry.Def.(float64), entry.Desc)
		case "int":
			entry.Value = fs.FlagSet.Int(key, entry.Def.(int), entry.Desc)
		case "string":
			entry.Value = fs.FlagSet.String(key, entry.Def.(string), entry.Desc)
		default:
			return fmt.Errorf("invalid flag type: %s", entry.Type)
		}

		cli[key] = entry
	}

	fs.FlagSet.SetOutput(fs.Output)
	fs.FlagSet.Usage = fs.Usage
	fs.FlagSet.Parse(fs.Arguments)

	// endregion: cli
	// region: 3. merge `env` and `cli`

	// 	a. copy `env` entries to `fs.Entries` where `env.Value` is not null

	for key := range env {
		if env[key].Value != nil {
			entry := env[key]
			entry.Source = EntrySourceEnv
			fs.Entries[key] = entry
		}
	}

	// 	b. copy `env` entries w/ null value into `fs.Entries` w/ default value

	for key := range env {
		if env[key].Value == nil {
			entry := env[key]
			entry.Value = entry.Def
			entry.Source = EntrySourceDef
			fs.Entries[key] = entry
		}
	}

	// 	c. flag.Visit parsed cli flags and overwrite `fs.Entries[key]`

	err = nil
	fs.FlagSet.Visit(func(f *flag.Flag) {
		entry := cli[f.Name]
		entry.Source = EntrySourceCli

		switch entry.Value.(type) {
		case *bool:
			entry.Value = *entry.Value.(*bool)
		case *float64:
			entry.Value = *entry.Value.(*float64)
		case *int:
			entry.Value = *entry.Value.(*int)
		case *string:
			entry.Value = *entry.Value.(*string)
		case *time.Duration:
			entry.Value = *entry.Value.(*time.Duration)
		default:
			err = fmt.Errorf("invalid flag value %s", reflect.TypeOf(entry.Value).String())
		}

		fs.Entries[f.Name] = entry
	})
	if err != nil {
		return
	}

	// endregion: merge
	// region: 4. resolve flags w/ FlagSetFileSuffix if value != "" or null

	for key := range fs.Entries {
		if strings.HasSuffix(key, fs.FileSuffix) && fs.Entries[key].Value != "" && fs.Entries[key].Value != nil {
			data, err := os.ReadFile(fs.Entries[key].Value.(string))
			if err != nil {
				return err
			}
			value := strings.TrimSuffix(string(data), "\n")
			linked := strings.TrimSuffix(key, fs.FileSuffix)
			entry := fs.Entries[linked]

			switch entry.Type {
			case "bool":
				entry.Value, err = strconv.ParseBool(value)
			case "time.Duration":
				entry.Value, err = time.ParseDuration(value)
			case "float64":
				entry.Value, err = strconv.ParseFloat(value, 64)
			case "int":
				entry.Value, err = strconv.Atoi(value)
			case "string":
				entry.Value = value
			default:
				return fmt.Errorf("invalid flag type: '%s' at %s/%s ", entry.Type, linked, key)
			}
			if err != nil {
				return err
			}
			fs.Entries[linked] = entry
		}
	}

	// endregion: file suffix
	// region: doublecheck types

	for key := range fs.Entries {
		entry := fs.Entries[key]
		typ := reflect.TypeOf(entry.Value).String()
		if typ != entry.Type {
			return fmt.Errorf("type mismatch for %s: %s vs %s", key, typ, entry.Type)
		}
	}

	// endregion: types

	return
}

// endregion: flagset parse
// region: flagset (parse and) copy

// copy fs.Entries -> config.Entries

func (fs *FlagSet) Copy() {
	for key := range fs.Entries {
		fs.config.Entries[key] = fs.Entries[key]
	}
}

func (fs *FlagSet) ParseCopy() (err error) {
	err = fs.Parse()
	if err != nil {
		return
	}
	fs.Copy()
	return
}

// endregion: flagset copy
// region: dumps

func (fs *FlagSet) Sdump() string {
	return spew.Sdump(fs)
}

func (fs *FlagSet) MarshalJSON() ([]byte, error) {
	type FuncAlias FlagSet
	return json.Marshal(&struct {
		*FuncAlias
		FlagSet string `json:"FlagSet"`
	}{
		FuncAlias: (*FuncAlias)(fs),
		FlagSet:   "func () value is not supported therefore it is masked",
	})
}
func (fs *FlagSet) JSON(prefix string, indent string) ([]byte, error) {
	return json.MarshalIndent(fs, prefix, indent)
}

// endregion: dumps
