package lyric

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

var coloringOpacity = map[int]string{
	0: `style="opacity:0.6"`,
}

func getColoringStyle(ctx context.Context, verse, totalLyric int) string {
	// verse mode? & focus mode?
	// refrein / fine
	prm, _ := params.GetParamFromContext(ctx)
	if totalLyric == 1 && verse == 2 {
		return coloringOpacity[0]
	}

	if prm.Verse < 2 || prm.SingleVerseMode {
		return ""
	}

	return coloringOpacity[verse]

}

func PrintWithHighlight(canv canvas.Canvas, snippet string, spans []params.Span, offset int, x, y float64, secondary ...bool) {
	li := NewLyric()
	rawText, prefixes, words := params.Highlight(snippet, spans, offset, x, y, secondary...)
	for i, prefix := range prefixes {
		canv.Rect(int(x+li.CalculateLyricWidth(prefix))-2, int(y-15.5), int(li.CalculateLyricWidth(words[i])), 24, `fill="#E3EFEC"`)
	}
	canv.TextUnescaped(x, y, rawText)
}
