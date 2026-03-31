package staff

import "github.com/jodi-ivan/numbered-notation-xml/internal/entity"

type StaffInfo struct {
	Multiline        bool
	MarginBottom     int
	MarginLeft       int
	NextLineRenderer []*entity.NoteRenderer

	ForceNewLine bool
}

const (
	MEASURE_TEXT_REFREIN = "Refrein"
	MEASURE_TEXT_FINE    = "Fine"

	FIRST_STAFF_Y_POS   = 95
	MEASURE_TEXT_OFFSET = 15
)
