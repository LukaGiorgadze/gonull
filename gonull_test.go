package gonull

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNullable(t *testing.T) {
	value := "test"
	n := NewNullable(value)

	assert.True(t, n.Valid)
	assert.Equal(t, value, n.Val)
}

func TestNullableScan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		Valid   bool
		wantErr bool
	}{
		{
			name:  "nil value",
			value: nil,
			Valid: false,
		},
		{
			name:  "string value",
			value: "test",
			Valid: true,
		},
		{
			name:    "unsupported type",
			value:   []byte{1, 2, 3},
			wantErr: true,
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
		name          string
		jsonData      []byte
		expectedVal   int
		expectedValid bool
	}

	testCases := []testCase{
		{
			name:          "ValuePresent",
			jsonData:      []byte(`123`),
			expectedVal:   123,
			expectedValid: true,
		},
		{
			name:          "ValueNull",
			jsonData:      []byte(`null`),
			expectedVal:   0,
			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var nullable Nullable[int]

			err := nullable.UnmarshalJSON(tc.jsonData)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedVal, nullable.Val)
			assert.Equal(t, tc.expectedValid, nullable.Valid)
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
