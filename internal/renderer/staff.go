package renderer

import (
	"context"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/moveabledo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderStaff(ctx context.Context, canv canvas.Canvas, x, y int, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, measures []musicxml.Measure, prevNotes ...*entity.NoteRenderer) (staffInfo StaffInfo) {
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

			leftBarline := measure.Barline[0]
			if (leftBarline.Location == musicxml.BarlineLocationLeft) && (leftBarline.BarStyle != musicxml.BarLineStyleRegular) {
				pos := x
				if lastRightBarlinePosition != nil {
					pos = int(lastRightBarlinePosition.X)
				}
				alignMeasures = append(alignMeasures, &entity.NoteRenderer{
					PositionX:     pos,
					Width:         int(barlineWidth[leftBarline.BarStyle]),
					Barline:       &leftBarline,
					MeasureNumber: measure.Number,
				})

				x += 5

				if leftBarline.Repeat != nil {
					x += UPPERCASE_LENGTH
				}
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
			additionalRenderer := numbered.RenderLengthNote(rctx, timeSignature, measure.Number, noteLength)

			renderer := &entity.NoteRenderer{
				PositionX:     x,
				PositionY:     int(y),
				Note:          n,
				NoteLength:    note.Type,
				Octave:        octave,
				Striketrough:  strikethrough,
				IsRest:        (note.Rest != nil),
				Beam:          map[int]entity.Beam{},
				IsNewLine:     measure.NewLineIndex == notePos,
				MeasureNumber: measure.Number,

				TimeMofication: note.TimeModification,
			}

			for _, mt := range note.MeasureText {
				if renderer.MeasureText != nil {
					renderer.MeasureText = []musicxml.MeasureText{}
				}
				alignment := musicxml.TextAlignmentLeft
				if notePos == len(measure.Notes)-1 {
					alignment = musicxml.TextAlignmentRight
				}
				renderer.MeasureText = append(renderer.MeasureText, musicxml.MeasureText{
					Text:          mt.Text,
					RelativeY:     mt.RelativeY,
					TextAlignment: alignment,
				})
			}

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

			hasBreathMark := false

			if note.Notations != nil {

				for i, slur := range note.Notations.Slur {
					if i == 0 {
						renderer.Slur = map[int]entity.Slur{}
					}

					_, existing := renderer.Slur[slur.Number]
					if !existing {
						renderer.Slur[slur.Number] = entity.Slur{
							Number: slur.Number,
							Type:   slur.Type,
						}
					} else {
						renderer.Slur[slur.Number] = entity.Slur{
							Number: slur.Number,
							Type:   musicxml.NoteSlurTypeHop,
						}
					}

				}

				if note.Notations.Tied != nil {
					renderer.Tie = &entity.Slur{
						Number: n,
						Type:   note.Notations.Tied.Type,
					}
				}

				// breath mark
				hasBreathMark = note.Notations.Articulation != nil &&
					note.Notations.Articulation.BreathMark != nil

				renderer.Tuplet = note.Notations.Tuplet
			}

			if len(note.Beam) > 0 {
				if currTimesig.BeatType != 4 {
					for _, beam := range note.Beam {
						renderer.Beam[beam.Number] = entity.Beam{
							Number: beam.Number,
							Type:   beam.State,
						}
					}
				}
			}

			// lyric
			var lyricWidth, noteWidth int

			if len(note.Lyric) > 0 {
				staffInfo.MarginBottom = ((len(note.Lyric) - 1) * 25)
				renderer.Lyric = make([]entity.Lyric, len(note.Lyric))
				for i, currLyric := range note.Lyric {
					lyricText := ""
					l := entity.Lyric{
						Syllabic: currLyric.Syllabic,
					}

					texts := []entity.Text{}
					for _, t := range currLyric.Text {
						lyricText += t.Value
						texts = append(texts, entity.Text{
							Value:     t.Value,
							Underline: t.Underline,
						})
					}

					l.Text = texts

					renderer.Lyric[i] = l
					currWidth := int(math.Round(lyric.CalculateLyricWidth(lyricText)))
					if currLyric.Syllabic == musicxml.LyricSyllabicTypeEnd || currLyric.Syllabic == musicxml.LyricSyllabicTypeSingle {
						//FIXME: edge cases kj-101, [ki]dung ma[laikat] no spaces between them
						currWidth += LOWERCASE_LENGTH
					}
					currWidth += 4 // lyric padding

					lyricWidth = int(math.Max(float64(lyricWidth), float64(currWidth)))
				}

			}

			noteWidth = LOWERCASE_LENGTH

			if noteWidth > lyricWidth {
				renderer.Width = noteWidth
				renderer.IsLengthTakenFromLyric = false
			} else {
				renderer.Width = lyricWidth
				renderer.IsLengthTakenFromLyric = true
				if float64(lyricWidth) > float64(noteWidth+UPPERCASE_LENGTH) {
					renderer.Width = UPPERCASE_LENGTH * 1.7
				}
			}

			notes = append(notes, renderer)

			for i, additional := range additionalRenderer {
				if i == 0 {
					continue
				}
				additionalNote := &entity.NoteRenderer{
					PositionY:     int(y),
					Width:         LOWERCASE_LENGTH,
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
				notes = append(notes, additionalNote)

			}
			if hasBreathMark {
				// FIXME: the breath mark stopped the continuation of the beam
				notes = append(notes, &entity.NoteRenderer{
					Articulation: &entity.Articulation{
						BreathMark: &entity.ArticulationTypesBreathMark,
					},
					MeasureNumber: measure.Number,

					// move the new line indicator to this
					IsNewLine: renderer.IsNewLine,
				})

				if renderer.IsNewLine {
					// remove the new line, since it is transferrerd to the breathmark
					renderer.IsNewLine = false
				}
			}

		}

		xNotes := 0
		continueDot := false
		lastDotLoc := 0
		dotCount := 0

		var prev *entity.NoteRenderer
		revisionX := map[int]int{}
		for i, n := range notes {
			if n.IsDotted {
				dotCount++
				if continueDot {
					revisionX[i] = lastDotLoc + UPPERCASE_LENGTH
					lastDotLoc = lastDotLoc + UPPERCASE_LENGTH
				} else {
					revisionX[i] = xNotes + UPPERCASE_LENGTH
					lastDotLoc = xNotes + UPPERCASE_LENGTH
				}
				continueDot = true
			} else if n.Articulation != nil && n.Articulation.BreathMark != nil {
				if prev != nil && prev.IsLengthTakenFromLyric {
					x -= prev.Width - LOWERCASE_LENGTH
				}
			} else {
				if continueDot {
					x += LOWERCASE_LENGTH
				}
				xNotes = x
				continueDot = false
				dotCount = 0
			}

			n.PositionX = x
			n.PositionY = y
			x += n.Width
			if prev != nil && prev.IsLengthTakenFromLyric && n.IsDotted {
				x = x - n.Width
			}
			if n.IsNewLine {
				x = constant.LAYOUT_INDENT_LENGTH
				staffInfo.Multiline = staffInfo.Multiline || true
			}
			n.IndexPosition = i
			prev = n
			// FIXED: the dotted does not give proper space at the end of measure
			// FIXED: the one dot on the last measure give uncessary space
			if n.IsDotted && i == len(notes)-1 && dotCount > 1 {
				x += LOWERCASE_LENGTH
			}

		}

		barline := musicxml.Barline{
			BarStyle: musicxml.BarLineStyleRegular,
		}

		if len(measure.Barline) == 1 {
			if measure.Barline[0].Location == musicxml.BarlineLocationRight {
				barline = measure.Barline[0]
			}
		} else if len(measure.Barline) > 1 {
			if measure.Barline[1].Location == musicxml.BarlineLocationRight {
				barline = measure.Barline[1]
			}
		}
		if barline.Repeat != nil && barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
			x += 5
		}
		barlineX := x

		if staffInfo.Multiline {
			staffInfo.MarginLeft = int(x) + LOWERCASE_LENGTH
		}

		x += LOWERCASE_LENGTH

		for i, rev := range revisionX {
			note := notes[i]

			note.PositionX = rev
			notes[i] = note
		}

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
					break
				}
				if i > indexNewLine {
					staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, note)
				}
			}
			staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, &entity.NoteRenderer{
				Barline:       &barline,
				PositionX:     barlineX,
				MeasureNumber: measure.Number,
			})

		} else {
			barlineRenderer := &entity.NoteRenderer{
				Barline:       &barline,
				PositionX:     barlineX,
				MeasureNumber: measure.Number,
			}
			lastRightBarlinePosition = &entity.Coordinate{
				X: float64(barlineX),
				Y: float64(y),
			}
			if measure.RightMeasureText != nil {
				barlineRenderer.MeasureText = []musicxml.MeasureText{
					musicxml.MeasureText{
						Text:          measure.RightMeasureText.Text,
						RelativeY:     measure.RightMeasureText.RelativeY,
						TextAlignment: musicxml.TextAlignmentRight,
					},
				}

			}
			alignMeasures = append(alignMeasures, barlineRenderer)
		}

		align = append(align, alignMeasures)
	}
	RenderWithAlign(ctx, canv, y, align)

	return
}
