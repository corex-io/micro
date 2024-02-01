package db

import (
	"bytes"
	"fmt"
	"github.com/corex-io/micro/datatypes"
	"strconv"
	"strings"
)

type SearchBody struct {
	Raw          string   `json:"Raw,omitempty"`
	Namespace    string   `json:"Namespace,omitempty"`
	ResultColumn []string `json:"ResultColumn,omitempty"`
	Condition    []M      `json:"Condition,omitempty"`
	Offset       int      `json:"Offset,omitempty"`
	Limit        int      `json:"Limit,omitempty"`
	Comment      string   `json:"Comment,omitempty"`
}

func tov[T int | string | int64 | float64](v []T) []any {
	var s = make([]any, 0, len(v))
	for _, i := range v {
		s = append(s, i)
	}
	return s
}

func Format(v any) []any {
	switch s := v.(type) {
	case string, int, uint, int64, float64:
		return []any{v}
	case []string:
		return tov(s)
	case []int:
		return tov(s)
	case []any:
		return s
	case []float64:
		return tov(s)
	case *datatypes.Time:
		return []any{s.Format("2006-01-02 15:04:05")}
	default:
		return []any{v}
	}
}

type M map[string]any

func (m M) SubSQL() (string, []any, error) {
	var condition []string
	var args []any
	for k, v := range m { // name@in
		if k == "" {
			return "", nil, fmt.Errorf("key is NULL")
		}
		vs := Format(v)
		mask := strings.TrimRight(strings.Repeat("?, ", len(vs)), ", ")
		ko := strings.SplitN(k, "$", 2)
		if len(ko) == 1 {
			ko = append(ko, "IN")
		}
		key, op := ko[0], ko[1]

		switch op {
		case "=", "!=", ">", ">=", "<", "<=":
			if len(vs) != 1 {
				return "", nil, fmt.Errorf("key=%s, value=%v", k, v)
			}
			condition = append(condition, fmt.Sprintf("%s %s %s", key, op, mask))
		case "IN", "In", "in", "NOT IN":
			condition = append(condition, fmt.Sprintf("%s %s ( %s )", key, op, mask))
		case "=~", "LIKE":
			if len(vs) != 1 {
				return "", nil, fmt.Errorf("key=%s, value=%v", k, v)
			}
			condition = append(condition, fmt.Sprintf("%s LIKE %s", key, mask))
		case "<>", "RANGE":
			if len(vs) != 2 {
				return "", nil, fmt.Errorf("key=%s, value=%v", k, v)
			}
			condition = append(condition, fmt.Sprintf("%s >= ? AND %s < ?", key, key))
		default:
			return "", nil, fmt.Errorf("op=%v unknown", op)
		}
		args = append(args, vs...)
	}
	return strings.Join(condition, " AND "), args, nil
}

func (s *SearchBody) SQL() (string, []any, error) {
	if s.Raw != "" {
		return s.Raw + "/*" + s.Comment + "*/", nil, nil
	}
	if len(s.ResultColumn) == 0 || s.Namespace == "" {
		return "", nil, fmt.Errorf("namespace or resultColumn is empty")
	}
	var buf bytes.Buffer
	buf.WriteString("SELECT ")
	buf.WriteString(strings.Join(s.ResultColumn, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(s.Namespace)

	var conditions []string
	var args []any
	for _, cond := range s.Condition {
		if len(cond) == 0 {
			continue
		}
		raw, _args, err := cond.SubSQL()
		if err != nil {
			return "", nil, err
		}
		conditions = append(conditions, "( "+raw+" )")
		args = append(args, _args...)
	}

	if len(conditions) != 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(conditions, " OR "))
	}

	if s.Offset > 0 {
		buf.WriteString(" OFFSET " + strconv.Itoa(s.Offset))
	}

	if s.Limit > 0 {
		buf.WriteString(" LIMIT " + strconv.Itoa(s.Limit))
	}

	buf.WriteString("/*" + s.Comment + "*/")
	return buf.String(), args, nil
}
