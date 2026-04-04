package footnote

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

// var li = lyric.NewLyric()

type VerseLineCursor struct {
	VerseNo    int
	LinePos    int
	Leftmargin int
	LineText   string
}

func (fi *footnoteInteractor) AssignFootnotesMarker(canv canvas.Canvas, pos entity.Coordinate, defaultX int, cursor VerseLineCursor, verseFootnote map[int]map[int]repository.VerseFootNotes) {

	footnotes, hasFootnotes := verseFootnote[cursor.VerseNo]
	if !hasFootnotes {
		return
	}

	currentLine, lineHasFootnotes := footnotes[cursor.LinePos]

	verseStyle := VerseNoteStyle(currentLine.MarkerStyle.Int32)
	if verseStyle == VerseNoteStyleHeadless || !lineHasFootnotes {
		return
	}

	xPos := cursor.Leftmargin + int(fi.li.CalculateLyricWidth(cursor.LineText)+pos.X)
	styleFontSize := BASE_FOOTNOTES_STYLE
	switch verseStyle {
	case VerseNoteStyleAlignRight:
		styleFontSize = BASE_STYLE_GROUP
		approxLineLength := constant.LAYOUT_WIDTH - (2 * defaultX)
		xPos = int(pos.X) + cursor.Leftmargin + approxLineLength
	case VerseNoteStyleHeadonly, VerseNoteStyleDirectAppendText:
		styleFontSize = BASE_VERSE_STYLE
		xPos -= int(fi.li.CalculateLyricWidth(" ")) // lyric on db is just white spaces
	}
	canv.Group(CLASSNAME_GROUP, fmt.Sprintf(STYLE_VERSE_FOOTNOTES_GROUP_CUSTOM, styleFontSize))
	canv.Text(xPos, int(pos.Y), currentLine.FootnoteMarker.String)
	canv.Gend()
}

func (fi *footnoteInteractor) RenderVerseFootnotes(canv canvas.Canvas, y *int, footnotes map[int]map[int]repository.VerseFootNotes) {
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
