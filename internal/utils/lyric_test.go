package utils_test

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestCalculateSecondaryLyricWidth(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		text string
		want float64
	}{
		{
			name: "all of them",
			text: "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwYyZz1.2,3 4(5)67890",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.CalculateSecondaryLyricWidth(tt.text)
			assert.Equal(t, tt.want, got)
		})
	}
}
