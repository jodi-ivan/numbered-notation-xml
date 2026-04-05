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

func alignJustify(measure []*entity.NoteRenderer, y int, addedSpace float64, count *int, measureIndex int, lastMeasure bool) {
	dotPositioner := numbered.DotPosition{}

	for i, note := range measure { // measure

		note.PositionY = y

		// do not add left spacing on first note  on the first measure
		if i == 0 && measureIndex == 0 {
			dotPositioner.Reset(note.PositionX)
			continue
		}

		// don't add to the end either
		if lastMeasure && i == len(measure)-1 {
			dotPositioner.Render(note.PositionX)
			continue
		}

		note.PositionX += int(addedSpace * float64(*count))
		*count = *count + 1

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
}

func (rsa *renderStaffAlign) getAddedSpace(lastNote *entity.NoteRenderer, rightAlignOffset *int, totalNotes int) (float64, int) {

	remaining := float64((constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH) - lastNote.PositionX)

	lastPos := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH
	lastNote.PositionX = lastPos + 4
	if lastNote.Barline != nil {
		remaining -= barline.GetBarlineWidth(lastNote.Barline.BarStyle)
	} else if len(lastNote.Lyric) > 0 {
		lyricWidth := int(math.Round(rsa.Lyric.CalculateOverallWidth(lastNote.Lyric)))
		*rightAlignOffset = lyricWidth / 2
		lastNote.PositionX -= lyricWidth
		remaining -= float64(lyricWidth)
	}

	return remaining / float64(totalNotes), lastPos
}

func (rsa *renderStaffAlign) RenderWithAlign(ctx context.Context, canv canvas.Canvas, y int, ts timesig.TimeSignature, noteRenderer [][]*entity.NoteRenderer) {

	if len(noteRenderer) == 0 {
		return
	}
	flatten := []*entity.NoteRenderer{}

	count := 1
	slurTiesNote := []*entity.NoteRenderer{}
	rightAlignOffset := 0

	// proprocessing
	totalNotes := 0
	for _, measure := range noteRenderer {
		totalNotes += len(measure)
		flatten = append(flatten, measure...)

		for _, note := range measure {
			if note.Tie != nil || note.Slur != nil {
				slurTiesNote = append(slurTiesNote, note)
			}
		}
	}
	lastMeasure := noteRenderer[len(noteRenderer)-1]
	lastNote := lastMeasure[len(lastMeasure)-1]
	added, lastPos := rsa.getAddedSpace(lastNote, &rightAlignOffset, totalNotes)

	canv.Group("class='staff'")
	for mi, measure := range noteRenderer { // staff
		alignJustify(measure, y, added, &count, mi, mi == len(noteRenderer)-1)

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
