package lyric

import (
	"context"
	"math"
	"regexp"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

var (
	numberedLyric *regexp.Regexp
	baitPrefix    *regexp.Regexp
)

func init() {
	if numberedLyric == nil {
		numberedLyric, _ = regexp.Compile(`^\d*\.\s{0,1}`)
	}

	if baitPrefix == nil {
		baitPrefix, _ = regexp.Compile(`^(bait)\d*\:`)
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
	RenderLyrics(ctx context.Context, canv canvas.Canvas, measure []*entity.NoteRenderer)
}

type lyricInteractor struct{}

func NewLyric() Lyric {
	return &lyricInteractor{}
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
	}
	noteRenderer.Width += constant.LOWERCASE_LENGTH / 2 // lyric padding

	return VerseInfo{
		MarginBottom: marginBottom,
	}
}

func (li *lyricInteractor) RenderLyrics(ctx context.Context, canv canvas.Canvas, measure []*entity.NoteRenderer) {
	prefixes := map[string]LyricPosition{}

	canv.Group("class='lyric'", "style='font-family:Caladea'")
	var prev *entity.NoteRenderer
	minPrefix := float64(constant.LAYOUT_WIDTH)

	for _, n := range measure {
		yPos := float64(n.PositionY)
		for i, l := range n.Lyric {
			if len(l.Text) == 0 {
				continue
			}

			xPos := n.PositionX
			yPos = float64(n.PositionY + DISTANCE_NOTE_TO_LYRIC + (i * LINE_BETWEEN_LYRIC))
			text := l.Text

			prefix := li.SplitLyricPrefix(n, i, prev)
			if len(prefix) == 1 {
				text = prefix[0].Lyrics.Text
				xPos, yPos = int(prefix[0].Coordinate.X), prefix[0].Coordinate.Y
			} else if len(prefix) == 2 {
				if n.LeadingHeader != "" {
					prefixWidth := li.CalculateLyricWidth(n.LeadingHeader)
					minPrefix = math.Min(minPrefix, prefix[0].Coordinate.X-prefixWidth)
					prefixes[n.LeadingHeader] = LyricPosition{
						Coordinate: entity.NewCoordinate(float64(n.PositionX), float64(n.PositionY)),
						Lyrics:     entity.Lyric{Text: []entity.Text{{Value: n.LeadingHeader}}},
					}
				}
				minPrefix = math.Min(minPrefix, prefix[0].Coordinate.X)
				text = prefix[1].Lyrics.Text
				l.Text = text
				xPos, yPos = int(prefix[1].Coordinate.X), prefix[1].Coordinate.Y

				if len(n.Lyric) > MAX_VERSE_IN_MUSIC {
					prefix[0].Coordinate.Y = yPos + (math.Trunc(float64(i)/MAX_LINE_PER_VERSE_IN_MUSIC) * LINE_BETWEEN_LYRIC)
				}
				prefixes[entity.LyricVal(prefix[0].Lyrics.Text).String()] = prefix[0]
			}

			lyricVal := entity.LyricVal(text).String()
			if len(n.Lyric) > MAX_VERSE_IN_MUSIC {
				yPos = yPos + (math.Trunc(float64(i)/MAX_LINE_PER_VERSE_IN_MUSIC) * LINE_BETWEEN_LYRIC)
			}
			if strings.HasPrefix(lyricVal, "*") {
				xPos -= int(li.CalculateLyricWidth("*"))
			}
			canv.Text(xPos, int(yPos), lyricVal)
			li.RenderElision(ctx, canv, text, i, entity.Coordinate{X: float64(xPos), Y: yPos})
			n.Lyric[i] = l
		}

		if len(prefixes) > 0 {
			canv.Group("class='leading-header'")
			for _, p := range prefixes {
				prefixVal := entity.LyricVal(p.Lyrics.Text).String()
				style := []string{}
				if p.Lyrics.Text[0].Italic {
					style = []string{
						`font-style="italic"`,
						`font-size="90%"`,
					}
					minPrefix += li.CalculateLyricWidth(prefixVal) * 0.1
				}
				canv.Text(int(minPrefix), int(p.Coordinate.Y), prefixVal, style...)
			}
			prefixes = map[string]LyricPosition{}
			canv.Gend()
		}
		prev = n
	}

	canv.Gend()
}
