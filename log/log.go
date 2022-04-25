// region: packages

package log

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/syslog"
	"os"
	"strings"
)

// endregion: packages

// region: types

type Encoder func(c *Ch, n ...interface{}) (string, error)

type ChType int

type SeverityLabels map[syslog.Priority]string

type ChConfig struct {
	// Db      interface{}
	Bye            *string
	Delimiter      *string
	Depth          *int
	Encoder        *Encoder
	Facility       *syslog.Priority
	File           interface{}
	FileFlags      *int
	FilePerm       *int
	Flags          *int
	Mark           *string
	Prefix         *string
	Severity       *syslog.Priority
	SeverityLabels *SeverityLabels
	Type           ChType
	Welcome        *string
}

type Ch struct {
	Config  ChConfig
	Encoder *Encoder
	File    *os.File
	Inst    *log.Logger
	// Inst interface{}
	Type ChType
}

type LoggerConfig struct {
	Ch []ChConfig
}

type Logger struct {
	Ch []*Ch
}

// endregion: types
// region: constants

// channel types
const (
	ChUndefined ChType = iota
	ChDb
	ChFile
	ChSyslog
)

// syslog priority
const (
	LOG_EMERG   = syslog.LOG_EMERG
	LOG_ALERT   = syslog.LOG_ALERT
	LOG_CRIT    = syslog.LOG_CRIT
	LOG_ERR     = syslog.LOG_ERR
	LOG_WARNING = syslog.LOG_WARNING
	LOG_NOTICE  = syslog.LOG_NOTICE
	LOG_INFO    = syslog.LOG_INFO
	LOG_DEBUG   = syslog.LOG_DEBUG

	LOG_KERN     = syslog.LOG_KERN
	LOG_USER     = syslog.LOG_USER
	LOG_MAIL     = syslog.LOG_MAIL
	LOG_DAEMON   = syslog.LOG_DAEMON
	LOG_AUTH     = syslog.LOG_AUTH
	LOG_SYSLOG   = syslog.LOG_SYSLOG
	LOG_LPR      = syslog.LOG_LPR
	LOG_NEWS     = syslog.LOG_NEWS
	LOG_UUCP     = syslog.LOG_UUCP
	LOG_CRON     = syslog.LOG_CRON
	LOG_AUTHPRIV = syslog.LOG_AUTHPRIV
	LOG_FTP      = syslog.LOG_FTP

	LOG_LOCAL0 = syslog.LOG_LOCAL0
	LOG_LOCAL1 = syslog.LOG_LOCAL1
	LOG_LOCAL2 = syslog.LOG_LOCAL2
	LOG_LOCAL3 = syslog.LOG_LOCAL3
	LOG_LOCAL4 = syslog.LOG_LOCAL4
	LOG_LOCAL5 = syslog.LOG_LOCAL5
	LOG_LOCAL6 = syslog.LOG_LOCAL6
	LOG_LOCAL7 = syslog.LOG_LOCAL7
)

// endregion: constants
// region: defaults

// region: ChDefaults

var bye = os.Args[0] + " logger is leaving..."
var delimiter = " -> "
var depth = 0
var facility = syslog.LOG_LOCAL0
var fileflags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
var fileperm = 0640
var flags = log.Ldate | log.Ltime | log.LUTC | log.Lshortfile
var mark = "logger was here..."
var prefix = "==> "
var severity = syslog.LOG_INFO
var severityLabels SeverityLabels = map[syslog.Priority]string{
	LOG_EMERG:   "__EMERG__: ",
	LOG_ALERT:   "__ALERT__: ",
	LOG_CRIT:    "__CRIT__: ",
	LOG_ERR:     "__ERR__: ",
	LOG_WARNING: "__WARNING__: ",
	LOG_NOTICE:  "__NOTICE__: ",
	LOG_INFO:    "__INFO__: ",
	LOG_DEBUG:   "__DEBUG__: ",
}
var welcome = os.Args[0] + " logger has been initiated\n"

