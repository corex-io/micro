package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/corex-io/micro/datatypes"
	"strconv"
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
			if v[i] == nil {
				continue
			}
			switch types[i].ScanType().Name() {
			case "RawBytes":
				row[i] = string(v[i]) //*(*string)(unsafe.Pointer(&v[i]))， 这里要深度copy, 以免数据覆盖
			case "NullInt64", "uint8", "int32", "uint64":
				if row[i], err = strconv.Atoi(string(v[i])); err != nil {
					return cnt, pErr(types[i].Name(), types[i].ScanType().Name(), v[i], v, err)
				}
			case "NullTime":
				if *(*string)(unsafe.Pointer(&v[i])) == "0001-01-01T00:00:00Z" {
					continue
				}
				if row[i], err = datatypes.ParseInLocation(timeFormat, string(v[i]), time.Local); err != nil {
					return cnt, pErr(types[i].Name(), types[i].ScanType().Name(), v[i], v, err)
				}

			default:
				return cnt, pErr(types[i].Name(), types[i].ScanType().Name(), v[i], row, nil)
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
