package lyric

import (
	"context"
	"math"
	"regexp"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

var numberedLyric *regexp.Regexp

func init() {
	if numberedLyric == nil {
		numberedLyric, _ = regexp.Compile(`^\d*\.\s{0,1}`)
	}
}

type Lyric interface {
	CalculateLyricWidth(string) float64
	SetLyricRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note) VerseInfo
	CalculateHypen(ctx context.Context, prevLyric, currentLyric *LyricPosition) (location []entity.Coordinate)
	RenderHypen(ctx context.Context, canv canvas.Canvas, measure []*entity.NoteRenderer)
	RenderElision(ctx context.Context, canv canvas.Canvas, text []entity.Text, lyricPart int, pos entity.Coordinate)
	CalculateMarginLeft(txt string) float64
	CalculateOverallWidth(ls []entity.Lyric) float64
	SplitLyricPrefix(note *entity.NoteRenderer, part int, leftBarline *entity.NoteRenderer) []LyricPosition
}

type lyricInteractor struct{}

func NewLyric() Lyric {
	return &lyricInteractor{}
}

func (li *lyricInteractor) CalculateOverallWidth(ls []entity.Lyric) float64 {
	result := 0.0

	for _, v := range ls {
		result = math.Max(result, li.CalculateLyricWidth(entity.LyricVal(v.Text).String()))
	}
	return result
}

func (li *lyricInteractor) CalculateLyricWidth(txt string) float64 {
	// Only for Caladea font
	width := map[string]float64{
		"A": 9.59, "a": 7.52, "1": 9.28,
		"B": 9.27, "b": 8.32, "2": 7.55,
		"C": 8.1, "c": 6.74, "3": 7.43,
		"D": 10, "d": 8.32, "4": 8.57,
		"E": 8.65, "e": 7.06, "5": 7.61,
		"F": 8.15, "f": 5.87, "6": 7.53,
		"G": 8.63, "g": 7.35, "7": 7.53,
		"H": 11.15, "h": 8.86, "8": 8,
		"I": 5.49, "i": 4.44, "9": 7.65,
		"J": 4.99, "j": 4.76, "0": 8.57,
		"K": 10.08, "k": 8.43, ",": 3.28,
		"L": 8.02, "l": 4.34, "'": 3.28,
		"M": 14.21, "m": 13.01, ".": 3.3,
		"N": 11.09, "n": 8.94, "!": 4.58,
		"O": 9.59, "o": 7.69, ";": 4.23,
		"P": 8.53, "p": 8.32, " ": 4,
		"Q": 9.59, "q": 8.02, "-": 5.27,
		"R": 9.81, "r": 6.34, "—": 16,
		"S": 7.25, "s": 6.28, "*": 6.83,
		"T": 8.92, "t": 5.21, "\"": 8,
		"U": 11, "u": 8.74, "+": 8.86,
		"V": 9.57, "v": 8.08,
		"W": 14.23, "w": 12.08,
		"X": 9.95, "x": 7.78,
		"Y": 8.92, "y": 8.18,
		"Z": 8.11, "z": 6.85,
	}
	res := 0.0

	for _, l := range txt {
		res += width[string(l)]
	}

	return res
}

// SetLyricRenderer prepares the renderer for lyrics, also calculate space underneath the note and after the note
func (li *lyricInteractor) SetLyricRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note) VerseInfo {
	// lyric
	var lyricWidth, noteWidth, marginBottom int

	if len(note.Lyric) > 0 {
		marginBottom = ((len(note.Lyric) - 1) * 25)
		if len(note.Lyric) > MAX_VERSE_IN_MUSIC {
			marginBottom += int(len(note.Lyric)/MAX_VERSE_IN_MUSIC) * LINE_BETWEEN_LYRIC
		}

		noteRenderer.Lyric = make([]entity.Lyric, len(note.Lyric))
		for _, currLyric := range note.Lyric {
			lyricText := ""
			l := entity.Lyric{
				Syllabic: currLyric.Syllabic,
			}

			texts := []entity.Text{}
			for _, t := range currLyric.Text {
				lyricText += t.Value
				texts = append(texts, entity.Text{
					Value:     t.Value,
					Underline: t.Underline,
				})
			}

			l.Text = texts
			if currLyric.Number > len(note.Lyric) {

				for i := len(note.Lyric); i < currLyric.Number; i++ {
					noteRenderer.Lyric = append(noteRenderer.Lyric, entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeMiddle,
						Text:     []entity.Text{{}},
					})
				}
			}
			noteRenderer.Lyric[currLyric.Number-1] = l
			currWidth := int(math.Round(li.CalculateLyricWidth(lyricText)))
			if currLyric.Syllabic == musicxml.LyricSyllabicTypeEnd || currLyric.Syllabic == musicxml.LyricSyllabicTypeSingle {
				currWidth += constant.LOWERCASE_LENGTH
			}

			lyricWidth = int(math.Max(float64(lyricWidth), float64(currWidth)))
		}

	}

	noteWidth = constant.LOWERCASE_LENGTH

	if noteWidth > lyricWidth {
		noteRenderer.Width = noteWidth + (constant.LOWERCASE_LENGTH / 2)
		noteRenderer.IsLengthTakenFromLyric = false
	} else {
		noteRenderer.Width = lyricWidth + 6
		noteRenderer.IsLengthTakenFromLyric = true
		// if float64(lyricWidth) < float64(noteWidth+constant.UPPERCASE_LENGTH) {
		// 	noteRenderer.Width = constant.UPPERCASE_LENGTH * 1.7
		// }
	}
	noteRenderer.Width += constant.LOWERCASE_LENGTH / 2 // lyric padding

	return VerseInfo{
		MarginBottom: marginBottom,
	}
}

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
					Coordinate: entity.Coordinate{
						X: float64(xPos) - li.CalculateLyricWidth(parts[0]+"."),
						Y: float64(note.PositionY + 25 + (part * LINE_BETWEEN_LYRIC)),
					},
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeSingle,
						Text: []entity.Text{
							entity.Text{
								Value: parts[0] + ".",
							},
						},
					},
				}

				mainLyric := LyricPosition{
					Coordinate: entity.Coordinate{
						X: float64(xPos),
						Y: float64(note.PositionY + 25 + (part * LINE_BETWEEN_LYRIC)),
					},
					Lyrics: note.Lyric[part],
				}
				mainLyric.Lyrics.Text[0] = entity.Text{Value: strings.ReplaceAll(mainLyric.Lyrics.Text[0].Value, parts[0]+".", "")}

				return []LyricPosition{header, mainLyric}
			}
		}

		xPos += int(li.CalculateMarginLeft(lyricVal))
		return []LyricPosition{
			LyricPosition{
				Coordinate: entity.Coordinate{
					X: float64(xPos),
					Y: float64(note.PositionY + 25 + (part * LINE_BETWEEN_LYRIC)),
				},
				Lyrics: note.Lyric[part],
			},
		}
	}

	return nil
}
