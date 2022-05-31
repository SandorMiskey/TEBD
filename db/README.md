# TEx-kit/log

## ToC

1. [ToC](#toc)
2. [Examples](#examples)
3. [Infra](#infra)
4. [Random improvements to be made](#random-improvements-to-be-made)

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

## Infra

Launch database instances (in docker, where applicable) for development:

```sh
docker  run                                   \
        --detach                              \
        --name tex-mariadb                    \
        --env MARIADB_DATABASE=tex            \
        --env MARIADB_USER=tex                \
        --env MARIADB_PASSWORD=<pw_here>      \
        --env MARIADB_ROOT_PASSWORD=<pw_here> \
        --publish 13306:3306                  \
        mariadb:latest
```

```sh
docker  run                                 \
        --detach                            \
        --name tex-mysql                    \
        --env MYSQL_DATABASE=tex            \
        --env MYSQL_USER=tex                \
        --env MYSQL_PASSWORD=<pw_here>      \
        --env MYSQL_ROOT_PASSWORD=<pw_here> \
        --publish 23306:3306                \
        mysql:latest
```

```sh
docker  run                               \
        --detach                          \
        --name tex-postgres               \
        --env POSTGRES_DB=tex             \
        --env POSTGRES_USER=tex           \
        --env POSTGRES_PASSWORD=TEx99! \
        --publish 15432:5432              \
        postgres:latest
```

```sh
sqlite3 tex.db
```

## Random improvements to be made

* Exec
  * add Statement.Db
  * Executable => interface
  * Exec() handle Tx
  * Tx.Exec
  * Statement.Exec() handle Tx
* add Commit func w/ history
* copy from trustone
  * Query to map
    * s.Result interface{} vs sql.Result
    * s.Rows?
    * s.Map, s.JSON?
  * Query to JSON
  * execTransaction
  * in all above: use type Statement, register last query(s) and result sets
  * in all above: prepared statements
  * do batches
  * typed result sets where applicable
* db.Statement initialization from JSON
* support context.WithTimeout() in queries
* setters (like Db.SetLogger()), reset (re-parse config)
* connection catalog w/ close all
* SQLite authentication
* ---
* [gorm?](https://gorm.io/index.html)
