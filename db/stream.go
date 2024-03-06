package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/corex-io/micro/datatypes"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const timeFormat = "2006-01-02T15:04:05+08:00"

func Stream(ctx context.Context, db *sql.DB, query string, args []any, handle func(index int64, row []any) error) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("stream: db is nil")
	}
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("query rows: %w", err)
	}
	defer rows.Close()

	types, err := rows.ColumnTypes()
	if err != nil {
		return 0, fmt.Errorf("column types: %w", err)
	}
	cnt, v, scans := int64(0), make([]sql.RawBytes, len(types), len(types)), make([]any, len(types), len(types))

	for i := range v {
		scans[i] = &v[i]
	}
	for rows.Next() {
		if err = rows.Scan(scans...); err != nil {
			return 0, fmt.Errorf("rows scan: %w", err)
		}
		cnt += 1
		row := make([]any, len(v), len(v))

		for i := 0; i <= len(v)-1; i++ {
			if row[i], err = format(*(*string)(unsafe.Pointer(&v[i])), types[i].ScanType().Name()); err != nil {
				//return cnt, pErr(types[i].Name(), types[i].ScanType().Name(), v[i], v, err)
				fmt.Println(pErr(types[i].Name(), types[i].ScanType().Name(), v[i], v, err))
			}
		}
		if err = handle(cnt, row); err != nil {
			return cnt, err
		}
	}
	return cnt, rows.Err()
}

func pErr(name, kind string, v sql.RawBytes, row any, err error) error {
	return fmt.Errorf("name=%s, kind=%s, v=%s, row=%#v, err=%w", name, kind, v, row, err)
}

func format(s, fmt string) (any, error) {

	switch fmt {
	case "RawBytes", "string":
		return strings.Clone(s), nil
	case "NullInt64", "uint8", "uint32", "uint64", "int", "int8", "int32", "int64":
		return strconv.Atoi(s)
	case "NullTime":
		if s == "0001-01-01T00:00:00Z" || s == "" {
			return &time.Time{}, nil
		}
		return datatypes.ParseInLocation(timeFormat, s, time.Local)
	case "float64":
		return strconv.ParseFloat(s, 64)
	case "bool":
		return strconv.ParseBool(s)
	default:
		return nil, errors.New("format unknown")
	}
}
