package datatypes

import "testing"

func TestLine_RangeMatch(t *testing.T) {
	v := NewLine([]float64{-1, 1, -1, 2, 3, 4, 5}...)
	t.Logf("%#v", v)
	t.Logf("%s", v.String())
	ret := v.RangeMatch(func(v *float64) bool {
		return v != nil && *v != -1
	}, 1)
	t.Logf("%#v", ret)
}
