package credits

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Credits interface {
	RenderCredits(ctx context.Context, canv canvas.Canvas, y int, metadata repository.HymnData, verseFootnotes map[int]map[int]repository.VerseFootNotes) int
	RenderForKidsFootnotes(ctx context.Context, canv canvas.Canvas, y int)
	RenderMuiscFootnotes(ctx context.Context, canv canvas.Canvas, metadata *repository.HymnMetadata, x, y int)
}

type creditsInteractor struct {
	Lyric lyric.Lyric
}

func NewCredits() Credits {
	return &creditsInteractor{
		Lyric: lyric.NewLyric(), // REFACTOR: reuse instance of the lyric
	}
}

func CalculateLyric(text string, italic bool) float64 {
	res := 0.0
	for _, l := range text {
		w, ok := charWidth[string(l)]
		if !ok {
			w = 6
		}
		res += w
	}

	return res
}

// autowrapText wrap the text based on the layout length
// replace the <i> with tspan italic. hence needs to be wrapped by text
// the pair <i> .. </i> maybe separated in the new line, needs to fill the pair tspan manually
// returns:
//   - the breakdown lines, with parse <i> to <tspan>
//   - the length for each line
//
// FIXME: error non character after </i> without spaces. ie. `</i>,` will return parsing error
func (ci *creditsInteractor) autoWrapText(text string, leftIndent int) ([]string, []int) {
	full := strings.Fields(text)
	result := []string{}
	length := 0

	lines := []string{}
	lenLines := []int{}
	available := constant.LAYOUT_WIDTH - (leftIndent + constant.LAYOUT_INDENT_LENGTH)

	italic := false
	for _, word := range full {
		word = strings.TrimSpace(word)
		length += spaceWidth + 4
		if strings.HasPrefix(word, "<i>") || italic {
			cleaned := strings.TrimSuffix(word, "</i>")
			cleaned = strings.TrimPrefix(cleaned, "<i>")
			length += int(CalculateLyric(cleaned, true))
			if !italic {
				result = append(result, fmt.Sprintf("<tspan font-style=\"italic\">%s", cleaned))
			} else {

				if !strings.HasSuffix(word, "</i>") {
					result = append(result, word)
				}
			}
			if !strings.HasSuffix(word, "</i>") {
				italic = true
			} else {
				italic = false
				result = append(result, fmt.Sprintf("%s</tspan>", cleaned))
			}
		} else {
			length += int(CalculateLyric(word, false))
			result = append(result, word)
		}

		if length >= available {
			if italic {
				result = append(result, "</tspan>")
			}

			lines = append(lines, strings.Join(result, " "))
			result = []string{}
			lenLines = append(lenLines, length)
			length = 0
		}
	}

	if len(result) > 0 {
		lines = append(lines, strings.Join(result, " "))
		lenLines = append(lenLines, length)
	}
	return lines, lenLines
}

// UNTESTED: needs more cases
// this will FORCE the text to be aligned
// added spaces between words, with assumed that the length of the space is 2px
func alignText(text string, textLength, targetLength int) string {
	// clean the tag
	text = strings.ReplaceAll(text, "tspan font-style", "tspan-font-style")
	words := strings.Fields(text)
	spaceLeft := targetLength - textLength
	if len(words) > 2 && (spaceLeft > (len(words)-2)*spaceWidth) {
		text = strings.Join(words, strings.Repeat("&#160;", (spaceLeft/(len(words)-2))))
	}
	return strings.ReplaceAll(text, "tspan-font-style", "tspan font-style")
}

