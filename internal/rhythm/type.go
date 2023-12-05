package rhythm

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
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
	Start CoordinateWithOctave
	End   CoordinateWithOctave
	Pull  CoordinateWithOctave
}

type BeamLine struct {
	Start entity.Coordinate
	End   entity.Coordinate
}

type beamSplitMarker struct {
	StartIndex int
	EndIndex   int
}
