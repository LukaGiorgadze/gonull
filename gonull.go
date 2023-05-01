// Package gonull provides a generic Nullable type for handling nullable values in a convenient way.
// This is useful when working with databases and JSON, where nullable values are common.
package gonull

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var (
	// ErrUnsupportedConversion is an error that occurs when attempting to convert a value to an unsupported type.
	// This typically happens when Scan is called with a value that cannot be converted to the target type T.
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
)

// Nullable is a generic struct that holds a nullable value of any type T.
// It keeps track of the value (Val) and a flag (IsSet) indicating whether the value has been set.
// This allows for better handling of nullable values, ensuring proper value management and serialization.
type Nullable[T any] struct {
	Val     T
	IsValid bool
}

// NewNullable creates a new Nullable with the given value and sets IsSet to true.
// This is useful when you want to create a Nullable with an initial value, explicitly marking it as set.
func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{Val: value, IsValid: true}
}

// Scan implements the sql.Scanner interface for Nullable, allowing it to be used as a nullable field in database operations.
// It is responsible for properly setting the IsSet flag and converting the scanned value to the target type T.
// This enables seamless integration with database/sql when working with nullable values.
func (n *Nullable[T]) Scan(value interface{}) error {
	if value == nil {
		n.IsValid = false
		return nil
	}

	var err error
	n.Val, err = convertToType[T](value)
	if err == nil {
		n.IsValid = true
	}
	return err
}

// Value implements the driver.Valuer interface for Nullable, enabling it to be used as a nullable field in database operations.
// This method ensures that the correct value is returned for serialization, handling unset Nullable values by returning nil.
func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.IsValid {
		return nil, nil
	}
	return n.Val, nil
}

// convertToType is a helper function that attempts to convert the given value to type T.
// This function is used by Scan to properly handle value conversion, ensuring that Nullable values are always of the correct type.
func convertToType[T any](value interface{}) (T, error) {
	switch v := value.(type) {
	case T:
		return v, nil
	default:
		var zero T
		return zero, ErrUnsupportedConversion
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for Nullable, allowing it to be used as a nullable field in JSON operations.
// This method ensures proper unmarshalling of JSON data into the Nullable value, correctly setting the IsSet flag based on the JSON data.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.IsValid = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Val = value
	n.IsValid = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Nullable, enabling it to be used as a nullable field in JSON operations.
// This method ensures proper marshalling of Nullable values into JSON data, representing unset values as null in the serialized output.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.IsValid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Val)
}
