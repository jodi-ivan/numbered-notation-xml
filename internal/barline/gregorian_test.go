package barline

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	ubermock "go.uber.org/mock/gomock"
)

func Test_barlineInteractor_RenderGregorian(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uberCtrl := ubermock.NewController(t)
	defer uberCtrl.Finish()

	// Mock setup helpers
	setupMock := func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff) {
		canvMock := canvas.NewMockCanvas(ctrl)
		staffMock := lines.NewMockLineStaff(uber)
		return canvMock, staffMock
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		initMock func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff)

		b          *musicxml.Barline
		isLastNote bool
		pos        entity.Coordinate
	}{
		{
			name: "Regular Barline (Not Last Note)",
			initMock: func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff) {
				canvMock, staffMock := setupMock(ctrl, uber)
				canvMock.EXPECT().Group(`class="barline"`)
				canvMock.EXPECT().Line(50, 100, 50, 150, `fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9`)
				canvMock.EXPECT().Gend()
				staffMock.EXPECT().GetTopLine().Return(100)
				staffMock.EXPECT().GetBottomLine().Return(150)
				staffMock.EXPECT().GetMiddleLine().Return(125)
				return canvMock, staffMock
			},
			b:          &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular},
			isLastNote: false,
			pos:        entity.Coordinate{X: 50, Y: 0},
		},
		{
			name: "Regular Barline ( Last Note)",
			initMock: func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff) {
				canvMock, staffMock := setupMock(ctrl, uber)
				canvMock.EXPECT().Group(`class="barline"`)
				canvMock.EXPECT().Line(54, 100, 54, 150, `fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9`)
				canvMock.EXPECT().Gend()
				staffMock.EXPECT().GetTopLine().Return(100)
				staffMock.EXPECT().GetBottomLine().Return(150)
				staffMock.EXPECT().GetMiddleLine().Return(125)
				return canvMock, staffMock
			},
			b:          &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular},
			isLastNote: true,
			pos:        entity.Coordinate{X: 50, Y: 0},
		},
		{
			name: "Barline LightLight",
			initMock: func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff) {
				canvMock, staffMock := setupMock(ctrl, uber)
				canvMock.EXPECT().Group(`class="barline"`)
				canvMock.EXPECT().Line(51, 100, 51, 150, `fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9`)
				canvMock.EXPECT().Line(54, 100, 54, 150, `fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9`)
				canvMock.EXPECT().Gend()
				staffMock.EXPECT().GetTopLine().Return(100)
				staffMock.EXPECT().GetBottomLine().Return(150)
				staffMock.EXPECT().GetMiddleLine().Return(125)
				return canvMock, staffMock
			},
			b:   &musicxml.Barline{BarStyle: musicxml.BarLineStyleLightLight},
			pos: entity.Coordinate{X: 50, Y: 0},
		},
		{
			name: "Barline LightHeavy",
			initMock: func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff) {
				canvMock, staffMock := setupMock(ctrl, uber)
				canvMock.EXPECT().Group(`class="barline"`)
				canvMock.EXPECT().Line(51, 100, 51, 150, `fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9`)
				canvMock.EXPECT().Line(56, 102, 56, 148, `fill:none;stroke:#000000;stroke-linecap:square;stroke-width:4.6`)

				canvMock.EXPECT().Gend()
				staffMock.EXPECT().GetTopLine().Return(100)
				staffMock.EXPECT().GetBottomLine().Return(150)
				staffMock.EXPECT().GetMiddleLine().Return(125)
				return canvMock, staffMock
			},
			b:   &musicxml.Barline{BarStyle: musicxml.BarLineStyleLightHeavy},
			pos: entity.Coordinate{X: 50, Y: 0},
		},
		{
			name: "Barline LightHeavy with backward repeat",
			initMock: func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff) {
				canvMock, staffMock := setupMock(ctrl, uber)
				canvMock.EXPECT().Group(`class="barline"`)
				canvMock.EXPECT().Line(51, 100, 51, 150, `fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9`)
				canvMock.EXPECT().Line(56, 102, 56, 148, `fill:none;stroke:#000000;stroke-linecap:square;stroke-width:4.6`)
				canvMock.EXPECT().Text(44, 130, ":", `style="font-family:Noto Music;font-size:19.2px"`)

				canvMock.EXPECT().Gend()
				staffMock.EXPECT().GetTopLine().Return(100)
				staffMock.EXPECT().GetBottomLine().Return(150)
				staffMock.EXPECT().GetMiddleLine().Return(125)
				return canvMock, staffMock
			},
			b: &musicxml.Barline{
				BarStyle: musicxml.BarLineStyleLightHeavy,
				Repeat: &musicxml.BarLineRepeat{
					Direction: musicxml.BarLineRepeatDirectionBackward,
				},
			},
			pos: entity.Coordinate{X: 50, Y: 0},
		},
		{
			name: "Barline HeavyLight",
			initMock: func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff) {
				canvMock, staffMock := setupMock(ctrl, uber)
				canvMock.EXPECT().Group(`class="barline"`)
				canvMock.EXPECT().Line(52, 102, 52, 148, `fill:none;stroke:#000000;stroke-linecap:square;stroke-width:4.6`)
				canvMock.EXPECT().Line(57, 100, 57, 150, `fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9`)
				// canvMock.EXPECT().Text(44, 130, ":", `style="font-family:Noto Music;font-size:19.2px"`)

				canvMock.EXPECT().Gend()
				staffMock.EXPECT().GetTopLine().Return(100)
				staffMock.EXPECT().GetBottomLine().Return(150)
				staffMock.EXPECT().GetMiddleLine().Return(125)
				return canvMock, staffMock
			},
			b: &musicxml.Barline{
				BarStyle: musicxml.BarLineStyleHeavyLight,
			},
			pos: entity.Coordinate{X: 50, Y: 0},
		},
		{
			name: "Barline HeavyLight wit repeat",
			initMock: func(ctrl *gomock.Controller, uber *ubermock.Controller) (*canvas.MockCanvas, *lines.MockLineStaff) {
				canvMock, staffMock := setupMock(ctrl, uber)
				canvMock.EXPECT().Group(`class="barline"`)
				canvMock.EXPECT().Line(52, 102, 52, 148, `fill:none;stroke:#000000;stroke-linecap:square;stroke-width:4.6`)
				canvMock.EXPECT().Line(57, 100, 57, 150, `fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9`)
				canvMock.EXPECT().Text(58, 130, ":", `style="font-family:Noto Music;font-size:19.2px"`)

				canvMock.EXPECT().Gend()
				staffMock.EXPECT().GetTopLine().Return(100)
				staffMock.EXPECT().GetBottomLine().Return(150)
				staffMock.EXPECT().GetMiddleLine().Return(125)
				return canvMock, staffMock
			},
			b: &musicxml.Barline{
				BarStyle: musicxml.BarLineStyleHeavyLight,
				Repeat: &musicxml.BarLineRepeat{
					Direction: musicxml.BarLineRepeatDirectionForward,
				},
			},
			pos: entity.Coordinate{X: 50, Y: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canvMock, staffMock := tt.initMock(ctrl, uberCtrl)
			var bi barlineInteractor
			bi.RenderGregorian(canvMock, tt.b, tt.isLastNote, staffMock, tt.pos)
		})
	}
}
