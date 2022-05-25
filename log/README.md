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

* Logger.HR [hint](https://stackoverflow.com/questions/16569433/get-terminal-size-in-go)
* max message width (in sample encoder)
* add taxonomy field
* extend file and line: func name(?), and full trace
* welcome/mark/bye severity (if severity present then use Out() otherwise c.Out())
* channel id/name, display like logLevel tags
* init by config json/struct (both Ch and Logger) (prerequisite: json/struct in cfg/)
* Ch.Type vs. Ch.Config.Type
* l.Out() parallel (goroutine) writes (w/ context and errGroup?)
* endpoints to change/reset config and level
* scheduled marker (after scheduler is implemented, use mark severity, could be a smart function)
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
