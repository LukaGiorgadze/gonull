// Package gonull provides a generic Nullable type for handling nullable values in a convenient way.
// This is useful when working with databases and JSON, where nullable values are common.
package gonull

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

var (
	// ErrUnsupportedConversion is an error that occurs when attempting to convert a value to an unsupported type.
	// This typically happens when Scan is called with a value that cannot be converted to the target type T.
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
)

// Nullable is a generic struct that holds a nullable value of any type T.
// It keeps track of the value (Val), a flag (Valid) indicating whether the value has been set and a flag (Present)
// indicating if the value is in the struct.
// This allows for better handling of nullable and undefined values, ensuring proper value management and serialization.
type Nullable[T any] struct {
	Val     T
	Valid   bool
	Present bool
}

// NewNullable creates a new Nullable with the given value and sets Valid to true.
// This is useful when you want to create a Nullable with an initial value, explicitly marking it as set.
func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{Val: value, Valid: true, Present: true}
}

// Scan implements the sql.Scanner interface for Nullable, allowing it to be used as a nullable field in database operations.
// It is responsible for properly setting the Valid flag and converting the scanned value to the target type T.
// This enables seamless integration with database/sql when working with nullable values.
func (n *Nullable[T]) Scan(value interface{}) error {
	n.Present = true

	if value == nil {
		n.Val = zeroValue[T]()
		n.Valid = false
		return nil
	}

	var err error
	n.Val, err = convertToType[T](value)
	if err == nil {
		n.Valid = true
	}
	return err
}

// Value implements the driver.Valuer interface for Nullable, enabling it to be used as a nullable field in database operations.
// This method ensures that the correct value is returned for serialization, handling unset Nullable values by returning nil.
func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Nullable, allowing it to be used as a nullable field in JSON operations.
// This method ensures proper unmarshalling of JSON data into the Nullable value, correctly setting the Valid flag based on the JSON data.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Present = true

	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Val = value
	n.Valid = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Nullable, enabling it to be used as a nullable field in JSON operations.
// This method ensures proper marshalling of Nullable values into JSON data, representing unset values as null in the serialized output.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Val)
}

// zeroValue is a helper function that returns the zero value for the generic type T.
// It is used to set the zero value for the Val field of the Nullable struct when the value is nil.
func zeroValue[T any]() T {
	var zero T
	return zero
}

// convertToType is a helper function that attempts to convert the given value to type T.
// This function is used by Scan to properly handle value conversion, ensuring that Nullable values are always of the correct type.
// ErrUnsupportedConversion is returned when a conversion cannot be made to the generic type T.
func convertToType[T any](value interface{}) (T, error) {
	var zero T
	if value == nil {
		return zero, nil
	}

	if reflect.TypeOf(value) == reflect.TypeOf(zero) {
		return value.(T), nil
	}

	// Check if the value is a numeric type and if T is also a numeric type.
	valueType := reflect.TypeOf(value)
	targetType := reflect.TypeOf(zero)
	if valueType.Kind() >= reflect.Int && valueType.Kind() <= reflect.Float64 &&
		targetType.Kind() >= reflect.Int && targetType.Kind() <= reflect.Float64 {
		if valueType.ConvertibleTo(targetType) {
			convertedValue := reflect.ValueOf(value).Convert(targetType)
			return convertedValue.Interface().(T), nil
		}
	}

	return zero, ErrUnsupportedConversion
}
