# TEx-kit/log

## ToC

1. [ToC](#toc)
2. [Random improvements to be made](#random-improvements-to-be-made)

## Random improvements to be made

* welcome/mark/bye severity (if severity present then use Out() otherwise c.Out())
* init by config json/struct (prerequisite: json/struct in cfg/)
* endpoints to change/reset config and level
* scheduled marker (after scheduler is implemented, use mark severity)
* hooks
* log rotation
* output destinations:
  * db
    * implement db/ first
    * Ch.File (?) and ChClose for all ChType
  * syslog local:
    * fix on mac
    * implement *Ch.Close()
  * syslog remote
  * net: nc, s3, nfs etc.
* output encoder
  * encoding/json
  * encoding/csv
  * encoding/xml
  * db
