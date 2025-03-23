package credits

import (
	"context"
	"fmt"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Credits interface {
	RenderCredits(ctx context.Context, canv canvas.Canvas, y int, metadata repository.HymnData)
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
	width := nonItalic
	if italic {
		width = italicWidth
	}
	for _, l := range text {
		if italic {
			res += (width[string(l)] * 0.6)
		} else {
			res += width[string(l)]
		}
	}

	return res
}

// autowrapText wrap the text based on the layout length
// replace the <i> with tspan italic. hence needs to be wrapped by text
// the pair <i> .. </i> maybe separated in the new line, needs to fill the pair tspan manually
// returns:
//   - the breakdown lines, with parse <i> to <tspan>
//   - the length for each line
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
		length += spaceWidth
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
			lines = append(lines, strings.Join(result, " "))
			result = []string{}
			lenLines = append(lenLines, length)
			length = 0
		}
	}

	lines = append(lines, strings.Join(result, " "))
	lenLines = append(lenLines, length)
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

// RenderCredits
// TODO: only supports wrapping int the lyric, the music does not have wrapping feature
func (ci *creditsInteractor) RenderCredits(ctx context.Context, canv canvas.Canvas, y int, metadata repository.HymnData) {
	leftIndent := indentLyric
	lyricMusicMerged := metadata.Lyric == metadata.Music
	copyrightY := y
	if lyricMusicMerged {
		leftIndent = indentMusicAndLyric
	}

	wrapped, lenLines := ci.autoWrapText(metadata.Lyric, leftIndent)
	canv.Group("class='credit'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
	if lyricMusicMerged {
		canv.Text(constant.LAYOUT_INDENT_LENGTH, y, fmt.Sprintf("Syair dan lagu : %s", metadata.Lyric))
	} else {

		canv.Text(constant.LAYOUT_INDENT_LENGTH, y, "Syair: ")

		for i, line := range wrapped {
			text := line
			hasBegin := strings.Contains(line, "<tspan font-style=")
			hasEnd := strings.Contains(line, "</tspan>")
			if hasBegin && !hasEnd {
				text = fmt.Sprintf("%s</tspan>", text)
			} else if !hasBegin && hasEnd {
				text = fmt.Sprintf("<tspan font-style=\"italic\">%s", text)
			}
			y += (i * newLineHeight)
			if len(wrapped) > 1 && i < len(wrapped)-1 {
				text = alignText(text, lenLines[i], constant.LAYOUT_WIDTH)
			}
			fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">%s</text>`, constant.LAYOUT_INDENT_LENGTH+leftIndent, y, text)
		}
		copyrightY = y
		y += newLineHeight

		musicCredit := strings.ReplaceAll(metadata.Music, "<i>", "<tspan font-style=\"italic\">")
		musicCredit = strings.ReplaceAll(musicCredit, "</i>", "</tspan>")
		fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">Lagu: %s</text>`, constant.LAYOUT_INDENT_LENGTH, y, musicCredit)

	}

	if metadata.Copyright.Valid {
		length := ci.Lyric.CalculateLyricWidth(metadata.Copyright.String) + constant.UPPERCASE_LENGTH
		canv.Text(constant.LAYOUT_WIDTH-int(length), copyrightY, fmt.Sprintf("© %s", metadata.Copyright.String))

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

}
