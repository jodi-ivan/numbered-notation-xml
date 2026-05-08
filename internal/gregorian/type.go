package gregorian

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type StemInfo struct {
	LengthCompensation float64
	ClampY1            float64
	ClampY2            float64
	LowestYPosition    float64
}

type CoordinateWithNoteLength struct {
	entity.Coordinate
	NoteLength musicxml.NoteLength
	Beam       map[int]entity.Beam
}
