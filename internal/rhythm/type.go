package rhythm

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type SlurTieType string

const (
	SlurTieTypeSlur SlurTieType = "slur"
	SlurTieTypeTie  SlurTieType = "tie"
)

type CoordinateWithOctave struct {
	entity.Coordinate
	Octave int
}

// for (svg.SVG).Qbez
type SlurBezier struct {
	Start    CoordinateWithOctave
	End      CoordinateWithOctave
	Pull     CoordinateWithOctave
	LineType musicxml.NoteSlurLineType

	SlurTieType SlurTieType
}

type BeamLine struct {
	Start  entity.Coordinate
	End    entity.Coordinate
	Number int
}
