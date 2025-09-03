package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidOID(t *testing.T) {
	tests := []struct {
		name string
		oid  string
		want bool
	}{
		{
			name: "valid oid",
			oid:  ".1.3.6.1.2.1.1.1.0",
			want: true,
		},
		{
			name: "invalid oid",
			oid:  "1.3.6.1.2.1.1.1.0.0.",
			want: false,
		},
		{
			name: "empty string",
			oid:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isValidOID(tt.oid))
		})
	}
}
