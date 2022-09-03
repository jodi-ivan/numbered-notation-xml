package renderer

import (
	"fmt"
	"net/http"

	svg "github.com/ajstarks/svgo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

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
	-7: "cis", //Cb
}

var minorKeySignature = map[int]string{
	7:  "ais",
	6:  "dis",
	5:  "gis",
	4:  "cis",
	3:  "fis",
	2:  "b",
	1:  "e",
	0:  "a",
	-1: "d",
	-2: "g",
	-3: "c",
	-4: "f",
	-5: "bis",
	-6: "es",
	-7: "as",
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
	beat := music.Part.Measures[0].Attribute.Time
	s.Text(LAYOUT_INDENT_LENGTH+(len(keySignature)*LOWERCASE_LENGTH), relativeY, fmt.Sprintf("%d ketuk", beat.Beats))
	s.End()
}
