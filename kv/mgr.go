package kv

import (
	"fmt"

	"go.etcd.io/bbolt"
)

var _db *DB

// Get get
func Get(bucket, key string, v interface{}, opts ...Option) error {
	return _db.Get(bucket, key, v, opts...)
}

// Put put
func Put(bucket, key string, v interface{}, opts ...Option) error {
	return _db.Put(bucket, key, v, opts...)
}

// First first
func First(bucket string, v interface{}, opts ...Option) (string, error) {
	return _db.First(bucket, v)
}

// Delete delete
func Delete(bucket, key string, opts ...Option) error {
	return _db.Delete(bucket, key, opts...)
}

// KeyList foreach
func KeyList(bucket string, filter func(key string) (bool, error)) ([]string, error) {
	return _db.KeyList(bucket, filter)
}

// Init init
func Init(path string, buckets ...string) error {
	var err error
	_db, err = New("bbolt", path)
	if err != nil {
		return err
	}

	_db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}

		return nil
	})

	_db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("InitSeed"))
		for i := 0; i < 10; i++ {
			seq, err := b.NextSequence()
			fmt.Println(seq, err)
		}

		return fmt.Errorf("ddd")
	})

	return err
}

// Stats stats
func Stats() {
	_db.Stats()
}

// Close close
func Close() error {
	return _db.Close()
}
