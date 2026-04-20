package staff

import (
	"context"
	"math"
	"slices"

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
	"github.com/jodi-ivan/numbered-notation-xml/internal/verse"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Staff interface {
	RenderStaff(ctx context.Context, canv canvas.Canvas, x, y int, metadata *repository.HymnMetadata, measures []musicxml.Measure, data StaffData) StaffInfo
	SplitLines(ctx context.Context, part musicxml.Part) [][]musicxml.Measure
	SetMeasureTextRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note, directionDashses map[int]musicxml.DirectionDashesType, isLastNote bool) bool
	Render(ctx context.Context, canv canvas.Canvas, part musicxml.Part, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, metadata *repository.HymnMetadata) int
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

func (si *staffInteractor) RenderStaff(ctx context.Context, canv canvas.Canvas, x, y int, metadata *repository.HymnMetadata, measures []musicxml.Measure, data StaffData) (staffInfo StaffInfo) {

	staffInfo.NextLineRenderer = []*entity.NoteRenderer{}

	var lastRightBarlinePosition *barline.CoordinateWithBarline
	yOffsetRepeat, yOffset := false, false
	refreinStartNote := false
	align := [][]*entity.NoteRenderer{}
	pos := 0
	startSyllable := data.SyllableCount
	if len(data.PrevNotes) > 0 {
		align, staffInfo = ProcessPreviousLines(data.PrevNotes, y)
		pos = data.PrevNotes[len(data.PrevNotes)-1].IndexPosition + 1
	}
	for mi, measure := range measures {
		mSyllcount := 0
		measure.Build()

		notes := []*entity.NoteRenderer{}

		currTimesig := data.TimeSig.GetTimesignatureOnMeasure(ctx, measure.Number)
		currKeySig := data.KeySig.GetKeyOnMeasure(ctx, measure.Number)

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

				yOffsetRepeat = yOffsetRepeat || leftBarlineRenderer.Barline.Ending != nil
			}
		}

		skipNote := map[int]bool{}
		for notePos, note := range measure.Notes {

			if skipNote[notePos] {
				continue
			}

			n, octave, strikethrough := moveabledo.GetNumberedNotation(currKeySig, note)
			noteLength := data.TimeSig.GetNoteLength(rctx, measure.Number, note)

			if rhythm.HasTies(note) && (notePos+1 < len(measure.Notes)) && currTimesig.IsCommonTime() {
				if mergedLength, mergedNote := rhythm.MergeNotes(ctx, note, measure.Notes[notePos+1], currTimesig); mergedLength > noteLength && mergedLength < 3 {
					note, noteLength = mergedNote, mergedLength
					skipNote[notePos+1] = true
				}

			}

			// additionalRenderer is all the new notes that needs represented in numbered when the original musicxml doesnot
			// for example a half note C have to be represented by following . next to number
			additionalNotes := si.Numbered.GetLengthNote(rctx, data.TimeSig, measure.Number, noteLength)
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
				IndexPosition: pos + note.IndexPosition + data.IndexStart,
			}

			// log.Println(measure.Number, pos+note.IndexPosition+data.IndexStart, data.IndexStart, pos, note.IndexPosition)
			staffInfo.EndIndex = pos + note.IndexPosition + data.IndexStart

			if note.Notations != nil && note.Notations.Fermata != nil {
				renderer.Fermata = note.Notations.Fermata
			}

			staffInfo.Multiline = staffInfo.Multiline || renderer.IsNewLine

			// text above the measure
			isLastNote := notePos == len(measure.Notes)-1 && mi == len(measures)-1
			hasMeasureText := si.SetMeasureTextRenderer(renderer, note, measure.DirectionDashes[notePos], isLastNote)
			if hasMeasureText || (len(note.MeasureText) > 0 && y == FIRST_STAFF_Y_POS) {
				yOffset = true

				isRefrein := slices.ContainsFunc(renderer.MeasureText, func(t musicxml.MeasureText) bool {
					return t.Text == "Refrein"
				})

				if hasMeasureText && isRefrein && renderer.MeasureNumber == 1 && notePos == 0 {
					staffInfo.StartRenderOtherNotes = false
				}

				refreinStartNote = isRefrein
			}

			si.Rhythm.SetRhythmNotation(renderer, note, n)

			// lyric
			verseInfo := si.Lyric.SetLyricRenderer(renderer, note.Lyric)
			if staffInfo.MarginBottom < verseInfo.MarginBottom {
				staffInfo.MarginBottom = verseInfo.MarginBottom
			}

			if verseInfo.HasLyric {
				mSyllcount++
				staffInfo.StartRenderOtherNotes = staffInfo.StartRenderOtherNotes || lyric.HasPrefix(renderer)
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

			if notePos == len(measure.Notes)-1 {
				pos = measure.Notes[len(measure.Notes)-1].IndexPosition + 1
			}

		}

		if (data.ReffAtStart || staffInfo.StartRenderOtherNotes) || (measures[0].Number == 1 && !refreinStartNote) {
			marginBottom := verse.LoadOtherVerse(notes, metadata, startSyllable)
			if staffInfo.MarginBottom < marginBottom {
				staffInfo.MarginBottom = marginBottom
			}
			startSyllable += mSyllcount
			staffInfo.StartRenderOtherNotes = true
		} else {
			mSyllcount = 0
			startSyllable = 0
			staffInfo.SyllableCount = 0
		}

		// data || staffinfo

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

			lastbarPos := entity.NewCoordinate(float64(rightBarlineRenderer.PositionX), float64(y))
			lastRightBarlinePosition = &barline.CoordinateWithBarline{
				Coordinate: lastbarPos,
				Barline:    *rightBarlineRenderer.Barline,
			}

			x += constant.LOWERCASE_LENGTH
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
			align = append(align, alignMeasures)
		}

	}
	if yOffset {
		y += MEASURE_TEXT_OFFSET
		staffInfo.MarginBottom = int(math.Max(MEASURE_TEXT_OFFSET, float64(staffInfo.MarginBottom)))

		if yOffsetRepeat {
			y += 10
			staffInfo.MarginBottom += 10
		}
	}

	staffInfo.SyllableCount += startSyllable

	si.RenderAlign.RenderWithAlign(ctx, canv, y, data.TimeSig, align)

	return
}
