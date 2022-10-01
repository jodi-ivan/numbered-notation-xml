package renderer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

var majorAccidental = map[int][]string{
	7:  []string{"C", "D", "E", "F", "G", "A", "B"},
	6:  []string{"F", "G", "A", "C", "D", "E"},
	5:  []string{"C", "D", "F", "G", "A"},
	4:  []string{"F", "G", "C", "D"},
	3:  []string{"C", "F", "G"},
	2:  []string{"F", "C"},
	1:  []string{"F"},
	0:  []string{},
	-1: []string{"B"},
	-2: []string{"B", "E"},
	-3: []string{"E", "A", "B"},
	-4: []string{"A", "B", "D", "E"},
	-5: []string{"D", "E", "G", "A", "B"},
	-6: []string{"G", "A", "B", "C", "D", "E"},
	-7: []string{"G", "A", "B", "C", "D", "E"},
}

var minorAccidental = map[int][]string{
	7:  []string{"A", "B", "C", "D", "E", "F", "G"},
	6:  []string{"D", "E", "F", "G", "A", "C"},
	5:  []string{"G", "A", "C", "D", "F"},
	4:  []string{"C", "D", "F", "G"},
	3:  []string{"F", "G", "C"},
	2:  []string{"C", "F", "G", "A"},
	1:  []string{"F#"},
	0:  []string{},
	-1: []string{"B"},
	-2: []string{"B", "E"},
	-3: []string{"E", "A", "B"},
	-4: []string{"A", "B", "D", "E"},
	-5: []string{"B", "D", "E", "G", "A"},
	-6: []string{"E", "G", "A", "B", "C", "D"},
	-7: []string{"A", "B", "C", "D", "E", "F", "G"},
}

const LOWERCASE_LENGTH = 15
const UPPERCASE_LENGTH = 20
const LAYOUT_INDENT_LENGTH = 50
const LAYOUT_WIDTH = 1000

var majorKeySignature = map[int]string{
	7:  "cis", // C#
	6:  "fis", // F#
	5:  "b",
	4:  "e",
	3:  "a",
	2:  "d",
	1:  "g",
	0:  "c",
	-1: "f",
	-2: "bes", // Bb
	-3: "es",  // Eb
	-4: "as",  // Ab
	-5: "des", // Db
	-6: "ges", //Gb
	-7: "ces", //Cb
}

var minorKeySignature = map[int]string{
	7:  "ais", // A#
	6:  "dis", // D#
	5:  "gis", // G#
	4:  "cis", // C#
	3:  "fis", // F#
	2:  "b",
	1:  "e",
	0:  "a",
	-1: "d",
	-2: "g",
	-3: "c",
	-4: "f",
	-5: "bes", //Bb
	-6: "es",  //Eb
	-7: "as",  // Ab
}

var majorLetteredKeySignature = map[int]string{
	7:  "C#",
	6:  "F#",
	5:  "B",
	4:  "E",
	3:  "A",
	2:  "D",
	1:  "G",
	0:  "C",
	-1: "F",
	-2: "Bb",
	-3: "Eb",
	-4: "Ab",
	-5: "Db",
	-6: "Gb",
	-7: "Cb",
}

var minorLetteredKeySignature = map[int]string{
	7:  "A#",
	6:  "D#",
	5:  "G#",
	4:  "C#",
	3:  "F#",
	2:  "B",
	1:  "E",
	0:  "A",
	-1: "D",
	-2: "G",
	-3: "C",
	-4: "F",
	-5: "Bb",
	-6: "Eb",
	-7: "Ab",
}

type KeySignatureMode int

const (
	KeySignatureModeMajor KeySignatureMode = 0
	KeySignatureModeMinor KeySignatureMode = 1
)

func (ksm KeySignatureMode) String() string {
	return []string{"major", "minor"}[int(ksm)]
}

type KeySignature struct {
	Key       string
	Mode      KeySignatureMode
	Humanized string
	Fifth     int
}

