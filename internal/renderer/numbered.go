package renderer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"

	svg "github.com/ajstarks/svgo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/moveabledo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

func googlefont(f string) []byte {
	empty := []byte{}
	r, err := http.Get(gwfURI + url.QueryEscape(f))
	log.Println("error call", err)
	if err != nil {
		return empty
	}
	defer r.Body.Close()
	b, rerr := ioutil.ReadAll(r.Body)
	log.Println(rerr, r.Status, string(b))
	if rerr != nil || r.StatusCode != http.StatusOK {
		return empty
	}

	return b
}

func RenderNumbered(w http.ResponseWriter, r *http.Request, music musicxml.MusicXML) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(200)
	s := svg.New(w)
	s.Start(LAYOUT_WIDTH, 1000)

	s.Def()
	fmt.Fprintf(s.Writer, fontfmt, string(googlefont("Caladea|Old Standard TT|Noto Music")))
	s.DefEnd()

	relativeY := 100
	// render title
	titleX := (LAYOUT_WIDTH / 2) - ((len(music.Credit.Words) * LOWERCASE_LENGTH) / 2) + (LAYOUT_INDENT_LENGTH * 2)
	s.Text(titleX, relativeY, music.Credit.Words)

	// render key signature
	relativeY += 25

	keySignature := keysig.NewKeySignature(music.Part.Measures[0].Attribute.Key)

	humanizedKeySignature := keySignature.String()

	s.Text(LAYOUT_INDENT_LENGTH, relativeY, keySignature.String())

	// render time signature
	// TODO: check the time signature on github issue
	// TODO: time signature changing happens on the top and not on the measure
	/*
		time signatures
		4/4
		3/4
		6/4
		1/4
		6/8 (shown as 3 x 2)
		2/4
	*/
	beat := music.Part.Measures[0].Attribute.Time
	s.Text(LAYOUT_INDENT_LENGTH+(len(humanizedKeySignature)*LOWERCASE_LENGTH), relativeY, fmt.Sprintf("%d ketuk", beat.Beats))
	relativeY += 50

	RenderMeasures(r.Context(), s, LAYOUT_INDENT_LENGTH, relativeY, music.Part)
	s.End()
}

