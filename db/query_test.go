package db

import (
	"strings"
	"testing"
)

func TestM_SubSQL(t *testing.T) {
	c := strings.SplitN("key", "@", 2)
	t.Logf("%#v", c)
	query := M{
		"key_string":   "123",
		"key_list@in":  []string{"s1", "s2"},
		"key_like1@=~": []string{"%123ddd%"},
		"key_like2@=~": "9999",
	}
	raw, args, err := query.SubSQL()
	t.Logf("raw=%s, args=%#v, err=%v", raw, args, err)
}

func TestSearch_SQL(t *testing.T) {
	search := SearchBody{
		Namespace:    "table",
		ResultColumn: []string{"*"},
		Condition: []M{
			{"1key_string": "123",
				"key1_list@in":  []string{"s1", "s2"},
				"key1_like1@=~": []string{"%123ddd%"},
				"key1_like2@=~": "9999",
			},
			{
				"key2_string":   "123",
				"key2_list@in":  []string{"s1", "s2"},
				"key2_like1@=~": []string{"%123ddd%"},
				"key2_like2@=~": "9999",
			}},
	}
	sql, args, err := search.SQL()
	t.Logf("sql=%s, args=%#v, err=%v", sql, args, err)

}