func contains(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}

	return -1
}

func NewKeySignature(key musicxml.KeySignature) KeySignature {
	keyMode := key.Mode
	fifths := key.Fifth

	var letterKey string

	modeMapper := map[string]KeySignatureMode{
		"major": KeySignatureModeMajor,
		"minor": KeySignatureModeMinor,
	}

	keySignature := ""
	if keyMode == "minor" {
		keySignature = fmt.Sprintf("la = %s", minorKeySignature[fifths])
		letterKey = minorKeySignature[fifths]
	} else if keyMode == "major" {
		keySignature = fmt.Sprintf("do = %s", majorKeySignature[fifths])
		letterKey = majorKeySignature[fifths]
	}

	return KeySignature{
		Fifth:     fifths,
		Key:       letterKey,
		Mode:      modeMapper[keyMode],
		Humanized: keySignature,
	}
}

func (ks *KeySignature) String() string {
	return ks.Humanized
}

func (ks KeySignature) GetNumberedNotation(note musicxml.Note) (numberedNote int, octave int, strikethrough bool) {
	octave = note.Pitch.Octave - 4

	pitch := note.Pitch.Step
	var accidental musicxml.NoteAccidental
	var accidentals []string

	// get the key signature
	if ks.Mode == KeySignatureModeMajor {
		accidentals = majorAccidental[ks.Fifth]
	} else if ks.Mode == KeySignatureModeMinor {
		accidentals = minorAccidental[ks.Fifth]
	}

	if contains(accidentals, pitch) >= 0 {
		if ks.Fifth > 0 {
			accidental = musicxml.NoteAccidentalSharp
		} else if ks.Fifth < 0 {
			accidental = musicxml.NoteAccidentalFlat
		}
	}

	if note.Accidental != "" {
		accidental = note.Accidental
	}

	pitch = fmt.Sprintf("%s%s", pitch, accidental.GetAccidental())

	numberedNote, strikethrough = ConvertPitchToNumbered(ks, pitch)
	return numberedNote, octave, strikethrough

}

func getNextHalfStep(pitch string) string {
	step := []string{"C", "D", "E", "F", "G", "A", "B"}

	wholeStep := func(p string) string {
		index := contains(step, p)
		if index == len(step)-1 {
			return step[0]
		}
		return step[index+1]
	}

	if strings.HasSuffix(pitch, "bb") { // double flat
		return fmt.Sprintf("%sb", string(pitch[0]))
	}

	if strings.HasSuffix(pitch, "b") { // flat, just remove the flat
		return string(pitch[0])
	}

	if strings.HasSuffix(pitch, "x") { // double sharp
		if string(pitch[0]) == "B" {
			return getNextHalfStep("C#")
		}

		if string(pitch[0]) == "E" {
			return getNextHalfStep("F#")
		}

		newPitch := wholeStep(string(pitch[0]))
		if newPitch == "B" || newPitch == "E" {
			return wholeStep(newPitch)
		}

		return fmt.Sprintf("%s#", newPitch)
	}

	if strings.HasSuffix(pitch, "#") {
		if string(pitch[0]) == "B" || string(pitch[0]) == "E" {
			return fmt.Sprintf("%s#", wholeStep(string(pitch[0])))
		}

		return wholeStep(string(pitch[0]))
	}

	if string(pitch[0]) == "B" || string(pitch[0]) == "E" {
		return wholeStep(string(pitch[0]))
	}

	return fmt.Sprintf("%s#", string(pitch[0]))

}

