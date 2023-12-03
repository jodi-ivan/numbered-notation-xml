package renderer

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

var barlineWidth = map[musicxml.BarLineStyle]float64{
	musicxml.BarLineStyleRegular:    4.16,
	musicxml.BarLineStyleLightHeavy: 7.7,
	musicxml.BarLineStyleLightLight: 6.28,
	musicxml.BarLineStyleHeavyHeavy: 8,
	musicxml.BarLineStyleHeavyLight: 7.7,
}

func RenderWithAlign(ctx context.Context, canv canvas.Canvas, y int, noteRenderer [][]*entity.NoteRenderer) {

	flatten := []*entity.NoteRenderer{}
	// get the note
	lastMeasure := noteRenderer[len(noteRenderer)-1]
	lastNote := lastMeasure[len(lastMeasure)-1]

	remaining := (constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH) - lastNote.PositionX

	lastPos := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH
	lastNote.PositionX = lastPos
	if lastNote.Barline != nil {
		remaining -= int(barlineWidth[lastNote.Barline.BarStyle])
	} else if lastNote.Articulation != nil && lastNote.Articulation.BreathMark != nil {
		// TODO: the breakmark last notess
		lastTwo := lastMeasure[len(lastMeasure)-2]

		remaining -= (lastTwo.Width - int(lyric.CalculateLyricWidth(",")))
	} else {
		lastNote.PositionX -= lastNote.Width
		remaining -= lastNote.Width
	}

	// get last remaining whitespace
	totalNotes := 0
	for _, measure := range noteRenderer {
		totalNotes += len(measure)
	}
	added := float64(remaining) / (float64(totalNotes) - 2)

	count := 1
	slurTiesNote := []*entity.NoteRenderer{}
	canv.Group("staff")
	for mi, measure := range noteRenderer {

		for i, note := range measure {

			note.PositionY = y
			flatten = append(flatten, note)

			if note.Tie != nil || note.Slur != nil {
				slurTiesNote = append(slurTiesNote, note)
			}

			// do not add left spacing on first not first measure
			if i == 0 && mi == 0 {
				continue
			}

			// don't add to the end either
			if mi == len(noteRenderer)-1 && i == len(measure)-1 {
				continue
			}

			note.PositionX += int(added * float64(count))
			count++

			//TODO: reposition if distance betweeb 2 lyrics are zless than 2spaces and 1 dashes space.
			// move the prev to -x distance
		}

		canv.Group("measure-align")
		canv.Group("class='note'", "style='font-family:Old Standard TT;font-weight:500'")
		for _, n := range measure {

			if n.IsDotted {
				canv.Text(n.PositionX, y, ".")
			} else if n.Articulation != nil && n.Articulation.BreathMark != nil {
				canv.Text(n.PositionX, y-10, ",")
			} else if n.Barline != nil {
				RenderBarline(ctx, canv, *n.Barline, entity.Coordinate{
					X: float64(n.PositionX),
					Y: float64(y),
				})
			} else {
				canv.Text(n.PositionX, y, fmt.Sprintf("%d", n.Note))
				if n.Striketrough {
					canv.Line(n.PositionX+10, y-16, n.PositionX, y+5, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.45")
				}
			}

		}

		canv.Group("class='lyric'", "style='font-family:Caladea'")
		for _, n := range measure {
			if len(n.Lyric) > 0 {
				for i, l := range n.Lyric {
					if len(l.Text) > 0 {
						lyricVal := entity.LyricVal(l.Text).String()
						xPos := n.PositionX
						if n.PositionX == constant.LAYOUT_INDENT_LENGTH {
							xPos += int(lyric.CalculateMarginLeft(lyricVal))
						}
						canv.Text(xPos, n.PositionY+25+(i*20), lyricVal)

						offsetLyric := ""
						for _, t := range l.Text {

							if t.Underline == 1 {
								currTextLength := lyric.CalculateLyricWidth(t.Value)
								offset := lyric.CalculateLyricWidth(offsetLyric)
								canv.Qbez(
									xPos+int(offset), n.PositionY+28,
									xPos+int(offset)+int(currTextLength/2), n.PositionY+28+6,
									xPos+int(offset)+int(currTextLength), n.PositionY+28,
									"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
								)
							} else {
								offsetLyric += t.Value
							}
						}
					}
				}
			}
		}

		canv.Gend()

		RenderOctave(ctx, canv, measure)
		RenderMeasureText(ctx, canv, measure)
		RenderBeam(ctx, canv, measure)
		RenderTuplet(ctx, canv, measure)
		canv.Gend()
		canv.Gend()

	}

	lyric.RenderHypen(ctx, canv, flatten)
	RenderMeasureTopping(ctx, canv, flatten)
	RenderSlurTies(ctx, canv, slurTiesNote, float64(lastPos))
	canv.Gend()

}
