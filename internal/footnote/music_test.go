package footnote

import (
	"context"
	"database/sql"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_footnoteInteractor_RenderMusicFootnotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv     func(*gomock.Controller) *canvas.MockCanvas
		metadata *repository.HymnMetadata
		y        int
	}{
		{
			name:     "",
			canv:     func(c *gomock.Controller) *canvas.MockCanvas { return nil },
			y:        100,
			metadata: &repository.HymnMetadata{},
		},
		{
			name: "usual",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().Group("class='footnotes'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
				canvMock.EXPECT().Text(616, 65, "* unit = tets")
				canvMock.EXPECT().Gend()
				return canvMock
			},
			y: 100,
			metadata: &repository.HymnMetadata{
				HymnData: repository.HymnData{
					Footnotes: sql.NullString{
						Valid:  true,
						String: "* unit = tets",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var fi footnoteInteractor
			fi.RenderMusicFootnotes(context.Background(), tt.canv(ctrl), tt.metadata, tt.y)
		})
	}
}