func ConvertPitchToNumbered(ks KeySignature, pitch string) (numbered int, strike bool) {

	//                                |     |   |   |   |     |     |   |   |   |   |   |    |     |   |   |   |
	major := []float64{ //            |     |   |   |   |     |     |   |   |   |   |   |    |     |   |   |   |
		1,   // do -> re              |     |   |   |   |     |     |   |   |   |   |   |    |     |   |   |   |
		1,   // re -> mi              |     |   |   |   |     |     |   |   |   |   |   |    |     |   |   |   |
		0.5, // mi -> fa              |     +---+   +---+     |     +---+   +---+   +---+    |     +---+   +---+
		1,   // fa -> sol             |       |       |       |       |       |       |      |       |       |
		1,   // sol -> la             |       |       |       |       |       |       |      |   *   |   *   |
		1,   // la -> si (ti)         |   1   |   2   |   3   |   4   |   5   |   6   |   7  |   1   |   2   |
		0.5, // si -> do              |       |       |       |       |       |       |      |       |       |
	} //                              +-------+-------+-------+-------+-------+-------+------+-------+-------+
	//                C major scale       C       D       E       F       G      A        B      C       D

	// minor step
	//                                  |   |   |   |     |     |   |   |   |     |    |   |   |   |   |   |
	minor := []float64{ //              |   |   |   |     |     |   |   |   |     |    |   |   |   |   |   |
		1,   // do -> re                |   |   |   |     |     |   |   |   |     |    |   |   |   |   |   |
		0.5, // re -> mi                |   |   |   |     |     |   |   |   |     |    |   |   |   |   |   |
		1,   // mi -> fa                +---+   +---+     |     +---+   +---+     |    +---+   +---+   +---+
		1,   // fa -> sol                 |       |       |       |       |       |      |       |       |
		0.5, // sol -> la                 |       |       |       |       |       |      |       |   *   |
		1,   // la -> si (ti)             |   1   |   2   |   3   |   4   |   5   |   6  |   7   |   1   |
		1,   // si -> do                  |       |       |       |       |       |      |       |       |
	} //                                --+-------+-------+-------+-------+-------+------+-------+-------+---
	//                A minor scale           A       B       C       D       E      F        G      A

	var current string
	var steps []float64

	if ks.Mode == KeySignatureModeMajor {
		current = majorLetteredKeySignature[ks.Fifth]
		steps = major
	} else if ks.Mode == KeySignatureModeMinor {
		current = minorLetteredKeySignature[ks.Fifth]
		steps = minor
	}

	counter := 0

	for !(isPitchEqual(current, pitch)) {
		current = getNextHalfStep(current)
		counter++
	}

	stepped := 0.0
	increase := 0

	for stepped < float64(counter)/2 {

		stepped += steps[increase]
		increase++
	}

	if stepped > (float64(counter) / 2) {
		return increase, true
	}
	return increase + 1, false

}

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

func RenderNumbered(w http.ResponseWriter, music musicxml.MusicXML) {
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

	keySignature := NewKeySignature(music.Part.Measures[0].Attribute.Key)

	humanizedKeySignature := keySignature.String()

	s.Text(LAYOUT_INDENT_LENGTH, relativeY, keySignature.String())

	// render time signature
	// TODO:
	// - check the time signature on github issue
	// - time signature changing happens on the top and not on the measure
	beat := music.Part.Measures[0].Attribute.Time
	s.Text(LAYOUT_INDENT_LENGTH+(len(humanizedKeySignature)*LOWERCASE_LENGTH), relativeY, fmt.Sprintf("%d ketuk", beat.Beats))
	relativeY += 50

	RenderMeasures(s, LAYOUT_INDENT_LENGTH, relativeY, music.Part)
	s.End()
}

