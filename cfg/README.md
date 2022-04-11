# TEx-kit/cfg

## ToC

1. [ToC](#toc)
2. [Examples](#examples)
3. [Random improvements to be made](#random-improvements-to-be-made)

## Examples

```go
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
```

## Random improvements to be made

* Trace() and `caller` in cfg/cfg.go and cfg/fs.go
* json type
* db
* set env. variables
* reload/dump function (maybe restart main?)
* logging?
