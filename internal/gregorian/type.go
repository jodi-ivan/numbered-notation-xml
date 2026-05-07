package gregorian

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type CoordinateWithNoteLength struct {
	entity.Coordinate
	NoteLength musicxml.NoteLength
	Beam       map[int]entity.Beam
}
