# TEx-kit/log

## ToC

1. [ToC](#toc)
2. [Examples](#examples)
3. [Random improvements to be made](#random-improvements-to-be-made)

## Examples

```go
dc, _ := log.NewCh(log.ChConfig{Type: log.ChSyslog})
defer dc.Close()
dc.Out(*log.ChDefaults.Mark)

Logger = *log.NewLogger()
lc, _ := Logger.NewCh()
defer Logger.Close()
_ = Logger.Ch[0].Out(*log.ChDefaults.Mark, "bar", 1, 1.1, true) // write direct to the first channel
_ = lc.Out(*log.ChDefaults.Mark)                                // write to identified channel
_ = Logger.Out(*log.ChDefaults.Mark)                            // write to all channels
_ = log.Out(lc, log.LOG_EMERG, "entry", "with", "severity")     // write to identified channel with severity
_ = log.Out(&Logger, log.LOG_EMERG, "foobar")                   // write to all logger channels with severity

```

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
