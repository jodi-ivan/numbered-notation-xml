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
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Staff interface {
	RenderStaff(ctx context.Context, canv canvas.Canvas, x, y int, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, measures []musicxml.Measure, prevNotes ...*entity.NoteRenderer) StaffInfo
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
		Rhythm:      rhythm.New(),
		RenderAlign: NewRenderAlign(),
	}
}

func (si *staffInteractor) RenderStaff(ctx context.Context, canv canvas.Canvas, x, y int, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, measures []musicxml.Measure, prevNotes ...*entity.NoteRenderer) (staffInfo StaffInfo) {
	restBeginning := false

	staffInfo.NextLineRenderer = []*entity.NoteRenderer{}

	var lastRightBarlinePosition *entity.Coordinate

	align := [][]*entity.NoteRenderer{}
	if len(prevNotes) > 0 {
		align = append(align, prevNotes)
	}
	for _, measure := range measures {
		measure.Build()
		notes := []*entity.NoteRenderer{}
		currTimesig := timeSignature.GetTimesignatureOnMeasure(ctx, measure.Number)
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
		for notePos, note := range measure.Notes {

			// FIXME: use the hidden attributes on the musicxml files instead of forcing hide every beginnning rest
			// don't print anything when rest on the beginning on the music
			if note.Rest != nil && measure.Number == 1 {

				if notePos == 0 {
					restBeginning = true
					continue
				}

				if restBeginning {
					continue
				}
			}

			restBeginning = false

			n, octave, strikethrough := moveabledo.GetNumberedNotation(keySignature, note)
			noteLength := timeSignature.GetNoteLength(rctx, measure.Number, note)

			additionalRenderer := si.Numbered.GetLengthNote(rctx, timeSignature, measure.Number, noteLength)
			renderer := &entity.NoteRenderer{
				PositionX:     x,
				PositionY:     int(y),
				Note:          n,
				NoteLength:    note.Type,
				Octave:        octave,
				Strikethrough: strikethrough,
				IsRest:        (note.Rest != nil),
				Beam:          map[int]entity.Beam{},
				IsNewLine:     measure.NewLineIndex == notePos,
				MeasureNumber: measure.Number,

				TimeModifications: note.TimeModification,
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
			staffInfo.MarginBottom = verseInfo.MarginBottom

			notes = append(notes, renderer)

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
				//FIXME: the new newline on current renderer has to be transferred to this new dotted renderer
				// currently usually handled by the breathmark, but still handled by the renderer without dotted
				// case test on kj 226. kj 309
				notes = append(notes, additionalNote)

			}
			breathPauseRenderer := si.BreathPause.SetAndGetBreathPauseRenderer(renderer, note)
			if breathPauseRenderer != nil {
				notes = append(notes, breathPauseRenderer)
			}

		}

		x, y = si.Rhythm.AdjustMultiDottedRenderer(notes, x, y)

		barlineX, rightBarlineRenderer := si.Barline.GetRendererRightBarline(measure, x)

		if staffInfo.Multiline {
			staffInfo.MarginLeft = int(x) + constant.LOWERCASE_LENGTH
		}

		x += constant.LOWERCASE_LENGTH

		filteredNotes := []*entity.NoteRenderer{}
		indexNewLine := -1
		for i, note := range notes {
			filteredNotes = append(filteredNotes, note)
			if note.IsNewLine {
				indexNewLine = i
				break
			}
		}

		alignMeasures = append(alignMeasures, filteredNotes...)
		if staffInfo.Multiline {
			for i, note := range notes {
				if i > 0 && note.IndexPosition == 0 {
					// TODO: check what the heck is this?
					break
				}
				if i > indexNewLine {
					staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, note)
				}
			}
			staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, rightBarlineRenderer)

		} else {

			lastRightBarlinePosition = &entity.Coordinate{
				X: float64(barlineX),
				Y: float64(y),
			}
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

		align = append(align, alignMeasures)
	}

	si.RenderAlign.RenderWithAlign(ctx, canv, y, align)

	return
}
