package kv

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/corex-io/micro/log"

	"go.etcd.io/bbolt"
)

var ctxLog = log.WithName("kv")

// DB db
type DB struct {
	*bbolt.DB
}

// New new
func New(driver, dsn string) (*DB, error) {
	ctxLog.Infof("connect[%s]: %s", driver, dsn)

	db, err := bbolt.Open(dsn, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("connect [%s] dsn=%s, err=%w", driver, dsn, err)
	}
	return &DB{db}, nil
}

// Close close connection
func (my *DB) Close() error {
	return my.DB.Close()
}

// Stats status
func (my *DB) Stats() {
	go func() {
		// Grab the initial stats.
		prev := my.DB.Stats()
		for {

			// Grab the current stats and diff them.
			stats := my.DB.Stats()
			diff := stats.Sub(&prev)
			// Encode stats to JSON and print to STDERR.
			json.NewEncoder(os.Stderr).Encode(diff)
			// Save stats for the next loop.
			prev = stats
			// Wait for 10s.
			time.Sleep(5 * time.Second)
		}
	}()
}

// Put save
func (my *DB) Put(bucket, key string, v interface{}, opts ...Option) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return my.DB.Update(func(tx *bbolt.Tx) error {
		_bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return _bucket.Put([]byte(key), b)
	})
}

// Get get
func (my *DB) Get(bucket, key string, v interface{}, opts ...Option) error {
	return my.DB.View(func(tx *bbolt.Tx) error {
		_bucket := tx.Bucket([]byte(bucket))
		if _bucket == nil {
			return fmt.Errorf("bucket[%s] NotFound", bucket)
		}

		b := _bucket.Get([]byte(key))
		if b == nil {
			return fmt.Errorf("bucket[%s]/key[%s] NotFound", bucket, key)
		}
		return json.Unmarshal(b, v)
	})
}

// Delete delete
func (my *DB) Delete(bucket string, key string, opts ...Option) error {
	return my.DB.Update(func(tx *bbolt.Tx) error {
		_bucket := tx.Bucket([]byte(bucket))
		if _bucket == nil {
			return fmt.Errorf("bucket[%s] NotFound", bucket)
		}
		return _bucket.Delete([]byte(key))
	})
}

// ForEach foreach
func (my *DB) ForEach(bucket string, fc func(k, v []byte) error) error {
	return my.DB.View(func(tx *bbolt.Tx) error {
		_bucket := tx.Bucket([]byte(bucket))
		if _bucket == nil {
			return fmt.Errorf("bucket[%s] NotFound", bucket)
		}
		return _bucket.ForEach(fc)
	})
}

// First 获取第一个元素
//db.View(func(tx *bolt.Tx) error {
// 	b := tx.Bucket([]byte("MyBucket"))
// 	c := b.Cursor()
// 	for k, v := c.First(); k != nil; k, v = c.Next() {
// 		fmt.Printf("key=%s, value=%s\n", k, v)
// 	}
// 	return nil
// })

// First first
func (my *DB) First(bucket string, v interface{}) (string, error) {
	var key string
	return key, my.DB.View(func(tx *bbolt.Tx) error {
		_bucket := tx.Bucket([]byte(bucket))
		k, b := _bucket.Cursor().First()
		if k == nil {
			return io.EOF
		}
		key = string(k)
		return json.Unmarshal(b, v)

	})
}

// KeyList return keys in bucket
func (my *DB) KeyList(bucket string, filter func(key string) (bool, error)) ([]string, error) {
	var keyList []string
	return keyList, my.DB.View(func(tx *bbolt.Tx) error {
		_bucket := tx.Bucket([]byte(bucket))
		if _bucket == nil {
			return fmt.Errorf("bucket[%s] NotFound", bucket)
		}

		return _bucket.ForEach(func(k, v []byte) error {
			key := string(k)
			ok, err := filter(key)
			if err != nil {
				return err
			}
			if ok {
				keyList = append(keyList, key)
			}
			return nil
		})
	})
}

// Dump dump
func (my *DB) Dump() error {
	return nil
}
