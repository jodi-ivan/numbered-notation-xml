package lyric

import (
	"context"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
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
	SetLyricRenderer(noteRenderer *entity.NoteRenderer, rawLyric []musicxml.Lyric) VerseInfo
	CalculateHypen(ctx context.Context, prevLyric, currentLyric *LyricPosition) (location []HyphenPosition)
	RenderHypen(ctx context.Context, y, offsetCenter int, canv canvas.Canvas, measure []*entity.NoteRenderer)
	RenderElision(ctx context.Context, canv canvas.Canvas, text []entity.Text, lyricPart int, pos entity.Coordinate, sty ...string)
	CalculateMarginLeft(txt string) float64
	CalculateOverallWidth(ls []entity.Lyric) float64
	SplitLyricPrefix(note *entity.NoteRenderer, y int, part int, leftBarline *entity.NoteRenderer) []LyricPosition
	RenderLyrics(ctx context.Context, y int, canv canvas.Canvas, measure []*entity.NoteRenderer, prevNote ...*entity.NoteRenderer) int
}

type lyricInteractor struct{}

func NewLyric() Lyric {
	return &lyricInteractor{}
}

// SetLyricRenderer prepares the renderer for lyrics, also calculate space underneath the note and after the note
func (li *lyricInteractor) SetLyricRenderer(noteRenderer *entity.NoteRenderer, rawLyric []musicxml.Lyric) VerseInfo {
	// lyric
	var lyricWidth, noteWidth, marginBottom int
	countedLyric := 0

	if len(rawLyric) > 0 {

		noteRenderer.Lyric = make([]entity.Lyric, len(rawLyric))
		for _, currLyric := range rawLyric {
			lyricText := ""
			l := entity.Lyric{
				Syllabic: currLyric.Syllabic,
				Verse:    currLyric.Verse,
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
			if lyricText != "" {
				countedLyric++
			}

			if currLyric.Number > len(rawLyric) {

				for i := len(rawLyric); i < currLyric.Number; i++ {
					noteRenderer.Lyric = append(noteRenderer.Lyric, entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeMiddle,
						Text:     []entity.Text{{}},
						Verse:    currLyric.Verse,
					})
				}
			}
			noteRenderer.Lyric[currLyric.Number-1] = l
			currWidth := int(math.Round(li.CalculateLyricWidth(lyricText)))

			lyricWidth = int(math.Max(float64(lyricWidth), float64(currWidth)))
		}

	}

	if countedLyric > 0 {
		marginBottom = ((countedLyric - 1) * 25)
	}

	if countedLyric > MAX_VERSE_IN_MUSIC {
		marginBottom += int(countedLyric/MAX_VERSE_IN_MUSIC) * LINE_BETWEEN_LYRIC
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
		HasLyric:     len(rawLyric) > 0,
	}
}

func (li *lyricInteractor) RenderLyrics(ctx context.Context, y int, canv canvas.Canvas, measure []*entity.NoteRenderer, prevNote ...*entity.NoteRenderer) int {

	prefixes := map[string]LyricPosition{}

	offsetCenterVal := 0
	canv.Group("class='lyric'", "style='font-family:Caladea;font-size:16px'")
	var prev *entity.NoteRenderer
	minPrefix := float64(constant.LAYOUT_WIDTH)
	offsetCenter := false
	for notePos, n := range measure {
		yPos := float64(y)
		for i, l := range n.Lyric {
			if len(l.Text) == 0 {
				continue
			}

			xPos := n.PositionX
			v := float64(l.Verse)
			if len(n.Lyric) == 1 {
				v = 0
			}
			yPos = float64(y+DISTANCE_NOTE_TO_LYRIC+(i*LINE_BETWEEN_LYRIC)) + (v * 5)
			prevNoteLen := len(n.Lyric)
			if notePos == 0 && len(prevNote) > 0 {
				prevNoteLen = len(prevNote[0].Lyric)
			} else if notePos > 0 && len(measure[notePos-1].Lyric) > 0 {
				prevNoteLen = len(measure[notePos-1].Lyric)
			}

			if !offsetCenter && prevNoteLen > len(n.Lyric) {
				offsetCenter = true
			}

			if offsetCenter {
				yPos += 10
				offsetCenterVal = 10
			}
			text := l.Text

			prefix := li.SplitLyricPrefix(n, y+(l.Verse*5), i, prev)
			if len(prefix) == 1 {
				text = prefix[0].Lyrics.Text
				xPos, yPos = int(prefix[0].Coordinate.X), prefix[0].Coordinate.Y
			} else if len(prefix) == 2 {

				minPrefix = math.Min(minPrefix, prefix[0].Coordinate.X)
				text = prefix[1].Lyrics.Text
				l.Text = text
				xPos, yPos = int(prefix[1].Coordinate.X), prefix[1].Coordinate.Y

				if len(n.Lyric) > MAX_VERSE_IN_MUSIC {
					prefix[0].Coordinate.Y = yPos + (math.Trunc(float64(i)/MAX_LINE_PER_VERSE_IN_MUSIC) * LINE_BETWEEN_LYRIC)
				}
				prefixes[entity.LyricVal(prefix[0].Lyrics.Text).String()] = prefix[0]
			}

			if n.LeadingHeader != "" && !unicode.IsNumber(rune(n.LeadingHeader[0])) {
				prefixWidth := li.CalculateLyricWidth(n.LeadingHeader)
				if notePos > 0 && measure[notePos-1].Barline != nil {
					prefixWidth += barline.GetBarlineWidth(measure[notePos-1].Barline.BarStyle)
				}
				minPrefix = math.Min(minPrefix, float64(n.PositionX)-prefixWidth-constant.LOWERCASE_LENGTH)
				prefixes[n.LeadingHeader] = LyricPosition{
					Coordinate: entity.NewCoordinate(float64(n.PositionX), float64(y)),
					Lyrics:     entity.Lyric{Text: []entity.Text{{Value: n.LeadingHeader}}, Verse: l.Verse},
				}
			}

			lyricVal := entity.LyricVal(text).String()
			if len(n.Lyric) > MAX_VERSE_IN_MUSIC {
				yPos = yPos + (math.Trunc(float64(i)/MAX_LINE_PER_VERSE_IN_MUSIC) * LINE_BETWEEN_LYRIC)
			}
			if strings.HasPrefix(lyricVal, "*") {
				xPos -= int(li.CalculateLyricWidth("*"))
			}
			if lyricVal == "" {
				n.Lyric[i] = l
				continue
			}
			styles := []string{}
			sty := getColoringStyle(ctx, l.Verse, len(n.Lyric))
			if sty != "" {
				styles = append(styles, sty)
			}
			canv.Text(xPos, int(yPos), lyricVal, styles...)
			elisionOpacity := "stroke-opacity:0.6"
			if sty == "" {
				elisionOpacity = ""
			}
			li.RenderElision(ctx, canv, text, i, entity.Coordinate{X: float64(xPos), Y: yPos}, elisionOpacity)
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
						`font-size="14.4px"`,
					}
					minPrefix += li.CalculateLyricWidth(prefixVal) * 0.1
				}
				sty := getColoringStyle(ctx, p.Lyrics.Verse, p.TotalLyric)
				style = append(style, sty)
				canv.Text(int(minPrefix), int(p.Coordinate.Y), prefixVal, style...)
			}
			prefixes = map[string]LyricPosition{}
			canv.Gend()
		}
		prev = n
	}

	canv.Gend()

	return offsetCenterVal
}