func (ci *creditsInteractor) RenderCredits(ctx context.Context, canv canvas.Canvas, y int, metadata repository.HymnData, verseFootnotes map[int]map[int]repository.VerseFootNotes) int {
	leftIndent := indentLyric
	lyricMusicMerged := metadata.Lyric == metadata.Music
	copyrightY := y
	if lyricMusicMerged {
		leftIndent = indentMusicAndLyric
	}

	if len(verseFootnotes) > 0 {
		flatten := []repository.VerseFootNotes{}
		versenoteHeadonlyCnt := 0
		hasInternalItalic := false
		for _, fn := range verseFootnotes {
			for _, t := range fn {
				hasInternalItalic = hasInternalItalic || strings.Contains(t.FootnoteMarker.String, "<i>") || strings.Contains(t.Footnote.String, "<i>")
				if lyric.VerseNoteStyle(t.MarkerStyle.Int32) == lyric.VerseNoteStyleHeadonly {
					versenoteHeadonlyCnt++
					continue
				}

				flatten = append(flatten, t)
			}
		}

		if versenoteHeadonlyCnt != len(flatten) {

			// Sort the footnotes by its markers
			sort.Slice(flatten, func(i, j int) bool {
				return flatten[i].FootnoteMarker.String < flatten[j].FootnoteMarker.String
			})

			y -= 10

			footnotesStyle := ""
			if !hasInternalItalic {
				footnotesStyle = "font-style:italic"
			}
			canv.Group("class='footnotes'", fmt.Sprintf(`style="font-size:60%%;font-family:'Figtree';font-weight:600;%s"`, footnotesStyle))
			totalLine := 0
			for i, fn := range flatten {
				lines := strings.Split(fn.Footnote.String, "<br/>")
				if len(lines) >= 2 {
					marker := strings.ReplaceAll(fn.FootnoteMarker.String, "<i>", "<tspan font-style=\"italic\">")
					marker = strings.ReplaceAll(marker, "</i>", "</tspan>")

					clean := strings.ReplaceAll(fn.FootnoteMarker.String, "<i>", "")
					clean = strings.ReplaceAll(clean, "</i>", "")

					xNotes := int(CalculateLyric(clean, true))
					fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">%s</text>`,
						constant.LAYOUT_INDENT_LENGTH+20, (15*i)+y,
						marker,
					)
					for li, line := range lines {
						line = strings.ReplaceAll(line, "<i>", "<tspan font-style=\"italic\">")
						line = strings.ReplaceAll(line, "</i>", "</tspan>")
						fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">%s</text>`,
							constant.LAYOUT_INDENT_LENGTH+20+xNotes, (15*(i+li))+y,
							line,
						)
					}
					y += (15 * (i + len(lines)))
				} else {
					totalLine++
					canv.Text(constant.LAYOUT_INDENT_LENGTH+20, (15*i)+y, fn.FootnoteMarker.String+fn.Footnote.String)
				}
			}
			canv.Gend()
			y += 15 + (totalLine * 15)
		}
	}
	wrapped, lenLines := ci.autoWrapText(metadata.Lyric, leftIndent)
	canv.Group("class='credit'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)

	prefix := "Syair: "
	if lyricMusicMerged {
		prefix = "Syair dan lagu :"
	}
	canv.Text(constant.LAYOUT_INDENT_LENGTH, y, prefix)

	for i, line := range wrapped {
		text := line
		hasBegin := strings.Contains(line, "<tspan font-style=")
		hasEnd := strings.Contains(line, "</tspan>")
		if hasBegin && !hasEnd {
			text = fmt.Sprintf("%s</tspan>", text)
		} else if !hasBegin && hasEnd {
			text = fmt.Sprintf("<tspan font-style=\"italic\">%s", text)
		}
		if len(wrapped) > 1 && i < len(wrapped)-1 {
			text = alignText(text, lenLines[i], constant.LAYOUT_WIDTH-(constant.LAYOUT_INDENT_LENGTH*2))
		}
		fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">%s</text>`, constant.LAYOUT_INDENT_LENGTH+leftIndent, y, text)
		y += newLineHeight
	}
	copyrightY = y

	if !lyricMusicMerged {
		musicCredit := strings.ReplaceAll(metadata.Music, "<i>", "<tspan font-style=\"italic\">")
		musicCredit = strings.ReplaceAll(musicCredit, "</i>", "</tspan>")
		fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">Lagu: %s</text>`, constant.LAYOUT_INDENT_LENGTH, y, musicCredit)
	}

	if metadata.Copyright.Valid {
		length := CalculateLyric(metadata.Copyright.String, false)
		lastMusicLen := CalculateLyric(wrapped[len(wrapped)-1], false)
		if (constant.LAYOUT_WIDTH - (leftIndent + int(lastMusicLen) + int(length))) < constant.LAYOUT_INDENT_LENGTH {
			copyrightY += newLineHeight
			y += newLineHeight

		}

		canv.Text(constant.LAYOUT_WIDTH-int(length)-constant.LAYOUT_INDENT_LENGTH+constant.UPPERCASE_LENGTH, copyrightY, fmt.Sprintf("© %s", metadata.Copyright.String))
		y += newLineHeight
	}

	ref := ""
	if metadata.RefBE.Valid {
		ref += fmt.Sprintf("BE %d", metadata.RefBE.Int16)
	}

	if metadata.RefNR.Valid {
		if ref != "" {
			ref += ", "
		}
		ref += fmt.Sprintf("NR %d", metadata.RefNR.Int16)

	}

	if ref != "" {
		l := CalculateLyric(ref, false)
		canv.Text(constant.LAYOUT_WIDTH-constant.UPPERCASE_LENGTH-int(l), y, ref)
	}

	if metadata.TitleFootnotes.Valid {
		notes := "<tspan font-style=\"italic\">*  " + metadata.TitleFootnotes.String + "</tspan>"
		y += 30
		fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">%s</text>`, constant.LAYOUT_INDENT_LENGTH, y, notes)

	}

	canv.Gend()

	return y
}

func (ci *creditsInteractor) RenderForKidsFootnotes(ctx context.Context, canv canvas.Canvas, y int) {
	canv.Group("class='credit'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
	fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">
				<tspan font-style="italic">Semua nyayian dengan tanda</tspan>
				<tspan font-style="bold" font-size="125%%">☆</tspan>
				<tspan font-style="italic">: khusus untuk anak-anak</tspan>
			</text>`, constant.LAYOUT_INDENT_LENGTH, y)
	canv.Gend()
}
