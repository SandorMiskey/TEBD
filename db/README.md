# TEx-kit/log

## ToC

1. [ToC](#toc)
2. [Examples](#examples)
3. [Infra](#infra)
4. [Random improvements to be made](#random-improvements-to-be-made)

## Examples

```go
...
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

* copy from trustone
* typed result sets where applicable
* support context.WithTimeout() in queries
* setters (like Db.SetLogger()), reset (re-parse config)
* SQLite authentication
* transactions
* ---
* [gorm?](https://gorm.io/index.html)
