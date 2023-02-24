package datatypes

import (
	"encoding/json"
	"time"
)

// Line Line
type Line []*float64

// PointAt PointAt
func (line Line) PointAt(startAt *Time, index, period int) *Time {
	return startAt.Add(time.Duration(index*period) * time.Second)
}

// String String
func (v Line) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

// RangeMatch 设置一个条件, 统计连续满足改条件的数据范围取值有哪些
// minRange int 最小的持续区间
func (v Line) RangeMatch(f func(v *float64) bool, minRange int) [][]int {
	s, e, k := 0, 0, "s"
	var ret [][]int
	for i, x := range v {
		if f(x) && k == "s" {
			s = i
			k = "e"
		} else if f(x) && k == "e" {
			e = i
		} else {
			if e != 0 && e-s >= minRange {
				ret = append(ret, []int{s, e})
			}
			s, e, k = 0, 0, "s"
		}
	}
	return ret
}

func (v Line) RangeMatchCount(f func(v *float64) bool, minRange int) int {
	return len(v.RangeMatch(f, minRange))
}
