package header

import (
	"context"
	"fmt"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/footnote"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Header interface {
	RenderSheetHeader(ctx context.Context, canv canvas.Canvas, credit []musicxml.Credit, metadata *repository.HymnMetadata)
	RenderKeyandTimeSignatures(ctx context.Context, canv canvas.Canvas, key keysig.KeySignature, timeSignature timesig.TimeSignature)
}

type headerInteractor struct {
	Lyric lyric.Lyric
}

func NewHeader(l lyric.Lyric) Header {
	return &headerInteractor{
		Lyric: l,
	}
}

func hasTitleNotes(verseNotes *repository.HymnMetadata) bool {
	for _, notes := range verseNotes.VerseFootNotes {
		for _, line := range notes {
			if footnote.VerseNoteStyle(line.MarkerStyle.Int32) == footnote.VerseNoteStyleForTitle {
				return true
			}
		}
	}
	return false
}
func (hi *headerInteractor) renderTitle(ctx context.Context, canv canvas.Canvas, credit []musicxml.Credit, metadata *repository.HymnMetadata) {
	relativeY := constant.TITLE_Y_POS

	workTitle := ""
	for _, v := range credit {
		if v.Type == musicxml.CreditTypeTitle {
			workTitle = strings.ToUpper(v.Words)
		}

	}
	if metadata != nil {
		workTitle = fmt.Sprintf(TITLE_FMT, metadata.Number, strings.ToUpper(metadata.Title))
		if metadata.Variant.Valid {
			workTitle = fmt.Sprintf(TITLE_VARIANT_FMT, metadata.Number, strings.ToLower(metadata.Variant.String), strings.ToUpper(metadata.Title))
		}
		if metadata.TitleFootnotes.Valid || hasTitleNotes(metadata) {
			workTitle += TITLE_FOOTNOTES
		}
		if metadata.IsForKids.Int16 == 1 {
			canv.TextUnescaped(
				constant.LAYOUT_INDENT_LENGTH, float64(relativeY),
				FOR_KIDS_ELMNT)
		}
	}
	titleWidth := hi.Lyric.CalculateLyricWidth(workTitle)
	titleX := (constant.LAYOUT_WIDTH / 2) - (titleWidth * 0.5)
	canv.Text(int(titleX), relativeY, workTitle)

}

func (hi *headerInteractor) renderSubtitle(ctx context.Context, canv canvas.Canvas, credit []musicxml.Credit, metadata *repository.HymnMetadata) {
	relativeY := constant.TITLE_Y_POS

	subtitle := ""
	for _, v := range credit {
		if v.Type == musicxml.CreditTypeSubtitle && (!strings.EqualFold(v.Words, EMPTY_SUBTITLE) || v.Words == "") {
			subtitle = v.Words
		}
	}

	if subtitle == "" {
		return
	}

	num := 0.0
	if metadata != nil {
		num = hi.Lyric.CalculateLyricWidth(fmt.Sprintf("%d. ", metadata.Number)) / 2
		if metadata.Variant.Valid {
			num = hi.Lyric.CalculateLyricWidth(fmt.Sprintf("%d%s. ", metadata.Number, metadata.Variant.String)) / 2
		}
	}
	subtitleWidth := (utils.CalculateSecondaryLyricWidth(subtitle) * SUBTITLE_TO_CREDITS_SIZE_RATIO)
	subtitleX := (constant.LAYOUT_WIDTH / 2) - (subtitleWidth * 0.5)
	canv.Text(int(subtitleX+num), relativeY+SUBTITLE_Y_POS, subtitle, SUBTITLE_ATTR)
}

func (hi *headerInteractor) RenderSheetHeader(ctx context.Context, canv canvas.Canvas, credit []musicxml.Credit, metadata *repository.HymnMetadata) {
	hi.renderTitle(ctx, canv, credit, metadata)
	hi.renderSubtitle(ctx, canv, credit, metadata)
}
