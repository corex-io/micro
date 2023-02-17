package datatypes

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/corex-io/codec"
)

type myTime struct {
	Timestamp int64     `json:"Timestamp,omitempty"`
	CreatedAt Time      `json:"CreatedAt,omitempty"`
	UpdatedAt time.Time `json:"UpdatedAt,omitempty"`
	ABC       string    `json:"ABC,omitempty"`
}

func TestTime(t *testing.T) {
	c := time.Now()
	my := myTime{Timestamp: time.Now().Unix(), UpdatedAt: c}
	t.Logf("%#v", my)
	data := make(map[string]interface{})
	if err := codec.Format(&data, my); err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", data)
	// t.Log(my.CreatedAt.Time().IsZero())

	b, err := json.Marshal(my)
	t.Logf("marshal: %s, %v", string(b), err)

	ts := Date(2020, time.Month(12), 1, 1, 1, 1, 1, time.Local)
	// var ts *Time
	bs := ts.After(Now())
	fmt.Printf("%s, %v\n", ts, bs)

	cc := NewTime(time.Time{})
	fmt.Printf("Zero: %s\n", cc)

}
