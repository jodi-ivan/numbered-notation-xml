package header

import (
	"context"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (hi *headerInteractor) RenderKeyandTimeSignatures(ctx context.Context, canv canvas.Canvas, key keysig.KeySignature, timeSignature timesig.TimeSignature) {

	relativeY := constant.TITLE_Y_POS + SIGNATURES_Y_POS
	currKeySig := key.GetKeyOnMeasure(ctx, 1)
	humanized := currKeySig.String()
	if len(key.Signatures) > 2 {
		for i, v := range key.Signatures {
			if i == 0 {
				continue
			}

			humanized += " - " + strings.Split(v.String(), "=")[1]
		}
	}
	canv.Text(constant.LAYOUT_INDENT_LENGTH, relativeY, humanized)

	if !timeSignature.IsEmpty() {
		humanizedTimeSignature := timeSignature.GetHumanized()
		canv.Text(constant.LAYOUT_INDENT_LENGTH+(3*constant.LOWERCASE_LENGTH)+int(hi.Lyric.CalculateLyricWidth(humanized)), relativeY, humanizedTimeSignature)
	}

}
