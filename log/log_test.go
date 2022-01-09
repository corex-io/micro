package log

import (
	"fmt"
	"os"
	"testing"
)

func Test_Log(t *testing.T) {
	log := NewLog("", os.Stdout, Format("json"))
	log1 := log.WithValues("222", "444")
	log1.Debugf("111")
	log2 := log1.WithValues("333", "555")
	log2.Debugf("2222")
	log3 := log.WithName("test", Format("Line"))

	log3.Infof("$$$")
}

func Benchmark_Log(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	log := NewLog("Benchmark", os.Stdout)

	for i := 0; i < b.N; i++ {
		log = log.WithName(fmt.Sprintf("bench-%d", i))
		log.Debugf("12345")
		// if err != nil {
		// 	b.Error(err)
		// }
	}
}

func Test_Appender(t *testing.T) {
	Infof("123")
	f, _ := os.OpenFile("t.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	SetWriter(f)
	Debugf("debug")

}
