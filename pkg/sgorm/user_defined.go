package sgorm

import (
	"database/sql/driver"
	"fmt"
)

// Bool is a custom Boolean type,compatible with MySQL bit(1), tinyint(1) and PostgreSQL bool types
type Bool bool

// NewBool create a new *Bool from bool or *bool
func NewBool(b interface{}) *Bool {
	if b == nil {
		return nil
	}
	switch v := b.(type) {
	case bool:
		val := Bool(v)
		return &val
	case *bool:
		if v == nil {
			return nil
		}
		val := Bool(*v)
		return &val
	default:
		return nil
	}
}

// Value implementing the driver.Valuer interface
func (b Bool) Value() (driver.Value, error) {
	switch currentDriver {
	case "postgres", "postgresql":
		return bool(b), nil
	default: // default MySQL processing
		if b {
			return []byte{1}, nil
		}
		return []byte{0}, nil
	}
}

// Scan implementing the Scanner interface
func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		*b = false
		return nil
	}

	switch v := value.(type) {
	case bool:
		*b = Bool(v)
	case []byte:
		*b = len(v) == 1 && v[0] == 1
	case int64:
		*b = v != 0
	case string:
		*b = v == "t" || v == "true" || v == "1"
	default:
		return fmt.Errorf("unsupported type: %T", value)
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

var currentDriver string

// SetDriver sets the name of the current database driver, such as "mysql" or "postgres"
func SetDriver(driver string) {
	currentDriver = driver
}
