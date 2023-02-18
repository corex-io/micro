package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/corex-io/codec"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
func New(driver, dsn string) (*DB, error) {
	ctxLog.Infof("connect: %s %s", driver, dsn)

	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("driver=%s is invalid", driver)
	}

	config := &gorm.Config{
		SkipDefaultTransaction: true,  // 不需要对单次写入操作使用事务
		PrepareStmt:            false, // 缓存预编译语句
		// Logger:                 &dbLog{Deepth: 5},
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

// Create 创建一条记录
func (my *DB) Create(v interface{}, opts ...Option) (int64, error) {
	options := newOptions(opts...)
	ret := my.WithOptions(options).Create(v)
	return ret.RowsAffected, ret.Error

}

// CreateMany 创建多条记录
func (my *DB) CreateMany(vs []interface{}, opts ...Option) (int64, error) {
	options := newOptions(opts...)
	ret := my.WithOptions(options).Create(vs)
	return ret.RowsAffected, ret.Error

}

// Remove remove
func (my *DB) Remove(cond interface{}, opts ...Option) (int64, error) {
	options := newOptions(opts...)
	db, err := my.withWhere(cond, options)
	if err != nil {
		return 0, err
	}
	ret := db.Delete(nil)
	return ret.RowsAffected, ret.Error
}

// Find 查询多条记录
// cond 查询条件
// results 返回结果
func (my *DB) Find(cond interface{}, results interface{}, opts ...Option) (int64, error) {
	options := newOptions(opts...)
	options.Omit = nil
	db, err := my.withWhere(cond, options)
	if err != nil {
		return 0, err
	}
	ret := db.Find(results)
	return ret.RowsAffected, ret.Error
}

// Update updateOrUpsert
func (my *DB) Update(cond, values interface{}, opts ...Option) (int64, error) {
	options := newOptions(opts...)
	db, err := my.withWhere(cond, options)
	if err != nil {
		return 0, err
	}

	ret := db.Updates(values)
	return ret.RowsAffected, ret.Error

}

// IsExist isExist
func (my *DB) IsExist(cond interface{}, opts ...Option) (bool, error) {
	cnt, err := my.Count(cond, opts...)
	if err != nil {
		return false, err
	}
	return cnt != 0, nil
}

// Upsert upset upsert
func (my *DB) Upsert(cond, values interface{}, opts ...Option) (int64, error) {
	exists, err := my.IsExist(cond, opts...)
	if err != nil {
		return 0, err
	}

	if v, ok := values.(Model); ok {
		opts = append(opts, Table(v.TableName()))
	}

	if !exists {

		if cnt, err := my.Create(cond, opts...); err != nil {
			return cnt, err
		}
	}
	return my.Update(cond, values, opts...)
}

// RawQuery raw qurey
func (my *DB) RawQuery(ctx context.Context,
	condFunc func() (string, []interface{}, error),
	results interface{},
	opts ...Option,
) (int64, error) {
	options := newOptions(opts...)
	sql, values, err := condFunc()
	if err != nil {
		return 0, fmt.Errorf("genSQL, err=%w", err)
	}
	ret := my.WithOptions(options).Raw(sql, values...).Find(results)
	return ret.RowsAffected, ret.Error
}

// Count count
func (my *DB) Count(cond interface{}, opts ...Option) (int64, error) {
	options := newOptions(opts...)
	db, err := my.withWhere(cond, options)
	if err != nil {
		return 0, nil
	}
	var cnt int64
	err = db.Count(&cnt).Error
	return cnt, err
}

// WithOptions WithOptions
func (my *DB) WithOptions(options Options, opts ...Option) *gorm.DB {
	db := my.db.Omit(options.Omit...)

	if options.Table != "" {
		db = db.Table(options.Table)
	}
	if options.Limit != 0 {
		db = db.Limit(options.Limit)
	}
	if options.Debug {
		db = db.Debug()
	}

	return db
}

// withWhere 生成where条件, 已经包含了withTable
func (my *DB) withWhere(cond interface{}, options Options) (*gorm.DB, error) {
	query := make(map[string]interface{})
	if err := codec.Format(&query, cond); err != nil {
		return nil, fmt.Errorf("withWhere: %w", err)
	}
	db := my.WithOptions(options)
	for k, v := range query {
		sv := fmt.Sprintf("%v", v)
		switch {
		case strings.HasPrefix(sv, "~"):
			db = db.Where(fmt.Sprintf("`%s` Like ?", k), "%"+sv[1:]+"%")
		case strings.HasPrefix(sv, ">"):
			db = db.Where(fmt.Sprintf("`%s` > ?", k), sv[1:])
		case strings.HasPrefix(sv, "<"):
			db = db.Where(fmt.Sprintf("`%s` < ?", k), sv[1:])
		case strings.HasPrefix(sv, "!"):
			db = db.Where(fmt.Sprintf("`%s` != ?", k), sv[1:])
		default:
			db = db.Where(fmt.Sprintf("`%s` = ?", k), v)
		}
	}
	return db, nil
}

// Transaction 事务
func (my *DB) Transaction(fc func(my *DB) error, opts ...Option) error {
	db := my.WithOptions(newOptions(opts...))
	return db.Transaction(func(tx *gorm.DB) error {
		return fc(WithDB(tx))
	})
}

// Model model
type Model interface {
	TableName() string
}
