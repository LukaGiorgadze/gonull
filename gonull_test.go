package gonull

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNullable(t *testing.T) {
	value := "test"
	n := NewNullable(value)

	assert.True(t, n.IsSet)
	assert.Equal(t, value, n.Val)
}

func TestNullableScan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		isSet   bool
		wantErr bool
	}{
		{
			name:  "nil value",
			value: nil,
			isSet: false,
		},
		{
			name:  "string value",
			value: "test",
			isSet: true,
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
				assert.Equal(t, tt.isSet, n.IsSet)
				if tt.isSet {
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
			nullable:  Nullable[string]{IsSet: false},
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
