package log

import (
	"fmt"
	"github.com/corex-io/micro/datatypes"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func openFile(elems ...string) (*os.File, error) {
	path := filepath.Join(elems...)
	dirname := filepath.Dir(path)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		panic(fmt.Sprintf("create directory %s, err=%v", dirname, err))
	}
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
}

type Rotate struct {
	Path     string
	truncate time.Duration
	w        io.Writer
	m        sync.RWMutex
}

func NewRotate(path string, truncate time.Duration) (*Rotate, error) {
	fpath := datatypes.Now().Format(path)
	f, err := openFile(fpath)
	if err != nil {
		return nil, err
	}
	rotate := &Rotate{Path: path, truncate: truncate, w: f}
	go rotate.RotateSched()
	return rotate, nil
}
func (r *Rotate) Write(p []byte) (n int, err error) {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.w.Write(p)
}

func (r *Rotate) RotateSched() {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-timer.C: // 按时间滚动
			if err := r.Rotate(); err != nil {
				panic(err)
			}
			interval := datatypes.Now().Truncate(r.truncate).Add(r.truncate).Sub(datatypes.Now())
			timer.Reset(interval)
		}
	}
}

func (r *Rotate) Rotate() error {
	r.m.Lock()
	defer r.m.Unlock()

	f, ok := r.w.(*os.File)
	if !ok {
		return fmt.Errorf("not file, no need ratate")
	}
	stat, err := f.Stat()
	if err != nil {
		return err
	}

	if datatypes.NewTime(stat.ModTime()).Truncate(r.truncate) == datatypes.Now().Truncate(r.truncate) {
		return nil
	}

	defer f.Close()

	fpath := datatypes.Now().Format(r.Path)
	r.w, err = openFile(fpath)
	return err
}
