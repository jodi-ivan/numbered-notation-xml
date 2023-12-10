package barline

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Barline interface {
	GetRendererLeftBarline(measure musicxml.Measure, x int, lastRightBarlinePosition *entity.Coordinate) (*entity.NoteRenderer, *BarlineInfo)
	GetRendererRightBarline(measure musicxml.Measure, x int) (int, *entity.NoteRenderer)
	RenderBarline(ctx context.Context, canv canvas.Canvas, barline musicxml.Barline, coordinate entity.Coordinate)
}

type barlineInteractor struct{}

func NewBarline() Barline {
	return &barlineInteractor{}
}

func (bi *barlineInteractor) GetRendererLeftBarline(measure musicxml.Measure, x int, lastRightBarlinePosition *entity.Coordinate) (*entity.NoteRenderer, *BarlineInfo) {
	leftBarline := measure.Barline[0]
	if (leftBarline.Location == musicxml.BarlineLocationLeft) && (leftBarline.BarStyle != musicxml.BarLineStyleRegular) {
		pos := x
		if lastRightBarlinePosition != nil {
			pos = int(lastRightBarlinePosition.X)
		}
		result := &entity.NoteRenderer{
			PositionX:     pos,
			Width:         int(barlineWidth[leftBarline.BarStyle]),
			Barline:       &leftBarline,
			MeasureNumber: measure.Number,
		}

		incr := 5

		if leftBarline.Repeat != nil {
			incr += constant.UPPERCASE_LENGTH
		}

		return result, &BarlineInfo{
			XIncrement: incr,
		}

	}

	return nil, nil
}

func (bi *barlineInteractor) GetRendererRightBarline(measure musicxml.Measure, x int) (int, *entity.NoteRenderer) {
	barline := musicxml.Barline{
		BarStyle: musicxml.BarLineStyleRegular,
	}

	if len(measure.Barline) == 1 {
		if measure.Barline[0].Location == musicxml.BarlineLocationRight {
			barline = measure.Barline[0]
		}
	} else if len(measure.Barline) > 1 {
		if measure.Barline[1].Location == musicxml.BarlineLocationRight {
			barline = measure.Barline[1]
		}
	}
	if barline.Repeat != nil && barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
		x += 5
	}

	barlineRenderer := &entity.NoteRenderer{
		MeasureNumber: measure.Number,
		PositionX:     x,
		Barline:       &barline,
	}
	return x, barlineRenderer
}

func (bi *barlineInteractor) RenderBarline(ctx context.Context, canv canvas.Canvas, barline musicxml.Barline, coordinate entity.Coordinate) {
	forward := ""
	backward := ""

	if barline.Repeat != nil {
		if barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
			backward = fmt.Sprintf(`<tspan x="%f" y="%f">:</tspan>`, coordinate.X-5, coordinate.Y)
		} else if barline.Repeat.Direction == musicxml.BarLineRepeatDirectionForward {
			//FIXED: adjust the size and position of forward barline
			forward = fmt.Sprintf(`<tspan x="%f" y="%f">:</tspan>`, coordinate.X+10, coordinate.Y)
		}
	}
	fmt.Fprintf(canv.Writer(), `<text x="%f" y="%f" style="font-family:Noto Music">
		%s
		<tspan x="%f" y="%f" font-size="180%%"> %s </tspan>
		%s
		</text>`,
		coordinate.X,
		coordinate.Y+6,
		backward,
		coordinate.X,
		coordinate.Y+6, unicode[barline.BarStyle], forward)
}

func GetBarlineWidth(barlineStyle musicxml.BarLineStyle) float64 {
	return barlineWidth[barlineStyle]
}
