package entity

import (
	"fmt"

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
	Bold      bool
	Italic    bool
}
type Lyric struct {
	Text     []Text
	Syllabic musicxml.LyricSyllabic
}

type Slur struct {
	// Number attributes for slur
	Number   int
	Type     musicxml.NoteSlurType
	LineType musicxml.NoteSlurLineType

	NumberedOnly bool
}

type Beam struct {
	Number int
	Type   musicxml.NoteBeamType
	Locked bool
}

type Articulation struct {
	BreathMark *ArticulationTypes
}
type ArticulationTypes string

var (
	ArticulationTypesBreathMark ArticulationTypes = "breathMark"
)

type NoteRenderer struct {
	UUID string

	AbsoluteNote       string
	AbsoluteOctave     int
	AbsoluteAccidental musicxml.NoteAccidental

	Articulation  *Articulation
	Barline       *musicxml.Barline
	Beam          map[int]Beam
	Fermata       *musicxml.Femata
	IsDotted      bool
	IsRest        bool
	Lyric         []Lyric
	Note          int
	NoteLength    musicxml.NoteLength
	NoteValue     float64
	Octave        int
	PositionX     int
	PositionY     int
	Slur          map[int]Slur
	Strikethrough bool
	Tie           *Slur
	Width         int

	// internal use
	IndexPosition          int
	IsAdditional           bool
	IsLengthTakenFromLyric bool
	IsNewLine              bool
	MeasureNumber          int

	LeadingHeader     string
	MeasureDash       map[int]musicxml.DirectionDashesType
	MeasureText       []musicxml.MeasureText
	TimeModifications *musicxml.TimeModification
	Tuplet            *musicxml.Tuplet
}

func (nr *NoteRenderer) GetNonAccidentalAbsoluteNote() string {
	return fmt.Sprintf("%s%d", nr.AbsoluteNote, nr.AbsoluteOctave)
}

func (nr *NoteRenderer) UpdateBeamWithLock(beamNum int, beamType musicxml.NoteBeamType) {
	newBeam := nr.Beam
	if len(newBeam) == 0 {
		return
	}

	if b, ok := newBeam[beamNum]; !ok || b.Locked {
		return
	}

	newBeam[beamNum] = Beam{
		Number: beamNum,
		Type:   beamType,
		Locked: true,
	}

}

func (nr *NoteRenderer) UpdateBeam(beamNum int, beamType musicxml.NoteBeamType) {
	newBeam := nr.Beam
	if len(newBeam) == 0 {
		return
	}

	if b, ok := newBeam[beamNum]; !ok || b.Locked {
		return
	}

	newBeam[beamNum] = Beam{
		Number: beamNum,
		Type:   beamType,
	}

}
