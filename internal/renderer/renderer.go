package renderer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/credits"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Renderer interface {
	Render(ctx context.Context, music musicxml.MusicXML, canv canvas.Canvas, metadata *repository.HymnMetadata)
}

type rendererInteractor struct {
	Lyric   lyric.Lyric
	Staff   staff.Staff
	Credits credits.Credits
}

func NewRenderer() Renderer {
	return &rendererInteractor{
		Lyric:   lyric.NewLyric(),
		Staff:   staff.NewStaff(),
		Credits: credits.NewCredits(),
	}
}

func (ir *rendererInteractor) Render(ctx context.Context, music musicxml.MusicXML, canv canvas.Canvas, metadata *repository.HymnMetadata) {
	canv.Start(constant.LAYOUT_WIDTH, 1500)
	canv.Def()
	fmt.Fprintf(canv.Writer(), fontfmt, string(googlefont("Caladea|Old Standard TT|Noto Music|Figtree")))
	canv.DefEnd()

	relativeY := 100
	// render title

	workTitle := ""
	for _, v := range music.Credit {
		if v.Type == musicxml.CreditTypeTitle {
			workTitle = strings.ToUpper(v.Words)
		}
	}
	if metadata != nil {
		workTitle = fmt.Sprintf("%d. %s", metadata.Number, strings.ToUpper(metadata.Title))
	}
	titleWidth := ir.Lyric.CalculateLyricWidth(workTitle)
	titleX := (constant.LAYOUT_WIDTH / 2) - (titleWidth * 0.5)
	canv.Text(int(titleX), relativeY, workTitle)

	relativeY += 25

	keySignature := keysig.NewKeySignature(music.Part.Measures[0].Attribute.Key)
	timeSignature := timesig.NewTimeSignatures(ctx, music.Part.Measures)

	humanizedKeySignature := timeSignature.GetHumanized()

	canv.Text(constant.LAYOUT_INDENT_LENGTH, relativeY, keySignature.String())

	canv.Text(constant.LAYOUT_INDENT_LENGTH+(3*constant.LOWERCASE_LENGTH)+int(ir.Lyric.CalculateLyricWidth(keySignature.String())), relativeY, humanizedKeySignature)
	relativeY += 70

	staffes := ir.Staff.SplitLines(ctx, music.Part)
	x := constant.LAYOUT_INDENT_LENGTH
	info := staff.StaffInfo{
		NextLineRenderer: []*entity.NoteRenderer{},
	}

	for _, st := range staffes {
		info = ir.Staff.RenderStaff(ctx, canv, x, relativeY, keySignature, timeSignature, st, info.NextLineRenderer...)
		relativeY = relativeY + 80 + info.MarginBottom
		if info.Multiline {
			x = info.MarginLeft
		} else {
			x = constant.LAYOUT_INDENT_LENGTH
		}
	}

	if metadata != nil {
		verseInfo := ir.Lyric.RenderVerse(ctx, canv, relativeY, metadata.Verse)
		relativeY = verseInfo.MarginBottom

		ir.Credits.RenderCredits(ctx, canv, relativeY, metadata.HymnData)

	}
	canv.End()

}

func googlefont(f string) []byte {
	empty := []byte{}

	link := gwfURI + url.QueryEscape(f)
	r, err := http.Get(link)
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
