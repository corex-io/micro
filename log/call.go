package log

import (
	"runtime"
	"strings"
)

// Caller return the file name and line number of the current file
func Caller(extends ...string) (string, int) {
	// the second caller usually from gorm internal, so set i start from 2
	ignores := []string{"micro/log/log.go", "micro/log/mgr.go"}
	ignores = append(ignores, extends...)
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !hasSuffix(file, ignores...) {
			return file, line
		}
	}
	return "", 0
}

func hasSuffix(k string, vList ...string) bool {
	for _, v := range vList {
		if strings.HasSuffix(k, v) {
			return true
		}
	}
	return false
}
