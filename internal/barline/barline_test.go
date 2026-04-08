package barline

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

func Test_barlineInteractor_GetRendererLeftBarline(t *testing.T) {
	type args struct {
		measure                  musicxml.Measure
		x                        int
		lastRightBarlinePosition *entity.Coordinate
	}
	tests := []struct {
		name             string
		args             args
		wantNoteRenderer *entity.NoteRenderer
		wantBarlineInfo  *BarlineInfo
	}{
		{
			name: "no barline",
			args: args{
				measure: musicxml.Measure{},
			},
			wantNoteRenderer: nil,
			wantBarlineInfo:  nil,
		},
		{
			name: "no left barline",
			args: args{
				measure: musicxml.Measure{
					Barline: []musicxml.Barline{
						musicxml.Barline{
							Location: musicxml.BarlineLocationRight,
						},
					},
				},
			},
			wantNoteRenderer: nil,
			wantBarlineInfo:  nil,
		},
		{
			name: "there is barline, but it is just a regular barline",
			args: args{
				measure: musicxml.Measure{
					Barline: []musicxml.Barline{
						musicxml.Barline{
							Location: musicxml.BarlineLocationLeft,
							BarStyle: musicxml.BarLineStyleRegular,
						},
					},
				},
				x: 25,
			},
			wantNoteRenderer: nil,
			wantBarlineInfo:  nil,
		},
		{
			name: "everything went fine, double light, no repeat, no last bar location",
			args: args{
				measure: musicxml.Measure{
					Number: 1,
					Barline: []musicxml.Barline{
						musicxml.Barline{
							Location: musicxml.BarlineLocationLeft,
							BarStyle: musicxml.BarLineStyleLightLight,
						},
					},
				},
				x: 25,
			},
			wantNoteRenderer: &entity.NoteRenderer{
				PositionX: 25,
				Width:     10,
				Barline: &musicxml.Barline{
					Location: musicxml.BarlineLocationLeft,
					BarStyle: musicxml.BarLineStyleLightLight,
				},
				MeasureNumber: 1,
			},

			wantBarlineInfo: &BarlineInfo{
				XIncrement: 9,
			},
		},
		{
			name: "everything went fine, heavy light with repeat",
			args: args{
				measure: musicxml.Measure{
					Number: 1,
					Barline: []musicxml.Barline{
						musicxml.Barline{
							Location: musicxml.BarlineLocationLeft,
							BarStyle: musicxml.BarLineStyleHeavyLight,
							Repeat:   &musicxml.BarLineRepeat{},
						},
					},
				},
				x: 25,
				lastRightBarlinePosition: &entity.Coordinate{
					X: 30,
				},
			},
			wantNoteRenderer: &entity.NoteRenderer{
				PositionX: 30,
				Width:     11,
				Barline: &musicxml.Barline{
					Location: musicxml.BarlineLocationLeft,
					BarStyle: musicxml.BarLineStyleHeavyLight,
					Repeat:   &musicxml.BarLineRepeat{},
				},
				MeasureNumber: 1,
			},

			wantBarlineInfo: &BarlineInfo{
				XIncrement: 29,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := barlineInteractor{}
			got, got1 := bi.GetRendererLeftBarline(tt.args.measure, tt.args.x, &CoordinateWithBarline{Coordinate: *tt.args.lastRightBarlinePosition})
			if !assert.Equal(t, tt.wantNoteRenderer, got) {
				t.Errorf("barlineInteractor.GetRendererLeftBarline() got = %v, want %v", got, tt.wantNoteRenderer)
			}
			if !assert.Equal(t, tt.wantBarlineInfo, got1) {
				t.Errorf("barlineInteractor.GetRendererLeftBarline() got1 = %v, want %v", got1, tt.wantBarlineInfo)
			}
		})
	}
}

