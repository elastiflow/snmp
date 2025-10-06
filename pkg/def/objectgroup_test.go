package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectGroup_Validate(t *testing.T) {
	tests := []struct {
		name        string
		objectGroup *ObjectGroup
		expected    error
	}{
		{
			name: "empty object group",
			objectGroup: &ObjectGroup{
				Objects: []string{},
			},
			expected: nil,
		},
		{
			name: "object group with objects",
			objectGroup: &ObjectGroup{
				Objects: []string{"IF-MIB::ifEntry", "IF-MIB::system"},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.objectGroup.Validate())
		})
	}
}

func TestObjectGroup_Type(t *testing.T) {
	objectGroup := ObjectGroup{}
	assert.Equal(t, "object_group", objectGroup.Kind())
}
