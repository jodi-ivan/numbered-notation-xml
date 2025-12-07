package renderer

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/credits"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

func Test_rendererInteractor_Render(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	googleFont := `
	@font-face {
		font-family: 'Caladea';
		font-style: normal;
		font-weight: 400;
		src: url(https://fonts.gstatic.com/s/caladea/v1/default.ttf) format('truetype');
	  }
	  @font-face {
		font-family: 'Figtree';
		font-style: normal;
		font-weight: 400;
		src: url(https://fonts.gstatic.com/s/figtree/v1/default.ttf) format('truetype');
	  }
	  @font-face {
		font-family: 'Noto Music';
		font-style: normal;
		font-weight: 400;
		src: url(https://fonts.gstatic.com/s/notomusic/v1/default.ttf) format('truetype');
	  }
	  @font-face {
		font-family: 'Old Standard TT';
		font-style: normal;
		font-weight: 400;
		src: url(https://fonts.gstatic.com/s/oldstandardtt/v1/defailt.ttf) format('truetype');
	  }
	  `

	measures := []musicxml.Measure{
		musicxml.Measure{
			Number: 1,
			Attribute: &musicxml.Attribute{
				Key: musicxml.KeySignature{
					Fifth: 2, // D major
				},
				Time: &struct {
					Beats    int `xml:"beats"`
					BeatType int `xml:"beat-type"`
				}{
					Beats:    4,
					BeatType: 4,
				},
			},
		},
		musicxml.Measure{
			Number: 2,
			Print: &musicxml.Print{
				NewSystem: musicxml.PrintNewSystemTypeYes,
			},
		},
		musicxml.Measure{
			Number: 3,
		},
		musicxml.Measure{
			Number: 4,
			Print: &musicxml.Print{
				NewSystem: musicxml.PrintNewSystemTypeYes,
			},
		},
		musicxml.Measure{
			Number: 5,
		},
	}

	spilttedLine := [][]musicxml.Measure{
		[]musicxml.Measure{musicxml.Measure{
			Number: 1,
			Attribute: &musicxml.Attribute{
				Key: musicxml.KeySignature{
					Fifth: 2, // D major
				},
				Time: &struct {
					Beats    int `xml:"beats"`
					BeatType int `xml:"beat-type"`
				}{
					Beats:    4,
					BeatType: 4,
				},
			},
		}},
		[]musicxml.Measure{
			musicxml.Measure{
				Number: 2,
				Print: &musicxml.Print{
					NewSystem: musicxml.PrintNewSystemTypeYes,
				},
			},
			musicxml.Measure{
				Number: 3,
			},
		},
		[]musicxml.Measure{
			musicxml.Measure{
				Number: 4,
				Print: &musicxml.Print{
					NewSystem: musicxml.PrintNewSystemTypeYes,
				},
			},
			musicxml.Measure{
				Number: 5,
			},
		},
	}

	type args struct {
		music    musicxml.MusicXML
		metadata *repository.HymnMetadata
	}
	tests := []struct {
		name string
		args args

		lyricMock   func(ctrl *gomock.Controller) *lyric.MockLyric
		staffMock   func(ctrl *gomock.Controller) *staff.MockStaff
		creditsMock func(ctrl *gomock.Controller) *credits.MockCredits
		canvasMock  func(ctrl *gomock.Controller) *canvas.MockCanvas
	}{
		{
			name: "default",
			args: args{
				metadata: &repository.HymnMetadata{
					HymnData: repository.HymnData{
						HymnIndicator: repository.HymnIndicator{
							Number: 1,
						},
						Title: "Unittest",
					},
					Verse: []repository.HymnVerse{},
				},
				music: musicxml.MusicXML{
					Part: musicxml.Part{
						Measures: measures,
					},
				},
			},
			canvasMock: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)

				httpmock.RegisterResponder("GET", "https://fonts.googleapis.com/css?family=Caladea%7COld+Standard+TT%7CNoto+Music%7CFigtree",
					httpmock.NewStringResponder(200, googleFont))

				writerMock := canvas.NewMockWriter(ctrl)
				style := fmt.Sprintf(fontfmt, googleFont)
				writerMock.EXPECT().Write([]byte(style))

				canv.EXPECT().Start(int(720), int(2000))
				canv.EXPECT().Def()
				canv.EXPECT().Writer().Return(writerMock)
				canv.EXPECT().DefEnd()
				canv.EXPECT().Text(int(310), int(100), "1. UNITTEST")

				canv.EXPECT().Text(int(50), int(125), "do = d")
				canv.EXPECT().Text(int(195), int(125), "4 ketuk")
				canv.EXPECT().End()

				return canv
			},
			lyricMock: func(ctrl *gomock.Controller) *lyric.MockLyric {
				l := lyric.NewMockLyric(ctrl)

				l.EXPECT().CalculateLyricWidth("do = d").Return(float64(100))

				l.EXPECT().CalculateLyricWidth("1. UNITTEST").Return(float64(100))
				l.EXPECT().RenderVerse(gomock.Any(), gomock.Any(), int(490), []repository.HymnVerse{}).Return(lyric.VerseInfo{MarginBottom: 40})

				return l
			},
			staffMock: func(ctrl *gomock.Controller) *staff.MockStaff {
				mockStaff := staff.NewMockStaff(ctrl)
				mockStaff.EXPECT().SplitLines(gomock.Any(), musicxml.Part{Measures: measures}).Return(spilttedLine)

				currTimeSig := timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{Measure: 1, Beat: 4, BeatType: 4},
					},
				}

				currTimeSig.GetHumanized()

				mockStaff.EXPECT().RenderStaff(
					gomock.Any(),
					gomock.Any(),
					int(50),
					int(175),
					keysig.NewKeySignature(musicxml.KeySignature{Fifth: 2}),
					currTimeSig,
					spilttedLine[0],
				).Return(staff.StaffInfo{
					MarginBottom: 25,
				})
				mockStaff.EXPECT().RenderStaff(
					gomock.Any(),
					gomock.Any(),
					int(50),  // x
					int(280), // y
					keysig.NewKeySignature(musicxml.KeySignature{Fifth: 2}),
					currTimeSig,
					spilttedLine[1],
				).Return(staff.StaffInfo{
					MarginBottom: 25,
					Multiline:    true,
					MarginLeft:   100,
				})
				mockStaff.EXPECT().RenderStaff(
					gomock.Any(),
					gomock.Any(),
					int(100), // x
					int(385), // y
					keysig.NewKeySignature(musicxml.KeySignature{Fifth: 2}),
					currTimeSig,
					spilttedLine[2],
				).Return(staff.StaffInfo{
					MarginBottom: 25,
					Multiline:    false,
				})
				return mockStaff
			},
			creditsMock: func(ctrl *gomock.Controller) *credits.MockCredits {
				cMock := credits.NewMockCredits(ctrl)

				cMock.EXPECT().RenderCredits(gomock.Any(), gomock.Any(), int(40), repository.HymnData{
					HymnIndicator: repository.HymnIndicator{
						Number: 1,
					},
					Title: "Unittest",
				})

				return cMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ir := rendererInteractor{
				Lyric:   tt.lyricMock(ctrl),
				Staff:   tt.staffMock(ctrl),
				Credits: tt.creditsMock(ctrl),
			}
			ir.Render(context.Background(), tt.args.music, tt.canvasMock(ctrl), tt.args.metadata)
		})
	}
}

