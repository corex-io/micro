package common

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type key string

// WithRequestId append taskId to context
func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, key("RequestId"), requestId)
}

// GetRequestId GetTaskId
func GetRequestId(ctx context.Context) string {
	taskId, _ := ctx.Value(key("RequestId")).(string)
	return taskId
}

// WithSpanId append SpanId to context
func WithSpanId(ctx context.Context, spanId string) context.Context {
	return context.WithValue(ctx, key("SpanId"), spanId)
}

// GetSpanId GetSpanId
func GetSpanId(ctx context.Context) string {
	taskId, _ := ctx.Value(key("SpanId")).(string)
	return taskId
}

// GenId gen id
func GenId(id ...string) string {
	if len(id) != 0 && id[0] != "" {
		return id[0]
	}
	m.Lock()
	s := uint64(time.Now().UnixMicro()*1000 + source.Int63n(1000)) // % 4738381338321616895
	m.Unlock()
	return strings.ToUpper(strconv.FormatUint(s, 36))
}

func ParseId(id string) int64 {
	tsm, err := strconv.ParseUint(id, 36, 64)
	if err != nil {
		return 0
	}
	return int64(tsm / 1000)
}

func GenId2() string {
	return fmt.Sprintf("T%s%03d", time.Now().Format("20060102150405.000000"), source.Intn(1000))
}

var source = rand.New(rand.NewSource(time.Now().Unix()))
var m sync.Mutex
