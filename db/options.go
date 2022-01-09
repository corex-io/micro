package db

import (
	"time"
)

// Options options
type Options struct {
	Debug     bool
	Database  string // mysql: database; mongodb;database
	Table     string // mysql: table; mongodb;collection
	CreatedAt time.Time
	UpdatedAt time.Time
	Limit     int
	Omit      []string // 忽略字段; mysql:Omit(columns...)
	Tx        *DB
}

// Option func
type Option func(*Options)

// NewOptions new options
func NewOptions(opts ...Option) Options {
	return newOptions(opts...)
}

// NewOptions new Options
func newOptions(opts ...Option) Options {
	opt := Options{
		//Omit: []string{"CreatedAt", "UpdatedAt", "DeletedAt"},
	}
	for i := len(opts) - 1; i >= 0; i-- {
		opts[i](&opt)
	}
	return opt
}

// Tx tx
func Tx(tx *DB) Option {
	return func(o *Options) {
		o.Tx = tx
	}
}

// Debug debug
func Debug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}

// Database database
func Database(database string) Option {
	return func(o *Options) {
		o.Database = database
	}
}

// Table => mysql: tableName; mongodb;collection
func Table(table string) Option {
	return func(o *Options) {
		o.Table = table
	}
}

// CreatedAt Set createdAt
func CreatedAt() Option {
	return func(o *Options) {
		o.CreatedAt = time.Now()
	}
}

// UpdatedAt set updatedAt
func UpdatedAt() Option {
	return func(o *Options) {
		o.UpdatedAt = time.Now()
	}
}

// Limit set limit
func Limit(limit int) Option {
	return func(o *Options) {
		o.Limit = limit
	}
}

// Omit omit
func Omit(omit ...string) Option {
	return func(o *Options) {
		o.Omit = omit
	}
}
