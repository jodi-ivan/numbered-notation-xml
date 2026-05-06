package timesig

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderGregorian(ctx context.Context, canv canvas.Canvas, lines [5]int, timeSig TimeSignature, x float64) {

	canv.Group(`class="timesig"`, `style="font-size:2em"`)
	uniHex := map[int]string{
		1:  "&#xF031;",
		2:  "&#xF032;",
		3:  "&#xF033;",
		4:  "&#xF034;",
		5:  "&#xF035;",
		6:  "&#xF036;",
		7:  "&#xF037;",
		8:  "&#xF038;",
		9:  "&#xF039;",
		0:  "&#xF030;",
		12: "&#xF031;&#xF032;",
	}

	ts := timeSig.GetTimesignatureOnMeasure(ctx, 1)
	if timeSig.IsMixed {
		for i, t := range timeSig.UniqueSign {
			canv.Group(`class="time"`)
			canv.TextUnescaped(x+8+float64(16*i), float64(lines[0]+lines[2])/2, uniHex[t.Beat])
			offset := 0.0
			if t.Beat > 9 {
				offset += 7
			}
			canv.TextUnescaped(x+8+float64(16*i)+offset, float64(lines[2]+lines[4])/2, uniHex[t.BeatType])
			canv.Gend()

		}
	} else if ts.Beat == 1 {
		canv.TextUnescaped(x, float64(lines[2]), uniHex[ts.Beat])
	} else {
		offset := 0.0
		if ts.Beat > 9 {
			offset += 3.5
		}
		canv.Group(`class="time"`)
		canv.TextUnescaped(x, float64(lines[0]+lines[2])/2, uniHex[ts.Beat])
		canv.TextUnescaped(x+offset, float64(lines[2]+lines[4])/2, uniHex[ts.BeatType])
		canv.Gend()

	}
	canv.Gend()
}
