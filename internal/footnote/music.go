package footnote

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/layout"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (fi *footnoteInteractor) RenderMusicFootnotes(ctx context.Context, canv canvas.Canvas, metadata *repository.HymnMetadata, y int) {
	if !metadata.Footnotes.Valid {
		return
	}
	_, marginRight := layout.GetMargin(ctx)
	layoutWidth := layout.GetLayoutWidth(ctx)
	canv.Group(CLASSNAME_GROUP, STYLE_GROUP)
	xPos := layoutWidth - marginRight - int(utils.CalculateSecondaryLyricWidth(metadata.Footnotes.String))
	canv.Text(xPos, y-MUSIC_FOOTNOTES_Y_OFFSET, metadata.Footnotes.String)
	canv.Gend()

}
