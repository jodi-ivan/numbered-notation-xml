package footnote

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

var li = lyric.NewLyric()

type VerseLineCursor struct {
	VerseNo    int
	LinePos    int
	Leftmargin int
	LineText   string
}

func AssignFootnotesMarker(canv canvas.Canvas, pos entity.Coordinate, defaultX int, cursor VerseLineCursor, verseFootnote map[int]map[int]repository.VerseFootNotes) {

	if footnotes, hasFootnotes := verseFootnote[cursor.VerseNo]; hasFootnotes {
		currentLine, lineHasFootnotes := footnotes[cursor.LinePos]

		verseStyle := VerseNoteStyle(currentLine.MarkerStyle.Int32)
		if lineHasFootnotes && verseStyle != VerseNoteStyleHeadless {
			xPos := cursor.Leftmargin + int(li.CalculateLyricWidth(cursor.LineText)+pos.X)
			styleFontSize := "font-family:'Figtree';font-weight:600;"
			switch verseStyle {
			case VerseNoteStyleAlignRight:
				styleFontSize = "font-family:'Figtree';font-size:60%;font-weight:600;"
				approxLineLength := constant.LAYOUT_WIDTH - (2 * defaultX)
				xPos = int(pos.X) + cursor.Leftmargin + approxLineLength
			case VerseNoteStyleHeadonly, VerseNoteStyleDirectAppendText:
				styleFontSize = "font-family:'Caladea';font-size:90%;font-weight:600;"
				xPos -= int(li.CalculateLyricWidth(" ")) // lyric on db is just white spaces
			}
			canv.Group("class='footnotes'", fmt.Sprintf(`style="font-style:italic;%s"`, styleFontSize))
			canv.Text(xPos, int(pos.Y), currentLine.FootnoteMarker.String)
			canv.Gend()
		}
	}
}

func RenderVerseFootnotes(canv canvas.Canvas, y *int, footnotes map[int]map[int]repository.VerseFootNotes) {
	if len(footnotes) == 0 {
		return
	}
	flatten := []repository.VerseFootNotes{}
	versenoteHeadonlyCnt := 0
	hasInternalItalic := false
	for _, fn := range footnotes {
		for _, t := range fn {
			hasInternalItalic = hasInternalItalic || strings.Contains(t.FootnoteMarker.String, utils.ITALIC_OPENING) || strings.Contains(t.Footnote.String, utils.ITALIC_OPENING)
			if VerseNoteStyle(t.MarkerStyle.Int32) == VerseNoteStyleHeadonly {
				versenoteHeadonlyCnt++
				continue
			}

			flatten = append(flatten, t)
		}
	}

	if versenoteHeadonlyCnt == len(flatten) {
		return
	}

	// Sort the footnotes by its markers
	sort.Slice(flatten, func(i, j int) bool {
		return flatten[i].FootnoteMarker.String < flatten[j].FootnoteMarker.String
	})

	*y = *y + Y_ADJUSTED_POS

	footnotesStyle := ""
	if !hasInternalItalic {
		footnotesStyle = STYLE_ITALIC
	}
	canv.Group(CLASSNAME_GROUP, fmt.Sprintf(STYLE_GROUP_CUSTOM, footnotesStyle))
	totalLine := 0

	indent := float64(constant.LAYOUT_INDENT_LENGTH + 20)
	for i, fn := range flatten {
		lines := strings.Split(fn.Footnote.String, "<br/>")
		if len(lines) >= 2 {
			marker, clean := utils.ReplaceItalicToSpanWithClean(fn.FootnoteMarker.String)

			xNotes := int(utils.CalculateSecondaryLyricWidth(clean))
			yPosMarker := float64((NEWLINE_HEIGHT * i) + (*y))
			canv.TextUnescaped(indent, yPosMarker, marker)

			for li, line := range lines {
				line = utils.ReplaceItalicToSpan(line)
				yPos := float64((NEWLINE_HEIGHT * (i + li)) + (*y))
				canv.TextUnescaped(indent+float64(xNotes), yPos, line)
			}
			*y = *y + (NEWLINE_HEIGHT * (i + len(lines)))
		} else {
			totalLine++
			canv.Text(int(indent), (NEWLINE_HEIGHT*i)+(*y), fn.FootnoteMarker.String+fn.Footnote.String)
		}
	}
	canv.Gend()
	*y = *y + (NEWLINE_HEIGHT + (totalLine * NEWLINE_HEIGHT))

}
