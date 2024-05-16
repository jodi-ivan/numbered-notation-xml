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
					Coordinate: entity.Coordinate{
						X: 25,
						Y: 120,
					},
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeBegin,
						Text: []entity.Text{
							entity.Text{
								Value: "hel",
							},
						},
					},
				},
				currentLyric: &LyricPosition{
					Coordinate: entity.Coordinate{
						X: 45.26,
						Y: 120,
					},
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeEnd,
						Text: []entity.Text{
							entity.Text{
								Value: "lo",
							},
						},
					},
				},
			},
		},
		{
			name: "short distance between two lyrics",
			args: args{
				prevLyric: &LyricPosition{
					Coordinate: entity.Coordinate{
						X: 25,
						Y: 120,
					},
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeBegin,
						Text: []entity.Text{
							entity.Text{
								Value: "hel", // width 20.46
							},
						},
					},
				},
				currentLyric: &LyricPosition{
					Coordinate: entity.Coordinate{
						X: 53,
						Y: 120,
					},
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeEnd,
						Text: []entity.Text{
							entity.Text{
								Value: "lo",
							},
						},
					},
				},
			},
			wantLocation: []entity.Coordinate{
				entity.Coordinate{
					X: 45.26,
					Y: 120,
				},
			},
		},
		{
			name: "long distance between two lyrics",
			args: args{
				prevLyric: &LyricPosition{
					Coordinate: entity.Coordinate{
						X: 25,
						Y: 120,
					},
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeBegin,
						Text: []entity.Text{
							entity.Text{
								Value: "hel", // width 20.46
							},
						},
					},
				},
				currentLyric: &LyricPosition{
					Coordinate: entity.Coordinate{
						X: 350,
						Y: 120,
					},
					Lyrics: entity.Lyric{
						Syllabic: musicxml.LyricSyllabicTypeEnd,
						Text: []entity.Text{
							entity.Text{
								Value: "lo",
							},
						},
					},
				},
			},
			wantLocation: []entity.Coordinate{
				entity.Coordinate{
					X: 45.26,
					Y: 120,
				},
				entity.Coordinate{
					X: 96.05,
					Y: 120,
				},
				entity.Coordinate{
					X: 146.84,
					Y: 120,
				},
				entity.Coordinate{
					X: 197.63,
					Y: 120,
				},
				entity.Coordinate{
					X: 248.42,
					Y: 120,
				},
				entity.Coordinate{
					X: 299.21,
					Y: 120,
				},
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
					&entity.NoteRenderer{},
				},
			},
		},
		{
			name: "Positive case. begin-middle-end",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().Writer().Return(writerMock).Times(2)
				writerMock.EXPECT().Write([]byte(`<text x="78.8850" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="33.6500" y="25">-</text>`))
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "E",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 40,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "xam",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 100,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
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
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().Writer().Return(writerMock).Times(2)
				writerMock.EXPECT().Write([]byte(`<text x="78.8850" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="33.6500" y="25">-</text>`))
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "E",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 40,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "xam",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 100,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "ple",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 160,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "duh",
									},
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
				writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().Writer().Return(writerMock).Times(3)
				writerMock.EXPECT().Write([]byte(`<text x="120.7500" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="90.0900" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="55.0350" y="25">-</text>`))
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Na",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 77,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "si",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 103,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "go",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 134,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "reng",
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
			name: "Positive case. begin-end-begin-middle-end",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().Writer().Return(writerMock).Times(3)
				writerMock.EXPECT().Write([]byte(`<text x="161.1300" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="127.0650" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="49.0950" y="25">-</text>`))
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Ma",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 62,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "kan",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 102,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Ber",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 140,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "sa",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 179,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "ma",
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
			name: "Positive case. begin-middle-end-single-begin-end",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().Writer().Return(writerMock).Times(3)
				writerMock.EXPECT().Write([]byte(`<text x="243.0950" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="79.1300" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="50.0650" y="25">-</text>`))
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Ber",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 63,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "sa",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 92,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "ma",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 179,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "dan",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeSingle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 219,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Ma",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 256,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "kan",
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
			name: "Positive case. single-single-begin-middle-middle",
			initCanvMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().Writer().Return(writerMock).Times(2)
				writerMock.EXPECT().Write([]byte(`<text x="187.2050" y="25">-</text>`))
				writerMock.EXPECT().Write([]byte(`<text x="163.6300" y="25">-</text>`))
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Boom",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeSingle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 78,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Boom",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeSingle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 140,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Sha",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 171,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "ka",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 198,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "la",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
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
				writerMock := canvas.NewMockWriter(ctrl)
				canv.EXPECT().Group("hyphens")
				canv.EXPECT().Writer().Return(writerMock)
				writerMock.EXPECT().Write([]byte(`<text x="35.2300" y="25">-</text>`))
				canv.EXPECT().Gend()
				return canv
			},
			args: args{
				measure: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "ka",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 56,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "sum",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
					&entity.NoteRenderer{
						PositionX: 96,
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "Boom",
									},
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
