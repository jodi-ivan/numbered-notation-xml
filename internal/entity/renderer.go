package entity

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type LyricVal []Text

func (lv LyricVal) String() string {
	result := ""
	for _, l := range lv {
		result += l.Value
	}

	return result
}

type Text struct {
	Value     string
	Underline int
}
type Lyric struct {
	Text     []Text
	Syllabic musicxml.LyricSyllabic
}

type Slur struct {
	// Number attributes for slur
	Number int
	Type   musicxml.NoteSlurType
}

type Beam struct {
	Number int
	Type   musicxml.NoteBeamType
}

type Articulation struct {
	BreathMark *ArticulationTypes
}
type ArticulationTypes string

var (
	ArticulationTypesBreathMark ArticulationTypes = "breathMark"
)

type NoteRenderer struct {
	IsDotted     bool
	IsRest       bool
	PositionX    int
	PositionY    int
	Note         int
	Octave       int
	Striketrough bool
	NoteLength   musicxml.NoteLength
	// BarType      string
	Width        int
	Lyric        []Lyric
	Slur         map[int]Slur
	Beam         map[int]Beam
	Tie          *Slur
	Articulation *Articulation
	Barline      *musicxml.Barline
	// internal use
	IsLengthTakenFromLyric bool
	IndexPosition          int
	IsNewLine              bool
	MeasureNumber          int

	MeasureText    []musicxml.MeasureText
	Tuplet         *musicxml.Tuplet
	TimeMofication *musicxml.TimeModification
}

func (nr *NoteRenderer) UpdateBeam(beamNum int, beamType musicxml.NoteBeamType) {
	newBeam := nr.Beam

	newBeam[beamNum] = Beam{
		Number: beamNum,
		Type:   beamType,
	}

}
