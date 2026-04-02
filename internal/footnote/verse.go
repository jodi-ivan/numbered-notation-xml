package footnote

import (
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
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
