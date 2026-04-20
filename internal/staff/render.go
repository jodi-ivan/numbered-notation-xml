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
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (si *staffInteractor) Render(ctx context.Context, canv canvas.Canvas, part musicxml.Part, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, metadata *repository.HymnMetadata) int {

	relativeY := constant.TITLE_Y_POS + header.HEADER_OFFSET

	staffes := si.SplitLines(ctx, part)
	x := constant.LAYOUT_INDENT_LENGTH
	info := StaffInfo{
		NextLineRenderer: []*entity.NoteRenderer{},
	}
	oldMarginButtom := 0
	for _, st := range staffes {
		data := StaffData{
			TimeSig:       timeSignature,
			KeySig:        keySignature,
			PrevNotes:     info.NextLineRenderer,
			SyllableCount: info.SyllableCount,
			IndexStart:    info.EndIndex,
			ReffAtStart:   info.StartRenderOtherNotes,
		}

		info = si.RenderStaff(ctx, canv, x, relativeY, metadata, st, data)
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
		data := StaffData{
			TimeSig:       timeSignature,
			KeySig:        keySignature,
			PrevNotes:     info.NextLineRenderer,
			SyllableCount: info.SyllableCount,
			IndexStart:    info.EndIndex,
			ReffAtStart:   info.StartRenderOtherNotes,
		}
		x = constant.LAYOUT_INDENT_LENGTH
		info = si.RenderStaff(ctx, canv, x, relativeY, metadata, nil, data)
		relativeY += info.MarginBottom + 70
	}

	return relativeY
}
