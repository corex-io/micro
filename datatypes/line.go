package datatypes

import (
	"encoding/json"
	"time"
)

// Line Line
type Line []*float64

func NewLine(p ...float64) Line {
	var v Line
	for _, c := range p {
		c := c
		v = append(v, &c)
	}
	return v
}

// PointAt PointAt
func (v Line) PointAt(startAt *Time, index, period int) *Time {
	return startAt.Add(time.Duration(index*period) * time.Second)
}

// String String
func (v Line) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

// Cut 去除尾部的无效值，check函数判断是否有效，有效=true，无效=false
func (v Line) Cut(check func(v *float64) bool) Line {
	pos := len(v)
	for i := pos - 1; i >= 0; i-- {
		if check(v[i]) {
			break
		}
		pos = i

	}
	return v[:pos]
}

// RangeMatch 设置一个条件, 统计连续满足改条件的数据范围取值有哪些
// check函数判断是否有效，有效=true，无效=false
// minRange int 最小的持续区间
func (v Line) RangeMatch(check func(v *float64) bool, minRange int) [][]int {
	s, e, max, ret := 0, 0, len(v), make([][]int, 0, 0)
	for e < max {
		for e < max && !check(v[e]) {
			e++
		}
		s = e
		for e < max && check(v[e]) {
			e++
		}

		if e-s >= minRange {
			ret = append(ret, []int{s, e})
		}
	}
	return ret
}
