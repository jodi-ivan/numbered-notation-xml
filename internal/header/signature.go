package header

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (hi *headerInteractor) RenderKeyandTimeSignatures(ctx context.Context, canv canvas.Canvas, key keysig.KeySignature, timeSignature timesig.TimeSignature) {

	relativeY := constant.TITLE_Y_POS + SIGNATURES_Y_POS
	currKeySig := key.GetKeyOnMeasure(ctx, 1)
	humanizedKeySignature := timeSignature.GetHumanized()

	canv.Text(constant.LAYOUT_INDENT_LENGTH, relativeY, currKeySig.String())
	canv.Text(constant.LAYOUT_INDENT_LENGTH+(3*constant.LOWERCASE_LENGTH)+int(hi.Lyric.CalculateLyricWidth(currKeySig.String())), relativeY, humanizedKeySignature)

}
