package entity

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type Coordinate struct {
	X float64
	Y float64
}

func (c *Coordinate) IsEmpty() bool {
	return c.X == 0 && c.Y == 0
}

func NewCoordinate(x, y float64) Coordinate {
	return Coordinate{
		X: x,
		Y: y,
	}
}

type CoordinateWithNoteLength struct {
	Coordinate
	NoteLength musicxml.NoteLength
	Beam       map[int]Beam
	NoteID     string
	Tuplet     *musicxml.Tuplet
}
