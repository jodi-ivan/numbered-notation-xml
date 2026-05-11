package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/header"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (si *staffInteractor) Render(ctx context.Context, canv canvas.Canvas, part musicxml.Part, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) int {

	relativeY := constant.TITLE_Y_POS + header.HEADER_OFFSET

	staffes := si.SplitLines(ctx, part)

	staffLines := lines.NewLineStaff(timeSignature, keySignature)

	// TODO: config toggle here
	x := staffLines.GetLeftIndentWithTimeSignature()
	info := StaffInfo{
		NextLineRenderer: []*entity.NoteRenderer{},
	}
	oldMarginButtom := 0
	for i, st := range staffes {
		info = si.RenderStaff(ctx, canv, x, relativeY, i, i == len(staffes)-1, keySignature, timeSignature, st, info.NextLineRenderer...)
		relativeY = relativeY + STAFF_LINE_DISTANCE + 70 + info.MarginBottom

		nextMeasureNumber := 1 + len(st)
		if len(st) > 0 {
			measure := st[0]
			if i+1 < len(staffes) && len(staffes[i+1]) > 0 {
				measure = staffes[i+1][0]
			}
			nextMeasureNumber = measure.Number
		}

		if info.ForceNewLine {
			relativeY += oldMarginButtom
		}
		if info.Multiline {
			x = info.MarginLeft + barline.BARLINE_AFTER_SPACE
			info.Multiline = false
		} else {
			x = staffLines.GetLeftIndent(nextMeasureNumber)
		}
		oldMarginButtom = info.MarginBottom
	}

	for len(info.NextLineRenderer) > 0 {
		x = staffLines.GetLeftIndent(info.NextLineRenderer[0].MeasureNumber)
		idx := len(staffes) - 1
		if idx <= 0 {
			idx = 1
		}
		info = si.RenderStaff(ctx, canv, x, relativeY, idx, true, keySignature, timeSignature, nil, info.NextLineRenderer...)
		relativeY += info.MarginBottom + STAFF_LINE_DISTANCE + 70
	}

	return relativeY
}
