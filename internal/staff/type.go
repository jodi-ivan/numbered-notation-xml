package staff

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

type StaffInfo struct {
	Multiline        bool
	MarginBottom     int
	MarginLeft       int
	NextLineRenderer []*entity.NoteRenderer
	EndIndex         int // for staff level numbering index

	StartRenderOtherNotes bool
	ForceNewLine          bool
	SyllableCount         int
}

type StaffData struct {
	PrevNotes     []*entity.NoteRenderer
	SyllableCount int
	TimeSig       timesig.TimeSignature
	KeySig        keysig.KeySignature
	IndexStart    int
	ReffAtStart   bool
}
type CoordinateWithTuplet struct {
	entity.Coordinate
	Tuplet musicxml.Tuplet
}

const (
	MEASURE_TEXT_REFREIN = "Refrein"
	MEASURE_TEXT_FINE    = "Fine"

	FIRST_STAFF_Y_POS   = 95
	MEASURE_TEXT_OFFSET = 15
)
