package lines

import "github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"

const (
	CLEF_WIDTH    = 35
	PADDING_WIDTH = 8

	ACCIDENTAL_KEY_SIGNATURE_WIDTH = 8

	STAFF_OFFSET      = 65
	STAFF_SPACE_WIDTH = 8
)

const (
	TREBLE_CLEF_HEX = `&#xF026;`
)

var accidentalHex = map[musicxml.NoteAccidental]string{
	musicxml.NoteAccidentalNatural:     "&#xF02E;",
	musicxml.NoteAccidentalSharp:       "&#xF02B;",
	musicxml.NoteAccidentalFlat:        "&#xF02D;",
	musicxml.NoteAccidentalDoubleSharp: "&#xF02A;",
	musicxml.NoteAccidentalDoubleFlat:  "&#xF02C;",
}
