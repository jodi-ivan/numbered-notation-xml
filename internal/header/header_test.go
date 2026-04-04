package header

import (
	"context"
	"database/sql"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_headerInteractor_renderTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv      func(*gomock.Controller) *canvas.MockCanvas
		lyricMock func(*gomock.Controller) *lyric.MockLyric
		credit    []musicxml.Credit
		metadata  *repository.HymnMetadata
	}{
		{
			name: "no metadata",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeTitle,
					Words: "Unit Test",
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(320, 35, "UNIT TEST")
				return canv
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("UNIT TEST").Return(80.0)
				return li
			},
		},
		{
			name: "has metadata - title only",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeTitle,
					Words: "Unit Test",
				},
			},
			metadata: &repository.HymnMetadata{
				HymnData: repository.HymnData{
					Title: "Unit test title only",
					HymnIndicator: repository.HymnIndicator{
						Number: 1,
					},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(300, 35, "1. UNIT TEST TITLE ONLY")
				return canv
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("1. UNIT TEST TITLE ONLY").Return(120.0)
				return li
			},
		},
		{
			name: "has metadata - title only with varinat",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeTitle,
					Words: "Unit Test",
				},
			},
			metadata: &repository.HymnMetadata{
				HymnData: repository.HymnData{
					Title: "Unit test title with variant",
					HymnIndicator: repository.HymnIndicator{
						Number:  1,
						Variant: sql.NullString{Valid: true, String: "a"},
					},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(296, 35, "1a. UNIT TEST TITLE WITH VARIANT")
				return canv
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("1a. UNIT TEST TITLE WITH VARIANT").Return(128.0)
				return li
			},
		},
		{
			name: "has metadata - title only footnoes and for kids",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeTitle,
					Words: "Unit Test",
				},
			},
			metadata: &repository.HymnMetadata{
				HymnData: repository.HymnData{
					Title: "Unit test title only with footnotes and for kids",
					HymnIndicator: repository.HymnIndicator{
						Number: 1,
					},
					IsForKids:      sql.NullInt16{Valid: true, Int16: 1},
					TitleFootnotes: sql.NullString{Valid: true, String: "Bisa juga dinyanyikan dengan lagu unit test"},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(296, 35, "1. UNIT TEST TITLE ONLY WITH FOOTNOTES AND FOR KIDS *")
				canv.EXPECT().TextUnescaped(float64(50), float64(35), `<tspan font-style="bold" font-size="125%">☆</tspan>`)
				return canv
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("1. UNIT TEST TITLE ONLY WITH FOOTNOTES AND FOR KIDS *").Return(128.0)
				return li
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var hi headerInteractor
			if tt.lyricMock != nil {
				hi.Lyric = tt.lyricMock(ctrl)
			}
			canv := canvas.Canvas(nil)
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			hi.renderTitle(context.Background(), canv, tt.credit, tt.metadata)
		})
	}
}

func Test_headerInteractor_renderSubtitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv      func(*gomock.Controller) *canvas.MockCanvas
		lyricMock func(*gomock.Controller) *lyric.MockLyric
		credit    []musicxml.Credit
		metadata  *repository.HymnMetadata
	}{
		{
			name: "no subtitle - empty string",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeSubtitle,
					Words: "",
				},
			},
		},
		{
			name: "no subtitle - default string",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeSubtitle,
					Words: "Subtitle",
				},
			},
		},
		{
			name: "has subtitle, no metadata",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeSubtitle,
					Words: "(KANON)",
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(335, 53, "(KANON)", `style="font-size:70%;font-family:'Figtree';font-weight:600"`)
				return canv
			},
		},
		{
			name: "has subtitle, has metadata",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeSubtitle,
					Words: "(KANON)",
				},
			},
			metadata: &repository.HymnMetadata{
				HymnData: repository.HymnData{
					HymnIndicator: repository.HymnIndicator{
						Number: 2,
					},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(343, 53, "(KANON)", `style="font-size:70%;font-family:'Figtree';font-weight:600"`)
				return canv
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("2. ").Return(16.0)
				return li
			},
		},
		{
			name: "has subtitle, has metadata with variant",
			credit: []musicxml.Credit{
				{
					Type:  musicxml.CreditTypeSubtitle,
					Words: "(KANON)",
				},
			},
			metadata: &repository.HymnMetadata{
				HymnData: repository.HymnData{
					HymnIndicator: repository.HymnIndicator{
						Number:  2,
						Variant: sql.NullString{Valid: true, String: "a"},
					},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(347, 53, "(KANON)", `style="font-size:70%;font-family:'Figtree';font-weight:600"`)
				return canv
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("2. ").Return(18.0)
				li.EXPECT().CalculateLyricWidth("2a. ").Return(24.0)
				return li
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var hi headerInteractor
			if tt.lyricMock != nil {
				hi.Lyric = tt.lyricMock(ctrl)
			}
			canv := canvas.Canvas(nil)
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			hi.renderSubtitle(context.Background(), canv, tt.credit, tt.metadata)
		})
	}
}