func isPitchEqual(one, two string) bool {
	if one == two {
		return true
	}
	pitches := map[string][]string{
		"C":   []string{"B#", "Dbb"},
		"Cb":  []string{"B", "Ax"},
		"C#":  []string{"Db", "Bx"},
		"Cx":  []string{"B", "Ax"},
		"Cbb": []string{"B", "Ax"},
		"D":   []string{"Cx", "Ebb"},
		"Db":  []string{"C#"},
		"D#":  []string{"Eb", "Fbb"},
		"Dbb": []string{"C", "B#"},
		"Dx":  []string{"E", "Fb"},
		"E":   []string{"Dx", "Fb"},
		"Eb":  []string{"D#", "Fbb"},
		"E#":  []string{"F", "Gbb"},
		"Ebb": []string{"D", "Cx"},
		"Ex":  []string{"F#", "Gb"},
		"F":   []string{"E#", "Gbb"},
		"Fb":  []string{"E", "Dx"},
		"F#":  []string{"Gb", "Ex"},
		"Fbb": []string{"Eb", "D#"},
		"Fx":  []string{"G", "Abb"},
		"G":   []string{"Abb", "Fx"},
		"Gb":  []string{"F#", "Ex"},
		"G#":  []string{"Ab"},
		"Gbb": []string{"F", "E#"},
		"Gx":  []string{"A", "Bbb"},
		"A":   []string{"Gx", "Bbb"},
		"Ab":  []string{"G#"},
		"A#":  []string{"Bb", "Cbb"},
		"Abb": []string{"G", "Fx"},
		"Ax":  []string{"B", "Cb"},
		"B":   []string{"Cb", "Ax"},
		"Bb":  []string{"A#", "Cbb"},
		"B#":  []string{"C", "Dbb"},
		"Bbb": []string{"A", "Gx"},
		"Bx":  []string{"C#", "Db"},
	}
	result := contains(pitches[one], two) >= 0

	return result
}

/*                Numbered stuff                              */

type Lyric struct {
	Text     string
	Syllabic musicxml.LyricSyllabic
}

type Slur struct {
	Number int
	Type   musicxml.NoteSlurType
}
type NoteRenderer struct {
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
}

type CoordinateWithOctave struct {
	X      float32
	Y      float32
	Octave int
}

// for (svg.SVG).Qbez
type SlurBezier struct {
	Start CoordinateWithOctave
	End   CoordinateWithOctave
	Pull  CoordinateWithOctave
}

func RenderSlur(canvas *svg.SVG, notes []*NoteRenderer) {
	slurs := map[int]SlurBezier{}
	sets := []SlurBezier{}

	for _, note := range notes {
		for _, s := range note.Slur {
			if s.Type == musicxml.NoteSlurTypeStart {
				slurs[s.Number] = SlurBezier{
					Start: CoordinateWithOctave{
						X:      float32(note.PositionX),
						Y:      float32(note.PositionY),
						Octave: note.Octave,
					},
				}
			} else if s.Type == musicxml.NoteSlurTypeStop {
				temp := slurs[s.Number]
				temp.End = CoordinateWithOctave{
					X:      float32(note.PositionX),
					Y:      float32(note.PositionY),
					Octave: note.Octave,
				}
				slurs[s.Number] = temp

				sets = append(sets, slurs[s.Number])
				delete(slurs, s.Number)
			}
		}
	}

	for _, s := range sets {
		pull := CoordinateWithOctave{
			X: s.Start.X + ((s.End.X - s.Start.X) / 2) + 5,
			Y: s.Start.Y + 15,
		}

		slurResult := SlurBezier{
			Start: CoordinateWithOctave{
				X:      s.Start.X + 5,
				Y:      s.Start.Y + 5,
				Octave: s.Start.Octave,
			},
			End: CoordinateWithOctave{
				X:      s.End.X + 5,
				Y:      s.End.Y + 5,
				Octave: s.End.Octave,
			},
			Pull: pull,
		}

		if slurResult.Start.Octave < 0 {
			slurResult.Start = CoordinateWithOctave{
				X: slurResult.Start.X + 5,
				Y: slurResult.Start.Y + 5,
			}
		}

		if slurResult.End.Octave < 0 {
			slurResult.End = CoordinateWithOctave{
				X: slurResult.End.X - 5,
				Y: slurResult.End.Y - 5,
			}
		}

		canvas.Qbez(
			int(slurResult.Start.X),
			int(slurResult.Start.Y),
			int(pull.X),
			int(pull.Y),
			int(slurResult.End.X),
			int(slurResult.End.Y),
			"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.5",
		)
	}

}

