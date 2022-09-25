package renderer

import (
	"fmt"
	"net/http"
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
const LAYOUT_INDENT_LENGTH = 50

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
	} else {
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
	step := []string{"A", "B", "C", "D", "E", "F", "G"}

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

func RenderNumbered(w http.ResponseWriter, music musicxml.MusicXML) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(200)
	s := svg.New(w)
	s.Start(1000, 1000)

	relativeY := 100
	// render title
	titleX := (1000 / 2) - ((len(music.Credit.Words) * LOWERCASE_LENGTH) / 2) - (LAYOUT_INDENT_LENGTH * 2)
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
	relativeY += 30
	for _, measure := range music.Part.Measures {
		txt := ""
		for _, note := range measure.Notes {
			n, _, _ := keySignature.GetNumberedNotation(note)
			txt = fmt.Sprintf("%s %d", txt, n)
		}
		s.Text(LAYOUT_INDENT_LENGTH, relativeY, txt)
		relativeY += 25
	}

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

func RenderMeasure(s *svg.SVG, measure musicxml.Measure, attributes musicxml.Attribute) {

	// map the key signature

}
