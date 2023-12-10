package db

import (
	"context"
	"github.com/corex-io/micro/common"
	"github.com/corex-io/micro/log"
	gormLogger "gorm.io/gorm/logger"
	"path"
	"time"
)

// Log gormLogger
type dbLog struct {
	Lv gormLogger.LogLevel
	w  *log.Log
}

func newLog() *dbLog {
	return &dbLog{
		w: log.WithName("db"),
	}
}

// LogMode implement gorm/logger.Interface
func (l *dbLog) LogMode(lv gormLogger.LogLevel) gormLogger.Interface {
	newlogger := *l
	newlogger.Lv = lv
	return &newlogger
}

// Info implement gorm/logger.Interface
func (l *dbLog) Info(ctx context.Context, msg string, v ...interface{}) {
	if l.Lv >= gormLogger.Info {
		l.w.WithName(common.GetRequestId(ctx)).Infof(msg, v...)
	}
}

// Warn implement gorm/logger.Interface
func (l *dbLog) Warn(ctx context.Context, msg string, v ...interface{}) {
	if l.Lv >= gormLogger.Warn {
		l.w.WithName(common.GetRequestId(ctx)).Warnf(msg, v...)
	}
}

// Error implement gorm/logger.Interface
func (l *dbLog) Error(ctx context.Context, msg string, v ...interface{}) {
	if l.Lv >= gormLogger.Error {
		l.w.WithName(common.GetRequestId(ctx)).Errorf(msg, v...)
	}
}

// Trace implement gorm/logger.Interface
func (l *dbLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Lv <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	// if strings.HasPrefix(sql, "SELECT") && err == nil {
	// 	return
	// }
	file, line := log.Caller("micro/db/db.go", "micro/db/log.go", "micro/db/mgr.go")
	switch {
	case err != nil && l.Lv >= gormLogger.Error:
		l.w.WithName(common.GetRequestId(ctx)).Warnf("%s:%d %s, rows=%d, cost=%s, err=%v", path.Base(file), line, sql, rows, time.Since(begin), err)
	case elapsed > 3000*time.Microsecond && l.Lv >= gormLogger.Warn:
		l.w.WithName(common.GetRequestId(ctx)).Warnf("%s:%d %s, rows=%d, cost=%s", path.Base(file), line, sql, rows, time.Since(begin))
	case l.Lv >= gormLogger.Info:
		l.w.WithName(common.GetRequestId(ctx)).Infof("%s:%d %s, rows=%d, cost=%s", path.Base(file), line, sql, rows, time.Since(begin))
	}
}
