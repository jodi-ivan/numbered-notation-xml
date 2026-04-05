package utils_test

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseHymnWithVariant(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		raw     string
		want    int
		want2   string
		wantErr bool
	}{
		{
			name:    "error prsing -- empty string",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "error prsing -- invalid string",
			raw:     "n",
			wantErr: true,
		},
		{
			name:    "valid without variant",
			raw:     "1",
			wantErr: false,
			want:    1,
		},
		{
			name:    "valid with variant",
			raw:     "1a",
			wantErr: false,
			want:    1,
			want2:   "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := utils.ParseHymnWithVariant(tt.raw)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseHymnWithVariant() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseHymnWithVariant() succeeded unexpectedly")
			}
			assert.Equal(t, tt.want, got, "ParseHymnWithVariant()--> int part")
			assert.Equal(t, tt.want2, got2, "ParseHymnWithVariant()--> string part")
		})
	}
}
