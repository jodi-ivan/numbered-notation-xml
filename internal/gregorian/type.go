package gregorian

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type CoordinateWithNoteLength struct {
	entity.Coordinate
	NoteLength musicxml.NoteLength
	Beam       map[int]entity.Beam
}

var renderMap = map[int]func(canv canvas.Canvas, lines [5]int, pos ...CoordinateWithNoteLength){
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
