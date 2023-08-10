package renderer

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type StaffInfo struct {
	Multiline        bool
	MarginBottom     int
	MarginLeft       int
	NextLineRenderer []*entity.NoteRenderer
}

type Coordinate struct {
	X float64
	Y float64
}

type CoordinateWithOctave struct {
	Coordinate
	Octave int
}

// for (svg.SVG).Qbez
type SlurBezier struct {
	Start CoordinateWithOctave
	End   CoordinateWithOctave
	Pull  CoordinateWithOctave
}

type BeamLine struct {
	Start Coordinate
	End   Coordinate
}

type beamMarker struct {
	NoteBeamType   musicxml.NoteBeamType
	NoteBeginIndex int
}

type beamSplitMarker struct {
	StartIndex int
	EndIndex   int
}
