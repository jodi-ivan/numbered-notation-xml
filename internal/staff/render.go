package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/header"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/layout"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

func (si *staffInteractor) RenderWithWrap(ctx context.Context, canv canvas.Canvas, x, y, staffPos int, metadata *entity.HymnMetaData, measures []musicxml.Measure, data StaffData) StaffInfo {
	isGregorian := params.GetGregorianToggle(ctx)
	notes, staffInfo := si.RenderStaff(ctx, canv, x, y, staffPos, metadata, measures, data)

	di, ok := params.GetDeviceDimensions(ctx)
	if (ok && di == 0) || !ok {
		staffInfo.MarginBottom += si.RenderAlign.RenderWithAlign(ctx, canv, staffPos, y, data.TimeSig, data.KeySig, notes)
		return staffInfo
	}

	wrapped := Wrap(ctx, x, di, data.KeySig, notes[0])

	baseMargin := staffInfo.MarginBottom
	currentY := y
	lastMargin := 0

	for i, line := range wrapped {
		lineMargin := baseMargin // fresh per line, never accumulated across lines

		if i > 0 {
			_, prev := ProcessPreviousLines(ctx, wrapped[i-1], data.KeySig, currentY)
			if prev.MarginBottom > lineMargin {
				lineMargin = prev.MarginBottom
			}
			staffInfo.MarginLeft = prev.MarginLeft
		}

		lineMargin += si.RenderAlign.RenderWithAlign(ctx, canv, i, currentY, data.TimeSig, data.KeySig, [][]*entity.NoteRenderer{line})

		if i == len(wrapped)-1 {
			lastMargin = lineMargin // caller will add 70 + this
		} else {
			currentY += 70 + lineMargin
			if isGregorian {
				currentY += STAFF_LINE_DISTANCE
			}
		}
	}

	// Caller does: relativeY += 70 + MarginBottom
	// We want:     total = (currentY - y) + 70 + lastMargin
	staffInfo.MarginBottom = (currentY - y) + lastMargin
	return staffInfo
}

func (si *staffInteractor) Render(ctx context.Context, canv canvas.Canvas, part musicxml.Part, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, metadata *entity.HymnMetaData) int {

	relativeY := constant.TITLE_Y_POS + header.HEADER_OFFSET

	staffes := si.SplitLines(ctx, part)

	if params.IgnoreNewLine(ctx) {
		// remove the classical (the book marking) new line breaks to one line.
		//  hopefuly we found which measure to break later on
		flattenedStaffes := make([]musicxml.Measure, 0)
		for _, staff := range staffes {
			for _, note := range staff {
				flattenedStaffes = append(flattenedStaffes, note)
			}
		}

		staffes = [][]musicxml.Measure{flattenedStaffes}

	}

	staffLines := lines.NewLineStaff(timeSignature, keySignature)

	isGregorian := params.GetGregorianToggle(ctx)
	x, _ := layout.GetMargin(ctx)
	if isGregorian {
		x = staffLines.GetLeftIndentWithTimeSignature(ctx)
	}
	info := StaffInfo{
		NextLineRenderer: []*entity.NoteRenderer{},
		SyllableTracker: entity.SyllableCounterTracker{
			SyllableOffset: map[int]int{},
			RepeatInfo:     []*musicxml.RepeatInfo{},
		},
	}
	oldMarginButtom := 0
	for i, st := range staffes {
		data := StaffData{
			TimeSig:       timeSignature,
			KeySig:        keySignature,
			PrevNotes:     info.NextLineRenderer,
			SyllableCount: info.SyllableTracker.SyllableCount,
			IndexStart:    info.EndIndex,
			ReffAtStart:   info.SyllableTracker.StartRenderOtherNotes,
			RepeatInfo:    info.SyllableTracker.RepeatInfo,

			SyllableOffset: info.SyllableTracker.SyllableOffset,
		}
		info = si.RenderWithWrap(ctx, canv, x, relativeY, i, metadata, st, data)
		info.SyllableTracker.RepeatInfo = append(data.RepeatInfo, info.SyllableTracker.RepeatInfo...)
		relativeY = relativeY + 70 + info.MarginBottom
		if isGregorian {
			relativeY += STAFF_LINE_DISTANCE
		}

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
			x, _ = layout.GetMargin(ctx)
			if isGregorian {
				x = staffLines.GetLeftIndent(ctx, nextMeasureNumber)
			}
		}
		oldMarginButtom = info.MarginBottom
	}

	for len(info.NextLineRenderer) > 0 {
		data := StaffData{
			TimeSig:       timeSignature,
			KeySig:        keySignature,
			PrevNotes:     info.NextLineRenderer,
			SyllableCount: info.SyllableTracker.SyllableCount,
			IndexStart:    info.EndIndex,
			ReffAtStart:   info.SyllableTracker.StartRenderOtherNotes,
			RepeatInfo:    info.SyllableTracker.RepeatInfo,

			SyllableOffset: info.SyllableTracker.SyllableOffset,
		}
		x, _ = layout.GetMargin(ctx)
		if isGregorian {
			x = staffLines.GetLeftIndent(ctx, info.NextLineRenderer[0].MeasureNumber)
		}
		idx := len(staffes) - 1
		if idx <= 0 {
			idx = 1
		}
		info = si.RenderWithWrap(ctx, canv, x, relativeY, idx, metadata, nil, data)
		relativeY += info.MarginBottom + 70
		if isGregorian {
			relativeY += STAFF_LINE_DISTANCE
		}
	}

	return relativeY
}
