package sgorm

import (
	"database/sql/driver"
	"fmt"
)

// Bool is a custom type for bit fields in MySQL.
type Bool bool

// NewBool creates a new *Bool from a bool or *bool.
func NewBool(b interface{}) *Bool {
	if b == nil {
		return nil
	}
	var val Bool
	if v, ok := b.(bool); ok {
		val = Bool(v)
		return &val
	}
	if v, ok := b.(*bool); ok {
		if v == nil {
			return nil
		}
		val = Bool(*v)
		return &val
	}
	return nil
}

// Value implements the driver Valuer interface.
func (b Bool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// Scan implements the Scanner interface.
func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		*b = false
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type: %T", value)
	}
	if len(v) == 1 && v[0] == 1 {
		*b = true
	} else {
		*b = false
	}

	return nil
}

func (b *Bool) String() string {
	if b == nil {
		return "false"
	}
	if *b {
		return "true"
	}
	return "false"
}
