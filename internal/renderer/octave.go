package renderer

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderOctave(ctx context.Context, canv canvas.Canvas, notes []*NoteRenderer) {
	canv.Group("class='octaves'")
	for _, note := range notes {
		if note.Octave < 0 {
			canv.Circle(note.PositionX+5, note.PositionY+5, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
		}

		if note.Octave > 0 {
			canv.Circle(note.PositionX+5, note.PositionY-15, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
		}
	}
	canv.Gend()
}
