package footnote

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (fi *footnoteInteractor) RenderTitleFootnotes(canv canvas.Canvas, y int, metadata repository.HymnData) {
	if !metadata.TitleFootnotes.Valid && metadata.IsForKids.Int16 != 1 {
		return
	}
	canv.Group(CLASSNAME_GROUP, STYLE_GROUP)
	if metadata.TitleFootnotes.Valid {
		notes := utils.TSPAN_OPENING + "* " + metadata.TitleFootnotes.String + utils.TSPAN_CLOSING
		y += 30
		canv.TextUnescaped(constant.LAYOUT_INDENT_LENGTH, float64(y), notes)
	}
	if metadata.IsForKids.Int16 == 1 {
		canv.TextUnescaped(constant.LAYOUT_INDENT_LENGTH, float64(y+25),
			`<tspan font-style="italic">Semua nyayian dengan tanda</tspan>
			<tspan font-style="bold" font-size="125%%">☆</tspan>
			<tspan font-style="italic">: khusus untuk anak-anak</tspan>`,
		)
	}
	canv.Gend()
}
