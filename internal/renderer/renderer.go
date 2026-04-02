package renderer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/credits"
	"github.com/jodi-ivan/numbered-notation-xml/internal/footnote"
	"github.com/jodi-ivan/numbered-notation-xml/internal/header"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/verse"
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
	canv.Start(constant.LAYOUT_WIDTH, 2000)
	canv.Def()
	fmt.Fprintf(canv.Writer(), fontfmt, string(googlefont("Caladea|Old Standard TT|Noto Music|Figtree")))
	canv.DefEnd()

	keySignature := keysig.NewKeySignature(ctx, music.Part.Measures)
	timeSignature := timesig.NewTimeSignatures(ctx, music.Part.Measures)

	header.RenderSheetHeader(ctx, canv, music.Credit, metadata)
	header.RenderSignatures(ctx, canv, keySignature, timeSignature)

	relativeY := ir.Staff.Render(ctx, canv, music.Part, keySignature, timeSignature)

	if metadata != nil {
		footnote.RenderMusicFootnotes(ctx, canv, metadata, relativeY)
		verseInfo := verse.RenderVerse(ctx, canv, relativeY, metadata.Verse, metadata.VerseFootNotes)

		if verseInfo.MarginBottom != 0 {
			relativeY = verseInfo.MarginBottom
		}

		relativeY = ir.Credits.RenderCredits(ctx, canv, relativeY, metadata.HymnData, metadata.VerseFootNotes)

		if metadata.IsForKids.Int16 == 1 {
			ir.Credits.RenderForKidsFootnotes(ctx, canv, relativeY+25)
		}
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
