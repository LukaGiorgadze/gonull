package gonull

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNullable(t *testing.T) {
	value := "test"
	n := NewNullable(value)

	assert.True(t, n.Valid)
	assert.Equal(t, value, n.Val)
}

type NullableInt struct {
	Int  int
	Null bool
}

func TestNullableScan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		Valid   bool
		Present bool
		wantErr bool
	}{
		{
			name:    "nil value",
			value:   nil,
			Valid:   false,
			Present: true,
		},
		{
			name:    "string value",
			value:   "test",
			Valid:   true,
			Present: true,
		},
		{
			name:    "unsupported type",
			value:   []byte{1, 2, 3},
			wantErr: true,
			Present: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n Nullable[string]
			err := n.Scan(tt.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Valid, n.Valid)
				assert.Equal(t, tt.Present, n.Present)
				if tt.Valid {
					assert.Equal(t, tt.value, n.Val)
				}
			}
		})
	}
}

func TestNullableValue(t *testing.T) {
	tests := []struct {
		name      string
		nullable  Nullable[string]
		wantValue driver.Value
		wantErr   error
	}{
		{
			name:      "valid value",
			nullable:  NewNullable("test"),
			wantValue: "test",
			wantErr:   nil,
		},
		{
			name:      "unset value",
			nullable:  Nullable[string]{Valid: false},
			wantValue: nil,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.nullable.Value()

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantValue, value)
		})
	}
}

func TestNullableUnmarshalJSON(t *testing.T) {
	type testCase struct {
		name            string
		jsonData        []byte
		expectedVal     any
		expectedValid   bool
		expectedPresent bool
	}

	testCases := []testCase{
		{
			name:            "ValuePresent",
			jsonData:        []byte(`123`),
			expectedVal:     123,
			expectedValid:   true,
			expectedPresent: true,
		},
		{
			name:            "ValueNull",
			jsonData:        []byte(`null`),
			expectedVal:     0,
			expectedValid:   false,
			expectedPresent: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var nullable Nullable[int]

			err := nullable.UnmarshalJSON(tc.jsonData)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedVal, nullable.Val)
			assert.Equal(t, tc.expectedValid, nullable.Valid)
			assert.Equal(t, tc.expectedPresent, nullable.Present)
		})
	}
}

func TestNullableUnmarshalJSON_Error(t *testing.T) {
	jsonData := []byte(`"invalid_number"`)

	var nullable Nullable[int]
	err := nullable.UnmarshalJSON(jsonData)

	assert.Error(t, err)
	assert.False(t, nullable.Valid)
}

func TestNullableMarshalJSON(t *testing.T) {
	type testCase struct {
		name         string
		nullable     Nullable[int]
		expectedJSON []byte
	}

	testCases := []testCase{
		{
			name:         "ValuePresent",
			nullable:     NewNullable[int](123),
			expectedJSON: []byte(`123`),
		},
		{
			name:         "ValueNull",
			nullable:     Nullable[int]{Val: 0, Valid: false},
			expectedJSON: []byte(`null`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, err := tc.nullable.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedJSON, jsonData)
		})
	}
}

func TestNullableScan_UnconvertibleFromInt64(t *testing.T) {
	value := int64(123456789012345)

	var n Nullable[string]
	err := n.Scan(value)

	assert.Error(t, err)
	assert.False(t, n.Valid)
}

