package datatypes

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type Timestamp int64

func (ts *Timestamp) Value() (driver.Value, error) {
	return time.Unix(int64(*ts), 0).Format("2006-01-02 15:04:03"), nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (ts *Timestamp) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if !ok {
		return fmt.Errorf("can not convert %v to timestamp, %#v", v, v)
	}
	*ts = Timestamp(value.Unix())
	return nil
}

func (ts *Timestamp) Time() *Time {
	return NewTime(time.Unix(int64(*ts), 0))
}

func (ts *Timestamp) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(*ts), 10)), nil
}

// UnmarshalBinary implement encoding.BinaryUnmarshaler.
func (ts *Timestamp) UnmarshalBinary(b []byte) error {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*ts = Timestamp(i)
	return nil
}
