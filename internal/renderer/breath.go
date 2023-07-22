package renderer

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderBreath(ctx context.Context, canv canvas.Canvas, notes []*NoteRenderer) {
	for _, note := range notes {
		if note.Articulation != nil && note.Articulation.BreathMark != nil {
			canv.Text(note.PositionX+5, note.PositionY-10, ",")
		}
	}
}
