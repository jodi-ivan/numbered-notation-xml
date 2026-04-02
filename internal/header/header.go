package header

import (
	"context"
	"fmt"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/credits"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

var ir = lyric.NewLyric()

func renderTitle(ctx context.Context, canv canvas.Canvas, credit []musicxml.Credit, metadata *repository.HymnMetadata) {
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
		if metadata.TitleFootnotes.Valid {
			workTitle += TITLE_FOOTNOTES
		}
		if metadata.IsForKids.Int16 == 1 {
			canv.TextUnescaped(
				constant.LAYOUT_INDENT_LENGTH, float64(relativeY),
				FOR_KIDS_ELMNT)
		}
	}
	titleWidth := ir.CalculateLyricWidth(workTitle)
	titleX := (constant.LAYOUT_WIDTH / 2) - (titleWidth * 0.5)
	canv.Text(int(titleX), relativeY, workTitle)

}

func renderSubtitle(ctx context.Context, canv canvas.Canvas, credit []musicxml.Credit, metadata *repository.HymnMetadata) {
	relativeY := constant.TITLE_Y_POS

	subtitle := ""
	for _, v := range credit {
		if v.Type == musicxml.CreditTypeSubtitle && v.Words != EMPTY_SUBTITLE {
			subtitle = v.Words
		}
	}

	if subtitle != "" {
		num := 0.0
		if metadata != nil {
			num = ir.CalculateLyricWidth(fmt.Sprintf("%d. ", metadata.Number)) / 2
		}
		subtitleWidth := (credits.CalculateLyric(subtitle, false) * SUBTITLE_TO_CREDITS_SIZE_RATIO)
		subtitleX := (constant.LAYOUT_WIDTH / 2) - (subtitleWidth * 0.5)
		canv.Text(int(subtitleX+num), relativeY+SUBTITLE_Y_POS, subtitle, SUBTITLE_ATTR)
	}
}

func RenderSheetHeader(ctx context.Context, canv canvas.Canvas, credit []musicxml.Credit, metadata *repository.HymnMetadata) {
	renderTitle(ctx, canv, credit, metadata)
	renderSubtitle(ctx, canv, credit, metadata)

}
