package common

import (
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGenId(t *testing.T) {
	//var maxInt64 int64 = math.MaxInt64
	//var maxUint64 uint64 = math.MaxUint64
	//t.Logf("maxInt64=%d,  maxUint64=%d", maxInt64, maxUint64)
	//now := time.Now()
	//random := source.Int63n(1000)
	//t.Logf("now=%s, ts=%d, random=%d", now, now.UnixMicro(), random)
	//var sd uint64
	//sd = uint64(now.UnixMicro()) * 1000
	//t.Logf("sd__unit=%d", sd)
	//t.Logf("max_uint=%d", maxUint64)
	//id := uint64(now.UnixMicro())*1000 + uint64(random)
	//
	//t.Logf("id=%d", id)
	//genId := GenId()
	//id1 := strings.ToUpper(strconv.FormatUint(id, 16))
	//id2 := fmt.Sprintf("T%s%03d", now.Format("20060102150405.000000"), random)
	//t.Logf("%s, GenId=%s, id1=%s,  id2=%s", now, genId, id1, id2)
	//
	//i, err := strconv.ParseUint(genId, 16, 64)
	//timestamp := time.UnixMicro(int64(i / 1000))
	//t.Logf("%d, err=%v, timestamp=%s", i, err, timestamp)
	now := time.Now()
	ms := now.UnixMicro()
	t.Logf("unix_micro=%d", ms) // 16位
	ms_16, ms_32, ms_max := strconv.FormatUint(uint64(ms), 16), strconv.FormatUint(uint64(ms), 32), strconv.FormatUint(uint64(math.MaxInt64), 16)
	t.Logf("ms_16=%s, ms_32=##%019s##, ms_max=%s", ms_16, strings.ToUpper(ms_32), ms_max) // 16位

}

func TestGenId2(t *testing.T) {
	//id := GenId()
	//t.Logf("%s", id)
	now := time.Now()
	ts := uint64(now.UnixMicro()*1000+source.Int63n(1000)) % 4738381338321616895
	id1 := strings.ToUpper(strconv.FormatUint(ts, 16))
	id2 := strings.ToUpper(strconv.FormatUint(ts, 32))
	id3 := strings.ToUpper(strconv.FormatUint(ts, 34))
	id4 := strings.ToUpper(strconv.FormatUint(ts, 36))

	t.Logf("%s", id1)
	t.Logf("%s", id2)
	t.Logf("%s", id3)
	t.Logf("%s", id4)

	kk, err := strconv.ParseUint(id4, 36, 64)
	t.Logf("ts1=%d, kk=%d, %v", ts, kk, err)

	kk2, err2 := strconv.ParseUint("ZZZZZZZZZZZZ", 36, 64)
	t.Logf("kk2=%d, %v", kk2, err2)

	kk3, err3 := strconv.ParseUint("FFFFFFFF", 16, 64)
	t.Logf("kk3=%d, %v", kk3, err3)

}

// go test -v -bench='BenchmarkGenId' -benchmem .
func BenchmarkGenId1(b *testing.B) {

	for i := 0; i < b.N; i++ {
		GenId()
	}
}

// go test -v -bench='BenchmarkGenId' -benchmem .
func BenchmarkGenId2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenId11()
	}
}
