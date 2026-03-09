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

type beamMarker struct {
	NoteBeamType   musicxml.NoteBeamType
	NoteBeginIndex int
}

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
	Start entity.Coordinate
	End   entity.Coordinate
}

type beamSplitMarker struct {
	StartIndex int
	EndIndex   int
}

type Interval []beamSplitMarker

func (c Interval) Len() int           { return len(c) }
func (c Interval) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Interval) Less(i, j int) bool { return c[i].StartIndex < c[j].StartIndex }
