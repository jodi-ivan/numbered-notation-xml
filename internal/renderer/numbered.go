package renderer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Delegator interface {
	Render(ctx context.Context, music musicxml.MusicXML, canv canvas.Canvas, metadata *repository.HymnMetadata)
}

type delegator struct{}

func NewDelegator() Delegator {
	return &delegator{}
}

func (d *delegator) Render(ctx context.Context, music musicxml.MusicXML, canv canvas.Canvas, metadata *repository.HymnMetadata) {
	canv.Start(constant.LAYOUT_WIDTH, 1000)
	canv.Def()
	fmt.Fprintf(canv.Writer(), fontfmt, string(googlefont("Caladea|Old Standard TT|Noto Music|Figtree")))
	canv.DefEnd()

	relativeY := 100
	// render title

	workTitle := ""
	if metadata != nil {
		workTitle = fmt.Sprintf("%d. %s", metadata.Number, strings.ToUpper(metadata.Title))
	}
	titleWidth := lyric.CalculateLyricWidth(workTitle)
	titleX := (constant.LAYOUT_WIDTH / 2) - (titleWidth * 0.5)
	canv.Text(int(titleX), relativeY, workTitle)

	relativeY += 25

	keySignature := keysig.NewKeySignature(music.Part.Measures[0].Attribute.Key)
	timeSignature := timesig.NewTimeSignatures(ctx, music.Part.Measures)

	humanizedKeySignature := keySignature.String()

	canv.Text(constant.LAYOUT_INDENT_LENGTH, relativeY, keySignature.String())

	beat := music.Part.Measures[0].Attribute.Time
	canv.Text(constant.LAYOUT_INDENT_LENGTH+(len(humanizedKeySignature)*LOWERCASE_LENGTH), relativeY, fmt.Sprintf("%d ketuk", beat.Beats))
	relativeY += 70

	staff := SplitLines(ctx, music.Part)
	x := constant.LAYOUT_INDENT_LENGTH
	info := StaffInfo{
		NextLineRenderer: []*entity.NoteRenderer{},
	}
	for _, st := range staff {
		info = RenderStaff(ctx, canv, x, relativeY, keySignature, timeSignature, st, info.NextLineRenderer...)
		relativeY = relativeY + 80 + info.MarginBottom
		if info.Multiline {
			x = info.MarginLeft
		} else {
			x = constant.LAYOUT_INDENT_LENGTH
		}
	}

	if metadata != nil {
		verseInfo := RenderVerse(ctx, canv, relativeY+20, metadata.Verse)
		relativeY = verseInfo.MarginBottom

		RenderCredits(ctx, canv, relativeY+50, metadata.HymnData)

	}
	canv.End()

}

func googlefont(f string) []byte {
	empty := []byte{}
	r, err := http.Get(gwfURI + url.QueryEscape(f))
	if err != nil {
		return empty
	}
	defer r.Body.Close()
	b, rerr := io.ReadAll(r.Body)
	if rerr != nil || r.StatusCode != http.StatusOK {
		return empty
	}

	return b
}

// SplitLines split the measure in the lines manner
func SplitLines(ctx context.Context, part musicxml.Part) [][]musicxml.Measure {
	result := [][]musicxml.Measure{}
	currentLine := []musicxml.Measure{}
	for _, measure := range part.Measures {

		if measure.Print != nil && measure.Print.NewSystem == musicxml.PrintNewSystemTypeYes {
			finishLine := make([]musicxml.Measure, len(currentLine))
			copy(finishLine, currentLine)

			result = append(result, finishLine)

			currentLine = []musicxml.Measure{}
		}
		currentLine = append(currentLine, measure)

	}

	return append(result, currentLine)
}
