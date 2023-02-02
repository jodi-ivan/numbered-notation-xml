package renderer

import svg "github.com/ajstarks/svgo"

func RenderOctave(canvas *svg.SVG, notes []*NoteRenderer) {
	canvas.Group("class='octaves'")
	for _, note := range notes {
		if note.Octave < 0 {
			canvas.Circle(note.PositionX+5, note.PositionY+5, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
		}

		if note.Octave > 0 {
			canvas.Circle(note.PositionX+5, note.PositionY-15, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
		}
	}
	canvas.Gend()
}
