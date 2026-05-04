package lyric

import (
	"fmt"
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

func (li *lyricInteractor) SplitLyricPrefix(note *entity.NoteRenderer, y int, part int, leftBarline *entity.NoteRenderer) []LyricPosition {
	lyricVal := entity.LyricVal(note.Lyric[part].Text).String()

	xPos := float64(note.PositionX)
	yPos := float64(y + 25 + (part * LINE_BETWEEN_LYRIC))
	parts := []string{}
	needStylized := false
	if baitPrefix.Match([]byte(lyricVal)) {
		lead := strings.Split(lyricVal, ":")
		parts = []string{
			fmt.Sprintf("bait %s: ", strings.TrimPrefix(lead[0], "bait")),
			lead[1],
		}
		needStylized = true
	} else {
		if !numberedLyric.Match([]byte(lyricVal)) {
			return nil
		}
		parts = strings.Split(lyricVal, ".")
		parts[0] += "."
	}

	if len(parts) == 2 {

		text := entity.Text{Value: parts[0]}
		if needStylized {
			text.Bold = true
			text.Italic = true
		}
		header := LyricPosition{
			Coordinate: entity.NewCoordinate(xPos-li.CalculateLyricWidth(parts[0]+"."), yPos),
			Lyrics: entity.Lyric{
				Syllabic: musicxml.LyricSyllabicTypeSingle,
				Text:     []entity.Text{text},
			},
		}

		mainLyric := LyricPosition{
			Coordinate: entity.NewCoordinate(xPos, yPos),
			Lyrics: entity.Lyric{
				Syllabic: musicxml.LyricSyllabicTypeBegin,
				Text:     []entity.Text{{Value: parts[1]}},
			},
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
