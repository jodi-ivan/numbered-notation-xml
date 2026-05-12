package gregorian

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

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

var renderStemAndBeamMap = map[int]func(canv canvas.Canvas, lines [5]int, pos ...entity.CoordinateWithNoteLength) StemInfo{
	-1: RenderStemDown,
	0:  RenderStemDown,
	1:  RenderStemUp,
}

var restHex = map[musicxml.NoteLength]string{
	musicxml.NoteLengthQuarter: "&#xF074;",
	musicxml.NoteLengthEighth:  "&#xF075;",
	musicxml.NoteLength16th:    "&#xF076;",
}

var beanNoteHex = map[musicxml.NoteLength]string{
	musicxml.NoteLength16th:    `&#xF064;`,
	musicxml.NoteLengthEighth:  `&#xF064;`,
	musicxml.NoteLengthQuarter: `&#xF064;`,
	musicxml.NoteLengthHalf:    `&#xF063;`,
	musicxml.NoteLengthWhole:   `&#xF062;`,
}

var singleFlagHex = map[int]map[musicxml.NoteLength]string{
	-1: {
		musicxml.NoteLengthEighth: "&#xF06D;",
		musicxml.NoteLength16th:   "&#xF06E;",
	},
	0: {
		musicxml.NoteLengthEighth: "&#xF06D;",
		musicxml.NoteLength16th:   "&#xF06E;",
	},
	1: {
		musicxml.NoteLengthEighth: "&#xF069;",
		musicxml.NoteLength16th:   "&#xF06A;",
	},
}

var stemPosOffset = map[int][2]entity.Coordinate{
	-1: {
		entity.NewCoordinate(0.5, 2),
		entity.NewCoordinate(0.5, 28),
	},
	0: {
		entity.NewCoordinate(0.5, 2),
		entity.NewCoordinate(0.5, 28),
	},
	1: {
		entity.NewCoordinate(9, 0),
		entity.NewCoordinate(9, -24),
	},
}

var accidentalHex = map[musicxml.NoteAccidental]string{
	musicxml.NoteAccidentalNatural:     "&#xF02E;",
	musicxml.NoteAccidentalSharp:       "&#xF02B;",
	musicxml.NoteAccidentalFlat:        "&#xF02D;",
	musicxml.NoteAccidentalDoubleSharp: "&#xF02A;",
	musicxml.NoteAccidentalDoubleFlat:  "&#xF02C;",
}

var accidentalHexWithParentheses = map[musicxml.NoteAccidental]string{
	musicxml.NoteAccidentalNatural:     "&#xF0B2;",
	musicxml.NoteAccidentalSharp:       "&#xF0B1;",
	musicxml.NoteAccidentalFlat:        "&#xF0B3;",
	musicxml.NoteAccidentalDoubleSharp: "&#xF0B0;",
	musicxml.NoteAccidentalDoubleFlat:  "&#xF0B4;",
}
