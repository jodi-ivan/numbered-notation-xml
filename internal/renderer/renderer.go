package renderer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/credits"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/header"
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
	canv.Start(constant.LAYOUT_WIDTH, 2000)
	canv.Def()
	fmt.Fprintf(canv.Writer(), fontfmt, string(googlefont("Caladea|Old Standard TT|Noto Music|Figtree")))
	canv.DefEnd()
	keySignature := keysig.NewKeySignature(ctx, music.Part.Measures)
	timeSignature := timesig.NewTimeSignatures(ctx, music.Part.Measures)

	header.RenderSheetHeader(ctx, canv, music.Credit, metadata)
	header.RenderSignatures(ctx, canv, keySignature, timeSignature)

	relativeY := constant.TITLE_Y_POS + header.HEADER_OFFSET

	staffes := ir.Staff.SplitLines(ctx, music.Part)
	x := constant.LAYOUT_INDENT_LENGTH
	info := staff.StaffInfo{
		NextLineRenderer: []*entity.NoteRenderer{},
	}
	oldMarginButtom := 0
	for i, st := range staffes {
		info = ir.Staff.RenderStaff(ctx, canv, x, relativeY, i == len(staffes)-1, keySignature, timeSignature, st, info.NextLineRenderer...)
		relativeY = relativeY + 70 + info.MarginBottom
		if info.ForceNewLine {
			relativeY += oldMarginButtom
		}
		if info.Multiline {
			x = info.MarginLeft + barline.BARLINE_AFTER_SPACE
			info.Multiline = false
		} else {
			x = constant.LAYOUT_INDENT_LENGTH
		}
		oldMarginButtom = info.MarginBottom
	}

	for len(info.NextLineRenderer) > 0 {
		x = constant.LAYOUT_INDENT_LENGTH
		info = ir.Staff.RenderStaff(ctx, canv, x, relativeY, true, keySignature, timeSignature, nil, info.NextLineRenderer...)
		relativeY += info.MarginBottom + 70
	}

	if metadata != nil {
		ir.Credits.RenderMuiscFootnotes(ctx, canv, metadata, x, relativeY)
		verseInfo := ir.Lyric.RenderVerse(ctx, canv, relativeY, metadata.Verse, metadata.VerseFootNotes)

		// FIXED: Y is 0 value (at the top of pages) when there is no verses
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
