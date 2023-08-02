package renderer

import (
	"context"
	"fmt"
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

func RenderStaff(ctx context.Context, canv canvas.Canvas, x, y int, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature, measures []musicxml.Measure) (multiline bool, marginBottom, marginLeft int) {
	restBeginning := false

	slurTiesRenderer := []*entity.NoteRenderer{}
	lastXCoordinate := float64(0)

	var nextMeasure musicxml.Measure
	canv.Group("class='staff'")
	for measureIndex, measure := range measures {
		measure.Build()
		if measureIndex < len(measures)-1 {
			nextMeasure = measures[measureIndex+1]
		}
		currTimesig := timeSignature.GetTimesignatureOnMeasure(ctx, measure.Number)
		rctx := context.WithValue(ctx, constant.CtxKeyMeasureNum, measure.Number)
		rctx = context.WithValue(rctx, constant.CtxKeyTimeSignature, currTimesig)
		notes := []*entity.NoteRenderer{}

		canv.Group("class='measure'", fmt.Sprintf("id='measure-%d'", measure.Number))

		for notePos, note := range measure.Notes {

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
				PositionX:    x,
				PositionY:    int(y),
				Note:         n,
				NoteLength:   note.Type,
				Octave:       octave,
				Striketrough: strikethrough,
				IsRest:       (note.Rest != nil),
				Beam:         map[int]entity.Beam{},
				IsNewLine:    measure.NewLineIndex == notePos,
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

				slurTiesRenderer = append(slurTiesRenderer, renderer)

				// breath mark
				hasBreathMark = note.Notations.Articulation != nil &&
					note.Notations.Articulation.BreathMark != nil
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
				marginBottom = ((len(note.Lyric) - 1) * 25)
				renderer.Lyric = make([]entity.Lyric, len(note.Lyric))
				for i, currLyric := range note.Lyric {
					renderer.Lyric[i] = entity.Lyric{
						Text:     currLyric.Text.Value,
						Syllabic: currLyric.Syllabic,
					}
					currWidth := int(math.Round(lyric.CalculateLyricWidth(currLyric.Text.Value)))
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
					PositionY:  int(y),
					Width:      LOWERCASE_LENGTH,
					IsDotted:   additional.IsDotted,
					NoteLength: additional.Type,
					Beam:       map[int]entity.Beam{},
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
				})
			}

		}

		if len(measure.Barline) > 0 && measure.Barline[0].Location == musicxml.BarlineLocationLeft {
			RenderBarline(ctx, canv, measure.Barline[0], Coordinate{
				X: float64(x),
				Y: float64(y),
			})

			x += 5

			if measure.Barline[0].Repeat != nil {
				//FIXED: the x value does not add
				x += UPPERCASE_LENGTH
			}
		}

		// part x
		canv.Group("class='note'", "style='font-family:Old Standard TT;font-weight:500'")
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
					canv.Text(lastDotLoc+UPPERCASE_LENGTH, y, ".")
					revisionX[i] = lastDotLoc + UPPERCASE_LENGTH
					lastDotLoc = lastDotLoc + UPPERCASE_LENGTH
				} else {
					canv.Text(xNotes+UPPERCASE_LENGTH, y, ".")
					revisionX[i] = xNotes + UPPERCASE_LENGTH
					lastDotLoc = xNotes + UPPERCASE_LENGTH
				}
				continueDot = true
			} else if n.Articulation != nil && n.Articulation.BreathMark != nil {
				if prev != nil && prev.IsLengthTakenFromLyric {
					x -= prev.Width - LOWERCASE_LENGTH
				}
				x += 5
				canv.Text(x, y-10, ",")
				x += LOWERCASE_LENGTH
			} else {
				if continueDot {
					// FIXED:the dotted does not adding pad to the next notes
					x += LOWERCASE_LENGTH
				}
				canv.Text(x, y, fmt.Sprintf("%d", n.Note))
				xNotes = x
				continueDot = false
				if n.Striketrough {
					canv.Line(x+10, y-16, x, y+5, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.45")
				}
				dotCount = 0
			}

			n.PositionX = x
			n.PositionY = y
			x += n.Width
			if prev != nil && prev.IsLengthTakenFromLyric && n.IsDotted {
				x = x - n.Width
			}
			if n.IsNewLine {
				x = LAYOUT_INDENT_LENGTH
				multiline = multiline || true
				y += 70 + marginBottom
			}
			n.IndexPosition = i
			prev = n
			// FIXED: the dotted does not give proper space at the end of measure
			// FIXED: the one dot on the last measure give uncessary space
			if n.IsDotted && i == len(notes)-1 && dotCount > 1 {
				x += LOWERCASE_LENGTH
			}

		}

		canv.Gend() // note group

		var skipPrintBarline bool

		// FIXED: Print it as glyph
		// FIXED: skip if it has next forward bar line
		for _, barlineNextMeasure := range nextMeasure.Barline {
			if barlineNextMeasure.Location == musicxml.BarlineLocationLeft {
				// there is left barline on the next measure, skipp the regular barline
				skipPrintBarline = true
				break
			}
		}

		if !skipPrintBarline {
			barline := musicxml.Barline{
				BarStyle: musicxml.BarLineStyleRegular,
			}

			if len(measure.Barline) > 0 {
				if measure.Barline[0].Location == musicxml.BarlineLocationRight {
					barline = measure.Barline[0]
				}
			}

			if barline.Repeat != nil && barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
				x += 5
			}
			RenderBarline(ctx, canv, barline, Coordinate{
				X: float64(x),
				Y: float64(y),
			})
		}

		lastXCoordinate = math.Max(lastXCoordinate, float64(x))
		if multiline {
			marginLeft = int(x) + LOWERCASE_LENGTH
		}

		canv.Group("class='lyric'", "style='font-family:Caladea'")
		for _, n := range notes {
			if len(n.Lyric) > 0 {
				for i, l := range n.Lyric {
					if l.Text != "" {
						xPos := n.PositionX
						if n.PositionX == LAYOUT_INDENT_LENGTH {
							xPos += int(lyric.CalculateMarginLeft(l.Text))
						}
						canv.Text(xPos, n.PositionY+25+(i*20), l.Text)
					}

				}
			}
		}
		canv.Gend()

		x += LOWERCASE_LENGTH

		for i, rev := range revisionX {
			note := notes[i]

			note.PositionX = rev
			notes[i] = note
		}
		RenderOctave(rctx, canv, notes)
		RenderBeam(rctx, canv, notes, measure.Number)

		canv.Gend() // measure group
	}
	// align shit here
	RenderSlurTies(ctx, canv, slurTiesRenderer, lastXCoordinate)
	canv.Gend() // staff group

	return
}
