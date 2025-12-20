package staff

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type RenderStaffWithAlign interface {
	RenderWithAlign(ctx context.Context, canv canvas.Canvas, y int, noteRenderer [][]*entity.NoteRenderer)
}

func NewRenderAlign() RenderStaffWithAlign {
	return &renderStaffAlign{
		Barline:  barline.NewBarline(),
		Numbered: numbered.New(),
		Rhythm:   rhythm.New(),
		Lyric:    lyric.NewLyric(),
	}
}

type renderStaffAlign struct {
	Barline  barline.Barline
	Numbered numbered.Numbered
	Rhythm   rhythm.Rhythm
	Lyric    lyric.Lyric
}

type dotPosition struct {
	beforeXpos int
	afterXPos  int
	address    []*int
}

func (dt *dotPosition) Reset(startPosition int) {
	dt.beforeXpos = startPosition
	dt.address = []*int{}
}

func (dt *dotPosition) Render(endPosition int) {
	if len(dt.address) > 0 {
		dt.afterXPos = endPosition
		space := (dt.afterXPos - dt.beforeXpos) / (len(dt.address) + 1)
		for i, d := range dt.address {
			*d = (dt.beforeXpos + (space * (i + 1)))
		}

		// reset here
		dt.beforeXpos = endPosition
		dt.address = []*int{}
	}
}

func (rsa *renderStaffAlign) RenderWithAlign(ctx context.Context, canv canvas.Canvas, y int, noteRenderer [][]*entity.NoteRenderer) {

	if len(noteRenderer) == 0 {
		return
	}
	flatten := []*entity.NoteRenderer{}
	// get the note
	lastMeasure := noteRenderer[len(noteRenderer)-1]
	lastNote := lastMeasure[len(lastMeasure)-1]

	remaining := (constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH) - lastNote.PositionX

	lastPos := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH
	lastNote.PositionX = lastPos
	if lastNote.Barline != nil {
		remaining -= int(barline.GetBarlineWidth(lastNote.Barline.BarStyle))
	}

	// get last remaining whitespace
	totalNotes := 0
	for _, measure := range noteRenderer {
		totalNotes += len(measure)
	}
	added := float64(remaining) / (float64(totalNotes) - 2)

	count := 1
	slurTiesNote := []*entity.NoteRenderer{}
	dotPositioner := dotPosition{}
	canv.Group("staff")
	for mi, measure := range noteRenderer {

		for i, note := range measure {

			note.PositionY = y
			flatten = append(flatten, note)

			if note.Tie != nil || note.Slur != nil {
				slurTiesNote = append(slurTiesNote, note)
			}

			// do not add left spacing on first note  on the first measure
			if i == 0 && mi == 0 {
				dotPositioner.Reset(note.PositionX)
				continue
			}

			// don't add to the end either
			if mi == len(noteRenderer)-1 && i == len(measure)-1 {
				dotPositioner.Render(note.PositionX)
				continue
			}

			note.PositionX += int(added * float64(count))
			count++

			if note.IsDotted {
				if dotPositioner.address == nil {
					dotPositioner.address = []*int{}
				}
				dotPositioner.address = append(dotPositioner.address, &note.PositionX)
			} else {
				if len(dotPositioner.address) > 0 {
					dotPositioner.Render(note.PositionX)
				} else {
					dotPositioner.Reset(note.PositionX)
				}
			}

			//TODO: reposition if distance betweeb 2 lyrics are zless than 2spaces and 1 dashes space.
			// move the prev to -x distance
		}

		canv.Group("measure-align")
		canv.Group("class='note'", "style='font-family:Old Standard TT;font-weight:500'")
		for _, n := range measure {
			canv.Group("titled-group")
			if n.IsDotted {
				canv.Text(n.PositionX, y, ".")
			} else if n.Articulation != nil && n.Articulation.BreathMark != nil {
				canv.Text(n.PositionX, y-10, ",")
			} else if n.Barline != nil {
				rsa.Barline.RenderBarline(ctx, canv, *n.Barline, entity.Coordinate{
					X: float64(n.PositionX),
					Y: float64(y),
				})
			} else {
				canv.Text(n.PositionX, y, fmt.Sprintf("%d", n.Note))
				if n.Strikethrough {
					canv.Line(n.PositionX+10, y-16, n.PositionX, y+5, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.45")
				}
			}
			fmt.Fprintf(canv.Writer(), `<title>Width: %d</title>`, n.Width)
			canv.Gend()

		}

		canv.Group("class='lyric'", "style='font-family:Caladea'")
		for _, n := range measure {
			if len(n.Lyric) > 0 {
				canv.Group("titled-group")

				for i, l := range n.Lyric {
					if len(l.Text) > 0 {
						lyricVal := entity.LyricVal(l.Text).String()
						xPos := n.PositionX
						if n.PositionX == constant.LAYOUT_INDENT_LENGTH {
							xPos += int(rsa.Lyric.CalculateMarginLeft(lyricVal))
						}
						canv.Text(xPos, n.PositionY+25+(i*20), lyricVal)

						offsetLyric := ""
						for _, t := range l.Text {

							if t.Underline == 1 {
								currTextLength := rsa.Lyric.CalculateLyricWidth(t.Value)
								offset := rsa.Lyric.CalculateLyricWidth(offsetLyric)
								canv.Qbez(
									xPos+int(offset), n.PositionY+28+(i*20),
									xPos+int(offset)+int(currTextLength/2), n.PositionY+28+(i*20)+6,
									xPos+int(offset)+int(currTextLength), n.PositionY+28+(i*20),
									"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
								)
							} else {
								offsetLyric += t.Value
							}
						}
					}
				}
				fmt.Fprintf(canv.Writer(), `<title>Width: %d</title>`, n.Width)
				canv.Gend()
			}
		}

		canv.Gend()

		rsa.Numbered.RenderOctave(ctx, canv, measure)
		rsa.Rhythm.RenderBeam(ctx, canv, measure)

		rsa.RenderMeasureText(ctx, canv, measure)
		RenderTuplet(ctx, canv, measure)

		canv.Gend()
		canv.Gend()

	}

	rsa.Lyric.RenderHypen(ctx, canv, flatten)
	rsa.Rhythm.RenderSlurTies(ctx, canv, slurTiesNote, float64(lastPos))
	RenderMeasureTopping(ctx, canv, flatten)
	canv.Gend()

}