func TestNewRenderer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRenderer()

			if assert.IsType(t, got, &rendererInteractor{}) {
				cast := got.(*rendererInteractor)
				assert.NotNil(t, cast.Lyric)
				assert.NotNil(t, cast.Staff)
				assert.NotNil(t, cast.Credits)

			}

		})
	}
}

func Test_googlefont(t *testing.T) {

	calibriFont := `
		@font-face {
			font-family: 'Caladea';
			font-style: normal;
			font-weight: 400;
			src: url(https://fonts.gstatic.com/s/calibri/v1/default.ttf) format('truetype');
			}
	`

	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		init func()
		want []byte
	}{
		{
			name: "failed to get the font",
			args: args{
				f: "Calibri",
			},
			init: func() {
				httpmock.RegisterResponder("GET", "https://fonts.googleapis.com/css?family=Calibri",
					httpmock.NewErrorResponder(errors.New("nope")))
			},
			want: []byte{},
		},
		{
			name: "non 200",
			args: args{
				f: "Calibri",
			},
			init: func() {
				httpmock.RegisterResponder("GET", "https://fonts.googleapis.com/css?family=Calibri",
					httpmock.NewStringResponder(999, "what is this response?"))
			},
			want: []byte{},
		},
		{
			name: "Everything went fine",
			args: args{
				f: "Calibri",
			},
			init: func() {
				httpmock.RegisterResponder("GET", "https://fonts.googleapis.com/css?family=Calibri",
					httpmock.NewStringResponder(200, calibriFont))
			},
			want: []byte(calibriFont),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			tt.init()

			got := googlefont(tt.args.f)
			assert.Equal(t, string(tt.want), string(got))
		})
	}
}
