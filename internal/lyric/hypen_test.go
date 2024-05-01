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
		/*
			Syllabic:
			 - Begin
			 - Middle
			 - End

			Example:

			E ....... xam .... ple
			^         ^        ^
			begin     middle   end
		*/
		// case 3a: positive case. order Begin-middle-end
		{},
		// case 3: positive case 1: all in the same line
		// -- syllabic order:
		// case 3a: begin - middle - middle - end
		// case 3b: begin - end - begin - middle - end
		// case 3c: begin - middle - end - single - begin - end
		// case 3d:	single - single - begin - middle - middle
		// case 4: positive case 2: there is a new line
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := &lyricInteractor{}
			li.RenderHypen(context.Background(), tt.initCanvMock(ctrl), tt.args.measure)
		})
	}
}
