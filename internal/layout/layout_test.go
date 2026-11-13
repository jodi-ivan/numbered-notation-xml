package layout

import (
	"context"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
	"github.com/stretchr/testify/assert"
)

func TestMargins(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		newWidth     int
		defaultWidth int
		want         int
		want2        int
	}{
		// Default width (800) - should result in no clamping
		{
			name:         "DefaultWidth_NoClamping",
			newWidth:     800,
			defaultWidth: 800,
			want:         50, // 800 * 0.0625 = 50
			want2:        32, // 800 * 0.04 = 32
		},
		// Small width (Clamping to minimum)
		{
			name:         "SmallWidth_ClampingToMin",
			newWidth:     100,
			defaultWidth: 800,
			want:         20, // 100 * 0.0625 = 6.25 -> clamped to 20
			want2:        12, // 100 * 0.04 = 4.0 -> clamped to 12
		},
		// Width near minimum threshold (320)
		{
			name:         "WidthNearMinThreshold",
			newWidth:     319,
			defaultWidth: 800,
			want:         20, // 319 * 0.0625 = 19.9375 -> clamped to 20
			want2:        12,
		},
		// Large width (Clamping to maximum)
		{
			name:         "LargeWidth_ClampingToMax",
			newWidth:     1000,
			defaultWidth: 800,
			want:         50, // 1000 * 0.0625 = 62.5 -> clamped to 50
			want2:        32,
		},
		// Very large width (Clamping both to maximum)
		{
			name:         "VeryLargeWidth_ClampingBothToMax",
			newWidth:     2000,
			defaultWidth: 800,
			want:         50, // 2000 * 0.0625 = 125.0 -> clamped to 50
			want2:        50, // 2000 * 0.04 = 80.0 -> clamped to 50
		},
		// Width exactly at minimum threshold (320)
		{
			name:         "WidthAtMinThreshold",
			newWidth:     320,
			defaultWidth: 800,
			want:         20, // 320 * 0.0625 = 20.0 -> no clamp (or clamped to 20)
			want2:        12,
		},
		// Width exactly at maximum threshold (800)
		{
			name:         "WidthAtMaxThreshold",
			newWidth:     800,
			defaultWidth: 800,
			want:         50, // 800 * 0.0625 = 50.0 -> no clamp (or clamped to 50)
			want2:        32, // 800 * 0.04 = 32.0 -> no clamp (or clamped to 32)
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := Margins(tt.newWidth, tt.defaultWidth)
			assert.Equal(t, tt.want, got, "got")
			assert.Equal(t, tt.want2, got2, "got2")

		})
	}
}

func TestGetMargin(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		ctx   context.Context
		want  int
		want2 int
	}{
		// No context provided (Fallback)
		{
			name:  "NoContext",
			ctx:   context.Background(),
			want:  50, // Assuming LAYOUT_INDENT_LENGTH is 50
			want2: 50,
		},
		// DeviceType 0 (Uses LAYOUT_INDENT_LENGTH)
		{
			name: "DeviceTypeZero",
			ctx: params.NewParamContext(context.Background(), &params.Param{
				Render: &params.RenderParam{
					DeviceType: 0,
				},
			}),
			want:  50,
			want2: 50,
		},
		{
			name: "DeviceTypeOne_DefaultWidth",
			ctx: params.NewParamContext(context.Background(), &params.Param{
				Render: &params.RenderParam{
					DeviceType: 1,
				},
			}),
			want:  20,
			want2: 12,
		},
		{
			name: "DeviceTypeTwo_CustomWidth",
			ctx: params.NewParamContext(context.Background(), &params.Param{
				Render: &params.RenderParam{
					DeviceType: 2,
				},
			}),
			want:  20,
			want2: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := GetMargin(tt.ctx)
			assert.Equal(t, tt.want, got, "got")
			assert.Equal(t, tt.want2, got2, "got2")
		})
	}
}
