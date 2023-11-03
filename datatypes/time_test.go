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

func TestTime_Truncate(t *testing.T) {
	now := Now()
	truncated := now.Truncate(time.Hour)
	fmt.Printf("%s\n", truncated)
}

func TestTime_IsZero(t *testing.T) {
	c, err := Parse("", "0000-21-01 00:00:00")
	t.Logf("%v, err=%v", c, err)
	z := Time{}
	err = z.UnmarshalJSON([]byte(`"0000-00-00 00:00:00"`))
	t.Logf("%v, %v", z, err)
}
