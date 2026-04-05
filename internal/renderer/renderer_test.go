package renderer

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/credits"
	"github.com/jodi-ivan/numbered-notation-xml/internal/footnote"
	"github.com/jodi-ivan/numbered-notation-xml/internal/header"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/verse"
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
	creditsData := []musicxml.Credit{
		{Type: musicxml.CreditTypeTitle, Words: "Unit Test"},
	}

	measures := []musicxml.Measure{
		musicxml.Measure{
			Number: 1,
			Attribute: &musicxml.Attribute{
				Key: &musicxml.KeySignature{
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

	keySignature := keysig.NewKeySignature(context.Background(), measures)
	timeSignature := timesig.NewTimeSignatures(context.Background(), measures)

	metadata := &repository.HymnMetadata{
		HymnData: repository.HymnData{
			HymnIndicator: repository.HymnIndicator{
				Number: 1,
			},
			Title: "Unittest",
		},
		Verse: map[int]repository.HymnVerse{},
	}
	/*
		httpmock.RegisterResponder("GET", "https://fonts.googleapis.com/css?family=Caladea%7COld+Standard+TT%7CNoto+Music%7CFigtree",
					httpmock.NewStringResponder(200, googleFont))
	*/

	type args struct {
		music    musicxml.MusicXML
		metadata *repository.HymnMetadata
	}
	tests := []struct {
		name string
		args args

		canvasMock func(ctrl *gomock.Controller) *canvas.MockCanvas
		headerMock func(ctrl *gomock.Controller) *header.MockHeader
		staffMock  func(ctrl *gomock.Controller) *staff.MockStaff

		// optional
		creditsMock  func(ctrl *gomock.Controller) *credits.MockCredits
		footnoteMock func(ctrl *gomock.Controller) *footnote.MockFootnote
		verseMock    func(ctrl *gomock.Controller) *verse.MockVerse
	}{
		{
			name: "empty metadata",
			args: args{
				music: musicxml.MusicXML{
					Part: musicxml.Part{
						Measures: measures,
					},
					Credit: creditsData,
				},
			},
			canvasMock: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				writerMock := canvas.NewMockWriter(c)

				httpmock.RegisterResponder("GET", "https://fonts.googleapis.com/css?family=Caladea%7COld+Standard+TT%7CNoto+Music%7CFigtree",
					httpmock.NewStringResponder(200, googleFont))
				style := fmt.Sprintf(fontfmt, googleFont)

				canv.EXPECT().Start(720, 2000)
				canv.EXPECT().Def()
				canv.EXPECT().Writer().Return(writerMock)
				writerMock.EXPECT().Write([]byte(style))
				canv.EXPECT().DefEnd()
				canv.EXPECT().End()

				return canv
			},

			headerMock: func(c *gomock.Controller) *header.MockHeader {
				mockHeader := header.NewMockHeader(c)

				mockHeader.EXPECT().RenderSheetHeader(gomock.Any(), gomock.Any(), creditsData, nil)
				mockHeader.EXPECT().RenderKeyandTimeSignatures(gomock.Any(), gomock.Any(), keySignature, timeSignature)
				return mockHeader
			},
			staffMock: func(c *gomock.Controller) *staff.MockStaff {
				mockStaff := staff.NewMockStaff(c)
				mockStaff.EXPECT().Render(gomock.Any(), gomock.Any(), musicxml.Part{Measures: measures}, keySignature, timeSignature)
				return mockStaff
			},
		},
		{
			name: "with metadata",
			args: args{
				metadata: &repository.HymnMetadata{
					HymnData: repository.HymnData{
						HymnIndicator: repository.HymnIndicator{
							Number: 1,
						},
						Title: "Unittest",
					},
					Verse: map[int]repository.HymnVerse{},
				},
				music: musicxml.MusicXML{
					Part: musicxml.Part{
						Measures: measures,
					},
					Credit: creditsData,
				},
			},
			canvasMock: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				writerMock := canvas.NewMockWriter(c)

				httpmock.RegisterResponder("GET", "https://fonts.googleapis.com/css?family=Caladea%7COld+Standard+TT%7CNoto+Music%7CFigtree",
					httpmock.NewStringResponder(200, googleFont))
				style := fmt.Sprintf(fontfmt, googleFont)

				canv.EXPECT().Start(720, 2000)
				canv.EXPECT().Def()
				canv.EXPECT().Writer().Return(writerMock)
				writerMock.EXPECT().Write([]byte(style))
				canv.EXPECT().DefEnd()
				canv.EXPECT().End()

				return canv
			},

			headerMock: func(c *gomock.Controller) *header.MockHeader {
				mockHeader := header.NewMockHeader(c)

				mockHeader.EXPECT().RenderSheetHeader(gomock.Any(), gomock.Any(), creditsData, metadata)
				mockHeader.EXPECT().RenderKeyandTimeSignatures(gomock.Any(), gomock.Any(), keySignature, timeSignature)
				return mockHeader
			},
			staffMock: func(c *gomock.Controller) *staff.MockStaff {
				mockStaff := staff.NewMockStaff(c)
				mockStaff.EXPECT().Render(gomock.Any(), gomock.Any(), musicxml.Part{Measures: measures}, keySignature, timeSignature).Return(100)
				return mockStaff
			},

			footnoteMock: func(c *gomock.Controller) *footnote.MockFootnote {
				fm := footnote.NewMockFootnote(c)
				pos := 150
				fm.EXPECT().RenderMusicFootnotes(gomock.Any(), gomock.Any(), metadata, 100)
				fm.EXPECT().RenderVerseFootnotes(gomock.Any(), &pos, metadata.VerseFootNotes)
				fm.EXPECT().RenderTitleFootnotes(gomock.Any(), 150, metadata.HymnData)
				return fm
			},

			verseMock: func(c *gomock.Controller) *verse.MockVerse {
				vm := verse.NewMockVerse(c)
				vm.EXPECT().RenderVerse(gomock.Any(), gomock.Any(), 100, metadata.Verse, metadata.VerseFootNotes).Return(verse.VerseInfo{MarginBottom: 150})
				return vm

			},

			creditsMock: func(c *gomock.Controller) *credits.MockCredits {
				cm := credits.NewMockCredits(c)
				pos := 150
				cm.EXPECT().RenderCredits(gomock.Any(), gomock.Any(), &pos, metadata.HymnData)
				return cm
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ir := rendererInteractor{
				Staff:  tt.staffMock(ctrl),
				Header: tt.headerMock(ctrl),
			}
			if tt.creditsMock != nil {
				ir.Credits = tt.creditsMock(ctrl)

			}
			if tt.footnoteMock != nil {
				ir.Footnote = tt.footnoteMock(ctrl)

			}
			if tt.verseMock != nil {
				ir.Verse = tt.verseMock(ctrl)
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
				assert.NotNil(t, cast.Footnote)
				assert.NotNil(t, cast.Verse)
				assert.NotNil(t, cast.Header)

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
