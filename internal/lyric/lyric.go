package lyric

import (
	"context"
	"math"
	"regexp"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
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
	RenderVerse(ctx context.Context, canv canvas.Canvas, y int, verses []repository.HymnVerse) VerseInfo
	SetLyricRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note) VerseInfo
}

type lyricInteractor struct{}

func (li *lyricInteractor) CalculateLyricWidth(txt string) float64 {
	return CalculateLyricWidth(txt)
}

func (li *lyricInteractor) RenderVerse(ctx context.Context, canv canvas.Canvas, y int, verses []repository.HymnVerse) VerseInfo {
	return RenderVerse(ctx, canv, y, verses)
}

func (li *lyricInteractor) SetLyricRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note) VerseInfo {
	// lyric
	var lyricWidth, noteWidth, marginBottom int

	if len(note.Lyric) > 0 {
		marginBottom = ((len(note.Lyric) - 1) * 25)
		noteRenderer.Lyric = make([]entity.Lyric, len(note.Lyric))
		for i, currLyric := range note.Lyric {
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

			noteRenderer.Lyric[i] = l
			currWidth := int(math.Round(li.CalculateLyricWidth(lyricText)))
			if currLyric.Syllabic == musicxml.LyricSyllabicTypeEnd || currLyric.Syllabic == musicxml.LyricSyllabicTypeSingle {
				currWidth += constant.LOWERCASE_LENGTH
			}
			currWidth += 4 // lyric padding

			lyricWidth = int(math.Max(float64(lyricWidth), float64(currWidth)))
		}

	}

	noteWidth = constant.LOWERCASE_LENGTH

	if noteWidth > lyricWidth {
		noteRenderer.Width = noteWidth
		noteRenderer.IsLengthTakenFromLyric = false
	} else {
		noteRenderer.Width = lyricWidth
		noteRenderer.IsLengthTakenFromLyric = true
		if float64(lyricWidth) > float64(noteWidth+constant.UPPERCASE_LENGTH) {
			noteRenderer.Width = constant.UPPERCASE_LENGTH * 1.7
		}
	}

	return VerseInfo{
		MarginBottom: marginBottom,
	}
}

func NewLyric() Lyric {
	return &lyricInteractor{}
}

func CalculateMarginLeft(txt string) float64 {
	if numberedLyric.Match([]byte(txt)) {
		subStr := numberedLyric.FindStringSubmatch(txt)
		if len(subStr) == 0 {
			return 0
		}

		return CalculateLyricWidth(subStr[0]) * -1
	}
	return 0
}

func CalculateLyricWidth(txt string) float64 {
	// Only for Caladea font
	width := map[string]float64{
		"1": 9.28,
		"2": 7.55,
		"3": 7.43,
		"4": 8.57,
		"5": 7.61,
		"6": 7.53,
		"7": 7.53,
		"8": 8,
		"9": 7.65,
		"0": 8.57,
		"A": 9.59,
		"B": 9.27,
		"C": 8.1,
		"D": 10,
		"E": 8.65,
		"F": 8.15,
		"G": 8.63,
		"H": 11.15,
		"I": 5.49,
		"J": 4.99,
		"K": 10.08,
		"L": 8.02,
		"M": 14.21,
		"N": 11.09,
		"O": 9.59,
		"P": 8.53,
		"Q": 9.59,
		"R": 9.81,
		"S": 7.25,
		"T": 8.92,
		"U": 11,
		"V": 9.57,
		"W": 14.23,
		"X": 9.95,
		"Y": 8.92,
		"Z": 8.11,
		"a": 7.52,
		"b": 8.32,
		"c": 6.74,
		"d": 8.32,
		"e": 7.06,
		"f": 5.87,
		"g": 7.35,
		"h": 8.86,
		"i": 4.44,
		"j": 4.76,
		"k": 8.43,
		"l": 4.34,
		"m": 13.01,
		"n": 8.94,
		"o": 7.69,
		"p": 8.32,
		"q": 8.02,
		"r": 6.34,
		"s": 6.28,
		"t": 5.21,
		"u": 8.74,
		"v": 8.08,
		"w": 12.08,
		"x": 7.78,
		"y": 8.18,
		"z": 6.85,
		",": 3.28,
		"'": 3.28,
		".": 3.3,
		"!": 4.58,
		";": 4.23,
		" ": 4,
		"-": 5.27,
	}
	res := 0.0

	for _, l := range txt {
		res += width[string(l)]
	}

	return res
}
