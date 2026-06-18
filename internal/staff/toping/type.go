package toping

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type CoordinateWithTuplet struct {
	entity.Coordinate
	Tuplet musicxml.Tuplet
}
