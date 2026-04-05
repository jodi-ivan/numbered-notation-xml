package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/moveabledo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm/splitter"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Staff interface {
	RenderStaff(ctx context.Context, canv canvas.Canvas, x, y int, isLastStaff bool, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, measures []musicxml.Measure, prevNotes ...*entity.NoteRenderer) StaffInfo
	SplitLines(ctx context.Context, part musicxml.Part) [][]musicxml.Measure
	SetMeasureTextRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note, directionDashses map[int]musicxml.DirectionDashesType, isLastNote bool) bool
	Render(ctx context.Context, canv canvas.Canvas, part musicxml.Part, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) int
}

type staffInteractor struct {
	Barline     barline.Barline
	Lyric       lyric.Lyric
	Numbered    numbered.Numbered
	BreathPause breathpause.BreathPause
	Rhythm      rhythm.Rhythm
	RenderAlign RenderStaffWithAlign
}

func NewStaff() Staff {
	barlineInteractor := barline.NewBarline()
	lyricInteractor := lyric.NewLyric()
	return &staffInteractor{
		Barline:     barlineInteractor,
		Lyric:       lyricInteractor,
		Numbered:    numbered.New(lyricInteractor, barlineInteractor),
		BreathPause: breathpause.New(),
		Rhythm:      rhythm.New(splitter.New()),
		RenderAlign: NewRenderAlign(),
	}
}

