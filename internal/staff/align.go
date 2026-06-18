package staff

import (
	"context"
	"fmt"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/gregorian"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm/splitter"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/text"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type RenderStaffWithAlign interface {
	RenderWithAlign(ctx context.Context, canv canvas.Canvas, staffPos, y int, ts timesig.TimeSignature, ks keysig.KeySignature, noteRenderer [][]*entity.NoteRenderer) int
}

func NewRenderAlign() RenderStaffWithAlign {

	barlineInteractor := barline.NewBarline()
	lyricInteractor := lyric.NewLyric()
	return &renderStaffAlign{
		Numbered: numbered.New(lyricInteractor, barlineInteractor),
		Rhythm:   rhythm.New(splitter.New()),
		Barline:  barlineInteractor,
		Lyric:    lyricInteractor,
		Text:     text.NewText(lyricInteractor),
	}
}

type renderStaffAlign struct {
	Barline  barline.Barline
	Numbered numbered.Numbered
	Rhythm   rhythm.Rhythm
	Lyric    lyric.Lyric
	Text     text.Text
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

func (rsa *renderStaffAlign) RenderWithAlign(ctx context.Context, canv canvas.Canvas, staffPos, y int, ts timesig.TimeSignature, ks keysig.KeySignature, noteRenderer [][]*entity.NoteRenderer) int {

	if len(noteRenderer) == 0 {
		return 0
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

	additionalMarginBottom := 0

	for mi, measure := range noteRenderer { // preparation adn precalculate
		alignJustify(measure, y, added, &count, mi, mi == len(noteRenderer)-1)
		rsa.Rhythm.Split(ctx, ts, measure)

		for notePos, note := range measure {
			if (len(note.Lyric) > 0 || note.Note > 0 || note.IsRest) && notePos == len(measure)-1 {
				xPos := note.PositionX
				noteStr := fmt.Sprintf("%d", note.Note)

				noteWidth := rsa.Lyric.CalculateLyricWidth(noteStr)
				xPos = xPos + rightAlignOffset - int(math.Round(noteWidth))

				note.PositionX = xPos
			}

			if (note.Barline != nil && note.Barline.Ending != nil) || note.Fermata != nil {
				additionalMarginBottom = 6
			}
		}
	}

	canv.Group(`class="gregorian"`, "style='font-family:mozart11'")
	margin := gregorian.RenderStaffLine(ctx, staffPos, y, canv, flatten, ks, ts)
	RenderMeasureTopping(ctx, y+10, canv, flatten, true)
	canv.Gend()

	stafflines := lines.NewLineStaffWithLines(ts, ks, y)

	yPos := y + gregorian.STAFF_OFFSET + (int(margin.Bottom.Y) - margin.DefaultBottom) + additionalMarginBottom

	canv.Group(`class="numbered"`)
	offsetLyric := 0
	for mi, measure := range noteRenderer {

		canv.Group("class='measure-align'", fmt.Sprintf("number='%d'", measure[0].MeasureNumber))

		prev := []*entity.NoteRenderer{}
		if mi > 0 {
			prevMeasure := noteRenderer[mi-1]

			idx := -1
			for i := len(prevMeasure) - 1; i >= 0; i-- {
				if len(prevMeasure[i].Lyric) > 0 {
					idx = i
					break
				}
			}

			if idx >= 0 {
				prev = append(prev, prevMeasure[idx])

			}
		}
		newOffsetLyric := rsa.Lyric.RenderLyrics(ctx, yPos+offsetLyric, canv, measure, prev...)
		if newOffsetLyric > 0 && offsetLyric == 0 {
			offsetLyric = newOffsetLyric
		}

		canv.Group("class='note'", "style='font-family:Old Standard TT;font-weight:500'")
		rsa.Numbered.RenderNote(ctx, canv, measure, yPos, rightAlignOffset)
		rsa.Rhythm.RenderBeam(ctx, yPos, canv, ts, measure)

		canv.Group("class='staff-text'")

		rsa.Text.RenderMeasureText(ctx, y+10, canv, measure, stafflines)
		RenderStaffLineDash(measure, canv, y+10, stafflines)
		RenderTuplet(ctx, yPos, canv, measure)

		canv.Gend()

		canv.Gend()
		canv.Gend()

	}

	rsa.Lyric.RenderHypen(ctx, yPos, offsetLyric, canv, flatten)
	rsa.Rhythm.RenderSlurTies(ctx, yPos, canv, slurTiesNote, float64(lastPos))
	RenderMeasureTopping(ctx, yPos, canv, flatten)
	canv.Gend()
	canv.Gend()

	// canv.Circle(int(margin.Top.X), int(margin.Top.Y), 2, "stroke-width:1;fill:none;stroke:#FF0000")
	// canv.Circle(int(margin.Bottom.X), int(margin.Bottom.Y), 2, "stroke-width:1;fill:none;stroke:#FF0000")

	return int(margin.Bottom.Y) - margin.DefaultBottom

}
