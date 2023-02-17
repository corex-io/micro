package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// M defiend JSON data type, need to implements driver.Valuer, sql.Scanner interface
type M map[string]interface{}

// MarshalJSON to output non base64 encoded []byte
func (m M) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	v := (map[string]interface{})(m)
	return json.Marshal(v)
}

// UnmarshalJSON to deserialize []byte
func (m *M) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		// *m = M(map[string]interface{}{})
		*m = M(nil)
		return nil
	}
	v := map[string]interface{}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		err = fmt.Errorf("err=%w, string(b)=%s", err, string(b))
	}
	*m = M(v)
	return err
	// return json.Unmarshal(b, m)
}

// Value return json value, implement driver.Valuer interface
func (m M) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	b, err := m.MarshalJSON()
	return string(b), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *M) Scan(v interface{}) error {
	switch vt := v.(type) {
	case []byte:
		return m.UnmarshalJSON(vt)
	case string:
		return m.UnmarshalJSON([]byte(vt))
	default:
		return fmt.Errorf("Failed to unmarshal JSONB value: %#v", v)
	}
}

// Map map
func (m M) Map() map[string]interface{} {
	return m
}

// GormDataType gorm common data type
func (M) GormDataType() string {
	return "Map"
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
