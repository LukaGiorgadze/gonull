package gonull

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNullableScan_String(t *testing.T) {
	tests := []struct {
		name            string
		value, expected any
		Valid           bool
		Present         bool
		wantErr         bool
	}{
		{
			name:     "nil value",
			value:    nil,
			expected: nil,
			Valid:    false,
			Present:  true,
		},
		{
			name:     "string value",
			value:    "test",
			expected: "test",
			Valid:    true,
			Present:  true,
		},
		{
			name:     "[]byte type",
			value:    []byte{116, 101, 115, 116},
			expected: "test",
			Valid:    true,
			Present:  true,
		},
		{
			name:     "[]uint8 type",
			value:    []byte{116, 101, 115, 116},
			expected: "test",
			Valid:    true,
			Present:  true,
		},
		{
			name:    "unsupported type",
			value:   []int64{1, 2, 3},
			wantErr: true,
			Present: true,
		},
		{
			name:     "empty []uint8 type",
			value:    []byte{},
			expected: "",
			Valid:    true,
			Present:  true,
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
					assert.Equal(t, tt.expected, n.Val)
				}
			}
		})
	}
}

func TestNullableScan_Bool(t *testing.T) {
	tests := []struct {
		name            string
		value, expected any
		Valid           bool
		Present         bool
		wantErr         bool
	}{
		{
			name:     "bool type",
			value:    true,
			expected: true,
			Valid:    true,
			Present:  true,
		},
		{
			name:     "int64 true type",
			value:    int64(1),
			expected: true,
			Valid:    true,
			Present:  true,
		},
		{
			name:     "int64 false type",
			value:    int64(0),
			expected: false,
			Valid:    true,
			Present:  true,
		},
		{
			name:    "unsupported type",
			value:   int64(100),
			wantErr: true,
			Present: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n Nullable[bool]
			err := n.Scan(tt.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Valid, n.Valid)
				assert.Equal(t, tt.Present, n.Present)
				if tt.Valid {
					assert.Equal(t, tt.expected, n.Val)
				}
			}
		})
	}
}

