package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/header"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (si *staffInteractor) Render(ctx context.Context, canv canvas.Canvas, part musicxml.Part, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) int {

	relativeY := constant.TITLE_Y_POS + header.HEADER_OFFSET

	staffes := si.SplitLines(ctx, part)
	x := constant.LAYOUT_INDENT_LENGTH
	info := StaffInfo{
		NextLineRenderer: []*entity.NoteRenderer{},
	}
	oldMarginButtom := 0
	for i, st := range staffes {
		info = si.RenderStaff(ctx, canv, x, relativeY, i == len(staffes)-1, keySignature, timeSignature, st, info.NextLineRenderer...)
		relativeY = relativeY + 70 + info.MarginBottom
		if info.ForceNewLine {
			relativeY += oldMarginButtom
		}
		if info.Multiline {
			x = info.MarginLeft + barline.BARLINE_AFTER_SPACE
			info.Multiline = false
		} else {
			x = constant.LAYOUT_INDENT_LENGTH
		}
		oldMarginButtom = info.MarginBottom
	}

	for len(info.NextLineRenderer) > 0 {
		x = constant.LAYOUT_INDENT_LENGTH
		info = si.RenderStaff(ctx, canv, x, relativeY, true, keySignature, timeSignature, nil, info.NextLineRenderer...)
		relativeY += info.MarginBottom + 70
	}

	return relativeY
}
