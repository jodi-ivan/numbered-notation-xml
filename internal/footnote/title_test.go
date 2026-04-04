package footnote

import (
	"database/sql"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_footnoteInteractor_RenderTitleFootnotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv     func(*gomock.Controller) *canvas.MockCanvas
		y        int
		metadata repository.HymnData
	}{
		{
			name: "",
			canv: func(c *gomock.Controller) *canvas.MockCanvas { return nil },
			y:    100,
		},
		{
			name: "title only",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().Group("class='footnotes'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
				// canv.TextUnescaped(constant.LAYOUT_INDENT_LENGTH, float64(y), notes)
				canvMock.EXPECT().TextUnescaped(50.0, 130.0, `<tspan font-style="italic">* Bisa juga di unittest</tspan>`)
				canvMock.EXPECT().Gend()
				return canvMock
			},
			y: 100,
			metadata: repository.HymnData{
				TitleFootnotes: sql.NullString{
					Valid:  true,
					String: "Bisa juga di unittest",
				},
			},
		},
		{
			name: "kids mark only",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().Group("class='footnotes'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
				// canv.TextUnescaped(constant.LAYOUT_INDENT_LENGTH, float64(y), notes)
				canvMock.EXPECT().TextUnescaped(50.0, 125.0,
					`<tspan font-style="italic">Semua nyayian dengan tanda</tspan>
			<tspan font-style="bold" font-size="125%%">☆</tspan>
			<tspan font-style="italic">: khusus untuk anak-anak</tspan>`,
				)
				canvMock.EXPECT().Gend()
				return canvMock
			},
			y: 100,
			metadata: repository.HymnData{
				IsForKids: sql.NullInt16{
					Int16: 1,
					Valid: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var fi footnoteInteractor
			fi.RenderTitleFootnotes(tt.canv(ctrl), tt.y, tt.metadata)
		})
	}
}
