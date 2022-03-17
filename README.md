# TE-Commerce Backend

## ToC

1. [ToC](#toc)
2. [Random improvements to be made](#random-improvements-to-be-made)


## Random improvements to be made

* logging
  * review existing logger functions
  * structured (json, go.uber.org/zap + lumberjack?) + flat output
  * log levels
  * db (+ stderr/stdout, db, s3, url) output
  * endpoints to change config and level
  * [zap](https://pkg.go.dev/go.uber.org/zap#pkg-examples)
  * [log](https://pkg.go.dev/go.uber.org/zap#pkg-examples)
  * [syslog](https://pkg.go.dev/log/syslog)
* db
  * review - manage transactions during database inserts and updates
  * [gorm?](https://gorm.io/index.html)
* config
  * app -> db
  * reload/dump function
* http
  * fasthttprouter -> github.com/fasthttp/router
  * fiber
    * [fiber](https://github.com/gofiber/fiber)
    * <https://docs.gofiber.io>
  * <https://github.com/fasthttp/session>
  * ~~<https://github.com/fasthttp/websocket>~~ or <https://github.com/fasthttp/fastws> (w/ WASM client)
* redis support, with cache worker / automatic refresh
* scheduler, startup modules
* github.com/dgrijalva/jwt-go -> github.com/golang-jwt/jwt (or Fiber!)
* javascript support: <https://github.com/rogchap/v8go> (<https://esbuild.github.io/>)
* wasm
* swagger ui:
  * startup script: docker run -d -p 1081:8080 -e SWAGGER_JSON=/openapi.yaml -v repo root/???/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui
  * <https://github.com/arsmn/fiber-swagger>
  * <https://github.com/swaggo/swag>
* <https://gqlgen.com>
