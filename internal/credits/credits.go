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
		wordLength := int(utils.CalculateSecondaryLyricWidth(word))
		if length+wordLength < available && strings.HasPrefix(word, utils.ITALIC_OPENING) && strings.HasSuffix(word, utils.ITALIC_CLOSING) {
			result = append(result, utils.ReplaceItalicToSpan(word))
			continue
		}
		if strings.HasPrefix(word, utils.ITALIC_OPENING) || italic {
			cleaned := strings.TrimSuffix(word, utils.ITALIC_CLOSING)
			cleaned = strings.TrimPrefix(cleaned, utils.ITALIC_OPENING)
			length += int(utils.CalculateSecondaryLyricWidth(cleaned))
			if !italic {
				result = append(result, fmt.Sprintf("%s%s", utils.TSPAN_OPENING, cleaned))
			} else {

				if !strings.HasSuffix(word, utils.ITALIC_CLOSING) {
					result = append(result, word)
				}
			}
			if !strings.HasSuffix(word, utils.ITALIC_CLOSING) {
				italic = true
			} else {
				italic = false
				result = append(result, fmt.Sprintf("%s%s", cleaned, utils.TSPAN_CLOSING))
			}
		} else {
			length += wordLength
			result = append(result, word)
		}

		if length >= available {
			if italic {
				result = append(result, utils.TSPAN_CLOSING)
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
	escaped := "tspan-font-style"
	unescaped := "tspan font-style"
	// clean the tag
	text = strings.ReplaceAll(text, unescaped, escaped)
	words := strings.Fields(text)
	spaceLeft := targetLength - textLength
	if len(words) > 2 && (spaceLeft > (len(words)-2)*spaceWidth) {
		text = strings.Join(words, strings.Repeat("&#160;", (spaceLeft/(len(words)-2))))
	}
	return strings.ReplaceAll(text, escaped, unescaped)
}

func formatAndRenderText(canv canvas.Canvas, y, leftIndent int, text string) []string {
	wrapped, lenLines := autoWrapText(text, leftIndent)
	for i, line := range wrapped {
		text := line
		hasBegin := strings.Contains(line, TSPAN_CONTAINS_CHECK)
		hasEnd := strings.Contains(line, utils.TSPAN_CLOSING)
		if hasBegin && !hasEnd {
			text = fmt.Sprintf("%s %s", text, utils.TSPAN_CLOSING)
		} else if !hasBegin && hasEnd {
			text = fmt.Sprintf("%s %s", utils.TSPAN_OPENING, text)
		}
		if len(wrapped) > 1 && i < len(wrapped)-1 {
			text = alignText(text, lenLines[i], constant.LAYOUT_WIDTH-(constant.LAYOUT_INDENT_LENGTH*2))
		}
		canv.TextUnescaped(float64(constant.LAYOUT_INDENT_LENGTH+leftIndent), float64(y+(i*newLineHeight)), text)
		wrapped[i] = text
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
	lastLine := wrapped[len(wrapped)-1]
	if strings.Contains(lastLine, TSPAN_CONTAINS_CHECK) {
		lastLine = utils.CleanSpan(lastLine)
	}

	return float64(leftIndent) + utils.CalculateSecondaryLyricWidth(lastLine)
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

	if metadata.Lyric == "" {
		return
	}
	canv.Group(GROUP_CLASSNAME, GROUP_STYLE)

	lastLineIndent := renderMusicAndLyric(canv, y, metadata)

	renderCopyright(canv, y, lastLineIndent, metadata)
	renderReferences(canv, *y, metadata)

	canv.Gend()

}
