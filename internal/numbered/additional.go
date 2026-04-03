package numbered

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func RendererFromAdditional(note musicxml.Note, header *entity.NoteRenderer, additionals []NoteLength) []*entity.NoteRenderer {

	result := []*entity.NoteRenderer{}
	for i, additional := range additionals {

		var additionalNote *entity.NoteRenderer
		if i == 0 {
			additionalNote = header
		} else {
			additionalNote = &entity.NoteRenderer{
				PositionY:     header.PositionY,
				Width:         constant.LOWERCASE_LENGTH,
				IsDotted:      additional.IsDotted,
				NoteLength:    additional.Type,
				Beam:          map[int]entity.Beam{},
				MeasureNumber: header.MeasureNumber,
				IsNewLine:     header.IsNewLine && (i == len(additionals)-1) && !note.IsBreathMark(),
			}
			if additionalNote.IsNewLine {
				header.IsNewLine = !additionalNote.IsNewLine
			}
		}
		switch additional.Type {
		case musicxml.NoteLength16th:
			additionalNote.Beam[2] = entity.Beam{
				Number: 2,
				Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
			}
			fallthrough
		case musicxml.NoteLengthEighth:
			additionalNote.Beam[1] = entity.Beam{
				Number: 1,
				Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
			}
		}

		result = append(result, additionalNote)

	}

	return result
}
