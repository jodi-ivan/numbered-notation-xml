package numbered

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (ni *numberedInteractor) RenderOctave(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer) {
	hasOctave := false
	for _, note := range notes {
		if !hasOctave && (note.Octave != 0) {
			canv.Group("class='octaves'")
		}
		hasOctave = hasOctave || (note.Octave != 0)
		if note.Octave < 0 {
			canv.Circle(note.PositionX+5, note.PositionY+5, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
			continue
		}

		if note.Octave > 0 {
			canv.Circle(note.PositionX+5, note.PositionY-15, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
		}
	}
	if hasOctave {
		canv.Gend()
	}
}
