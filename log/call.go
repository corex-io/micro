package log

import (
	"fmt"
	"runtime"
	"strings"
)

// call return the file name and line number of the current file
func call() (string, int) {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		fmt.Println(file, line, ok)
		if ok && !(strings.HasSuffix(file, "micro/log/log.go") || strings.HasSuffix(file, "micro/log/mgr.go")) {
			return file, line
		}
	}
	return "", 0
}
