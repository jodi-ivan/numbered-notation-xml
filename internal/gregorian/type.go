package gregorian

import (
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type VMargin struct {
	Top           entity.Coordinate
	Bottom        entity.Coordinate
	DefaultTop    int
	DefaultBottom int
}

func (m *VMargin) SetTop(margin entity.Coordinate) {
	if m.Top.Y > margin.Y {
		m.Top = margin
	}
}

func (m *VMargin) SetBottom(margin entity.Coordinate) {
	if m.Bottom.Y < margin.Y {
		m.Bottom = margin
	}
}

func (m *VMargin) Set(margin entity.Coordinate) {
	m.SetTop(margin)
	m.SetBottom(margin)
}

func (m *VMargin) Merge(margin VMargin) {
	m.SetTop(margin.Top)
	m.SetBottom(margin.Bottom)
}

func (m *VMargin) GetTopDiffDelta(top int) int {
	return top - int(math.Floor(math.Min(m.Top.Y, float64(top))))
}

func (m *VMargin) GetBottomDiffDelta(bottom int) int {
	return int(math.Ceil(math.Max(m.Bottom.Y, float64(bottom)))) - bottom
}

type SlurWithCoordinates struct {
	Slur   entity.Slur
	Start  entity.Coordinate
	Ending entity.Coordinate
}

type StemInfo struct {
	LengthCompensation float64
	ClampY1            float64
	ClampY2            float64
	HighestYPosition   entity.Coordinate
	LowestYPosition    entity.Coordinate
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
