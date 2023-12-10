package datatypes

import "testing"

func TestLine_RangeMatch(t *testing.T) {
	v := NewLine([]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}...)
	t.Logf("%#v", v)
	t.Logf("%s", v.String())
	ret := v.RangeMatch(func(v *float64) bool {
		return v != nil && *v == 0
	}, 1)
	t.Logf("%#v", ret)
}

func TestLine_Cut(t *testing.T) {
	v := NewLine([]float64{-1, -1, 2, 3, 4, 5, -1, -1}...)
	t.Logf("%#v", v)
	t.Logf("%s", v.String())
	ret := v.Cut(func(v *float64) bool {
		return v != nil && *v != -1
	})
	t.Logf("%#v", ret.String())
}
