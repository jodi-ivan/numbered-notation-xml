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
	Lyric    lyric.Lyric
	Staff    staff.Staff
	Credits  credits.Credits
	Footnote footnote.Footnote
	Verse    verse.Verse
	Header   header.Header
}

func NewRenderer() Renderer {
	l := lyric.NewLyric()
	f := footnote.New(l)
	return &rendererInteractor{
		Lyric:    l,
		Staff:    staff.NewStaff(),
		Credits:  credits.NewCredits(),
		Footnote: f,
		Verse:    verse.New(f, l),
		Header:   header.NewHeader(l),
	}
}

func (ir *rendererInteractor) Render(ctx context.Context, music musicxml.MusicXML, canv canvas.Canvas, metadata *repository.HymnMetadata) {
	canv.Start(constant.LAYOUT_WIDTH, 2000)
	canv.Def()
	fmt.Fprintf(canv.Writer(), fontfmt, string(googlefont("Caladea|Old Standard TT|Noto Music|Figtree")))
	canv.DefEnd()

	keySignature := keysig.NewKeySignature(ctx, music.Part.Measures)
	timeSignature := timesig.NewTimeSignatures(ctx, music.Part.Measures)

	ir.Header.RenderSheetHeader(ctx, canv, music.Credit, metadata)
	ir.Header.RenderKeyandTimeSignatures(ctx, canv, keySignature, timeSignature)

	relativeY := ir.Staff.Render(ctx, canv, music.Part, keySignature, timeSignature)

	if metadata != nil {
		ir.Footnote.RenderMusicFootnotes(ctx, canv, metadata, relativeY)
		verseInfo := ir.Verse.RenderVerse(ctx, canv, relativeY, metadata.Verse, metadata.VerseFootNotes)

		if verseInfo.MarginBottom != 0 {
			relativeY = verseInfo.MarginBottom
		}
		ir.Footnote.RenderVerseFootnotes(canv, &relativeY, metadata.VerseFootNotes)
		ir.Credits.RenderCredits(ctx, canv, &relativeY, metadata.HymnData)
		ir.Footnote.RenderTitleFootnotes(canv, relativeY, metadata.HymnData)

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
