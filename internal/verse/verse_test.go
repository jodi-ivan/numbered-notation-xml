package verse

import (
	"context"
	"database/sql"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/footnote"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

func Test_verseInteractor_elisionPosition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		p                     LyricPartVerse
		y                     int
		lineBeforeWord        string
		syllableBeforeElision string
		lyricMock             func(c *gomock.Controller) *lyric.MockLyric
		want                  [2]entity.Coordinate
	}{
		{
			name: "with underline",
			p: LyricPartVerse{
				Breakdown: []LyricStylePart{
					{
						Text: "t",
					},
					{
						Text:      "ahu",
						Underline: true,
					},
				},
			},
			lineBeforeWord: "unittest",
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("t").Return(3.0)
				li.EXPECT().CalculateLyricWidth("u").Return(3.0)
				li.EXPECT().CalculateLyricWidth("ahu").Return(12.0)
				li.EXPECT().CalculateLyricWidth("unittest").Return(24.0)
				li.EXPECT().CalculateLyricWidth("").Return(0.0)

				return li
			},
			want: [2]entity.Coordinate{
				entity.NewCoordinate(27, 0),
				entity.NewCoordinate(37.5, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v verseInteractor
			v.Lyric = tt.lyricMock(ctrl)
			got := v.elisionPosition(tt.p, tt.y, tt.lineBeforeWord, tt.syllableBeforeElision)
			assert.Equal(t, tt.want, got, "elisionPosition()")

		})
	}
}

func Test_verseInteractor_parse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		y         int
		verses    map[int]repository.HymnVerse
		lyricMock func(c *gomock.Controller) *lyric.MockLyric

		want ParsedVerseWithInfo
	}{
		{
			name: "default",
			y:    300,
			verses: map[int]repository.HymnVerse{
				2: {
					Content:  sql.NullString{Valid: true, String: `[[{"word":"Dalam","breakdown":[{"text":"Da","type":"begin","combine":false,"breakdown":null},{"text":"lam","type":"end","combine":false,"breakdown":null}],"dash":false},{"word":"dunia","breakdown":[{"text":"du","type":"begin","combine":false,"breakdown":null},{"text":"nia","type":"end","combine":true,"breakdown":[{"text":"n","underline":false},{"text":"ia","underline":true}]}],"dash":false},{"word":"'ku","breakdown":[{"text":"'ku","type":"single","combine":false,"breakdown":null}],"dash":false},{"word":"dikawal","breakdown":[{"text":"di","type":"begin","combine":false,"breakdown":null},{"text":"ka","type":"middle","combine":false,"breakdown":null},{"text":"wal","type":"end","combine":false,"breakdown":null}],"dash":false}]]`},
					Row:      sql.NullInt16{Int16: 1, Valid: true},
					StyleRow: sql.NullInt32{Int32: 6, Valid: true},
					Col:      sql.NullInt16{Int16: 2, Valid: true},
					VerseNum: sql.NullInt32{Int32: 2, Valid: true},
				},
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("n").Return(3.0)
				li.EXPECT().CalculateLyricWidth("ia").Return(6.0)
				li.EXPECT().CalculateLyricWidth("a").Return(3.0)

				li.EXPECT().CalculateLyricWidth(" Dalam").Return(18.0)
				li.EXPECT().CalculateLyricWidth("du").Return(6.0)
				li.EXPECT().CalculateLyricWidth(" Dalam dunia 'ku dikawal").Return(90.0)

				return li
			},
			want: ParsedVerseWithInfo{
				Verses: map[int]ParsedVerse{
					2: {
						Verse: []string{" Dalam dunia 'ku dikawal"},
						ElisionMarks: [][2]entity.Coordinate{
							{
								entity.NewCoordinate(27, 0),
								entity.NewCoordinate(31.5, 0),
							},
						},
						Position: versePosition{
							Col:      2,
							Row:      1,
							RowWidth: 6,
							Style:    6,
						},
					},
				},
				IsMultiColumn: true,
				MaxLineWidth:  90,
				RowPositionY: map[int]int{
					1: 300,
				},
				MaxRightPos: 90,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v verseInteractor
			v.Lyric = tt.lyricMock(ctrl)
			got := v.parse(tt.y, tt.verses)

			assert.Equal(t, tt.want, got, "parse()")
		})
	}
}

