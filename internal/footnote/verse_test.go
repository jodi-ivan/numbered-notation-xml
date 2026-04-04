package footnote

import (
	"database/sql"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

func Test_footnoteInteractor_AssignFootnotesMarker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv      func(*gomock.Controller) *canvas.MockCanvas
		lyricMock func(*gomock.Controller) *lyric.MockLyric

		cursor        VerseLineCursor
		pos           entity.Coordinate
		defaultX      int
		verseFootnote map[int]map[int]repository.VerseFootNotes
	}{
		{
			name: "no footnotes",
		},
		{
			name: "has footnotes, but no line",
			verseFootnote: map[int]map[int]repository.VerseFootNotes{
				2: {
					3: repository.VerseFootNotes{
						FootnoteMarker: sql.NullString{Valid: true, String: "*)"},
					},
				},
			},
			cursor: VerseLineCursor{
				VerseNo: 2,
				LinePos: 2,
			},
		},
		{
			name: "has footnotes, has line, but headless style",
			verseFootnote: map[int]map[int]repository.VerseFootNotes{
				2: {
					3: repository.VerseFootNotes{
						MarkerStyle:    sql.NullInt32{Valid: true, Int32: int32(VerseNoteStyleHeadless)},
						FootnoteMarker: sql.NullString{Valid: true, String: "*)"},
					},
				},
			},
			cursor: VerseLineCursor{
				VerseNo: 2,
				LinePos: 3,
			},
		},
		{
			name: "Right align style",
			verseFootnote: map[int]map[int]repository.VerseFootNotes{
				2: {
					3: repository.VerseFootNotes{
						MarkerStyle:    sql.NullInt32{Valid: true, Int32: int32(VerseNoteStyleAlignRight)},
						FootnoteMarker: sql.NullString{Valid: true, String: "1)"},
					},
				},
			},
			cursor: VerseLineCursor{
				VerseNo:    2,
				LinePos:    3,
				LineText:   "this is unit",
				Leftmargin: 20,
			},
			defaultX: 50,
			pos:      entity.NewCoordinate(50, 100),
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='footnotes'", `style="font-style:italic;font-size:60%;font-family:'Figtree';font-weight:600"`)
				canv.EXPECT().Text(50+620+20, int(100), "1)")
				canv.EXPECT().Gend()
				return canv
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("this is unit").Return(100.0)
				return li
			},
		},
		{
			name: "Direct append",
			verseFootnote: map[int]map[int]repository.VerseFootNotes{
				2: {
					3: repository.VerseFootNotes{
						MarkerStyle:    sql.NullInt32{Valid: true, Int32: int32(VerseNoteStyleDirectAppendText)},
						FootnoteMarker: sql.NullString{Valid: true, String: "*)"},
					},
				},
			},
			cursor: VerseLineCursor{
				VerseNo:    2,
				LinePos:    3,
				LineText:   "this is unit",
				Leftmargin: 20,
			},
			defaultX: 50,
			pos:      entity.NewCoordinate(50, 100),
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='footnotes'", `style="font-style:italic;font-family:'Caladea';font-size:90%;font-weight:600;"`)
				canv.EXPECT().Text(162, int(100), "*)")
				canv.EXPECT().Gend()
				return canv
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("this is unit").Return(100.0)
				li.EXPECT().CalculateLyricWidth(" ").Return(8.0)
				return li
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var fi footnoteInteractor
			if tt.lyricMock != nil {
				fi.li = tt.lyricMock(ctrl)
			}
			canv := canvas.Canvas(nil)
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			fi.AssignFootnotesMarker(canv, tt.pos, tt.defaultX, tt.cursor, tt.verseFootnote)
		})
	}
}

func Test_footnoteInteractor_RenderVerseFootnotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	getY := func() *int {
		y := 100

		return &y
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv func(*gomock.Controller) *canvas.MockCanvas

		y         *int
		wantY     int
		footnotes map[int]map[int]repository.VerseFootNotes
	}{
		{
			name: "No footnotes",
		},
		{
			name: "All of them are head only",
			footnotes: map[int]map[int]repository.VerseFootNotes{
				2: {
					1: repository.VerseFootNotes{
						MarkerStyle: sql.NullInt32{Valid: true, Int32: int32(VerseNoteStyleHeadonly)},
					},
				},
			},
		},
		{
			name: "No new line br, and has no italic",
			footnotes: map[int]map[int]repository.VerseFootNotes{
				2: {
					1: repository.VerseFootNotes{
						MarkerStyle:    sql.NullInt32{Valid: true, Int32: int32(VerseNoteStyleAlignRight)},
						FootnoteMarker: sql.NullString{Valid: true, String: "unit"},
						Footnote:       sql.NullString{Valid: true, String: " = test"},
					},
				},
			},
			y: getY(),
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='footnotes'", `style="font-size:60%;font-family:'Figtree';font-weight:600;font-style:italic"`)
				canv.EXPECT().Text(70, 95, "unit = test")
				canv.EXPECT().Gend()
				return canv
			},
			wantY: 125,
		},
		{
			name: "has new line br, and has italic",
			footnotes: map[int]map[int]repository.VerseFootNotes{
				2: {
					1: repository.VerseFootNotes{
						MarkerStyle:    sql.NullInt32{Valid: true, Int32: int32(VerseNoteStyleAlignRight)},
						FootnoteMarker: sql.NullString{Valid: true, String: "this is <i>unit:</i> "},
						Footnote:       sql.NullString{Valid: true, String: "test1<br/>test2<br/>test3"},
					},
				},
			},
			y: getY(),
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='footnotes'", `style="font-size:60%;font-family:'Figtree';font-weight:600;"`)
				canv.EXPECT().TextUnescaped(float64(70), float64(95), `this is <tspan font-style="italic">unit:</tspan> `)
				canv.EXPECT().TextUnescaped(float64(123), float64(95), "test1")
				canv.EXPECT().TextUnescaped(float64(123), float64(110), "test2")
				canv.EXPECT().TextUnescaped(float64(123), float64(125), "test3")
				canv.EXPECT().Gend()
				return canv
			},
			wantY: 155,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi footnoteInteractor
			canv := canvas.Canvas(nil)
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			fi.RenderVerseFootnotes(canv, tt.y, tt.footnotes)
			if tt.y != nil {
				assert.Equal(t, tt.wantY, *tt.y)

			}
		})
	}
}
