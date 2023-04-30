package gonull

import (
	"database/sql/driver"
	"errors"
)

// Nullable wraps a generic nullable type that can be used with Go's database/sql package.
type Nullable[T any] struct {
	Val   T
	IsSet bool
}

// NewNullable returns a new Nullable with the given value set and Valid set to true.
func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{Val: value, IsSet: true}
}

// Scan implements the sql.Scanner interface.
func (n *Nullable[T]) Scan(value interface{}) error {
	if value == nil {
		n.IsSet = false
		return nil
	}

	var err error
	n.Val, err = convertToType[T](value)
	if err == nil {
		n.IsSet = true
	}
	return err
}

// Value implements the driver.Valuer interface.
func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.IsSet {
		return nil, nil
	}
	return n.Val, nil
}

func convertToType[T any](value interface{}) (T, error) {
	switch v := value.(type) {
	case T:
		return v, nil
	default:
		var zero T
		return zero, errors.New("unsupported type conversion")
	}
}
