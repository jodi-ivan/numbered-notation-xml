package staff

import (
	"context"
	"fmt"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type RenderStaffWithAlign interface {
	RenderWithAlign(ctx context.Context, canv canvas.Canvas, y int, ts timesig.TimeSignature, noteRenderer [][]*entity.NoteRenderer)
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

// TODO: right align text on the last node.
func (rsa *renderStaffAlign) RenderWithAlign(ctx context.Context, canv canvas.Canvas, y int, ts timesig.TimeSignature, noteRenderer [][]*entity.NoteRenderer) {

	if len(noteRenderer) == 0 {
		return
	}
	flatten := []*entity.NoteRenderer{}
	// get the note
	lastMeasure := noteRenderer[len(noteRenderer)-1]
	lastNote := lastMeasure[len(lastMeasure)-1]

	count := 1
	slurTiesNote := []*entity.NoteRenderer{}
	dotPositioner := numbered.DotPosition{}
	rightAlignOffset := 0

	// proprocessing
	totalNotes := 0
	for _, measure := range noteRenderer {
		var prev, next *entity.NoteRenderer

		for i, note := range measure {

			if i < len(measure)-1 {
				next = measure[i+1]
			}

			//clean up breathmark pause
			if note.Articulation != nil && note.Articulation.BreathMark != nil {
				breathpause.AdjustBreathmarkBeamCont(ctx, note, prev, next)
			}

			prev = note
		}
		totalNotes += len(measure)
	}

	// get last remaining whitespace
	remaining := (constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH) - lastNote.PositionX

	lastPos := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH
	lastNote.PositionX = lastPos
	if lastNote.Barline != nil {
		remaining -= int(barline.GetBarlineWidth(lastNote.Barline.BarStyle))
	} else if len(lastNote.Lyric) > 0 {
		lyricWidth := int(math.Round(rsa.Lyric.CalculateOverallWidth(lastNote.Lyric)))
		rightAlignOffset = lyricWidth / 2
		lastNote.PositionX -= lyricWidth
		remaining -= lyricWidth
	}

	added := float64(remaining) / (float64(totalNotes) - 2)

	canv.Group("staff")
	for mi, measure := range noteRenderer { // staff
		for i, note := range measure { // measure

			flatten = append(flatten, note)
			note.PositionY = y

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
				if dotPositioner.Address == nil {
					dotPositioner.Address = []*int{}
				}
				dotPositioner.Address = append(dotPositioner.Address, &note.PositionX)
			} else {
				if len(dotPositioner.Address) > 0 {
					dotPositioner.Render(note.PositionX)
				} else {
					dotPositioner.Reset(note.PositionX)
				}
			}
		}

		canv.Group("measure-align")
		// if len(measure) > 0 {
		// 	fmt.Fprintf(canv.Writer(), `<title>Measure %d</title>`, measure[0].MeasureNumber)
		// }

		canv.Group("class='lyric'", "style='font-family:Caladea'")
		for _, n := range measure {
			for i, l := range n.Lyric {
				if len(l.Text) == 0 {
					continue
				}
				lyricVal := entity.LyricVal(l.Text).String()
				xPos := n.PositionX
				if n.PositionX == constant.LAYOUT_INDENT_LENGTH {
					xPos += int(rsa.Lyric.CalculateMarginLeft(lyricVal))
				}
				canv.Text(xPos, n.PositionY+25+(i*20), lyricVal)
				rsa.Lyric.RenderElision(ctx, canv, l.Text, i, entity.Coordinate{X: float64(xPos), Y: float64(n.PositionY)})
				n.Lyric[i] = l

			}
		}

		canv.Gend()
		canv.Group("class='note'", "style='font-family:Old Standard TT;font-weight:500'")
		for notePos, n := range measure {
			if n.IsDotted {
				canv.Text(n.PositionX, y, ".")
			} else if n.Articulation != nil && n.Articulation.BreathMark != nil {
				canv.Text(n.PositionX-5, y-10, ",")
			} else if n.Barline != nil {
				rsa.Barline.RenderBarline(ctx, canv, *n.Barline, entity.Coordinate{
					X: float64(n.PositionX),
					Y: float64(y),
				})
			} else {
				xPos := n.PositionX
				noteStr := fmt.Sprintf("%d", n.Note)
				noteWidth := rsa.Lyric.CalculateLyricWidth(noteStr)
				if notePos == len(measure)-1 {
					xPos = xPos + rightAlignOffset - int(math.Round(noteWidth))
				}
				canv.Text(xPos, y, noteStr)

				coordinate := entity.Coordinate{X: float64(xPos), Y: float64(n.PositionY)}
				rsa.Numbered.RenderStrikethrough(ctx, canv, n.Strikethrough, coordinate)
				breathpause.RenderFermata(ctx, canv, n.Fermata, coordinate)
				rsa.Numbered.RenderOctave(ctx, canv, n.Octave, coordinate)
				n.PositionX = xPos

			}

		}

		rsa.Rhythm.RenderBeam(ctx, canv, ts, measure)
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
