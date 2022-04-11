# TEx-kit

## ToC

1. [ToC](#toc)
2. [cfg/](#cfg)
3. [log/](#log)
4. [Random improvements to be made](#random-improvements-to-be-made)

## cfg/

Read config from cli flags, environment or db...

## log/

Wrapper around `log` and `syslog`

## Random improvements to be made

* db
  * review - manage transactions during database inserts and updates
  * [gorm?](https://gorm.io/index.html)
* http
  * runtime.GOMAXPROCS(runtime.NumCPU())
  * expvar, runtime variables and functions, log status/count/file size
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
