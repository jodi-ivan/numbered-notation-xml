package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

func getTotalChar(note *entity.NoteRenderer) int {
	if len(note.Lyric) == 0 {
		return 0
	}
	total := len(entity.LyricVal(note.Lyric[0].Text).String())

	if note.Lyric[0].Syllabic == musicxml.LyricSyllabicTypeEnd || note.Lyric[0].Syllabic == musicxml.LyricSyllabicTypeSingle {
		total++
	}
	return total
}

func Highligts(ctx context.Context, notes []*entity.NoteRenderer) {
	param, ok := params.GetParamFromContext(ctx)
	if !ok || param.Hightliter == nil {
		return
	}

	spans := param.Hightliter["1"]
	if len(spans) == 0 {
		return
	}

	noteIndex := 0

	for _, span := range spans {
		start, end := span[0], span[1]

		currentWord := 0

		// Move forward until we reach the first note that
		// could overlap this span.
		for noteIndex < len(notes) {
			note := notes[noteIndex]
			noteLength := getTotalChar(note)

			noteStart := param.HightliterCtr
			noteEnd := param.HightliterCtr + noteLength

			// This note is completely before the span.
			if noteEnd <= start {
				param.HightliterCtr = noteEnd
				noteIndex++
				continue
			}

			// This note starts at or after the end of the span.
			// Since spans are ordered, we're done with this span.
			if noteStart >= end {
				break
			}

			// We have an overlap.
			if len(note.Lyric) > 0 {
				endSubSpan := currentWord
				if endSubSpan == 0 {
					endSubSpan += noteLength
				}
				note.Lyric[0].Highlights = params.Span{0, end - start - endSubSpan}

				// log.Printf(
				// 	"=================FOUND INTERSECTION span [%d,%d) note [%d,%d) lyric=%q",
				// 	start,
				// 	end,
				// 	param.HightliterCtr,
				// 	noteEnd,
				// 	entity.LyricVal(note.Lyric[0].Text).String(),
				// )
				currentWord += noteLength

			}

			// Move to the next note.
			param.HightliterCtr = noteEnd
			noteIndex++
		}
	}
}
