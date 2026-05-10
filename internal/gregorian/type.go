package gregorian

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type SlurWithCoordinates struct {
	Slur   entity.Slur
	Start  entity.Coordinate
	Ending entity.Coordinate
}

type StemInfo struct {
	LengthCompensation float64
	ClampY1            float64
	ClampY2            float64
	LowestYPosition    float64
	Flip               bool
}

type CoordinateWithNoteLength struct {
	entity.Coordinate
	NoteLength musicxml.NoteLength
	Beam       map[int]entity.Beam
	Direction  *int
	NoteID     string
}

type SlurTieGroup struct {
	AccumulativeDirection int

	Start entity.Coordinate
	End   entity.Coordinate

	NoteMember []string
	Ties       *entity.Slur
	Slur       *entity.Slur

	MaxY float64
	MinY float64
}
