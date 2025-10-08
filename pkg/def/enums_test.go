package def

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerEnum_Type(t *testing.T) {
	intEnum := IntegerEnum{}
	assert.Equal(t, "integer_enum", intEnum.Kind())
}

func TestIntegerEnum_Validate(t *testing.T) {
	tests := []struct {
		name     string
		enum     IntegerEnum
		expected error
	}{
		{
			name: "valid integer enum",
			enum: IntegerEnum{
				1: "Active",
				2: "Inactive",
				3: "Testing",
			},
			expected: nil,
		},
		{
			name:     "empty integer enum",
			enum:     IntegerEnum{},
			expected: nil,
		},
		{
			name: "integer enum with empty value",
			enum: IntegerEnum{
				1: "Active",
				2: "",
				3: "Testing",
			},
			expected: fmt.Errorf("integer enum value for 2 cannot be empty (null values must be stringified)"),
		},
		{
			name:     "nil value",
			expected: fmt.Errorf("integer enum cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.enum.Validate()
			if tt.expected == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expected.Error(), err.Error())
			}
		})
	}
}

func TestBitMapEnum_Type(t *testing.T) {
	bitMapEnum := BitMapEnum{}
	assert.Equal(t, "bit_map_enum", bitMapEnum.Kind())
}

func TestBitMapEnum_Validate(t *testing.T) {
	tests := []struct {
		name     string
		enum     BitMapEnum
		expected error
	}{
		{
			name: "valid bitmap enum",
			enum: BitMapEnum{
				0: "Read",
				1: "Write",
				2: "Execute",
			},
			expected: nil,
		},
		{
			name:     "empty bitmap enum",
			enum:     BitMapEnum{},
			expected: nil,
		},
		{
			name: "bitmap enum with empty value",
			enum: BitMapEnum{
				0: "Read",
				1: "",
				2: "Execute",
			},
			expected: fmt.Errorf("integer enum value for 1 cannot be empty (null values must be stringified)"),
		},
		{
			name:     "nil value",
			expected: fmt.Errorf("bit map enum cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.enum.Validate()
			if tt.expected == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expected.Error(), err.Error())
			}
		})
	}
}

func TestOidEnum_Type(t *testing.T) {
	oidEnum := OidEnum("")
	assert.Equal(t, "oid_enum", oidEnum.Kind())
}

func TestOidEnum_Validate(t *testing.T) {
	tests := []struct {
		name     string
		enum     OidEnum
		expected error
	}{
		{
			name:     "valid oid enum",
			enum:     OidEnum(".1.3.6.1.2.1"),
			expected: nil,
		},
		{
			name:     "empty oid enum",
			enum:     OidEnum(""),
			expected: fmt.Errorf("oid enum value cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.enum.Validate()
			if tt.expected == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expected.Error(), err.Error())
			}
		})
	}
}
