package barline

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Barline interface {
	GetRendererLeftBarline(measure musicxml.Measure, x int, lastRightBarlinePosition *CoordinateWithBarline) (*entity.NoteRenderer, *BarlineInfo)
	GetRendererRightBarline(measure musicxml.Measure, x int) (int, *entity.NoteRenderer)
	RenderBarline(ctx context.Context, canv canvas.Canvas, barline musicxml.Barline, coordinate entity.Coordinate, s ...string)
}

type barlineInteractor struct{}

func NewBarline() Barline {
	return &barlineInteractor{}
}

// GetRendererLeftBarline render the left side of the barline
// this is only utilize then the barline is not regular barline
// since the regular left barline is already added by default
func (bi *barlineInteractor) GetRendererLeftBarline(measure musicxml.Measure, x int, lastRightBarlinePosition *CoordinateWithBarline) (*entity.NoteRenderer, *BarlineInfo) {

	if len(measure.Barline) == 0 {
		return nil, nil
	}
	leftBarline := measure.Barline[0]
	if (leftBarline.Location == musicxml.BarlineLocationLeft) && (leftBarline.BarStyle != musicxml.BarLineStyleRegular) {
		pos := x
		lastBarlineRepeat := lastRightBarlinePosition != nil && lastRightBarlinePosition.Barline.Repeat != nil

		if lastRightBarlinePosition != nil {
			pos = int(lastRightBarlinePosition.X)
			if lastBarlineRepeat {
				pos += LEFT_BARLINE_RIGHT_AND_LEFT_REPEAT - 4
			}
		}
		result := &entity.NoteRenderer{
			PositionX:     pos,
			Width:         int(barlineWidth[leftBarline.BarStyle]) + BARLINE_AFTER_SPACE,
			Barline:       &leftBarline,
			MeasureNumber: measure.Number,
		}

		incr := 5
		if leftBarline.Repeat != nil && !lastBarlineRepeat {
			// HACK: do we need to check the direction == forward in this?
			incr += constant.UPPERCASE_LENGTH
		}

		return result, &BarlineInfo{
			XIncrement: incr + BARLINE_AFTER_SPACE,
		}

	}

	return nil, nil
}

func (bi *barlineInteractor) GetRendererRightBarline(measure musicxml.Measure, x int) (barlinePos int, renderer *entity.NoteRenderer) {
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
	barlineRenderer := &entity.NoteRenderer{
		MeasureNumber: measure.Number,
		PositionX:     x,
		Barline:       &barline,
		Width:         BARLINE_AFTER_SPACE,
		// HACK: why there is no width define here?
	}

	// x += BARLINE_AFTER_SPACE
	if barline.Repeat != nil && barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
		x -= 5 // the semicolon offset
	}

	return x, barlineRenderer
}

func (bi *barlineInteractor) RenderBarline(ctx context.Context, canv canvas.Canvas, barline musicxml.Barline, coordinate entity.Coordinate, styles ...string) {
	if _, ok := unicode[barline.BarStyle]; !ok {
		return
	}
	if barline.Repeat != nil {
		var x float64
		y := float64(coordinate.Y) - 3
		switch barline.Repeat.Direction {
		case musicxml.BarLineRepeatDirectionBackward:
			x = coordinate.X - 5
		case musicxml.BarLineRepeatDirectionForward:
			x = coordinate.X + 10
		}

		if x != 0 {
			canv.TextUnescaped(x, y, ":", `style="font-family:Noto Music;font-size:16px"`)
		}
	}

	canv.TextUnescaped(
		coordinate.X, coordinate.Y+6,
		unicode[barline.BarStyle],
		`style="font-family:Noto Music;font-size:28.8px"`,
	)

}

func GetBarlineWidth(barlineStyle musicxml.BarLineStyle) float64 {
	return barlineWidth[barlineStyle]
}
