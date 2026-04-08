package lyric

import (
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func (li *lyricInteractor) CalculateMarginLeft(txt string) float64 {
	if numberedLyric.Match([]byte(txt)) {
		subStr := numberedLyric.FindStringSubmatch(txt)
		if len(subStr) == 0 {
			return 0
		}

		return li.CalculateLyricWidth(strings.Join(subStr, "")) * -1
	}
	return 0
}

func (li *lyricInteractor) SplitLyricPrefix(note *entity.NoteRenderer, part int, leftBarline *entity.NoteRenderer) []LyricPosition {
	lyricVal := entity.LyricVal(note.Lyric[part].Text).String()

	if !numberedLyric.Match([]byte(lyricVal)) {
		return nil
	}
	xPos := float64(note.PositionX)
	yPos := float64(note.PositionY + 25 + (part * LINE_BETWEEN_LYRIC))

	parts := strings.Split(lyricVal, ".")
	if len(parts) == 2 {

		header := LyricPosition{
			Coordinate: entity.NewCoordinate(xPos-li.CalculateLyricWidth(parts[0]+"."), yPos),
			Lyrics: entity.Lyric{
				Syllabic: musicxml.LyricSyllabicTypeSingle,
				Text:     []entity.Text{{Value: parts[0] + "."}},
			},
		}

		mainLyric := LyricPosition{
			Coordinate: entity.NewCoordinate(xPos, yPos),
			Lyrics:     note.Lyric[part],
		}
		mainLyric.Lyrics.Text[0] = entity.Text{Value: strings.ReplaceAll(mainLyric.Lyrics.Text[0].Value, parts[0]+".", "")}
		return []LyricPosition{header, mainLyric}
	}
	xPos += li.CalculateMarginLeft(lyricVal)
	return []LyricPosition{
		{
			Coordinate: entity.NewCoordinate(xPos, yPos),
			Lyrics:     note.Lyric[part],
		},
	}

}
