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
	SetMeasureTextRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note, isLastNote bool)
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
	return &staffInteractor{
		Barline:     barline.NewBarline(),
		Lyric:       lyric.NewLyric(),
		Numbered:    numbered.New(),
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

			if rhythm.HasTies(note) {
				if notePos+1 < len(measure.Notes) && rhythm.HasTies(measure.Notes[notePos+1]) {

					endTieNote := measure.Notes[notePos+1]

					mergedNoteLegth := rhythm.MergeNote(ctx, note, endTieNote, currTimesig)
					if mergedNoteLegth < 3 {
						noteLength = mergedNoteLegth

						note = rhythm.TransferStopSlurAndBreathmark(endTieNote, note)

						// dont process next notes
						skipNote[notePos+1] = true
					}
				}
			}

			// additionalRenderer is all the new notes that needs represented in numbered when the original musicxml doesnot
			// for example a half note C have to be represented by following . next to number
			additionalRenderer := si.Numbered.GetLengthNote(rctx, timeSignature, measure.Number, noteLength)
			if skipNote[notePos+1] {
				// split notes by the beam. currently only happen when there is ties
				next := measure.Notes[notePos+1]
				if notePos+2 < len(measure.Notes) {
					next = measure.Notes[notePos+2]
				}
				additionalRenderer = si.Numbered.SplitNote(ctx, noteLength, currTimesig, note.Type, next.Type)
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

				// FIXME: duplicate if the barline exist
				LeadingHeader: measure.PrefixHeader,
			}

			if note.Notations != nil && note.Notations.Fermata != nil {
				renderer.Fermata = note.Notations.Fermata
			}

			staffInfo.Multiline = staffInfo.Multiline || renderer.IsNewLine

			// text above the measure
			si.SetMeasureTextRenderer(renderer, note, notePos == len(measure.Notes)-1)

			if len(additionalRenderer) > 0 {

				// the first additional notes is always altering the current note
				addRenderer := additionalRenderer[0]
				switch addRenderer.Type {
				case musicxml.NoteLength16th:
					renderer.Beam[2] = entity.Beam{
						Number: 2,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					}
					fallthrough
				case musicxml.NoteLengthEighth:
					renderer.Beam[1] = entity.Beam{
						Number: 1,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					}
				}
			}
			// set the beam, slur and ties
			si.Rhythm.SetRhythmNotation(renderer, note, n)

			// lyric
			verseInfo := si.Lyric.SetLyricRenderer(renderer, note)
			if staffInfo.MarginBottom < verseInfo.MarginBottom {
				staffInfo.MarginBottom = verseInfo.MarginBottom
			}

			notes = append(notes, renderer)

			hasBreathMark := note.Notations != nil &&
				note.Notations.Articulation != nil &&
				note.Notations.Articulation.BreathMark != nil

			// additional renderer is a several new renderer because of
			// the conversion to numbered
			// for example, a half note, means an additional note for the dot
			for i, additional := range additionalRenderer {
				if i == 0 {
					continue
				}
				additionalNote := &entity.NoteRenderer{
					PositionY:     int(y),
					Width:         constant.LOWERCASE_LENGTH,
					IsDotted:      additional.IsDotted,
					NoteLength:    additional.Type,
					Beam:          map[int]entity.Beam{},
					MeasureNumber: measure.Number,
					IsNewLine:     renderer.IsNewLine && (i == len(additionalRenderer)-1) && !hasBreathMark,
				}
				if additionalNote.IsNewLine {
					renderer.IsNewLine = !additionalNote.IsNewLine
				}

				shouldReplace := notePos+2 < len(measure.Notes) && note.Type == additional.Type
				if skipNote[notePos+1] && len(additionalRenderer) > 2 && i == len(additionalRenderer)-1 && shouldReplace {
					additionalNote.IsDotted = false
					additionalNote.Note = renderer.Note
					additionalNote.Octave = renderer.Octave
					additionalNote.Strikethrough = renderer.Strikethrough

					renderer.Tie = &entity.Slur{
						Number: 1,
						Type:   musicxml.NoteSlurTypeStart,
					}

					additionalNote.Tie = &entity.Slur{
						Number: 1,
						Type:   musicxml.NoteSlurTypeStop,
					}

				}

				switch additional.Type {
				case musicxml.NoteLength16th:
					additionalNote.Beam[2] = entity.Beam{
						Number: 2,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					}
					fallthrough
				case musicxml.NoteLengthEighth:
					additionalNote.Beam[1] = entity.Beam{
						Number: 1,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					}
				}

				notes = append(notes, additionalNote)

			}

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
					musicxml.MeasureText{
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
					renderer.MeasureText = []musicxml.MeasureText{musicxml.MeasureText{Text: keyChanges, TextAlignment: musicxml.TextAlignmentLeft}}
				}

				lastMeasure := mi == len(measures)-1
				noCarryOverLine := len(staffInfo.NextLineRenderer) == 0

				if isLastStaff && lastMeasure && noCarryOverLine {
					firstKeySig := keySignature.GetKeyOnMeasure(ctx, 1)
					indicator := keysig.TranstionFromTwoKeySignatures(currKeySig, firstKeySig)

					renderer := alignMeasures[len(alignMeasures)-1]
					renderer.MeasureText = []musicxml.MeasureText{musicxml.MeasureText{Text: indicator, TextAlignment: musicxml.TextAlignmentRight}}

				}
			}
			align = append(align, alignMeasures)
		}

	}

	si.RenderAlign.RenderWithAlign(ctx, canv, y, timeSignature, align)

	return
}
