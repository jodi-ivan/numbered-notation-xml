package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/gregorian"
	"github.com/jodi-ivan/numbered-notation-xml/internal/header"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (si *staffInteractor) Render(ctx context.Context, canv canvas.Canvas, part musicxml.Part, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) int {

	relativeY := constant.TITLE_Y_POS + header.HEADER_OFFSET

	staffes := si.SplitLines(ctx, part)
	key := keySignature.GetKeyOnMeasure(ctx, 1)
	x := gregorian.GetLeftIndentWithTimeSignature(key, timeSignature)
	info := StaffInfo{
		NextLineRenderer: []*entity.NoteRenderer{},
	}
	oldMarginButtom := 0
	for i, st := range staffes {
		info = si.RenderStaff(ctx, canv, x, relativeY, i, i == len(staffes)-1, keySignature, timeSignature, st, info.NextLineRenderer...)
		relativeY = relativeY + 130 + info.MarginBottom
		if info.ForceNewLine {
			relativeY += oldMarginButtom
		}
		if info.Multiline {
			x = info.MarginLeft + barline.BARLINE_AFTER_SPACE
			info.Multiline = false
		} else {
			measure := st[0]
			if i < len(staffes)-1 {
				measure = staffes[i+1][0]
			}
			key = keySignature.GetKeyOnMeasure(ctx, measure.Number)
			x = gregorian.GetLeftIndent(key)
		}
		oldMarginButtom = info.MarginBottom
	}

	for len(info.NextLineRenderer) > 0 {
		key = keySignature.GetKeyOnMeasure(ctx, info.NextLineRenderer[0].MeasureNumber)
		x = gregorian.GetLeftIndent(key)
		idx := len(staffes) - 1
		if idx == 0 {
			idx = 1
		}
		info = si.RenderStaff(ctx, canv, x, relativeY, idx, true, keySignature, timeSignature, nil, info.NextLineRenderer...)
		relativeY += info.MarginBottom + 130
	}

	return relativeY
}
