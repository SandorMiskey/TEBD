# TBD-Framework

## 1. Table of contents

1. [1. Table of contents](#1-table-of-contents)
2. [2. Shortcuts](#2-shortcuts)
3. [3. Work in progress](#3-work-in-progress)
   1. [3.1. Scheduled](#31-scheduled)
   2. [3.2. Improvements to be made later](#32-improvements-to-be-made-later)
4. [4. Build options](#4-build-options)
5. [5. Runtime configuration](#5-runtime-configuration)
6. [6. API](#6-api)
   1. [6.1. Protocols and communication](#61-protocols-and-communication)
   2. [6.2. Data and content types](#62-data-and-content-types)
   3. [6.3. Response status](#63-response-status)
   4. [6.4. Security measures](#64-security-measures)
   5. [6.5. Versioning](#65-versioning)

## 2. Shortcuts

* [Swagger UI](http://???:1081)
* [Demo API server](https://???:2443/)

## 3. Work in progress

### 3.1. Scheduled

### 3.2. Improvements to be made later

* github
* logging
  * review existing logger functions
  * db output
* db
  * review - manage transactions during database inserts and updates
* redis support, with cache worker / automatic refresh
* ethereum/hl support
* scheduler, startup modules
* javascript support: https://github.com/rogchap/v8go (https://esbuild.github.io/)
* wasm
* config
  * infra -> env
  * app -> db
  * reload/dump function
* config into database and provide endpoint to force update
* swagger ui:
  * startup script: docker run -d -p 1081:8080 -e SWAGGER_JSON=/openapi.yaml -v /srv/web/healthpass.ge-demo.te-food.com/src/healthpass/build/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui
  * remove search bar
  * link to swagger ui from README.md/Shortcuts
* are these actually used? (env.sh, global variables, database)
  * S3
  * PAGINATION
* containerization, make mainBuild available as an endpoint for diagnostic purposes
* superadmin ui and monitoring endpoint

## 4. Build options

For available build options see `build.sh -h` in the `build/` subdirectory.

## 5. Runtime configuration

For available runtime config parameters see `build.sh -h` or `./main --help`. Config parameters can be passed via cli flags (`-xxxxxx=yyyyyy`) or env. variables (`export XXXXXX=yyyyyy`). Order of precedence:

1. Command line options
2. Environment variables
3. Default values

## 6. API

Detailed description of API resources and endpoints can be found [here](http://???:1081/#/).

### 6.1. Protocols and communication

### 6.2. Data and content types

The API sends and receives all content in standard JSON string format. Some things must ensured when preparing a request:

* JSON strings are UTF-8 encoded.
* The MIME type of the request sent to the API is defined as 'application/json'.
* Content-length should also be specified explicitly in all HTTP request headers.

The data types used are the following:

* string: Standard character string, always UTF-8 encoded.
* int/integer: 32-bit signed integer.
* GUID (Global Unique Identifier): A 128-bit long hexadecimal identifier with 32 character string representation, like this: 21EC2020-3AEA-4069-A2DD-08002B30309D. Its value must be unique.
* DateTime: Date and time values, represented in a string format of the ISO-8601 standard, like this: 2014-12-06T08:35:46Z. All responses, where applicable, contain DateTimes as strings.
* bool: Boolean true or false values. (Note: avoid integer or byte evaluation when constructing the request JSON!)
* byte: Byte representation of an enumeration value. This is only for internal use only, API call responses always contain the string representation.

### 6.3. Response status

HTTP status indicates the nature of the outcome of the request:

* 200 for successful fulfillment.
* 4xx means error, presumably due to client-side data.
* 5xx means error, presumably due to server-side malfunction.

See details in the description of endpoints.

### 6.4. Security measures

### 6.5. Versioning