var ChDefaults = ChConfig{
	Bye:            &bye,            // default exit msg
	Delimiter:      &delimiter,      // default delimiter
	Depth:          &depth,          // default depth offset for Output()
	Encoder:        &EncoderFlat,    // default encoder
	Facility:       &facility,       // default syslog facility
	File:           os.Stdout,       // default file
	FileFlags:      &fileflags,      // default flags to OpenFile wrapping those of the underlying system
	FilePerm:       &fileperm,       // default permissions for logfiles
	Flags:          &flags,          // define which text to prefix to each log entry generated by the Logger
	Mark:           &mark,           // default mark msg
	Prefix:         &prefix,         // default output prefix
	Severity:       &severity,       // default syslog severity
	SeverityLabels: &severityLabels, // default labels for severities
	Type:           ChFile,          // default Ch.Type
	Welcome:        &welcome,        // default mark msg
}

// endregion: ChDefaults
// region: messages

var (
	ErrInvalidFile            = errors.New("invalid file")
	ErrInvalidLoggerOrChannel = errors.New("invalid logger or channel")
	ErrInvalidSeverity        = errors.New("invalid severity")
	ErrNotImplementedYet      = errors.New("not implemented yet")
	ErrTooManyParameters      = errors.New("too many parameters")
)

// endregion: messages

// endregion: defaults
// region: constructors and destructors

// region: logger

func NewLogger() (l *Logger) {
	return &Logger{
		Ch: make([]*Ch, 0),
	}
}

func (l *Logger) Close() (e error) {
	for _, ch := range l.Ch {
		e = ch.Close()
	}
	return e
}

// endregion: logger
// region: destinations

func NewCh(cs ...ChConfig) (*Ch, error) {

	// region: prepare input

	// not sure if this is idiomatic, but this way you can call log.NewCh() instead of log.NewCh(log.ChDefaults)

	if len(cs) == 0 {
		cs = append(cs, ChDefaults)
	}
	if len(cs) > 1 {
		return nil, ErrTooManyParameters
	}
	c := cs[0]

	// endregion: prepare
	// region: check/set defaults

	if c.Bye == nil {
		c.Bye = ChDefaults.Bye
	}
	if c.Delimiter == nil {
		c.Delimiter = ChDefaults.Delimiter
	}
	if c.Depth == nil {
		c.Depth = ChDefaults.Depth
	}
	if c.Encoder == nil {
		c.Encoder = ChDefaults.Encoder
	}
	if c.File == nil || c.File == "" {
		c.File = ChDefaults.File
	}
	if c.Facility == nil {
		c.Facility = ChDefaults.Facility
	}
	if c.FileFlags == nil {
		c.FileFlags = ChDefaults.FileFlags
	}
	if c.FilePerm == nil {
		c.FilePerm = ChDefaults.FilePerm
	}
	if c.Flags == nil {
		c.Flags = ChDefaults.Flags
	}
	if c.Mark == nil {
		c.Mark = ChDefaults.Mark
	}
	if c.Prefix == nil {
		c.Prefix = ChDefaults.Prefix
	}
	if c.Severity == nil {
		c.Severity = ChDefaults.Severity
	}
	if c.SeverityLabels == nil {
		c.SeverityLabels = ChDefaults.SeverityLabels
	}
	if c.Type == ChUndefined {
		c.Type = ChDefaults.Type
	}
	if c.Welcome == nil {
		c.Welcome = ChDefaults.Welcome
	}

	// endregion: defaults
	// region: create channel

	ch := Ch{
		Config:  c,
		Encoder: c.Encoder,
		Type:    c.Type,
	}

	switch c.Type {
	case ChDb:
		return nil, ErrNotImplementedYet
	case ChFile:
		if c.File == nil || c.File == "" {
			return nil, ErrInvalidFile
		}

		switch c.File.(type) {
		case *os.File:
			ch.Inst = log.New(c.File.(io.Writer), *c.Prefix, *c.Flags)
			ch.File = c.File.(*os.File)
		case string:
			f, err := os.OpenFile(c.File.(string), os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.FileMode(*c.FilePerm))
			if err != nil {
				return nil, err
			}
			if _, err := f.WriteString(""); err != nil {
				f.Close() // ignore error; Write error takes precedence
				return nil, err
			}
			ch.Inst = log.New(f, *c.Prefix, *c.Flags)
			ch.File = f
		default:
			return nil, fmt.Errorf("%s: c.File=%s, (%T)", ErrInvalidFile, c.File, c.File)
		}
	case ChSyslog:
		inst, err := syslog.NewLogger(*c.Severity|*c.Facility, *c.Flags)
		if err != nil {
			return nil, err
		}
		ch.Inst = inst
	default:
		return nil, fmt.Errorf("%s: %v", ErrInvalidLoggerOrChannel, c.Type)
	}

	// endregion: channel
	// region: welcome and back

	if c.Welcome != nil {
		ch.Out(*c.Welcome)
	}
	return &ch, nil

	// endregion: welcome and back

}

