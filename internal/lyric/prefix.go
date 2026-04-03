package lyric

import (
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
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
	if note.PositionX == constant.LAYOUT_INDENT_LENGTH || (leftBarline != nil && leftBarline.Barline != nil) {
		xPos := note.PositionX

		hasLeftBarline := leftBarline != nil && leftBarline.Barline != nil
		if len(note.Lyric) > 2 {

			if hasLeftBarline {
				xPos -= leftBarline.Width - barline.BARLINE_AFTER_SPACE + 5
			}

			parts := strings.Split(lyricVal, ".")
			if len(parts) == 2 {

				header := LyricPosition{
					Coordinate: entity.NewCoordinate(float64(xPos)-li.CalculateLyricWidth(parts[0]+"."), float64(note.PositionY+25+(part*LINE_BETWEEN_LYRIC))),
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeSingle,
						Text:     []entity.Text{{Value: parts[0] + "."}},
					},
				}

				mainLyric := LyricPosition{
					Coordinate: entity.NewCoordinate(float64(xPos), float64(note.PositionY+25+(part*LINE_BETWEEN_LYRIC))),
					Lyrics:     note.Lyric[part],
				}
				mainLyric.Lyrics.Text[0] = entity.Text{Value: strings.ReplaceAll(mainLyric.Lyrics.Text[0].Value, parts[0]+".", "")}
				return []LyricPosition{header, mainLyric}
			}
		}

		xPos += int(li.CalculateMarginLeft(lyricVal))
		return []LyricPosition{
			{
				Coordinate: entity.NewCoordinate(float64(xPos), float64(note.PositionY+25+(part*LINE_BETWEEN_LYRIC))),
				Lyrics:     note.Lyric[part],
			},
		}
	}

	return nil
}
