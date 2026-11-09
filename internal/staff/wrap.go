package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
)

func Wrap(ctx context.Context, x, maxPosX int, ks keysig.KeySignature, oneline []*entity.NoteRenderer) [][]*entity.NoteRenderer {
	// 1. Group consecutive notes into measures by MeasureNumber.
	var measures [][]*entity.NoteRenderer
	for _, n := range oneline {
		if n == nil {
			continue
		}
		last := len(measures) - 1
		if last >= 0 && measures[last][0].MeasureNumber == n.MeasureNumber {
			measures[last] = append(measures[last], n)
		} else {
			measures = append(measures, []*entity.NoteRenderer{n})
		}
	}

	// 2. Wrap whole measures into lines.
	var result [][]*entity.NoteRenderer
	var current []*entity.NoteRenderer
	posX := 0
	lineNum := 1
	// flush finalizes the current line and marks its last note as a line end.
	flush := func() {
		if len(current) == 0 {
			return
		}
		lineNum++
		current[len(current)-1].IsNewLine = true
		result = append(result, current)
		ri := rhythm.New(nil)
		ri.AdjustMultiDottedRenderer(ctx, current, x, (lineNum-1)*70, ks)

		current = nil
		posX = 0
	}

	for _, measure := range measures {

		mw := 0
		for _, n := range measure {
			mw += int(n.Width)
		}

		// The measure doesn't fit: move it entirely to a new line.
		if posX+mw > maxPosX {
			flush()
		}

		current = append(current, measure...)
		posX += mw
	}

	flush()

	return result
}
