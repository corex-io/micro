package db

import (
	"fmt"
	"sync"
)

// Config config
type Config struct {
	DriverName string `json:"DriverName" yaml:"DriverName"`
	UUID       string `json:"UUID" yaml:"UUID"`
	DSN        string `json:"DSN" yaml:"DSN"`
	Database   string `json:"Database" yaml:"Database"`
}

// 不采用懒加载方式的目的在于, 尽可能早的把问题暴露出来, 配置的数据库信息必须是正确的,否在程序会Panic
type dbmgr struct {
	mux *sync.RWMutex
	dbs map[string]*DB
}

// Load Load
func (mgr *dbmgr) Load(configs ...Config) error {
	mgr.mux.Lock()
	defer mgr.mux.Unlock()
	for _, config := range configs {
		if config.UUID == "" || config.DSN == "" {
			panic(fmt.Sprintf("config is invalid, config=%#v", config))
		}
		db, err := New(config.DriverName, config.DSN)
		if err != nil {
			return err
		}
		if err:=db.SetMaxOpenConns(200); err != nil {
			return err
		}
		mgr.dbs[config.UUID] = db
	}
	return nil
}

// Get get
func (mgr *dbmgr) Get(uuid string) *DB {
	mgr.mux.RLock()
	defer mgr.mux.RUnlock()
	storage, ok := mgr.dbs[uuid]
	if !ok {
		panic(fmt.Sprintf("uuid=%s not found", uuid))
	}
	return storage
}

var mgr *dbmgr

func init() {
	mgr = &dbmgr{
		mux: &sync.RWMutex{},
		dbs: make(map[string]*DB),
	}
}

// Load Load
func Load(config Config) error {
	return mgr.Load(config)
}

// Get Get
func Get(uuid string) *DB {
	return mgr.Get(uuid)
}

// Close close
func Close() error {
	for _, db := range mgr.dbs {
		db.Close()
	}
	return nil
}
