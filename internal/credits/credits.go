package credits

import (
	"context"
	"fmt"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Credits interface {
	RenderCredits(ctx context.Context, canv canvas.Canvas, y *int, metadata repository.HymnData)
}

type creditsInteractor struct {
	Lyric lyric.Lyric
}

func NewCredits() Credits {
	return &creditsInteractor{
		Lyric: lyric.NewLyric(), // REFACTOR: reuse instance of the lyric
	}
}

// autowrapText wrap the text based on the layout length
// replace the <i> with tspan italic. hence needs to be wrapped by text
// the pair <i> .. </i> maybe separated in the new line, needs to fill the pair tspan manually
// returns:
//   - the breakdown lines, with parse <i> to <tspan>
//   - the length for each line
func autoWrapText(text string, leftIndent int) ([]string, []int) {
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
			length += int(utils.CalculateSecondaryLyricWidth(cleaned))
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
			length += int(utils.CalculateSecondaryLyricWidth(word))
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

func formatAndRenderText(canv canvas.Canvas, y, leftIndent int, text string) []string {
	wrapped, lenLines := autoWrapText(text, leftIndent)
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
		canv.TextUnescaped(float64(constant.LAYOUT_INDENT_LENGTH+leftIndent), float64(y+(i*newLineHeight)), text)
	}

	return wrapped
}

func renderMusicAndLyric(canv canvas.Canvas, y *int, metadata repository.HymnData) (lastLineIndent float64) {
	leftIndent := indentLyric

	lyricMusicMerged := metadata.Lyric == metadata.Music
	if lyricMusicMerged {
		leftIndent = indentMusicAndLyric
	}

	prefix := PREFIX_LYRIC
	if lyricMusicMerged {
		prefix = PREFIX_MERGED_LYRIC_MUSIC
	}
	canv.Text(constant.LAYOUT_INDENT_LENGTH, *y, prefix)

	wrapped := formatAndRenderText(canv, *y, leftIndent, metadata.Lyric)
	*y += newLineHeight * len(wrapped)

	if !lyricMusicMerged {
		canv.Text(constant.LAYOUT_INDENT_LENGTH, *y, PREFIX_MUSIC)
		wrapped = formatAndRenderText(canv, *y, leftIndent, metadata.Music)
		*y += newLineHeight * len(wrapped)
	}

	*y = *y - newLineHeight
	return float64(leftIndent) + utils.CalculateSecondaryLyricWidth(wrapped[len(wrapped)-1])
}

func renderCopyright(canv canvas.Canvas, y *int, leftIndent float64, metadata repository.HymnData) {

	if !metadata.Copyright.Valid {
		return
	}

	copyrightY := *y
	length := utils.CalculateSecondaryLyricWidth(metadata.Copyright.String)
	if constant.LAYOUT_WIDTH-int(leftIndent+length) < constant.LAYOUT_INDENT_LENGTH {
		copyrightY += newLineHeight
		*y = *y + newLineHeight
	}

	canv.Text(constant.LAYOUT_WIDTH-int(length)-constant.LAYOUT_INDENT_LENGTH+constant.UPPERCASE_LENGTH, copyrightY, fmt.Sprintf("© %s", metadata.Copyright.String))
	*y = *y + newLineHeight

}

func renderReferences(canv canvas.Canvas, y int, metadata repository.HymnData) {
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
		l := utils.CalculateSecondaryLyricWidth(ref)
		canv.Text(constant.LAYOUT_WIDTH-constant.UPPERCASE_LENGTH-int(l), y, ref)
	}
}

func (ci *creditsInteractor) RenderCredits(ctx context.Context, canv canvas.Canvas, y *int, metadata repository.HymnData) {
	canv.Group(GROUP_CLASSNAME, GROUP_STYLE)

	lastLineIndent := renderMusicAndLyric(canv, y, metadata)

	renderCopyright(canv, y, lastLineIndent, metadata)
	renderReferences(canv, *y, metadata)

	canv.Gend()

}
