package renderer

import "github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"

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
