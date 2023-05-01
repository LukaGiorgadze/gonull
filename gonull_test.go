package gonull

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNullable(t *testing.T) {
	value := "test"
	n := NewNullable(value)

	assert.True(t, n.IsValid)
	assert.Equal(t, value, n.Val)
}

func TestNullableScan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		IsValid bool
		wantErr bool
	}{
		{
			name:    "nil value",
			value:   nil,
			IsValid: false,
		},
		{
			name:    "string value",
			value:   "test",
			IsValid: true,
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
				assert.Equal(t, tt.IsValid, n.IsValid)
				if tt.IsValid {
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
			nullable:  Nullable[string]{IsValid: false},
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
		expectedVal     int
		expectedIsValid bool
	}

	testCases := []testCase{
		{
			name:            "ValuePresent",
			jsonData:        []byte(`123`),
			expectedVal:     123,
			expectedIsValid: true,
		},
		{
			name:            "ValueNull",
			jsonData:        []byte(`null`),
			expectedVal:     0,
			expectedIsValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var nullable Nullable[int]

			err := nullable.UnmarshalJSON(tc.jsonData)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedVal, nullable.Val)
			assert.Equal(t, tc.expectedIsValid, nullable.IsValid)
		})
	}
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
			nullable:     Nullable[int]{Val: 0, IsValid: false},
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