func TestNullableScan_Float(t *testing.T) {
	tests := []struct {
		name            string
		value, expected any
		Valid           bool
		Present         bool
		wantErr         bool
	}{
		{
			name:     "float32 type",
			value:    float32(0.25),
			expected: float32(0.25),
			Valid:    true,
			Present:  true,
		},
		{
			name:     "float64 type",
			value:    float64(0.25),
			expected: float32(0.25),
			Valid:    true,
			Present:  true,
		},
		{
			name:     "[]uint8|[]byte type",
			value:    []byte{48, 46, 50, 53},
			expected: float32(0.25),
			Valid:    true,
			Present:  true,
		},
		{
			name:    "[]uint8|[]byte type empty",
			value:   []byte{},
			wantErr: true,
			Present: true,
		},
		{
			name:    "[]uint8|[]byte type non numbers",
			value:   []byte{1, 2, 3},
			wantErr: true,
			Present: true,
		},
		{
			name:    "unsupported type",
			value:   []int64{48, 46, 50, 53},
			wantErr: true,
			Present: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n Nullable[float32]
			err := n.Scan(tt.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Valid, n.Valid)
				assert.Equal(t, tt.Present, n.Present)
				if tt.Valid {
					assert.Equal(t, tt.expected, n.Val)
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
		expected any
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
		{
			name:     "Nil to time.Time",
			expected: time.Time{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var result any
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
			case time.Time:
				result, err = convertToType[time.Time](nil)
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

type testValuerScannerStruct struct {
	b []byte
}

func (t testValuerScannerStruct) Value() (driver.Value, error) {
	return t.b, nil
}

func (t *testValuerScannerStruct) Scan(src any) error {
	if src == nil {
		return nil
	}
	if str, ok := src.(string); ok && str == "error" {
		return errors.New("intentional error")
	}
	switch v := src.(type) {
	case string:
		t.b = []byte(v)
		return nil
	case []byte:
		t.b = v
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func TestValuerAndScanner(t *testing.T) {
	valueNullable1 := Nullable[testValuerScannerStruct]{
		Val:     testValuerScannerStruct{b: []byte("test output string")},
		Valid:   true,
		Present: true,
	}
	valueNullable2 := Nullable[testValuerScannerStruct]{
		Valid:   false,
		Present: true,
	}

	valueResult1, valueErr1 := valueNullable1.Value()
	assert.NoError(t, valueErr1)
	assert.Equal(t, []byte("test output string"), valueResult1)

	valueResult2, valueErr2 := valueNullable2.Value()
	assert.NoError(t, valueErr2)
	assert.Equal(t, nil, valueResult2)

	scannerData1 := []byte("test input string")

	var scannerNullable1 Nullable[testValuerScannerStruct]
	var scannerNullable2 Nullable[testValuerScannerStruct]

	scannerErr1 := scannerNullable1.Scan(scannerData1)
	assert.NoError(t, scannerErr1)
	assert.Equal(t, Nullable[testValuerScannerStruct]{
		Present: true,
		Valid:   true,
		Val: testValuerScannerStruct{
			b: []byte("test input string"),
		},
	}, scannerNullable1)

	scannerErr2 := scannerNullable2.Scan(nil)
	assert.NoError(t, scannerErr2)
	assert.Equal(t, Nullable[testValuerScannerStruct]{
		Present: true,
		Valid:   false,
		Val: testValuerScannerStruct{
			b: []byte(nil),
		},
	}, scannerNullable2)

	var scannerNullableUnsupported Nullable[testValuerScannerStruct]
	scannerErrUnsupported := scannerNullableUnsupported.Scan(123)
	assert.Error(t, scannerErrUnsupported)
	assert.Contains(t, scannerErrUnsupported.Error(), "unsupported type")
	assert.Equal(t, Nullable[testValuerScannerStruct]{
		Present: true,
		Valid:   false,
		Val:     testValuerScannerStruct{},
	}, scannerNullableUnsupported)
}

func TestNullableOrElse(t *testing.T) {
	value := "hello"
	nonEmpty := NewNullable(value)
	assert.Equal(t, value, nonEmpty.OrElse("world"))

	var empty Nullable[string]
	assert.Equal(t, "world", empty.OrElse("world"))
}

type customValuer struct {
	value any
	err   error
}

type unknowType interface{}

func (cv customValuer) Value() (driver.Value, error) {
	return cv.value, cv.err
}

func TestConvertToDriverValue(t *testing.T) {
	var (
		intVal           int          = 123
		int8Val          int8         = 12
		int16Val         int16        = 1234
		int32Val         int32        = 12345
		int64Val         int64        = 123456
		uintVal          uint         = 123
		uint8Val         uint8        = 12
		uint16Val        uint16       = 1234
		uint32Val        uint32       = 12345
		uint64Val        uint64       = 1 << 62
		float32Val       float32      = 12.34
		float64Val       float64      = 123.456
		boolVal          bool         = true
		stringVal        string       = "test"
		timeVal          time.Time    = time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
		byteSlice        []byte       = []byte("byte slice")
		ptrToInt         *int         = &intVal
		nilPtr           *int         = nil
		valuerSuccess    customValuer = customValuer{value: "valuer value", err: nil}
		valuerError      customValuer = customValuer{err: errors.New("valuer error")}
		unknowTypeError  unknowType   = map[bool]bool{}
		unsupportedSlice              = []int{1, 2, 3}
	)

	tests := []struct {
		name    string
		value   any
		want    driver.Value
		wantErr bool
	}{
		{"Int", intVal, int64(intVal), false},
		{"Int8", int8Val, int64(int8Val), false},
		{"Int16", int16Val, int64(int16Val), false},
		{"Int32", int32Val, int64(int32Val), false},
		{"Int64", int64Val, int64(int64Val), false},
		{"Uint", uintVal, int64(uintVal), false},
		{"Uint8", uint8Val, int64(uint8Val), false},
		{"Uint16", uint16Val, int64(uint16Val), false},
		{"Uint32", uint32Val, int64(uint32Val), false},
		{"Uint64", uint64Val, int64(uint64Val), false},
		{"Float32", float32Val, float64(float32Val), false},
		{"Float64", float64Val, float64(float64Val), false},
		{"Bool", boolVal, boolVal, false},
		{"String", stringVal, stringVal, false},
		{"ByteSlice", byteSlice, byteSlice, false},
		{"Time", timeVal, timeVal, false},
		{"PointerToInt", ptrToInt, int64(*ptrToInt), false},
		{"NilPointer", nilPtr, nil, false},
		{"UnsupportedType", struct{}{}, nil, true},
		{"Uint64HighBitSet", uint64(1 << 63), nil, true}, // Uint64 with high bit set
		{"ValuerInterfaceSuccess", valuerSuccess, "valuer value", false},
		{"ValuerInterfaceError", valuerError, nil, true},
		{"UnknowTypeError", unknowTypeError, nil, true},
		{"UnsupportedSliceType", unsupportedSlice, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToDriverValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToDriverValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToDriverValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableValue_Uint32(t *testing.T) {
	uint32Val := uint32(12345)
	nullableUint32 := NewNullable(uint32Val)

	convertedValue, err := nullableUint32.Value()

	if err != nil {
		t.Fatalf("Nullable[uint32].Value() returned an error: %v", err)
	}

	if _, ok := convertedValue.(int64); !ok {
		t.Fatalf("Nullable[uint32].Value() returned a non-int64 type: %T", convertedValue)
	}

	if int64(uint32Val) != convertedValue.(int64) {
		t.Errorf("Nullable[uint32].Value() returned %v, want %v", convertedValue, uint32Val)
	}
}

func Test_IsZero(t *testing.T) {
	type Foo struct {
		ID     Nullable[int64]  `json:"id,omitempty"`
		Name   Nullable[string] `json:"name,omitempty"`
		IsZero Nullable[bool]   `json:"is_zero,omitzero"`
	}

	foo1 := &Foo{}
	err := json.Unmarshal([]byte("{\"id\":0}"), foo1)
	require.NoError(t, err)
	assert.True(t, foo1.ID.Present)        // the value was passed
	assert.True(t, foo1.ID.Valid)          // the value was passed, the value is valid (0 is valid)
	assert.Equal(t, int64(0), foo1.ID.Val) // the value was passed, the value is valid
	assert.False(t, foo1.ID.IsZero())      // the value is not "zero"
	assert.False(t, foo1.Name.Present)     // name is not present
	assert.True(t, foo1.Name.IsZero())     // name is "zero" and will not be marshaled (Note, it needs go >= 1.24)
	hasZero, _ := json.Marshal(foo1)
	assert.Equal(t, `{"id":0,"name":null}`, string(hasZero))

	foo2 := &Foo{}
	err = json.Unmarshal([]byte("{\"id\":null,\"name\":\"foo\"}"), foo2)
	require.NoError(t, err)
	assert.True(t, foo2.ID.Present)       // the value was passed
	assert.False(t, foo2.ID.Valid)        // the value was passed, but it null, invalid (unset)
	assert.False(t, foo1.ID.IsZero())     // the value is not "zero"
	assert.True(t, foo2.Name.Present)     // the value was passed
	assert.True(t, foo2.Name.Valid)       // the value was passed
	assert.Equal(t, "foo", foo2.Name.Val) // the value was passed, the value is valid
	assert.False(t, foo1.ID.IsZero())     // the value is not "zero"
}

func TestNullableScan_Float64(t *testing.T) {
	tests := []struct {
		name            string
		value, expected any
		Valid           bool
		Present         bool
		wantErr         bool
	}{
		{
			name:     "float64 type",
			value:    float64(0.25),
			expected: float64(0.25),
			Valid:    true,
			Present:  true,
		},
		{
			name:     "float32 type",
			value:    float32(0.25),
			expected: float64(0.25),
			Valid:    true,
			Present:  true,
		},
		{
			name:     "[]uint8|[]byte type",
			value:    []byte("0.25"),
			expected: float64(0.25),
			Valid:    true,
			Present:  true,
		},
		{
			name:    "[]uint8|[]byte type empty",
			value:   []byte{},
			wantErr: true,
			Present: true,
		},
		{
			name:    "[]uint8|[]byte type non numbers",
			value:   []byte("not a number"),
			wantErr: true,
			Present: true,
		},
		{
			name:    "unsupported type",
			value:   []int64{48, 46, 50, 53},
			wantErr: true,
			Present: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n Nullable[float64]
			err := n.Scan(tt.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.Valid, n.Valid)
				assert.Equal(t, tt.Present, n.Present)
				if tt.Valid {
					assert.Equal(t, tt.expected, n.Val)
				}
			}
		})
	}
}