// TODO: breakline
func RenderMeasures(ctx context.Context, s *svg.SVG, x, y int, measures musicxml.Part) {

	keySignature := keysig.NewKeySignature(measures.Measures[0].Attribute.Key)
	x = LAYOUT_INDENT_LENGTH

	totalMeasure := len(measures.Measures)
	restBeginning := false

	timeSignature := timesig.NewTimeSignatures(ctx, measures.Measures)

	for _, measure := range measures.Measures {

		s.Group("class='measure'", fmt.Sprintf("id='measure-%d'", measure.Number))

		notes := []*NoteRenderer{}
		totalNotes := len(measure.Notes)
		currTimesig := timeSignature.GetTimesignatureOnMeasure(ctx, measure.Number)
		rctx := context.WithValue(ctx, constant.CtxKeyMeasureNum, measure.Number)
		rctx = context.WithValue(rctx, constant.CtxKeyTimeSignature, currTimesig)
		for notePosInMeasure, note := range measure.Notes {

			// don't print anything when rest on the beginning on the music
			if note.Rest != nil && measure.Number == 1 {

				if notePosInMeasure == 0 {
					restBeginning = true
					continue
				}

				if restBeginning {
					continue
				}
			}

			restBeginning = false

			// don't print when rest on the last of the music
			if measure.Number == totalMeasure && (notePosInMeasure+1) == totalNotes && note.Rest != nil {
				continue
			}

			n, octave, strikethrough := moveabledo.GetNumberedNotation(keySignature, note)

			noteLength := timeSignature.GetNoteLength(rctx, measure.Number, note)
			additionalRenderer := numbered.RenderLengthNote(rctx, timeSignature, measure.Number, noteLength)

			// aditiomal := &NoteRenderer{}
			renderer := &NoteRenderer{
				PositionX:    x,
				PositionY:    y,
				Note:         n,
				NoteLength:   note.Type,
				Octave:       octave,
				Striketrough: strikethrough,
				IsRest:       (note.Rest != nil),
				Beam:         map[int]Beam{},
			}

			if len(additionalRenderer) > 0 {

				// the first additional notes is always altering the current note
				addRenderer := additionalRenderer[0]
				switch addRenderer.Type {
				case musicxml.NoteLength16th:
					renderer.Beam[2] = Beam{
						Number: 2,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					}
					fallthrough
				case musicxml.NoteLengthEighth:
					renderer.Beam[1] = Beam{
						Number: 1,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					}
				}
			}

			hasBreathMark := false

			if note.Notations != nil {

				// slur
				for i, slur := range note.Notations.Slur {
					if i == 0 {
						renderer.Slur = map[int]Slur{}
					}
					renderer.Slur[slur.Number] = Slur{
						Number: slur.Number,
						Type:   slur.Type,
					}
				}

				if note.Notations.Tied != nil {
					renderer.Tie = &Slur{
						Number: n,
						Type:   note.Notations.Tied.Type,
					}
				}

				// breath mark
				hasBreathMark = note.Notations.Articulation != nil &&
					note.Notations.Articulation.BreathMark != nil

			}

			if len(note.Beam) > 0 {

				if currTimesig.BeatType != 4 {
					for _, beam := range note.Beam {
						renderer.Beam[beam.Number] = Beam{
							Number: beam.Number,
							Type:   beam.State,
						}
					}

				}

			}

			var lyricWidth, noteWidth int

			if len(note.Lyric) > 0 {
				renderer.Lyric = Lyric{
					Text:     note.Lyric[0].Text.Value,
					Syllabic: note.Lyric[0].Syllabic,
				}
				lyricWidth = int(math.Trunc(CalculateLyricWidth(note.Lyric[0].Text.Value)))
				if note.Lyric[0].Syllabic == musicxml.LyricSyllabicTypeEnd || note.Lyric[0].Syllabic == musicxml.LyricSyllabicTypeSingle {
					lyricWidth += SPACE_LENGTH
				}
				lyricWidth = lyricWidth + 4 // lyric padding
			}

			noteWidth = LOWERCASE_LENGTH

			if noteWidth > lyricWidth {
				renderer.Width = noteWidth
				renderer.isLengthTakenFromLyric = false
			} else {
				renderer.Width = lyricWidth
				renderer.isLengthTakenFromLyric = true
				if float64(lyricWidth) > float64(noteWidth+UPPERCASE_LENGTH) {
					renderer.Width = UPPERCASE_LENGTH * 1.7
				}
			}

			notes = append(notes, renderer)

			for i, additional := range additionalRenderer {
				if i == 0 {
					continue
				}
				additionalNote := &NoteRenderer{
					PositionY:  y,
					Width:      LOWERCASE_LENGTH,
					IsDotted:   additional.IsDotted,
					NoteLength: additional.Type,
					Beam:       map[int]Beam{},
				}

				switch additional.Type {
				case musicxml.NoteLength16th:
					additionalNote.Beam[2] = Beam{
						Number: 2,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					}
					fallthrough
				case musicxml.NoteLengthEighth:
					additionalNote.Beam[1] = Beam{
						Number: 1,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					}
				}
				notes = append(notes, additionalNote)

			}
			if hasBreathMark {
				// FIXME: the breath mark stopped the continuation of the beam
				notes = append(notes, &NoteRenderer{
					Articulation: &Articulation{
						BreathMark: &ArticulationTypesBreathMark,
					},
					Width: 5,
				})
			}

		}

		// TODO: align justify
		overallWidth := 0
		for _, n := range notes {
			overallWidth += n.Width
		}
		if (x + overallWidth) > (LAYOUT_WIDTH - LAYOUT_INDENT_LENGTH) {
			y = y + 70
			x = LAYOUT_INDENT_LENGTH

		}

		s.Group("class='note'", "style='font-family:Old Standard TT;font-weight:500'")
		xNotes := 0
		continueDot := false
		lastDotLoc := 0

		var prev *NoteRenderer
		revisionX := map[int]int{}
		for i, n := range notes {
			if n.IsDotted {
				if continueDot {
					s.Text(lastDotLoc+UPPERCASE_LENGTH, y, ".")
					revisionX[i] = lastDotLoc + UPPERCASE_LENGTH
					lastDotLoc = lastDotLoc + UPPERCASE_LENGTH
				} else {
					s.Text(xNotes+UPPERCASE_LENGTH, y, ".")
					revisionX[i] = xNotes + UPPERCASE_LENGTH
					lastDotLoc = xNotes + UPPERCASE_LENGTH
				}
				continueDot = true
			} else if n.Articulation != nil && n.Articulation.BreathMark != nil {
				// breath mark
			} else {
				if continueDot {
					// FIXED:the dotted does not adding pad to the next notes
					x += LOWERCASE_LENGTH
				}
				s.Text(x, y, fmt.Sprintf("%d", n.Note))
				xNotes = x
				continueDot = false
				if n.Striketrough {
					s.Line(x+10, y-16, x, y+5, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.45")
				}
			}

			n.PositionX = x
			n.PositionY = y
			x += n.Width
			if prev != nil && prev.isLengthTakenFromLyric && n.IsDotted {
				x = x - n.Width
			}
			n.indexPosition = i
			prev = n
			// FIXED: the dotted does not give proper space at the end of measure
			if n.IsDotted && i == len(notes)-1 {
				x += LOWERCASE_LENGTH / 2
			}

		}
		s.Gend()
		// FIXME: Print it as glyph
		s.Text(x, y, " | ", "font-family:Noto Music")

		s.Group("class='lyric'", "style='font-family:Caladea'")
		for _, n := range notes {
			if n.Lyric.Text != "" {
				s.Text(n.PositionX, n.PositionY+25, n.Lyric.Text)
			}
		}
		s.Gend()

		x += LOWERCASE_LENGTH

		for i, rev := range revisionX {
			note := notes[i]

			note.PositionX = rev
			notes[i] = note
		}
		RenderOctave(rctx, s, notes)
		RenderBreath(rctx, s, notes)
		RenderSlurAndBeam(rctx, s, notes, measure.Number)
		s.Gend()

	}
}