func (si *staffInteractor) RenderStaff(ctx context.Context, canv canvas.Canvas, x, y int, isLastStaff bool, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, measures []musicxml.Measure, prevNotes ...*entity.NoteRenderer) (staffInfo StaffInfo) {

	staffInfo.NextLineRenderer = []*entity.NoteRenderer{}

	var lastRightBarlinePosition *entity.Coordinate

	align := [][]*entity.NoteRenderer{}
	if len(prevNotes) > 0 {
		align, staffInfo = ProcessPreviousLines(prevNotes, y)
	}
	for mi, measure := range measures {
		measure.Build()
		notes := []*entity.NoteRenderer{}

		currTimesig := timeSignature.GetTimesignatureOnMeasure(ctx, measure.Number)
		currKeySig := keySignature.GetKeyOnMeasure(ctx, measure.Number)

		rctx := context.WithValue(ctx, constant.CtxKeyMeasureNum, measure.Number)
		rctx = context.WithValue(rctx, constant.CtxKeyTimeSignature, currTimesig)
		alignMeasures := []*entity.NoteRenderer{}

		// barline
		if len(measure.Barline) > 0 {
			leftBarlineRenderer, barlineInfo := si.Barline.GetRendererLeftBarline(measure, x, lastRightBarlinePosition)
			if leftBarlineRenderer != nil {
				if h, ok := measure.PrefixHeader[0]; ok {
					leftBarlineRenderer.LeadingHeader = h
				}
				alignMeasures = append(alignMeasures, leftBarlineRenderer)
				x += barlineInfo.XIncrement
			}
		}

		skipNote := map[int]bool{}
		for notePos, note := range measure.Notes {

			if skipNote[notePos] {
				continue
			}

			n, octave, strikethrough := moveabledo.GetNumberedNotation(currKeySig, note)
			noteLength := timeSignature.GetNoteLength(rctx, measure.Number, note)

			if rhythm.HasTies(note) && (notePos+1 < len(measure.Notes)) {
				if mergedLength, mergedNote := rhythm.MergeNotes(ctx, note, measure.Notes[notePos+1], currTimesig); mergedLength > noteLength {
					note, noteLength = mergedNote, mergedLength
					skipNote[notePos+1] = true
				}

			}

			// additionalRenderer is all the new notes that needs represented in numbered when the original musicxml doesnot
			// for example a half note C have to be represented by following . next to number
			additionalNotes := si.Numbered.GetLengthNote(rctx, timeSignature, measure.Number, noteLength)
			if skipNote[notePos+1] {
				// split notes by the beam. currently only happen when there is ties
				next := measure.Notes[notePos+1]
				if notePos+2 < len(measure.Notes) {
					next = measure.Notes[notePos+2]
				}
				additionalNotes = si.Numbered.SplitNote(ctx, noteLength, currTimesig, note.Type, next.Type)
			}
			renderer := &entity.NoteRenderer{
				PositionX:     x,
				PositionY:     int(y),
				Note:          n,
				NoteLength:    note.Type,
				Octave:        octave,
				Strikethrough: strikethrough,
				IsRest:        (note.Rest != nil),
				Beam:          map[int]entity.Beam{},
				IsNewLine:     measure.NewLineIndex[notePos],
				MeasureNumber: measure.Number,

				TimeModifications: note.TimeModification,

				LeadingHeader: measure.PrefixHeader[notePos],
			}

			if note.Notations != nil && note.Notations.Fermata != nil {
				renderer.Fermata = note.Notations.Fermata
			}

			staffInfo.Multiline = staffInfo.Multiline || renderer.IsNewLine

			// text above the measure
			isLastNote := notePos == len(measure.Notes)-1 && mi == len(measures)-1
			hasMeasureText := si.SetMeasureTextRenderer(renderer, note, measure.DirectionDashes[notePos], isLastNote)
			if hasMeasureText && (y == FIRST_STAFF_Y_POS || renderer.Fermata != nil) {
				y += MEASURE_TEXT_OFFSET
				staffInfo.MarginBottom = MEASURE_TEXT_OFFSET
			}

			si.Rhythm.SetRhythmNotation(renderer, note, n)

			// lyric
			verseInfo := si.Lyric.SetLyricRenderer(renderer, note)
			if staffInfo.MarginBottom < verseInfo.MarginBottom {
				staffInfo.MarginBottom = verseInfo.MarginBottom
			}

			additonalRenderer := si.Numbered.RendererFromAdditional(note, renderer, additionalNotes)
			if len(additonalRenderer) > 2 {
				additionalNote := additonalRenderer[len(additonalRenderer)-1]
				shouldReplace := notePos+2 < len(measure.Notes) && note.Type == additionalNotes[len(additionalNotes)-1].Type
				if skipNote[notePos+1] && measure.Notes[notePos+1].IsBreathMark() && shouldReplace {
					additonalRenderer[len(additonalRenderer)-1] = numbered.ReplaceDotWithNumbered(additionalNote, renderer)
				}

			}

			notes = append(notes, additonalRenderer...)

			breathPauseRenderer := si.BreathPause.SetAndGetBreathPauseRenderer(renderer, note)
			if breathPauseRenderer != nil {
				notes = append(notes, breathPauseRenderer)
			}

		}

		x, y = si.Rhythm.AdjustMultiDottedRenderer(notes, x, y)
		var rightBarlineRenderer *entity.NoteRenderer
		x, rightBarlineRenderer = si.Barline.GetRendererRightBarline(measure, x)

		if staffInfo.Multiline {
			staffInfo.MarginLeft = int(x) + constant.LOWERCASE_LENGTH

		}

		x += constant.LOWERCASE_LENGTH

		filteredNotes := []*entity.NoteRenderer{}
		for _, note := range notes {
			filteredNotes = append(filteredNotes, note)
			if note.IsNewLine {
				break
			}
		}

		// alignMeasures = append(alignMeasures, filteredNotes...)
		if staffInfo.Multiline {
			// // -------
			if len(staffInfo.NextLineRenderer) == 0 && len(align) > 0 && staffInfo.ForceNewLine {
				staffInfo.MarginLeft = constant.LAYOUT_INDENT_LENGTH
				si.Rhythm.AdjustMultiDottedRenderer(notes, constant.LAYOUT_INDENT_LENGTH, y)
				notes = append(notes, rightBarlineRenderer)
				staffInfo.NextLineRenderer = notes
			} else {
				nextstaffInfo := PrepareNextLines(staffInfo, notes, rightBarlineRenderer)
				staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, nextstaffInfo.NextLineRenderer...)
				alignMeasures = append(alignMeasures, filteredNotes...)
				staffInfo.MarginLeft = nextstaffInfo.MarginLeft
				if staffInfo.MarginBottom < nextstaffInfo.MarginBottom {
					staffInfo.MarginBottom = nextstaffInfo.MarginBottom
				}
			}
		} else {
			alignMeasures = append(alignMeasures, filteredNotes...)

			lastRightBarlinePosition = &entity.Coordinate{
				X: float64(rightBarlineRenderer.PositionX),
				Y: float64(y),
			}
			x += barline.BARLINE_AFTER_SPACE
			if measure.RightMeasureText != nil {
				rightBarlineRenderer.MeasureText = []musicxml.MeasureText{
					{
						Text:          measure.RightMeasureText.Text,
						RelativeY:     measure.RightMeasureText.RelativeY,
						TextAlignment: musicxml.TextAlignmentRight,
					},
				}

			}
			alignMeasures = append(alignMeasures, rightBarlineRenderer)
		}

		if len(alignMeasures) > 0 {
			if keySignature.IsMixed {
				if keyChanges, ok := keySignature.MeasureText[measure.Number]; ok {
					renderer := alignMeasures[0]
					renderer.MeasureText = []musicxml.MeasureText{{Text: keyChanges, TextAlignment: musicxml.TextAlignmentLeft}}
				}

				lastMeasure := mi == len(measures)-1
				noCarryOverLine := len(staffInfo.NextLineRenderer) == 0

				if isLastStaff && lastMeasure && noCarryOverLine {
					firstKeySig := keySignature.GetKeyOnMeasure(ctx, 1)
					indicator := keysig.TranstionFromTwoKeySignatures(currKeySig, firstKeySig)

					renderer := alignMeasures[len(alignMeasures)-1]
					renderer.MeasureText = []musicxml.MeasureText{{Text: indicator, TextAlignment: musicxml.TextAlignmentRight}}

				}
			}
			align = append(align, alignMeasures)
		}

	}

	si.RenderAlign.RenderWithAlign(ctx, canv, y, timeSignature, align)

	return
}