func TestConvertToTypeFromInt64(t *testing.T) {
	tests := []struct {
		name          string
		targetType    string
		value         int64
		expectedError error
	}{
		{name: "Convert int64 to int", targetType: "int", value: int64(1), expectedError: nil},
		{name: "Convert int64 to int8", targetType: "int8", value: int64(2), expectedError: nil},
		{name: "Convert int64 to int16", targetType: "int16", value: int64(3), expectedError: nil},
		{name: "Convert int64 to int32", targetType: "int32", value: int64(4), expectedError: nil},
		{name: "Convert int64 to uint", targetType: "uint", value: int64(5), expectedError: nil},
		{name: "Convert int64 to uint8", targetType: "uint8", value: int64(6), expectedError: nil},
		{name: "Convert int64 to uint16", targetType: "uint16", value: int64(7), expectedError: nil},
		{name: "Convert int64 to uint32", targetType: "uint32", value: int64(8), expectedError: nil},
		// Add more tests as necessary
		{name: "Convert int64 to string (expected to fail)", targetType: "string", value: int64(9), expectedError: ErrUnsupportedConversion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch tt.targetType {
			case "int":
				n := Nullable[int]{}
				err = n.Scan(tt.value)
			case "int8":
				n := Nullable[int8]{}
				err = n.Scan(tt.value)
			case "int16":
				n := Nullable[int16]{}
				err = n.Scan(tt.value)
			case "int32":
				n := Nullable[int32]{}
				err = n.Scan(tt.value)
			case "uint":
				n := Nullable[uint]{}
				err = n.Scan(tt.value)
			case "uint8":
				n := Nullable[uint8]{}
				err = n.Scan(tt.value)
			case "uint16":
				n := Nullable[uint16]{}
				err = n.Scan(tt.value)
			case "uint32":
				n := Nullable[uint32]{}
				err = n.Scan(tt.value)
			case "string":
				n := Nullable[string]{}
				err = n.Scan(tt.value)
			default:
				t.Fatalf("Unsupported type: %s", tt.targetType)
				return
			}

			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}

func TestNullableScanWithCustomEnum(t *testing.T) {
	type TestEnum float32

	const (
		TestEnumA TestEnum = iota
		TestEnumB
	)

	type TestModel struct {
		ID    int
		Field Nullable[TestEnum]
	}

	// Simulate the scenario where the SQL driver returns an int64
	// This is common as database integer types are usually scanned as int64 in Go
	//
	// sqlReturnedValue (int64(0)) is convertible to float32.
	// The converted value 0 (as float32) matches TestEnumA, which is also 0 when converted to float32.
	sqlReturnedValue := int64(0)

	model := TestModel{ID: 1, Field: NewNullable(TestEnumA)}

	err := model.Field.Scan(sqlReturnedValue)
	assert.NoError(t, err, "Scan failed with unsupported type conversion")
	assert.Equal(t, TestEnumA, model.Field.Val, "Scanned value does not match expected enum value")

}

func TestConvertToTypeWithNilValue(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
	}{
		{
			name:     "Nil to int",
			expected: int(0),
		},
		{
			name:     "Nil to int8",
			expected: int8(0),
		},
		{
			name:     "Nil to int16",
			expected: int16(0),
		},
		{
			name:     "Nil to int32",
			expected: int32(0),
		},
		{
			name:     "Nil to int64",
			expected: int64(0),
		},
		{
			name:     "Nil to uint",
			expected: uint(0),
		},
		{
			name:     "Nil to uint8 (byte)",
			expected: uint8(0),
		},
		{
			name:     "Nil to uint16",
			expected: uint16(0),
		},
		{
			name:     "Nil to uint32",
			expected: uint32(0),
		},
		{
			name:     "Nil to uint64",
			expected: uint64(0),
		},
		{
			name:     "Nil to float32",
			expected: float32(0),
		},
		{
			name:     "Nil to float64",
			expected: float64(0),
		},
		{
			name:     "Nil to bool",
			expected: bool(false),
		},
		{
			name:     "Nil to string",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var result interface{}
			var err error

			switch tc.expected.(type) {
			case int:
				result, err = convertToType[int](nil)
			case int8:
				result, err = convertToType[int8](nil)
			case int16:
				result, err = convertToType[int16](nil)
			case int32:
				result, err = convertToType[int32](nil)
			case int64:
				result, err = convertToType[int64](nil)
			case uint:
				result, err = convertToType[uint](nil)
			case uint8:
				result, err = convertToType[uint8](nil)
			case uint16:
				result, err = convertToType[uint16](nil)
			case uint32:
				result, err = convertToType[uint32](nil)
			case uint64:
				result, err = convertToType[uint64](nil)
			case float32:
				result, err = convertToType[float32](nil)
			case float64:
				result, err = convertToType[float64](nil)
			case bool:
				result, err = convertToType[bool](nil)
			case string:
				result, err = convertToType[string](nil)
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

type testStruct struct {
	Foo Nullable[*string] `json:"foo"`
}

func TestPresent(t *testing.T) {
	var nullable1 testStruct
	var nullable2 testStruct
	var nullable3 testStruct

	err := json.Unmarshal([]byte(`{"foo":"f"}`), &nullable1)
	assert.NoError(t, err)
	assert.Equal(t, true, nullable1.Foo.Valid)
	assert.Equal(t, true, nullable1.Foo.Present)

	err = json.Unmarshal([]byte(`{}`), &nullable2)
	assert.NoError(t, err)
	assert.Equal(t, false, nullable2.Foo.Valid)
	assert.Equal(t, false, nullable3.Foo.Present)
	assert.Nil(t, nullable2.Foo.Val)

	err = json.Unmarshal([]byte(`{"foo": null}`), &nullable3)
	assert.NoError(t, err)
	assert.Equal(t, false, nullable3.Foo.Valid)
	assert.Equal(t, true, nullable3.Foo.Present)
	assert.Nil(t, nullable3.Foo.Val)
}
