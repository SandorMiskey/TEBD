# TEx-kit/log

## ToC

1. [ToC](#toc)
2. [Examples](#examples)
3. [Random improvements to be made](#random-improvements-to-be-made)

## Examples

```go
dbConfig := db.Config{
        User:   Config.Entries["dbUser"].Value.(string),
        Passwd: Config.Entries["dbPasswd"].Value.(string),
        DBName: "tex",
        Logger: Logger,
}
// dbConfig.SetDefaults() // or db.SetDefaults(&dbConfig) is also available
// dbConfig.FormatDSN()   // or db.FormatDSN(&dbConfig)
// _ = dbConfig.ParseDSN("user:pass@tcp(host)/dbname?allowNativePasswords=true&checkConnLiveness=true&collation=utf8_general_ci&loc=UTC&maxAllowedPacket=4&foo=bar")

// MySQL
dbConfig.Type = db.DbMySQL
dbConfig.Addr = "localhost:23306"
_, _ = dbConfig.Open() // or db.Open(dbConfig)

// MariaDB
dbConfig.Type = db.DbMariaDB
dbConfig.Addr = "localhost:13306"
dbConfig.DSN = ""
_, _ = dbConfig.Open()

// PostgreSQL
db.DbDefaults = db.DbDefaultsPostgres
dbConfig.Type = db.DbPostgres
dbConfig.Addr = "localhost:15432"
dbConfig.DSN = ""
dbConfig.Params = nil
_, _ = dbConfig.Open()

// SQLite3
dbConfig.Type = db.DbSQLite3
dbConfig.Addr = "tex.db"
dbConfig.DSN = ""
dbConfig.Params = nil
_, _ = dbConfig.Open()
```

## Random improvements to be made

* exec
  * prepared statements:
    * ~~Db has []*Tx~~
      * ~~struct def~~
      * ~~append on new tx~~
      * ~~getter~~?
    * Db has []*Statement with prepared
    * prepare only for db
    * in case of i=tx use tx.Stmt
  * batch Args []interface{} -> [][]interface{} (in a Tx, prepared)
  * batch Statement -> []Statement
* copy from trustone
  * Query to map
    * s.Result interface{} vs sql.Result
    * s.Rows?
    * s.Map, s.JSON?
  * Query to JSON
  * in all above: use type Statement, register last query(s) and result sets
  * in all above: prepared statements
  * do batches
  * typed result sets where applicable
* db.Open from JSON/cfg
* Statement from/TO JSON
* Statement catalog from/to JSON (+cfg)
* driver independent insert/update/delete (shortcuts for exec)
* ---
* execTransaction from TO?
* support context(.WithTimeout) in queries
* setters (like Db.SetLogger()), reset (re-parse config)
* connection catalog w/ close all
* SQLite authentication
* ---
* [gorm?](https://gorm.io/index.html)