func (l *Logger) NewCh(cs ...ChConfig) (ch *Ch, e error) {
	ch, e = NewCh(cs...)
	if e != nil {
		return nil, e
	}
	l.Ch = append(l.Ch, ch)
	return
}

func (c *Ch) Close() (e error) {
	if c.Type != ChFile {
		return ErrNotImplementedYet
	}
	if c.Config.Bye != nil {
		c.Out(*c.Config.Bye)
	}
	if e := c.File.Close(); e != nil {
		return e
	}
	return nil
}

// endregion: destinations and destructors

// endregion: constructor
// region: encoders

var EncoderFlat Encoder = func(c *Ch, n ...interface{}) (s string, e error) {

	// prefix with severity label, if needed
	if severity, ok := n[0].(syslog.Priority); ok {
		labels := *c.Config.SeverityLabels
		label := labels[severity]
		s = label + s
		_, n = n[0], n[1:]
	}

	// encode
	for _, v := range n {
		s = fmt.Sprintf("%s%s%+v", s, *c.Config.Delimiter, v)
	}
	s = strings.Replace(s, *c.Config.Delimiter, "", 1)

	// done
	return s, nil
}

// endregion: encoders
// region: output

func (c *Ch) Out(s ...interface{}) (e error) {
	// set depth
	depth := 2 + *c.Config.Depth
	hood := Trace(depth)
	for k, v := range hood {
		if v.File != hood[0].File {
			depth = depth + k - 1
			break
		}
	}

	// check severity
	severity, severityOk := s[0].(syslog.Priority)
	if severityOk {
		if severity > LOG_DEBUG {
			return ErrInvalidSeverity
		}
		if *c.Config.Severity < severity {
			return nil
		}
	}

	// encode and out
	o, e := Encoder(*c.Encoder)(c, s...)
	if e != nil {
		c.Inst.Output(depth, e.Error())
	}

	switch c.Type {
	case ChDb:
		return ErrNotImplementedYet
	case ChFile:
		c.Inst.Output(depth, o)
	case ChSyslog:
		if severityOk {
			writer := c.Inst.Writer().(*syslog.Writer)
			switch severity {
			case LOG_EMERG:
				writer.Emerg(o)
			case LOG_ALERT:
				writer.Alert(o)
			case LOG_CRIT:
				writer.Crit(o)
			case LOG_ERR:
				writer.Err(o)
			case LOG_WARNING:
				writer.Warning(o)
			case LOG_NOTICE:
				writer.Notice(o)
			case LOG_INFO:
				writer.Info(o)
			case LOG_DEBUG:
				writer.Debug(o)
			default:
				return ErrInvalidSeverity
			}
		} else {
			c.Inst.Output(depth, o)
		}
	default:
		return fmt.Errorf("%s: %v", ErrInvalidLoggerOrChannel, c.Type)
	}
	return
}

func (l *Logger) Out(s ...interface{}) *[]error {
	es := make([]error, 0)
	for _, c := range l.Ch {
		e := c.Out(s...)
		if e != nil {
			es = append(es, e)
		}
	}
	if len(es) == 0 {
		return nil
	}
	return &es
}

func Out(c interface{}, p syslog.Priority, s ...interface{}) *[]error {

	// prepare return slice
	es := make([]error, 0)

	// validate severity
	if p > LOG_DEBUG {
		es = append(es, ErrInvalidSeverity)
		return &es
	}
	s = append([]interface{}{p}, s...)

	// switch output depending on channel type
	switch c.(type) {
	case *Ch:
		ch := c.(*Ch)
		e := ch.Out(s...)
		if e != nil {
			es = append(es, e)
		}
	case *Logger:
		l := c.(*Logger)
		e := l.Out(s...)
		if e != nil {
			es = append(es, *e...)
		}
	default:
		es = append(es, ErrInvalidLoggerOrChannel)
		return &es
	}

	// done
	if len(es) == 0 {
		return nil
	}
	return &es
}

// endregion: output
