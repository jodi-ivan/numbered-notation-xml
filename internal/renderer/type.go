package renderer

import "github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"

type Lyric struct {
	Text     string
	Syllabic musicxml.LyricSyllabic
}

type Slur struct {
	// Number attributes for slur
	// Pitch note for ties
	Number int
	Type   musicxml.NoteSlurType
}

type Beam struct {
	Number int
	Type   musicxml.NoteBeamType
}

type NoteRenderer struct {
	IsDotted     bool
	IsRest       bool
	PositionX    int
	PositionY    int
	Note         int
	Octave       int
	Striketrough bool
	NoteLength   musicxml.NoteLength
	BarType      string
	Width        int
	Lyric        []Lyric
	Slur         map[int]Slur
	Beam         map[int]Beam
	Tie          *Slur
	Articulation *Articulation

	// internal use
	isLengthTakenFromLyric bool
	indexPosition          int
	isNewLine              bool
}

func (nr *NoteRenderer) UpdateBeam(beamNum int, beamType musicxml.NoteBeamType) {
	newBeam := nr.Beam

	newBeam[beamNum] = Beam{
		Number: beamNum,
		Type:   beamType,
	}

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

type Articulation struct {
	BreathMark *ArticulationTypes
}
type ArticulationTypes string

var (
	ArticulationTypesBreathMark ArticulationTypes = "breathMark"
)