func Test_verseInteractor_RenderVerse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv         func(c *gomock.Controller) *canvas.MockCanvas
		lyricMock    func(c *gomock.Controller) *lyric.MockLyric
		footnoteMock func(c *gomock.Controller) *footnote.MockFootnote

		y             int
		verses        map[int]repository.HymnVerse
		verseFootnote map[int]map[int]repository.VerseFootNotes
		want          VerseInfo
	}{
		{
			name: "default",
			y:    300,
			verses: map[int]repository.HymnVerse{
				2: {
					Content:  sql.NullString{Valid: true, String: `[[{"word":"Dalam","breakdown":[{"text":"Da","type":"begin","combine":false,"breakdown":null},{"text":"lam","type":"end","combine":false,"breakdown":null}],"dash":false},{"word":"dunia","breakdown":[{"text":"du","type":"begin","combine":false,"breakdown":null},{"text":"nia","type":"end","combine":true,"breakdown":[{"text":"n","underline":false},{"text":"ia","underline":true}]}],"dash":false},{"word":"'ku","breakdown":[{"text":"'ku","type":"single","combine":false,"breakdown":null}],"dash":false},{"word":"dikawal","breakdown":[{"text":"di","type":"begin","combine":false,"breakdown":null},{"text":"ka","type":"middle","combine":false,"breakdown":null},{"text":"wal","type":"end","combine":false,"breakdown":null}],"dash":false}]]`},
					Row:      sql.NullInt16{Int16: 1, Valid: true},
					StyleRow: sql.NullInt32{Int32: 6, Valid: true},
					Col:      sql.NullInt16{Int16: 2, Valid: true},
					VerseNum: sql.NullInt32{Int32: 2, Valid: true},
				},
			},
			want: VerseInfo{MarginBottom: 330},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("n").Return(3.0)
				li.EXPECT().CalculateLyricWidth("ia").Return(6.0)
				li.EXPECT().CalculateLyricWidth("a").Return(3.0)
				li.EXPECT().CalculateLyricWidth("2. ").Return(9.0)

				li.EXPECT().CalculateLyricWidth(" Dalam").Return(18.0)
				li.EXPECT().CalculateLyricWidth("du").Return(6.0)
				li.EXPECT().CalculateLyricWidth(" Dalam dunia 'ku dikawal").Return(90.0)

				return li
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='verses'", "style='font-family:Caladea'")
				canv.EXPECT().Group("class='verse'", "number='2'")
				canv.EXPECT().Text(541, 300, "2. ")
				canv.EXPECT().Text(555, 300, " Dalam dunia 'ku dikawal")

				canv.EXPECT().Group()
				canv.EXPECT().Qbez(582, 302, 584, 307, 586, 302, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend().Times(3)
				return canv
			},
			footnoteMock: func(c *gomock.Controller) *footnote.MockFootnote {
				fn := footnote.NewMockFootnote(c)
				cursor := footnote.VerseLineCursor{
					VerseNo:    2,
					LinePos:    1,
					Leftmargin: 455,
					LineText:   " Dalam dunia 'ku dikawal",
				}
				fn.EXPECT().AssignFootnotesMarker(gomock.Any(), entity.NewCoordinate(100, 300), 315, cursor, nil)
				return fn
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v verseInteractor
			v.Lyric = tt.lyricMock(ctrl)
			v.Footnote = tt.footnoteMock(ctrl)
			got := v.RenderVerse(context.Background(), tt.canv(ctrl), tt.y, tt.verses, tt.verseFootnote)

			assert.Equal(t, tt.want, got)
		})
	}
}
