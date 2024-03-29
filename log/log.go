package log

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

type kv struct {
	key   string
	value any
}

// Log log
type Log struct {
	options Options
	name    string
	head    []kv
	w       io.Writer
}

// NewLog newLog
func NewLog(name string, writer io.Writer, opts ...Option) *Log {
	options := newOptions(opts...)
	return &Log{
		options: options,
		name:    name,
		w:       writer,
	}
}

// SetWriter set writer
func (log *Log) SetWriter(w io.Writer) *Log {
	newLog := log.copy()
	newLog.w = w
	return newLog
}

func (log *Log) copy() *Log {
	newLog := NewLog(log.name, log.w)
	newLog.options = log.options
	var head []kv
	head = append(head, log.head...)
	newLog.head = head
	return newLog
}

// WithName withName
func (log *Log) WithName(name string, opts ...Option) *Log {
	newLog := NewLog(name, log.w)
	newLog.options = withOptions(log.options, opts...)
	newLog.head = append(newLog.head, log.head...)
	return newLog
}

// WithValues withvalues
func (log *Log) WithValues(key string, value interface{}) *Log {
	newLog := log.copy()
	newLog.head = append(newLog.head, kv{key, value})
	return newLog
}

func (log *Log) outputLine(lv, format string, v ...interface{}) string {
	var msg []string
	fpath, no := Caller()

	msg = append(msg, fmt.Sprintf("%s[%s][%s][%s:%d]",
		lv,
		time.Now().Format(log.options.timeFormat),
		log.name,
		filepath.Base(fpath), no,
	))
	for _, kv := range log.head {
		msg = append(msg, fmt.Sprintf("%s=%v", kv.key, kv.value))
	}
	msg = append(msg, fmt.Sprintf(format, v...))
	return strings.Join(msg, " ")
}

type jsonLog struct {
	LogAt  string         `json:"LogAt"`
	Name   string         `json:"Name"`
	Lv     string         `json:"Lv"`
	Fields map[string]any `json:"Fields,omitempty"`
	Msg    string         `json:"Msg,omitempty"`
}

func (log *Log) outputJSON(lv, format string, v ...any) string {
	fields := make(map[string]any, len(log.head))
	for _, kv := range log.head {
		fields[kv.key] = kv.value
	}
	rets := jsonLog{
		LogAt:  time.Now().Format("2006/01/02 15:04:05.000"),
		Name:   log.name,
		Lv:     lv,
		Fields: fields,
		Msg:    fmt.Sprintf(format, v...),
	}

	msg, _ := json.Marshal(rets)
	return string(msg)
}

func (log *Log) output(lv, format string, v ...any) {
	var msg string
	switch log.options.msgFormat {
	case "json", "JSON":
		msg = log.outputJSON(lv, format, v...)
	default:
		msg = log.outputLine(lv, format, v...)
	}
	_, _ = fmt.Fprintln(log.w, msg)
}

// Close close
func (log *Log) Close() error {
	if w, ok := log.w.(io.WriteCloser); ok {
		return w.Close()
	}
	return nil
}

// Debugf debugf
func (log *Log) Debugf(format string, v ...any) {
	log.output("D", format, v...)
}

// Infof Infof
func (log *Log) Infof(format string, v ...any) {
	log.output("I", format, v...)
}

// Warnf Warnf
func (log *Log) Warnf(format string, v ...any) {
	log.output("W", format, v...)
}

// Errorf Errorf
func (log *Log) Errorf(format string, v ...any) {
	log.output("E", format, v...)
	debug.PrintStack()
}
