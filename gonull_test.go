package gonull_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/LukaGiorgadze/gonull"
	"github.com/stretchr/testify/assert"
)

func TestNewNullable(t *testing.T) {
	value := "test"
	n := gonull.NewNullable(value)

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
			var n gonull.Nullable[string]
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
		nullable  gonull.Nullable[string]
		wantValue driver.Value
		wantErr   error
	}{
		{
			name:      "valid value",
			nullable:  gonull.NewNullable("test"),
			wantValue: "test",
			wantErr:   nil,
		},
		{
			name:      "unset value",
			nullable:  gonull.Nullable[string]{Valid: false},
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
			var nullable gonull.Nullable[int]

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

	var nullable gonull.Nullable[int]
	err := nullable.UnmarshalJSON(jsonData)

	assert.Error(t, err)
	assert.False(t, nullable.Valid)
}

func TestNullableMarshalJSON(t *testing.T) {
	type testCase struct {
		name         string
		nullable     gonull.Nullable[int]
		expectedJSON []byte
	}

	testCases := []testCase{
		{
			name:         "ValuePresent",
			nullable:     gonull.NewNullable[int](123),
			expectedJSON: []byte(`123`),
		},
		{
			name:         "ValueNull",
			nullable:     gonull.Nullable[int]{Val: 0, Valid: false},
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

	var n gonull.Nullable[string]
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
		{name: "Convert int64 to string (expected to fail)", targetType: "string", value: int64(9), expectedError: gonull.ErrUnsupportedConversion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch tt.targetType {
			case "int":
				n := gonull.Nullable[int]{}
				err = n.Scan(tt.value)
			case "int8":
				n := gonull.Nullable[int8]{}
				err = n.Scan(tt.value)
			case "int16":
				n := gonull.Nullable[int16]{}
				err = n.Scan(tt.value)
			case "int32":
				n := gonull.Nullable[int32]{}
				err = n.Scan(tt.value)
			case "uint":
				n := gonull.Nullable[uint]{}
				err = n.Scan(tt.value)
			case "uint8":
				n := gonull.Nullable[uint8]{}
				err = n.Scan(tt.value)
			case "uint16":
				n := gonull.Nullable[uint16]{}
				err = n.Scan(tt.value)
			case "uint32":
				n := gonull.Nullable[uint32]{}
				err = n.Scan(tt.value)
			case "string":
				n := gonull.Nullable[string]{}
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
		Field gonull.Nullable[TestEnum]
	}

	// Simulate the scenario where the SQL driver returns an int64
	// This is common as database integer types are usually scanned as int64 in Go
	//
	// sqlReturnedValue (int64(0)) is convertible to float32.
	// The converted value 0 (as float32) matches TestEnumA, which is also 0 when converted to float32.
	sqlReturnedValue := int64(0)

	model := TestModel{ID: 1, Field: gonull.NewNullable(TestEnumA)}

	err := model.Field.Scan(sqlReturnedValue)
	assert.NoError(t, err, "Scan failed with unsupported type conversion")
	assert.Equal(t, TestEnumA, model.Field.Val, "Scanned value does not match expected enum value")

}

type testStruct struct {
	Foo gonull.Nullable[*string] `json:"foo"`
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
