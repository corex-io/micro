package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGenId(t *testing.T) {
	var maxInt64 int64 = math.MaxInt64
	var maxUint64 uint64 = math.MaxUint64
	t.Logf("maxInt64=%d,  maxUint64=%d", maxInt64, maxUint64)
	now := time.Now()
	random := source.Int63n(1000)
	t.Logf("now=%s, ts=%d, random=%d", now, now.UnixMicro(), random)
	var sd uint64
	sd = uint64(now.UnixMicro()) * 1000
	t.Logf("sd__unit=%d", sd)
	t.Logf("max_uint=%d", maxUint64)
	id := uint64(now.UnixMicro())*1000 + uint64(random)

	t.Logf("id=%d", id)
	genId := GenId()
	id1 := strings.ToUpper(strconv.FormatUint(id, 16))
	id2 := fmt.Sprintf("T%s%03d", now.Format("20060102150405.000000"), random)
	t.Logf("%s, GenId=%s, id1=%s,  id2=%s", now, genId, id1, id2)

	i, err := strconv.ParseUint(genId, 16, 64)
	timestamp := time.UnixMicro(int64(i / 1000))
	t.Logf("%d, err=%v, timestamp=%s", i, err, timestamp)

}

// go test -v -bench='BenchmarkGenId' -benchmem .
func BenchmarkGenId1(b *testing.B) {
	now := time.Now()
	random := source.Int63n(1000000)
	for i := 0; i < b.N; i++ {
		_ = strconv.FormatUint(uint64(now.UnixMicro()*1000000+random), 16)
	}
}

// go test -v -bench='BenchmarkGenId' -benchmem .
func BenchmarkGenId2(b *testing.B) {
	now := time.Now()
	random := source.Int63n(1000000)
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%x", uint64(now.UnixMicro()*1000000+random))
	}
}
