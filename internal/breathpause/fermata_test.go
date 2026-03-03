package breathpause_test

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func TestRenderFermata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string
		canv func(ctrl *gomock.Controller) *canvas.MockCanvas

		fermata *musicxml.Femata
		pos     entity.Coordinate
	}{
		{
			name: "nope",
		},
		{
			name: "everything went fine",
			canv: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				writerMock := canvas.NewMockWriter(ctrl)
				writerMock.EXPECT().Write([]byte(`<text x="-4.000" y="-17.500" style="font-family:Noto Music;font-size:200%"> &#x1D110; </text>`))
				canv.EXPECT().Writer().Return(writerMock)
				return canv
			},
			fermata: &musicxml.Femata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var canv *canvas.MockCanvas
			if tt.canv != nil {
				canv = tt.canv(ctrl)

			}
			breathpause.RenderFermata(context.Background(), canv, tt.fermata, tt.pos)
		})
	}
}
