package layout

import (
	"context"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

const (
	defaultMarginLeft = 50.0 // px at defaultWidth

	leftRatio  = 0.0625 // 50 / 800
	rightRatio = 0.04   // ~3:2 vs left

	minLeft  = 20.0 // room for "1." and repeat marks
	minRight = 12.0
	maxLeft  = 50.0
	maxRight = 50.0
)

func clamp(v, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, v))
}

// Margins returns left and right margins in px for a given layout width.
// defaultWidth is accepted so the scale stays relative to the reference layout.
func Margins(newWidth, defaultWidth int) (mLeft, mRight int) {
	scale := newWidth / defaultWidth
	base := defaultWidth * scale // == newWidth; kept explicit for clarity

	marginLeft := clamp(float64(base)*leftRatio, minLeft, maxLeft)
	marginRight := clamp(float64(base)*rightRatio, minRight, maxRight)
	return int(math.Round(marginLeft)), int(math.Round(marginRight))
}

func GetLayoutWidth(ctx context.Context) int {
	layoutWidth := constant.LAYOUT_WIDTH

	if p, ok := params.GetParamFromContext(ctx); ok && p.Render != nil {
		// one line
		if p.Render.DeviceType == 0 {
			layoutWidth = params.GetMaxPositionX(ctx)
		} else {
			layoutWidth, _ = params.GetDeviceDimensions(ctx)
		}
	}

	return layoutWidth
}

// GetMargin returns left and right margins in px based on context parameters.
func GetMargin(ctx context.Context) (left, right int) {
	if p, ok := params.GetParamFromContext(ctx); ok && p.Render != nil {
		if p.Render.DeviceType == 0 {
			return constant.LAYOUT_INDENT_LENGTH, constant.LAYOUT_INDENT_LENGTH
		}
		return Margins(GetLayoutWidth(ctx), constant.LAYOUT_WIDTH)
	}

	return constant.LAYOUT_INDENT_LENGTH, constant.LAYOUT_INDENT_LENGTH
}
