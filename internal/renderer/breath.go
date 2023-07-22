package renderer

import (
	"context"

	svg "github.com/ajstarks/svgo"
)

func RenderBreath(ctx context.Context, canvas *svg.SVG, notes []*NoteRenderer) {
	for _, note := range notes {
		if note.Articulation != nil && note.Articulation.BreathMark != nil {
			canvas.Text(note.PositionX+2, note.PositionY-10, ",")
		}
	}
}