func RenderOctave(canvas *svg.SVG, notes []*NoteRenderer) {
	for _, note := range notes {
		if note.Octave < 0 {
			canvas.Circle(note.PositionX+5, note.PositionY+20, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:1.3")
		}

		if note.Octave > 0 {
			canvas.Circle(note.PositionX+5, note.PositionY-15, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:1.3")
		}
	}
}

func CalculateLyricWidth(txt string) float64 {
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

func RenderMeasures(s *svg.SVG, x, y int, measures musicxml.Part) {

	keySignature := NewKeySignature(measures.Measures[0].Attribute.Key)

	var locationX int
	locationX = x

	for _, measure := range measures.Measures {

		notes := []*NoteRenderer{}
		for _, note := range measure.Notes {

			n, octave, strikethrough := keySignature.GetNumberedNotation(note)

			renderer := &NoteRenderer{
				PositionX:    x,
				PositionY:    y,
				Note:         n,
				NoteLength:   note.Type,
				Octave:       octave,
				Striketrough: strikethrough,
			}

			if note.Notations != nil && len(note.Notations.Slur) > 0 {
				renderer.Slur = map[int]Slur{}
				for _, slur := range note.Notations.Slur {
					renderer.Slur[slur.Number] = Slur{
						Number: slur.Number,
						Type:   slur.Type,
					}
				}
			}

			var lyricWidth, noteWidth int

			if len(note.Lyric) > 0 {
				renderer.Lyric = Lyric{
					Text:     note.Lyric[0].Text.Value,
					Syllabic: note.Lyric[0].Syllabic,
				}

				lyricWidth = int(math.Ceil(CalculateLyricWidth(note.Lyric[0].Text.Value)))
				if note.Lyric[0].Syllabic == musicxml.LyricSyllabicTypeEnd || note.Lyric[0].Syllabic == musicxml.LyricSyllabicTypeSingle {
					lyricWidth += 7 // space
				}
			}

			noteWidth = LOWERCASE_LENGTH
			switch note.Type {
			case musicxml.NoteLengthWhole:

				// whole note in musical number notation will add 3 dots in front of the note
				// C whole note will represent as
				//      1 . . . |
				noteWidth = 3 * LOWERCASE_LENGTH * 3

			case musicxml.NoteLengthHalf:
				// half  note in musical number notation will add 1 dots in front of the note
				// C half note will represent as
				//      1 . |
				noteWidth = LOWERCASE_LENGTH * 2
			}

			if noteWidth > lyricWidth {
				x += noteWidth
				renderer.Width = noteWidth
			} else {
				x += lyricWidth
				renderer.Width = lyricWidth

			}

			raw, _ := json.Marshal(renderer)
			log.Println(string(raw))

			notes = append(notes, renderer)
		}

		//    rough calculation of measure
		if x+((len(notes)+1)*LOWERCASE_LENGTH) > (LAYOUT_WIDTH - LAYOUT_INDENT_LENGTH) {
			y = y + 70
			locationX = LAYOUT_INDENT_LENGTH
			x = LAYOUT_INDENT_LENGTH
		}

		s.Gstyle("font-family:Old Standard TT;font-weight:600")
		for _, n := range notes {

			s.Text(locationX, y, fmt.Sprintf("%d", n.Note))

			n.PositionX = locationX
			n.PositionY = y
			locationX += n.Width

		}
		s.Gend()
		s.Text(locationX, y, " | ", "font-family:Noto Music")

		s.Gstyle("font-family:Caladea")
		for _, n := range notes {
			if n.Lyric.Text != "" {
				s.Text(n.PositionX, n.PositionY+25, n.Lyric.Text)
			}
		}
		s.Gend()
		locationX += LOWERCASE_LENGTH
		RenderOctave(s, notes)
		RenderSlur(s, notes)

	}

}
