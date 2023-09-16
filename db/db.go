package db

import (
	"context"
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// DB my
type DB struct {
	// dsn string
	db *gorm.DB
}

// WithDB with db
func WithDB(db *gorm.DB) *DB {
	return &DB{db}
}

// New new
func New(driver, dsn string, opts ...Option) (*DB, error) {
	ctxLog.Infof("connect: %s %s", driver, dsn)

	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "clickhouse":
		dialector = clickhouse.Open(dsn)
	default:
		return nil, fmt.Errorf("driver=%s is invalid", driver)
	}

	config := &gorm.Config{
		SkipDefaultTransaction: true,  // 不需要对单次写入操作使用事务
		PrepareStmt:            false, // 缓存预编译语句
		Logger:                 newLog(),
	}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, fmt.Errorf("connect DB fail, dsn=%s, err=%w", dsn, err)
	}
	return WithDB(db), nil
}

// Close connection
func (my *DB) Close() error {
	return nil
}

// DB db
func (my *DB) DB(ctx context.Context) *gorm.DB {
	return my.db.WithContext(ctx)
}

// SetMaxOpenConns SetMaxOpenConns
func (my *DB) SetMaxOpenConns(max int) error {
	db, err := my.db.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(max)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(max)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Hour)

	return nil
}

// Model model
type Model interface {
	TableName() string
}
