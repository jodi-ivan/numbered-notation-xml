package renderer

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/credits"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
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
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

type Renderer interface {
	Render(ctx context.Context, music musicxml.MusicXML, canv canvas.Canvas, metadata *entity.HymnMetaData)
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

func (ir *rendererInteractor) Render(ctx context.Context, music musicxml.MusicXML, canv canvas.Canvas, metadata *entity.HymnMetaData) {
	canvHeight := 3000
	canv.Def()
	fmt.Fprintf(canv.Writer(), fontfmt, string(googlefont()))
	canv.DefEnd()

	keySignature := keysig.NewKeySignature(ctx, music.Part.Measures)
	timeSignature := timesig.NewTimeSignatures(ctx, music.Part.Measures)

	ir.Header.RenderSheetHeader(ctx, canv, music.Credit, metadata)
	ir.Header.RenderKeyandTimeSignatures(ctx, canv, keySignature, timeSignature)

	relativeY := ir.Staff.Render(ctx, canv, music.Part, keySignature, timeSignature, metadata)
	if metadata != nil {
		prm, _ := params.GetParamFromContext(ctx)
		if prm.Verse > 2 || (prm.Verse > 1 && prm.SingleVerseMode) {
			firstVerse := verse.BuildContent(music, metadata)
			metadata.ParsedVerse[1] = firstVerse
			verseInfo := repository.HymnVerse{}
			if otherVerse, ok := metadata.Verse[2]; ok {
				verseInfo.StyleRow = otherVerse.StyleRow
			}
			metadata.Verse[1] = verseInfo
		}

		ir.Footnote.RenderMusicFootnotes(ctx, canv, metadata.HymnMetadata, relativeY)
		verseInfo := ir.Verse.RenderVerse(ctx, canv, relativeY, metadata)

		if verseInfo.MarginBottom != 0 {
			relativeY = verseInfo.MarginBottom
		}
		ir.Footnote.RenderVerseFootnotes(canv, &relativeY, metadata.VerseFootNotes)
		ir.Credits.RenderCredits(ctx, canv, &relativeY, metadata.HymnData)

		canvHeight = relativeY + 50
		ir.Footnote.RenderTitleFootnotes(canv, relativeY, metadata.HymnData)
	}
	canv.Start(constant.LAYOUT_WIDTH, canvHeight)
	canv.End()

}

func googlefont() []byte {
	return []byte(`@font-face {
         font-family: 'Caladea';
         font-style: normal;
         font-weight: 400;
         src: url(/assets/fonts/caladea.ttf) format('truetype');
       }
       @font-face {
         font-family: 'Figtree';
         font-style: normal;
         font-weight: 400;
         src: url(/assets/fonts/figtree.ttf) format('truetype');
        }
       @font-face {
         font-family: 'Noto Music';
         font-style: normal;
         font-weight: 400;
         src: url(/assets/fonts/noto-music.ttf) format('truetype');
       }
       @font-face {
         font-family: 'Old Standard TT';
         font-style: normal;
         font-weight: 400;
         src: url(/assets/fonts/old-standard-tt.ttf) format('truetype');
       }
	   @font-face {
         font-family: 'mozart11';
         font-style: normal;
         font-weight: 400;
         src: url(/assets/fonts/mozart11.ttf) format('truetype');
       }`)
}
