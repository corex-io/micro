package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// List defiend JSON data type, need to implements driver.Valuer, sql.Scanner interface
type List []string

// MarshalJSON to output non base64 encoded []byte
func (m List) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	v := ([]string)(m)

	return json.Marshal(v)
}

// UnmarshalJSON to deserialize []byte
func (m *List) UnmarshalJSON(b []byte) error {
	var v []string
	err := json.Unmarshal(b, &v)
	*m = List(v)
	return err
	// return json.Unmarshal(b, m)
}

// Value return json value, implement driver.Valuer interface
func (m List) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	b, err := m.MarshalJSON()
	return string(b), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *List) Scan(v interface{}) error {
	switch vt := v.(type) {
	case []byte:
		return m.UnmarshalJSON(vt)
	case string:
		return m.UnmarshalJSON([]byte(vt))
	default:
		return fmt.Errorf("Failed to unmarshal JSONB value: %#v", v)
	}
}

// GormDataType gorm common data type
func (List) GormDataType() string {
	return "List"
}

// GormDBDataType gorm db data type
// func (M) GormDBDataType(db *gorm.DB, field *schema.Field) string {
// 	switch db.Dialector.Name() {
// 	case "sqlite":
// 		return "text"
// 	case "mysql":
// 		return "text"
// 	case "postgres":
// 		return "text"
// 	}
// 	return ""
// }
