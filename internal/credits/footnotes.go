package credits

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (ci *creditsInteractor) RenderMuiscFootnotes(ctx context.Context, canv canvas.Canvas, metadata *repository.HymnMetadata, x, y int) {
	if !metadata.Footnotes.Valid {
		return
	}
	canv.Group("class='footnotes'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
	xPos := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH - int(CalculateLyric(metadata.Footnotes.String, true))
	canv.Text(xPos, y-45, metadata.Footnotes.String)
	canv.Gend()

}
