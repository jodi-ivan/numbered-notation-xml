package renderer

import (
	"context"
	"fmt"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderCredits(ctx context.Context, canv canvas.Canvas, y int, metadata repository.HymnData) {
	canv.Group("class='credit'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
	if metadata.Lyric == metadata.Music {
		canv.Text(constant.LAYOUT_INDENT_LENGTH, y, fmt.Sprintf("Syair dan lagu : %s", metadata.Lyric))
	} else {

		text := strings.ReplaceAll(metadata.Lyric, "<i>", `<tspan font-style="italic">`)
		text = strings.ReplaceAll(text, "</i>", "</tspan>")
		fmt.Fprintf(canv.Writer(), `<text x="%d" y="%d">Syair : %s </text>`, constant.LAYOUT_INDENT_LENGTH, y, text)

		y += 15
		canv.Text(constant.LAYOUT_INDENT_LENGTH, y, fmt.Sprintf("Lagu : %s", metadata.Music))

	}

	ref := ""
	if metadata.RefBE.Valid {
		ref += fmt.Sprintf("BE  %d", metadata.RefBE.Int16)
	}

	if metadata.RefNR.Valid {
		if ref != "" {
			ref += ", "
		}

		ref += fmt.Sprintf("NR  %d", metadata.RefNR.Int16)
	}

	if ref != "" {
		l := lyric.CalculateLyricWidth(ref)
		canv.Text(constant.LAYOUT_WIDTH-constant.LAYOUT_INDENT_LENGTH-int(l), y, ref)
	}

	canv.Gend()

}
