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
func WithSpanId(ctx context.Context, sessionId string) context.Context {
	return context.WithValue(ctx, key("SpanId"), sessionId)
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
	defer m.Unlock()
	return strings.ToUpper(strconv.FormatUint(uint64(time.Now().UnixMicro()*1000+source.Int63n(1000)), 16))
}

func GenId2() string {
	return fmt.Sprintf("T%s%03d", time.Now().Format("20060102150405.000000"), source.Intn(1000))
}

var source = rand.New(rand.NewSource(time.Now().Unix()))
var m sync.Mutex
