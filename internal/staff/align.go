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
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm/splitter"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type RenderStaffWithAlign interface {
	RenderWithAlign(ctx context.Context, canv canvas.Canvas, y int, ts timesig.TimeSignature, noteRenderer [][]*entity.NoteRenderer)
}

func NewRenderAlign() RenderStaffWithAlign {

	barlineInteractor := barline.NewBarline()
	lyricInteractor := lyric.NewLyric()
	return &renderStaffAlign{
		Numbered: numbered.New(lyricInteractor, barlineInteractor),
		Rhythm:   rhythm.New(splitter.New()),
		Barline:  barlineInteractor,
		Lyric:    lyricInteractor,
	}
}

type renderStaffAlign struct {
	Barline  barline.Barline
	Numbered numbered.Numbered
	Rhythm   rhythm.Rhythm
	Lyric    lyric.Lyric
}

func (rsa *renderStaffAlign) RenderWithAlign(ctx context.Context, canv canvas.Canvas, y int, ts timesig.TimeSignature, noteRenderer [][]*entity.NoteRenderer) {

	if len(noteRenderer) == 0 {
		return
	}
	flatten := []*entity.NoteRenderer{}

	count := 1
	slurTiesNote := []*entity.NoteRenderer{}
	dotPositioner := numbered.DotPosition{}
	rightAlignOffset := 0

	// proprocessing
	totalNotes := 0
	for _, measure := range noteRenderer {
		totalNotes += len(measure)
	}

	// get last remaining whitespace
	// get the note
	lastMeasure := noteRenderer[len(noteRenderer)-1]
	lastNote := lastMeasure[len(lastMeasure)-1]

	remaining := (constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH) - lastNote.PositionX

	lastPos := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH
	lastNote.PositionX = lastPos + 4
	if lastNote.Barline != nil {
		remaining -= int(barline.GetBarlineWidth(lastNote.Barline.BarStyle))
	} else if len(lastNote.Lyric) > 0 {
		lyricWidth := int(math.Round(rsa.Lyric.CalculateOverallWidth(lastNote.Lyric)))
		rightAlignOffset = lyricWidth / 2
		lastNote.PositionX -= lyricWidth
		remaining -= lyricWidth
	}

	added := float64(remaining) / (float64(totalNotes))

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

			if breathpause.IsBreathMark(note) {
				note.PositionX -= 5
			}
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

		canv.Group("class='measure-align'", fmt.Sprintf("number='%d'", measure[0].MeasureNumber))

		rsa.Lyric.RenderLyrics(ctx, canv, measure)

		canv.Group("class='note'", "style='font-family:Old Standard TT;font-weight:500'")

		rsa.Numbered.RenderNote(ctx, canv, measure, y, rightAlignOffset)
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
