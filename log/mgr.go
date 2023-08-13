package log

import (
	"io"
	"os"
	"sync"
)

type _mgr struct {
	mux sync.RWMutex
	m   map[string]*Log
}

var mgr *_mgr

func init() {
	mgr = &_mgr{
		mux: sync.RWMutex{},
		m: map[string]*Log{
			"_": NewLog("root", os.Stdout),
		},
	}
}

// Regist regist
func Regist(l *Log) {
	mgr.mux.Lock()
	defer mgr.mux.Lock()
	mgr.m[l.name] = l
}

// Get get
func Get(name string) *Log {
	mgr.mux.RLock()
	defer mgr.mux.RUnlock()
	return mgr.m[name]
}

// SetWriter setwriter
func SetWriter(w io.Writer) *Log {
	return Get("_").SetWriter(w)
}

// WithName  withName
func WithName(name string, opts ...Option) *Log {
	return Get("_").WithName(name, opts...)
}

// WithValues with values
func WithValues(k string, v interface{}) *Log {
	return Get("_").WithValues(k, v)
}

// Debugf debugf
func Debugf(format string, v ...interface{}) {
	Get("_").Debugf(format, v...)
}

// Infof infof
func Infof(format string, v ...interface{}) {
	Get("_").Infof(format, v...)
}

// Warnf warnf
func Warnf(format string, v ...interface{}) {
	Get("_").Warnf(format, v...)
}

// Errorf errorf
func Errorf(format string, v ...interface{}) {
	Get("_").Errorf(format, v...)
}

// Close if need
func Close() error {
	return Get("_").Close()
}
