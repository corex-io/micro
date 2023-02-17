package db

import (
	"context"
	"path"
	"runtime"
	"time"

	"github.com/corex-io/micro/log"

	gormLogger "gorm.io/gorm/logger"
)

var ctxLog = log.WithName("db")

// Log gormLogger
type dbLog struct {
	Lv     gormLogger.LogLevel
	Deepth int
}

// LogMode implement gorm/logger.Interface
func (log *dbLog) LogMode(lv gormLogger.LogLevel) gormLogger.Interface {
	// log.Lv = lv data race
	return log
}

// Info implement gorm/logger.Interface
func (log *dbLog) Info(ctx context.Context, msg string, v ...interface{}) {
	if log.Lv >= gormLogger.Info {
		ctxLog.Infof(msg, v...)
	}
}

// Warn implement gorm/logger.Interface
func (log *dbLog) Warn(ctx context.Context, msg string, v ...interface{}) {
	if log.Lv >= gormLogger.Warn {
		ctxLog.Warnf(msg, v...)
	}
}

// Error implement gorm/logger.Interface
func (log *dbLog) Error(ctx context.Context, msg string, v ...interface{}) {
	if log.Lv >= gormLogger.Error {
		ctxLog.Errorf(msg, v)
	}
}

// Trace implement gorm/logger.Interface
func (log *dbLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if log.Lv <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	// if strings.HasPrefix(sql, "SELECT") && err == nil {
	// 	return
	// }
	_, file, line, _ := runtime.Caller(log.Deepth)
	switch {
	case err != nil && log.Lv >= gormLogger.Error:
		ctxLog.Warnf("%s:%d %s, rows=%d, cost=%s, err=%v", path.Base(file), line, sql, rows, time.Since(begin), err)
	case elapsed > 200*time.Microsecond && log.Lv >= gormLogger.Warn:
		ctxLog.Warnf("%s:%d %s, rows=%d, cost=%s", path.Base(file), line, sql, rows, time.Since(begin))
	case log.Lv >= gormLogger.Info:
		ctxLog.Infof("%s:%d %s, rows=%d, cost=%s", path.Base(file), line, sql, rows, time.Since(begin))
	}
}
