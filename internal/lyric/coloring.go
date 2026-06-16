package lyric

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

var coloringOpacity = map[int]string{
	0: `style="opacity:0.6"`,
}

func getColoringStyle(ctx context.Context, verse, totalLyric int) string {
	prm, _ := params.GetParamFromContext(ctx)
	if totalLyric == 1 && verse == 2 {
		return coloringOpacity[0]
	}

	if prm.Verse < 2 || prm.SingleVerseMode {
		return ""
	}

	return coloringOpacity[verse]

}
