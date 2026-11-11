package footnote

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/layout"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (fi *footnoteInteractor) RenderTitleFootnotes(ctx context.Context, canv canvas.Canvas, y int, metadata repository.HymnData) {
	if !metadata.TitleFootnotes.Valid && metadata.IsForKids.Int16 != 1 {
		return
	}
	marginLeft, _ := layout.GetMargin(ctx)

	canv.Group(CLASSNAME_GROUP, STYLE_GROUP)
	if metadata.TitleFootnotes.Valid {
		notes := utils.TSPAN_OPENING + "* " + metadata.TitleFootnotes.String + utils.TSPAN_CLOSING
		y += 30
		canv.TextUnescaped(float64(marginLeft), float64(y), notes)
	}
	if metadata.IsForKids.Int16 == 1 {
		canv.TextUnescaped(float64(marginLeft), float64(y+25),
			`<tspan font-style="italic">Semua nyayian dengan tanda</tspan>
			<tspan font-style="bold" font-size="12px">☆</tspan>
			<tspan font-style="italic">: khusus untuk anak-anak</tspan>`,
		)
	}
	canv.Gend()
}
