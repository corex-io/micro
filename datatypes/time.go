package datatypes

import (
	"database/sql/driver"
	"fmt"
	"time"
	"unsafe"
)

// Since returns the time elapsed since t. It is shorthand for time.Now().Sub(t).
func Since(t *Time) time.Duration {
	return Now().Time().Sub(t.Time())
}

// Until returns the duration until t. It is shorthand for t.Sub(time.Now()).
func Until(t *Time) time.Duration {
	return t.Time().Sub(Now().Time())
}

const timeFormat = `"2006-01-02 15:04:05"`

const (
	Day   = 24 * time.Hour
	Week  = 7 * Day
	Month = 30 * Day
	Year  = 365 * Day
)

// Time time: 2006-01-02 15:04:05
type Time time.Time

// Now now
func Now() *Time {
	return NewTime(time.Now())
}

// NewTime newTime
func NewTime(t time.Time, layout ...string) *Time {
	ts := Time(t)
	return &ts
}

// Parse parses a formatted string and returns the time value it represents.
// The layout defines the format by showing )how the reference time, defined to be
func Parse(layout, value string) (*Time, error) {
	ts, err := time.Parse(layout, value)
	if err != nil {
		return nil, err
	}
	return NewTime(ts), nil
}

// ParseInLocation is like Parse but differs in two important ways.
// First, in the absence of time zone information, Parse interprets a time as UTC;
// ParseInLocation interprets the time as in the given location.
// Second, when given a zone offset or abbreviation, Parse tries to match it against the Local location;
// ParseInLocation uses the given location.
func ParseInLocation(layout, value string, loc *time.Location) (*Time, error) {
	ts, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return nil, err
	}
	return NewTime(ts), nil
}

func FromTimestamp(ts int64) *Time {
	return NewTime(time.Unix(ts/1000, (ts%1000)*int64(time.Millisecond)).UTC())
}

// Date returns the Time corresponding to
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) *Time {
	return NewTime(time.Date(year, month, day, hour, min, sec, nsec, loc))
}

// Add returns the time t+d.
func (t *Time) Add(d time.Duration) *Time {
	return NewTime(t.Time().Add(d))
}

// UTC returns t with the location set to UTC.
func (t *Time) UTC() *Time {
	return NewTime(t.Time().UTC())
}

// AddDate returns the time corresponding to adding the given number of years, months, and days to t.
// For example, AddDate(-1, 2, 3) applied to January 1, 2011 returns March 4, 2010.
// AddDate normalizes its result in the same way that Date does, so,
// for example, adding one month to October 31 yields December 1, the normalized form for November 31.
func (t *Time) AddDate(years int, months int, days int) *Time {
	return NewTime(t.Time().AddDate(years, months, days))
}

// After reports whether the time instant t is after u.
func (t *Time) After(u *Time) bool {
	return t.Time().After(u.Time())
}

// Before reports whether the time instant t is before u.
func (t *Time) Before(u *Time) bool {
	return t.Time().Before(u.Time())
}

// Clock returns the hour, minute, and second within the day specified by t.
func (t *Time) Clock() (hour, min, sec int) {
	return t.Time().Clock()
}

// Date returns the year, month, and day in which t occurs.
func (t *Time) Date() (year int, month time.Month, day int) {
	return t.Time().Date()
}

// Day returns the day of the month specified by t.
func (t *Time) Day() int {
	return t.Time().Day()
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with Time values;
// most code should use Equal instead.
func (t *Time) Equal(u *Time) bool {
	return t.Time().Equal(u.Time())
}

// Format returns a textual representation of the time value formatted according to layout,
// which defines the format by showing how the reference time, defined to be
func (t *Time) Format(layout string) string {
	return t.Time().Format(layout)
}

// Hour returns the hour within the day specified by t, in the range [0, 23].
func (t *Time) Hour() int {
	return t.Time().Hour()
}

// IsZero isZero
func (t *Time) IsZero() bool {
	return t.Time().IsZero()
}

// Minute returns the minute offset within the hour specified by t, in the range [0, 59].
func (t *Time) Minute() int {
	return t.Time().Minute()
}

// Month returns the month of the year specified by t.
func (t *Time) Month() time.Month {
	return t.Time().Month()
}

// Nanosecond returns the nanosecond offset within the second specified by t, in the range [0, 999999999].
func (t *Time) Nanosecond() int {
	return t.Time().Nanosecond()
}

// Second returns the second offset within the minute specified by t, in the range [0, 59].
func (t *Time) Second() int {
	return t.Time().Second()
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum) value that can be stored in a Duration,
// the maximum (or minimum) duration will be returned. To compute t-d for a duration d, use t.Add(-d).
func (t *Time) Sub(u *Time) time.Duration {
	return t.Time().Sub(u.Time())
}

// Year returns the year in which t occurs.
func (t *Time) Year() int {
	return t.Time().Year()
}

// Weekday returns the day of the week specified by t.
func (t *Time) Weekday() time.Weekday {
	return t.Time().Weekday()
}

// Location returns the time zone information associated with t.
func (t *Time) Location() *time.Location {
	return t.Time().Location()
}

// String returns the time formatted using the format string
func (t *Time) String() string {
	return t.Time().Format("2006-01-02 15:04:05")
}

// UnmarshalJSON unmarshaler interface
func (t *Time) UnmarshalJSON(b []byte) error {
	s := *(*string)(unsafe.Pointer(&b))
	if s == "null" || s == `""` || s == `"0000-00-00 00:00:00"` {
		*t = Time{}
		return nil
	}
	now, err := time.ParseInLocation(timeFormat, s, time.Local)
	*t = Time(now)
	return err
}

// MarshalJSON marshaler interface
func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil || t.Time().IsZero() {
		return []byte("null"), nil
	}
	s := t.Time().Format(timeFormat)
	return []byte(s), nil
}

// Value return json value, implement driver.Valuer interface
func (t *Time) Value() (driver.Value, error) {
	if t == nil || t.IsZero() {
		return nil, nil
	}
	// return t.Time().Format("2006-01-02 15:04:05"), nil
	return t.Time(), nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *Time) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)

}

// Time time
func (t *Time) Time() time.Time {
	return time.Time(*t)
}

// Unix unix
func (t *Time) Unix() int64 {
	return t.Time().Unix()
}

// Timestamp returns a new millisecond timestamp from a time.
func (t *Time) Timestamp() int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/int64(time.Millisecond)
}

// MarshalBinary implement encoding.BinaryMarshaler.
func (t Time) MarshalBinary() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalBinary implement encoding.BinaryUnmarshaler.
func (t *Time) UnmarshalBinary(b []byte) error {
	return t.UnmarshalJSON(b)
}

// Truncate returns a new time with the time truncated to the given number of seconds.
func (t *Time) Truncate(d time.Duration) *Time {
	if d/Year > 0 {
		return Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	} else if d/Month > 0 {
		return Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	} else if d/Day > 0 {
		return Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	} else if d/time.Hour > 0 {
		return Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local)
	} else if d/time.Minute > 0 {
		return Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
	}
	return t
}