func Test_barlineInteractor_GetRendererRightBarline(t *testing.T) {
	type args struct {
		measure musicxml.Measure
		x       int
	}
	tests := []struct {
		name           string
		args           args
		wantBarlinePos int
		wantRenderer   *entity.NoteRenderer
	}{
		{
			name: "One barline in the measure without repeat",
			args: args{
				x: 25,
				measure: musicxml.Measure{
					Number: 1,
					Barline: []musicxml.Barline{
						{
							Location: musicxml.BarlineLocationRight,
						},
					},
				},
			},
			wantBarlinePos: 25,
			wantRenderer: &entity.NoteRenderer{
				MeasureNumber: 1,
				PositionX:     25,
				Barline: &musicxml.Barline{
					Location: musicxml.BarlineLocationRight,
				},
				Width: 4,
			},
		},
		{
			name: "Two barlines in the measure with repeat",
			args: args{
				x: 25,
				measure: musicxml.Measure{
					Number: 1,
					Barline: []musicxml.Barline{
						{
							Location: musicxml.BarlineLocationLeft,
						},
						{
							Location: musicxml.BarlineLocationRight,
							BarStyle: musicxml.BarLineStyleLightHeavy,
							Repeat: &musicxml.BarLineRepeat{
								Direction: musicxml.BarLineRepeatDirectionBackward,
							},
						},
					},
				},
			},
			wantBarlinePos: 35,
			wantRenderer: &entity.NoteRenderer{
				MeasureNumber: 1,
				PositionX:     25,
				Width:         4,
				Barline: &musicxml.Barline{
					Location: musicxml.BarlineLocationRight,
					BarStyle: musicxml.BarLineStyleLightHeavy,
					Repeat: &musicxml.BarLineRepeat{
						Direction: musicxml.BarLineRepeatDirectionBackward,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := barlineInteractor{}
			barlinePos, renderer := bi.GetRendererRightBarline(tt.args.measure, tt.args.x)
			if barlinePos != tt.wantBarlinePos {
				t.Errorf("barlineInteractor.GetRendererRightBarline() got = %v, want %v", barlinePos, tt.wantBarlinePos)
			}
			if !assert.Equal(t, tt.wantRenderer, renderer) {
				t.Errorf("barlineInteractor.GetRendererRightBarline() got1 = %v, want %v", renderer, tt.wantRenderer)
			}
		})
	}
}

func Test_barlineInteractor_RenderBarline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		barline    musicxml.Barline
		coordinate entity.Coordinate
	}
	tests := []struct {
		name       string
		initCanvas func(ctrl *gomock.Controller) *canvas.MockCanvas
		args       args
	}{
		// case 1 : no repeat
		{
			name: "no repeat barline",
			args: args{
				barline: musicxml.Barline{
					Location: musicxml.BarlineLocationRight,
					BarStyle: musicxml.BarLineStyleLightHeavy,
				},
				coordinate: entity.Coordinate{
					X: 25,
					Y: 125,
				},
			},
			initCanvas: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(ctrl)

				canvMock.EXPECT().TextUnescaped(float64(25), float64(131),
					`<tspan x="25.00" y="131.00" font-size="180%">&#x01D102;</tspan>`,
					`style="font-family:Noto Music"`)

				return canvMock
			},
		},
		//case 2: repeat left - forward
		{
			name: "repeat left forward",
			args: args{
				barline: musicxml.Barline{
					Location: musicxml.BarlineLocationLeft,
					BarStyle: musicxml.BarLineStyleHeavyLight,
					Repeat: &musicxml.BarLineRepeat{
						Direction: musicxml.BarLineRepeatDirectionForward,
					},
				},
				coordinate: entity.Coordinate{
					X: 25,
					Y: 125,
				},
			},
			initCanvas: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(ctrl)

				canvMock.EXPECT().TextUnescaped(
					float64(25), float64(131),
					`<tspan x="25.00" y="131.00" font-size="180%">&#x01D103;</tspan><tspan x="35.00" y="121.00">:</tspan>`,
					`style="font-family:Noto Music"`,
				)

				return canvMock
			},
		},
		// case 3: repeat right - backward
		{
			name: "repeat right backward",
			args: args{
				barline: musicxml.Barline{
					Location: musicxml.BarlineLocationRight,
					BarStyle: musicxml.BarLineStyleLightHeavy,
					Repeat: &musicxml.BarLineRepeat{
						Direction: musicxml.BarLineRepeatDirectionBackward,
					},
				},
				coordinate: entity.Coordinate{
					X: 25,
					Y: 125,
				},
			},
			initCanvas: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(ctrl)

				canvMock.EXPECT().TextUnescaped(float64(25), float64(131),
					`<tspan x="20.00" y="121.00">:</tspan><tspan x="25.00" y="131.00" font-size="180%">&#x01D102;</tspan>`,
					`style="font-family:Noto Music"`)

				return canvMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			bi := barlineInteractor{}
			bi.RenderBarline(context.Background(), tt.initCanvas(ctrl), tt.args.barline, tt.args.coordinate)
		})
	}
}

func TestNewBarline(t *testing.T) {

	t.Run("everything is went fine", func(t *testing.T) {
		if got := NewBarline(); !assert.NotNil(t, got) {
			t.Fail()
		}
	})

}

func TestGetBarlineWidth(t *testing.T) {

	t.Run("GetBarLineWidth", func(t *testing.T) {
		if got := GetBarlineWidth(musicxml.BarLineStyleHeavyHeavy); !assert.Equal(t, float64(8), got) {
			t.Fail()
		}
	})
}
