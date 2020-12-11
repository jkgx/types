package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// NullTime implements sql.NullTime functionality.
type NullTime time.Time

// Scan implements the Scanner interface.
func (ns *NullTime) Scan(value interface{}) error {
	var v sql.NullTime
	if err := (&v).Scan(value); err != nil {
		return err
	}
	*ns = NullTime(v.Time)
	return nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (ns NullTime) MarshalJSON() ([]byte, error) {
	var t *time.Time
	if !time.Time(ns).IsZero() {
		tt := time.Time(ns)
		t = &tt
	}
	return json.Marshal(t)
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *NullTime) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*ns = NullTime(t)
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullTime) Value() (driver.Value, error) {
	return sql.NullTime{Valid: !time.Time(ns).IsZero(), Time: time.Time(ns)}.Value()
}

// JSONRawMessage represents a json.RawMessage that works well with JSON, SQL, and Swagger.
type JSONRawMessage json.RawMessage

// Scan implements the Scanner interface.
func (m *JSONRawMessage) Scan(value interface{}) error {
	*m = []byte(fmt.Sprintf("%s", value))
	return nil
}

// Value implements the driver Valuer interface.
func (m JSONRawMessage) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "null", nil
	}
	return string(m), nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (m JSONRawMessage) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *JSONRawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

// NullJSONRawMessage represents a json.RawMessage that works well with JSON, SQL, and Swagger and is NULLable-
type NullJSONRawMessage json.RawMessage

// Scan implements the Scanner interface.
func (m *NullJSONRawMessage) Scan(value interface{}) error {
	if value == nil {
		value = "null"
	}
	*m = []byte(fmt.Sprintf("%s", value))
	return nil
}

// Value implements the driver Valuer interface.
func (m NullJSONRawMessage) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return string(m), nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (m NullJSONRawMessage) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *NullJSONRawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

// JSONScan is a generic helper for storing a value as a JSON blob in SQL.
func JSONScan(dst interface{}, value interface{}) error {
	if value == nil {
		value = "null"
	}
	if err := json.Unmarshal([]byte(fmt.Sprintf("%s", value)), &dst); err != nil {
		return fmt.Errorf("unable to decode payload to: %s", err)
	}
	return nil
}

// JSONValue is a generic helper for retrieving a SQL JSON-encoded value.
func JSONValue(src interface{}) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&src); err != nil {
		return nil, err
	}
	return b.String(), nil
}