package timesig_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func TestRenderGregorian(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv    func(c *gomock.Controller) *canvas.MockCanvas
		lines   [5]int
		timeSig timesig.TimeSignature
		x       float64
	}{
		// Mixed Time Signature
		{
			name: "Mixed Time Signature",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				mockCanv := canvas.NewMockCanvas(c)
				// Expect calls for mixed signature logic
				mockCanv.EXPECT().Group("class=\"timesig\"", "style=\"font-size:32px\"")
				// Expect calls for unique signatures (e.g., 3/4 and 5/8)
				mockCanv.EXPECT().Group("class=\"time\"").Times(2)
				mockCanv.EXPECT().TextUnescaped(100.0, 20.0, "&#xF033;")
				mockCanv.EXPECT().TextUnescaped(100.0, 40.0, "&#xF034;")
				mockCanv.EXPECT().TextUnescaped(116.0, 20.0, "&#xF035;")
				mockCanv.EXPECT().TextUnescaped(116.0, 40.0, "&#xF038;")

				mockCanv.EXPECT().Gend().Times(2)
				mockCanv.EXPECT().Gend()
				return mockCanv
			},
			lines: [5]int{10, 20, 30, 40, 50},
			timeSig: timesig.TimeSignature{
				IsMixed: true,
				UniqueSign: []timesig.Time{
					{Beat: 3, BeatType: 4},
					{Beat: 5, BeatType: 8},
				},
			},
			x: 100.0,
		},
		// Simple Time Signature (Beat 1)
		{
			name: "Simple Time Signature Beat 1",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				mockCanv := canvas.NewMockCanvas(c)
				mockCanv.EXPECT().Group("class=\"timesig\"", "style=\"font-size:32px\"")
				// Expect TextUnescaped call for beat number at lines[2]
				mockCanv.EXPECT().TextUnescaped(100.0, 30.0, gomock.Any(), gomock.Any())
				mockCanv.EXPECT().Gend().Return()
				return mockCanv
			},
			lines: [5]int{10, 20, 30, 40, 50},
			timeSig: timesig.TimeSignature{
				IsMixed:    false,
				Signatures: []timesig.Time{{Beat: 1, BeatType: 4}},
			},
			x: 100.0,
		},
		// // Simple Time Signature (Beat > 1, Beat <= 9)
		{
			name: "Simple Time Signature Beat > 1",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				mockCanv := canvas.NewMockCanvas(c)
				mockCanv.EXPECT().Group("class=\"timesig\"", "style=\"font-size:32px\"")
				// Expect Group call for simple time signature
				mockCanv.EXPECT().Group("class=\"time\"")
				// Expect TextUnescaped calls for beat number and beat type
				mockCanv.EXPECT().TextUnescaped(100.0, 20.0, "&#xF034;")
				mockCanv.EXPECT().TextUnescaped(100.0, 40.0, "&#xF034;")
				mockCanv.EXPECT().Gend()
				mockCanv.EXPECT().Gend()

				return mockCanv
			},
			lines: [5]int{10, 20, 30, 40, 50},
			timeSig: timesig.TimeSignature{
				IsMixed:    false,
				Signatures: []timesig.Time{{Beat: 4, BeatType: 4}},
			},
			x: 100.0,
		},
		// // Simple Time Signature (Beat > 9)
		{
			name: "Simple Time Signature Beat > 9",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				mockCanv := canvas.NewMockCanvas(c)
				mockCanv.EXPECT().Group("class=\"timesig\"", "style=\"font-size:32px\"")
				// Expect Group call for simple time signature
				mockCanv.EXPECT().Group("class=\"time\"")
				// Expect TextUnescaped calls for beat number and beat type with offset
				mockCanv.EXPECT().TextUnescaped(100.0, 20.0, "&#xF031;&#xF032;")
				mockCanv.EXPECT().TextUnescaped(103.5, 40.0, "&#xF034;")
				mockCanv.EXPECT().Gend()
				mockCanv.EXPECT().Gend()

				return mockCanv
			},
			lines: [5]int{10, 20, 30, 40, 50},
			timeSig: timesig.TimeSignature{
				IsMixed:    false,
				Signatures: []timesig.Time{{Beat: 12, BeatType: 4}},
			},
			x: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timesig.RenderGregorian(context.Background(), tt.canv(ctrl), tt.lines, tt.timeSig, tt.x)
		})
	}
}
