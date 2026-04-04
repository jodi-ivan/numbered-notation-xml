package lyric

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

func Test_lyricInteractor_CalculateHypen(t *testing.T) {
	type args struct {
		prevLyric    *LyricPosition
		currentLyric *LyricPosition
	}
	tests := []struct {
		name         string
		args         args
		wantLocation []entity.Coordinate
	}{
		{
			name: "no need hypen between two lyrics",
			args: args{
				prevLyric: &LyricPosition{
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeEnd,
					},
				},
			},
		},
		{
			name: "gap too small to fit a hypen",
			args: args{
				prevLyric: &LyricPosition{
					Coordinate: entity.NewCoordinate(25, 120),

					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeBegin,
						Text: []entity.Text{
							{Value: "hel"},
						},
					},
				},
				currentLyric: &LyricPosition{
					Coordinate: entity.NewCoordinate(45.26, 120),
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeEnd,
						Text: []entity.Text{
							{Value: "lo"},
						},
					},
				},
			},
		},
		{
			name: "short distance between two lyrics",
			args: args{
				prevLyric: &LyricPosition{
					Coordinate: entity.NewCoordinate(25, 120),
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeBegin,
						Text: []entity.Text{
							{Value: "hel"}, // width 20.46

						},
					},
				},
				currentLyric: &LyricPosition{
					Coordinate: entity.NewCoordinate(53, 120),
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeEnd,
						Text: []entity.Text{
							{Value: "lo"},
						},
					},
				},
			},
			wantLocation: []entity.Coordinate{
				entity.NewCoordinate(45.26, 120),
			},
		},
		{
			name: "long distance between two lyrics",
			args: args{
				prevLyric: &LyricPosition{
					Coordinate: entity.NewCoordinate(25, 120),
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeBegin,
						Text: []entity.Text{
							{Value: "hel"}, // width 20.46
						},
					},
				},
				currentLyric: &LyricPosition{
					Coordinate: entity.NewCoordinate(350, 120),
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeEnd,
						Text: []entity.Text{
							{Value: "lo"},
						},
					},
				},
			},
			wantLocation: []entity.Coordinate{
				entity.NewCoordinate(96.05, 120),
				entity.NewCoordinate(146.84, 120),
				entity.NewCoordinate(197.63, 120),
				entity.NewCoordinate(248.42000000000002, 120),
				entity.NewCoordinate(299.21, 120),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := lyricInteractor{}
			if gotLocation := li.CalculateHypen(context.Background(), tt.args.prevLyric, tt.args.currentLyric); !assert.Equal(t, tt.wantLocation, gotLocation) {
				t.Errorf("lyricInteractor.CalculateHypen() = %v, want %v", gotLocation, tt.wantLocation)
			}
		})
	}
}

func Test_lyricInteractor_RenderHypen(t *testing.T) {
	type args struct {
		measure []*entity.NoteRenderer
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		initCanvMock func(*gomock.Controller) *canvas.MockCanvas
		args         args
	}{
		// case 1: empty measure
		{
			name: "empty measure",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				// writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				// canv.EXPECT().Writer().Return(writerMock)
				canv.EXPECT().Gend()
				return canv

			},
		},
		// case 2: empty lyric
		{
			name: "empty lyric",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				// writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				// canv.EXPECT().Writer().Return(writerMock)
				canv.EXPECT().Gend()
				return canv

			},
			args: args{
				measure: []*entity.NoteRenderer{
					{},
				},
			},
		},
		{
			name: "Positive case. begin-middle-end",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().TextUnescaped(78.8850, float64(25), "-")
				canv.EXPECT().TextUnescaped(33.6500, float64(25), "-")
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					{
						PositionX: 25,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									entity.Text{
										Value: "E",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					{
						PositionX: 40,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{
										Value: "xam",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					{
						PositionX: 100,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{
										Value: "ple",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
				},
			},
		},
		{
			name: "Positive case. begin-middle-end-single",
			initCanvMock: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().TextUnescaped(78.885, float64(25), "-")
				canv.EXPECT().TextUnescaped(33.6500, float64(25), "-")
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					{
						PositionX: 25,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "E"},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					{
						PositionX: 40,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "xam"},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					{
						PositionX: 100,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "ple"},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
					{
						PositionX: 160,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "duh"},
								},
								Syllabic: musicxml.LyricSyllabicTypeSingle,
							},
						},
					},
				},
			},
		},
		{
			name: "Positive case. begin-middle-middle-end",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().TextUnescaped(120.7500, float64(25), "-")
				canv.EXPECT().TextUnescaped(90.0900, float64(25), "-")
				canv.EXPECT().TextUnescaped(55.0350, float64(25), "-")
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					{
						PositionX: 25,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "Na"},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					{
						PositionX: 77,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "si"},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					{
						PositionX: 103,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "go"},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					{
						PositionX: 134,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "reng"},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
				},
			},
		},
		{
			name: "Positive case. begin-end-begin-middle-end",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().TextUnescaped(161.1300, float64(25), "-")
				canv.EXPECT().TextUnescaped(127.0650, float64(25), "-")
				canv.EXPECT().TextUnescaped(49.0950, float64(25), "-")
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					{
						PositionX: 25,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "Ma"},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					{
						PositionX: 62,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "kan"},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
					{
						PositionX: 102,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "Ber"},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					{
						PositionX: 140,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "sa"},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					{
						PositionX: 179,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "ma"},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
				},
			},
		},
		{
			name: "Positive case. begin-middle-end-single-begin-end",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().TextUnescaped(243.0950, float64(25), "-")
				canv.EXPECT().TextUnescaped(79.1300, float64(25), "-")
				canv.EXPECT().TextUnescaped(50.0650, float64(25), "-")
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					{
						PositionX: 25,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "Ber"},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					{
						PositionX: 63,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "sa"},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					{
						PositionX: 92,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "ma"},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
					{
						PositionX: 179,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "dan"},
								},
								Syllabic: musicxml.LyricSyllabicTypeSingle,
							},
						},
					},
					{
						PositionX: 219,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "Ma"},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					{
						PositionX: 256,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "kan"},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
				},
			},
		},
		{
			name: "Positive case. middle-middle-end",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().TextUnescaped(43.2050, float64(25), "-")
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					{
						PositionX: 25,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "ka"},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					{
						PositionX: 56,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "sum"},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
					{
						PositionX: 96,
						Lyric: []entity.Lyric{
							{
								Text: []entity.Text{
									{Value: "Boom"},
								},
								Syllabic: musicxml.LyricSyllabicTypeSingle,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := &lyricInteractor{}
			li.RenderHypen(context.Background(), tt.initCanvMock(ctrl), tt.args.measure)
		})
	}
}
