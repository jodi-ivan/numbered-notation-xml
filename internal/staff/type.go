package staff

import "github.com/jodi-ivan/numbered-notation-xml/internal/entity"

type StaffInfo struct {
	Multiline        bool
	MarginBottom     int
	MarginLeft       int
	NextLineRenderer []*entity.NoteRenderer
}

const MEASURE_TEXT_REFREIN = "Refrein"
