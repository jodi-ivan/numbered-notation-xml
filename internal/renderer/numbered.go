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
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/moveabledo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

const LOWERCASE_LENGTH = 15
const UPPERCASE_LENGTH = 20
const SPACE_LENGTH = 7
const LAYOUT_INDENT_LENGTH = 50
const LAYOUT_WIDTH = 1000

const (
	gwfURI  = "https://fonts.googleapis.com/css?family="
	fontfmt = "<style type=\"text/css\">\n<![CDATA[\n%s]]>\n</style>\n"
	gfmt    = "fill:white;font-size:36pt;text-anchor:middle"
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

/*                Numbered stuff                              */

type Lyric struct {
	Text     string
	Syllabic musicxml.LyricSyllabic
}

type Slur struct {
	// Number attributes for slur
	// Pitch note for ties
	Number int
	Type   musicxml.NoteSlurType
}

type Beam struct {
	Number int
	Type   musicxml.NoteBeamType
}

type NoteRenderer struct {
	IsDotted     bool
	IsRest       bool
	PositionX    int
	PositionY    int
	Note         int
	Octave       int
	Striketrough bool
	NoteLength   musicxml.NoteLength
	BarType      string
	Width        int
	Lyric        Lyric
	Slur         map[int]Slur
	Beam         map[int]Beam
	Tie          *Slur

	// internal use
	isLengthTakenFromLyric bool
}

type Coordinate struct {
	X float64
	Y float64
}

type CoordinateWithOctave struct {
	Coordinate
	Octave int
}

// for (svg.SVG).Qbez
type SlurBezier struct {
	Start CoordinateWithOctave
	End   CoordinateWithOctave
	Pull  CoordinateWithOctave
}

type BeamLine struct {
	Start Coordinate
	End   Coordinate
}

func cleanBeamByNumber(ctx context.Context, notes []*NoteRenderer, beamNumber int) []*NoteRenderer {

	switches := map[int]musicxml.NoteBeamType{}

	var prev *NoteRenderer

	for indexNote, note := range notes {

		if indexNote == len(notes)-1 { // ignore last note or no beam
			prev = note
			continue
		}

		if len(note.Beam) == 0 { // stopping the beam
			if indexNote == 0 {
				prev = note
				continue
			} else {

				if _, ok := switches[beamNumber]; !ok {
					prev = note
					continue
				}

				prev.Beam[beamNumber] = Beam{
					Number: beamNumber,
					Type:   musicxml.NoteBeamTypeEnd,
				}

				delete(switches, beamNumber)
			}
		}

		if t, ok := switches[beamNumber]; !ok {

			if _, hasBeam := note.Beam[beamNumber]; !hasBeam {
				prev = note
				continue
			}
			newBeam := map[int]Beam{}

			for k, v := range note.Beam {
				newBeam[k] = v
			}

			switches[beamNumber] = musicxml.NoteBeamTypeBegin
			newBeam[beamNumber] = Beam{
				Number: beamNumber,
				Type:   musicxml.NoteBeamTypeBegin,
			}
			note.Beam = newBeam
		} else {

			if prev == nil {
				prev = note
				continue
			}

			if _, hasBeam := note.Beam[beamNumber]; hasBeam {
				newBeam := map[int]Beam{}

				for k, v := range note.Beam {
					newBeam[k] = v
				}

				switches[beamNumber] = musicxml.NoteBeamTypeBegin
				newBeam[beamNumber] = Beam{
					Number: beamNumber,
					Type:   musicxml.NoteBeamTypeContinue,
				}
				note.Beam = newBeam
				prev = note
				continue
			}

			if t == musicxml.NoteBeamTypeBegin {
				if _, ok := prev.Beam[beamNumber]; !ok {
					prev = note
					continue
				}

				prev.Beam[beamNumber] = Beam{
					Number: beamNumber,
					Type:   musicxml.NoteBeamTypeEnd,
				}

				delete(switches, beamNumber)
			}

		}
		prev = note

	}

	if len(prev.Beam) > 0 {
		additional, ok := prev.Beam[beamNumber]
		if ok {
			if additional.Type != musicxml.NoteBeamTypeEnd {
				newBeam := prev.Beam

				newBeam[beamNumber] = Beam{
					Type:   musicxml.NoteBeamTypeEnd,
					Number: beamNumber,
				}

				prev.Beam = newBeam
			} else {
				if _, ok := switches[beamNumber]; !ok {
					newBeam := prev.Beam
					newBeam[beamNumber] = Beam{
						Type:   musicxml.NoteBeamTypeBackwardHook,
						Number: beamNumber,
					}
					prev.Beam = newBeam

				}
			}
		}
	}
	return notes
}

func RenderSlurAndBeam(ctx context.Context, canvas *svg.SVG, notes []*NoteRenderer, measureNo int) {
	slurs := map[int]SlurBezier{}
	slurSets := []SlurBezier{}

	// [ ] 6/8 time signature
	beams := map[int]BeamLine{}
	beamSets := []BeamLine{}

	// FIXED: support for multi-octave
	// currently it support multi-ties based on the pitch
	// since there is no indicator for what octave it could colliding with each other
	ties := map[int]SlurBezier{}
	tiesSet := []SlurBezier{}
	cleanedNote := cleanBeamByNumber(ctx, notes, 1)
	cleanedNote = cleanBeamByNumber(ctx, cleanedNote, 2)
	for _, note := range cleanedNote {

		for _, s := range note.Slur {
			if s.Type == musicxml.NoteSlurTypeStart {
				slurs[s.Number] = SlurBezier{
					Start: CoordinateWithOctave{
						Coordinate: Coordinate{
							X: float64(note.PositionX),
							Y: float64(note.PositionY),
						},
						Octave: note.Octave,
					},
				}
			} else if s.Type == musicxml.NoteSlurTypeStop {
				temp := slurs[s.Number]
				temp.End = CoordinateWithOctave{
					Coordinate: Coordinate{
						X: float64(note.PositionX),
						Y: float64(note.PositionY),
					},
					Octave: note.Octave,
				}
				slurs[s.Number] = temp

				slurSets = append(slurSets, slurs[s.Number])
				delete(slurs, s.Number)
			}
		}

		// TODO: team beam only support 2 notes grouping
		// TODO add support for backward hook and forward hook
		// TODO: add support for signular note notebeam type
		for _, b := range note.Beam {
			positionY := float64(note.PositionY - 20 + ((b.Number) * 3))

			switch b.Type {
			case musicxml.NoteBeamTypeBegin:
				beams[b.Number] = BeamLine{
					Start: Coordinate{
						X: float64(note.PositionX),
						Y: positionY,
					},
				}
			case musicxml.NoteBeamTypeEnd:

				beam := beams[b.Number]

				if beam.Start.X == 0 {
					beams[b.Number] = BeamLine{
						Start: Coordinate{
							X: float64(note.PositionX),
							Y: positionY,
						},
						End: Coordinate{
							X: float64(note.PositionX) + 8,
							Y: positionY,
						},
					}

				} else {
					beam.End = Coordinate{
						X: float64(note.PositionX) + 8,
						Y: beam.Start.Y,
					}
					beams[b.Number] = beam
				}

				beamSets = append(beamSets, beams[b.Number])

				delete(beams, b.Number)

			}
		}

		if note.Tie != nil {
			if note.Tie.Type == musicxml.NoteSlurTypeStart {
				ties[note.Note] = SlurBezier{
					Start: CoordinateWithOctave{
						Coordinate: Coordinate{
							X: float64(note.PositionX),
							Y: float64(note.PositionY),
						},
						Octave: note.Octave,
					},
				}
			} else if note.Tie.Type == musicxml.NoteSlurTypeStop {
				temp := ties[note.Note]
				temp.End = CoordinateWithOctave{
					Coordinate: Coordinate{
						X: float64(note.PositionX),
						Y: float64(note.PositionY),
					},
					Octave: note.Octave,
				}
				ties[note.Note] = temp

				tiesSet = append(tiesSet, ties[note.Note])
				delete(slurs, note.Note)
			}
		}

	}

	RenderBezier(slurSets, canvas)
	RenderBezier(tiesSet, canvas)

	canvas.Group("class='beam'")
	for _, b := range beamSets {
		canvas.Line(
			int(math.Round(b.Start.X)),
			int(math.Round(b.Start.Y)),
			int(math.Round(b.End.X)),
			int(math.Round(b.End.Y)),
			"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.2",
		)
	}
	canvas.Gend()
}

func RenderOctave(canvas *svg.SVG, notes []*NoteRenderer) {
	canvas.Group("class='octaves'")
	for _, note := range notes {
		if note.Octave < 0 {
			canvas.Circle(note.PositionX+5, note.PositionY+5, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
		}

		if note.Octave > 0 {
			canvas.Circle(note.PositionX+5, note.PositionY-15, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
		}
	}
	canvas.Gend()
}

func CalculateLyricWidth(txt string) float64 {
	// TODO: margin left of the lyric
	// TODO: continues syllable
	width := map[string]float64{
		"1": 9.28,
		"A": 9.59,
		"B": 9.27,
		"C": 8.1,
		"D": 10,
		"E": 8.65,
		"F": 8.15,
		"G": 8.63,
		"H": 11.15,
		"I": 5.49,
		"J": 4.99,
		"K": 10.08,
		"L": 8.02,
		"M": 14.21,
		"N": 11.09,
		"O": 9.59,
		"P": 8.53,
		"Q": 9.59,
		"R": 9.81,
		"S": 7.25,
		"T": 8.92,
		"U": 11,
		"V": 9.57,
		"W": 14.23,
		"X": 9.95,
		"Y": 8.92,
		"Z": 8.11,
		"a": 7.52,
		"b": 8.32,
		"c": 6.74,
		"d": 8.32,
		"e": 7.06,
		"f": 5.87,
		"g": 7.35,
		"h": 8.86,
		"i": 4.44,
		"j": 4.76,
		"k": 8.43,
		"l": 4.34,
		"m": 13.01,
		"n": 8.94,
		"o": 7.69,
		"p": 8.32,
		"q": 8.02,
		"r": 6.34,
		"s": 6.28,
		"t": 5.21,
		"u": 8.74,
		"v": 8.08,
		"w": 12.08,
		"x": 7.78,
		"y": 8.18,
		"z": 6.85,
		",": 3.28,
		"'": 3.28,
		".": 3.3,
		"!": 4.58,
		";": 4.23,
	}
	res := 0.0

	for _, l := range txt {
		res += width[string(l)]
	}

	return res
}

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

			noteLength := timeSignature.GetNoteLength(ctx, measure.Number, note)
			additionalRenderer := numbered.RenderLengthNote(ctx, timeSignature, measure.Number, noteLength)

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

			if note.Notations != nil {
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

			}

			if len(note.Beam) > 0 {
				currTimesig := timeSignature.GetTimesignatureOnMeasure(ctx, measure.Number)
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
				// FIXED: next position notes is wrong (there like more space than it should) when the width of lyric > note
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

		}

		// TODO: align justify
		overallWidth := 0
		for _, n := range notes {
			//FIXED [Maybe] width is unreliable, especially when the width is the lyric. calculation seems correct. but next note placement is wrong. makes the
			//	page calculation incorrect
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
			} else {
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
			prev = n

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
		RenderOctave(s, notes)
		RenderSlurAndBeam(ctx, s, notes, measure.Number)
		s.Gend()

	}
}

func RenderBezier(set []SlurBezier, canvas *svg.SVG) {
	canvas.Group("class='slurties'")
	// TODO: check ties across measure bar
	for _, s := range set {

		slurResult := SlurBezier{
			Start: CoordinateWithOctave{
				Coordinate: Coordinate{
					X: s.Start.X + 5,
					Y: s.Start.Y + 5,
				},
				Octave: s.Start.Octave,
			},
			End: CoordinateWithOctave{
				Coordinate: Coordinate{
					X: s.End.X + 5,
					Y: s.End.Y + 5,
				},
				Octave: s.End.Octave,
			},
		}

		if slurResult.Start.Octave < 0 {
			slurResult.Start = CoordinateWithOctave{
				Coordinate: Coordinate{
					X: slurResult.Start.X + 3,
					Y: slurResult.Start.Y + 3,
				},
			}
		}

		if slurResult.End.Octave < 0 {

			slurResult.End = CoordinateWithOctave{
				Coordinate: Coordinate{
					X: slurResult.End.X - 3,
					Y: slurResult.End.Y + 3,
				},
			}
		}

		pull := CoordinateWithOctave{
			Coordinate: Coordinate{
				X: slurResult.Start.X + ((slurResult.End.X - slurResult.Start.X) / 2),
				Y: slurResult.Start.Y + 7.5,
			},
		}
		slurResult.Pull = pull

		canvas.Qbez(
			int(math.Round(slurResult.Start.X)),
			int(math.Round(slurResult.Start.Y)),
			int(math.Round(pull.X)),
			int(math.Round(pull.Y)),
			int(math.Round(slurResult.End.X)),
			int(math.Round(slurResult.End.Y)),
			"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.5",
		)
	}
	canvas.Gend()
}
