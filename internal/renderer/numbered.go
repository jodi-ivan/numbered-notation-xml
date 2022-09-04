package renderer

import (
	"fmt"
	"net/http"

	svg "github.com/ajstarks/svgo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

var majorScale = map[string][]string{
	"cis": []string{"C#", "D#", "E#", "F#", "G#", "A#", "B#"},
	"fis": []string{"F#", "G#", "A#", "B", "C#", "D#", "E#"},
	"b":   []string{"B", "C#", "D#", "E", "F#", "G#", "A#"},
	"e":   []string{"E", "F#", "G#", "A", "B", "C#", "D#"},
	"a":   []string{"A", "B", "C#", "D", "E", "F#", "G#"},
	"d":   []string{"D", "E", "F#", "G", "A", "B", "C#"},
	"g":   []string{"G", "A", "B", "C", "D", "E", "F#"},
	"c":   []string{"C", "D", "E", "F", "G", "A", "B"},
	"f":   []string{"F", "G", "A", "Bb", "C", "D", "E"},
	"bes": []string{"Bb", "C", "D", "Eb", "F", "G", "A"},
	"es":  []string{"Eb", "F", "G", "Ab", "Bb", "C", "D"},
	"as":  []string{"Ab", "Bb", "C", "Db", "Eb", "F", "G"},
	"des": []string{"Db", "Eb", "F", "Gb", "Ab", "Bb", "C"},
	"ges": []string{"Gb", "Ab", "Bb", "Cb", "Db", "Eb", "F"},
	"ces": []string{"Gb", "Ab", "Bb", "Cb", "Db", "Eb", "F"},
}

var minorScale = map[string][]string{
	"ais": []string{"A#", "B#", "C#", "D#", "E#", "F#", "G#"},
	"dis": []string{},
	"gis": []string{},
	"cis": []string{},
	"fis": []string{},
	"b":   []string{},
	"e":   []string{},
	"a":   []string{},
	"d":   []string{},
	"g":   []string{},
	"c":   []string{},
	"f":   []string{},
	"bis": []string{},
	"es":  []string{},
	"as":  []string{},
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
	-5: "bis", //F#
	-6: "es",  //Eb
	-7: "as",  // Ab
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
	keyMode := music.Part.Measures[0].Attribute.Key.Mode
	fifths := music.Part.Measures[0].Attribute.Key.Fifth

	keySignature := ""
	if keyMode == "minor" {
		keySignature = fmt.Sprintf("la = %s", minorKeySignature[fifths])
	} else if keyMode == "major" {
		keySignature = fmt.Sprintf("do = %s", majorKeySignature[fifths])
	}

	s.Text(LAYOUT_INDENT_LENGTH, relativeY, keySignature)

	// render time signature
	// TODO:
	// - check the time signature on github issue
	beat := music.Part.Measures[0].Attribute.Time
	s.Text(LAYOUT_INDENT_LENGTH+(len(keySignature)*LOWERCASE_LENGTH), relativeY, fmt.Sprintf("%d ketuk", beat.Beats))

	s.End()
}

func RenderMeasure(s *svg.SVG, measure musicxml.Measure, attributes musicxml.Attribute) {

	// map the key signature

}
