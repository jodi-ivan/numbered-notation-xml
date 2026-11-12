package credits

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

const nbsp = "&#160;"

func rep(n int) string { return strings.Repeat(nbsp, n) }

func TestNewCredits(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "everything went fine",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCredits(); !assert.NotNil(t, got) {
				t.Fail()
			}
		})
	}
}

func Test_creditsInteractor_autoWrapText(t *testing.T) {
	type args struct {
		text       string
		leftIndent int
	}
	tests := []struct {
		name     string
		args     args
		lines    []string
		lenLines []int
	}{
		{
			name: "no italic, no new line",
			args: args{
				text:       "this is a simple text",
				leftIndent: constant.LAYOUT_INDENT_LENGTH,
			},
			lines:    []string{"this is a simple text"},
			lenLines: []int{99},
		},
		{
			name: "with italic terminated in the middle sentence, no new line",
			args: args{
				text:       "this is a simple text <i>with italic</i> added",
				leftIndent: constant.LAYOUT_INDENT_LENGTH,
			},
			lines:    []string{"this is a simple text <tspan font-style=\"italic\">with italic</tspan> added"},
			lenLines: []int{182},
		},
		{
			name: "with italic terminated in the end sentence, no new line",
			args: args{
				text:       "this is a simple text <i>with italic</i>",
				leftIndent: constant.LAYOUT_INDENT_LENGTH,
			},
			lines:    []string{"this is a simple text <tspan font-style=\"italic\">with italic</tspan>"},
			lenLines: []int{150},
		},
		{
			name: "with italic terminated is broken down to two lines",
			args: args{
				text:       "this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <i>Also added a long italic text for breaking down the text to the new line.</i>",
				leftIndent: constant.LAYOUT_INDENT_LENGTH,
			},
			lines: []string{
				"this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <tspan font-style=\"italic\">Also added a long italic text for breaking down the </tspan>",
				"text to the new line.</tspan>",
			},
			lenLines: []int{711, 106},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, lenLines := autoWrapText(context.Background(), tt.args.text, tt.args.leftIndent)
			if !assert.Equal(t, tt.lines, lines) {
				t.Errorf("creditsInteractor.autoWrapText() lines got = %v, want %v", lines, tt.lines)
			}
			if !assert.Equal(t, tt.lenLines, lenLines) {
				t.Errorf("creditsInteractor.autoWrapText() lenLines got = %v, want %v", lenLines, tt.lenLines)
			}
		})
	}
}

func Test_AlignText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		textLength   int
		targetLength int
		want         string
	}{
		{"single word", "hello", 50, 100, "hello"},
		{"two words (only 1 gap) are unchanged", "hello world", 50, 100, "hello world"},
		{"target equals text length", "a b c d", 100, 100, "a b c d"},
		{"target smaller than text length", "a b c d", 100, 90, "a b c d"},
		{"room smaller than one space", "a b c d", 100, 101, "a b c d"},
		{
			// 20/2 = 10 extra, 4 gaps -> 2 each, remainder 2 -> 3,3,2,2
			"remainder spread over first gaps",
			"a b c d e", 100, 120,
			"a " + rep(3) + "b " + rep(3) + "c " + rep(2) + "d " + rep(2) + "e",
		},
		{
			// 16/2 = 8 extra, 4 gaps -> 2 each
			"even distribution",
			"a b c d e", 100, 116,
			"a " + rep(2) + "b " + rep(2) + "c " + rep(2) + "d " + rep(2) + "e",
		},
		{
			// 9/2 = 4 extra (floor), 3 gaps -> 1 each, remainder 1 -> 2,1,1
			"odd px rounds down",
			"a b c d", 100, 109,
			"a " + rep(2) + "b " + rep(1) + "c " + rep(1) + "d",
		},
		{
			// 6/2 = 3 extra, 2 gaps -> 1 each, remainder 1 -> 2,1
			"multiple spaces in input are collapsed",
			"a   b    c", 50, 56,
			"a " + rep(2) + "b " + rep(1) + "c",
		},
		{
			// the tspan tag counts as one word; 3 extra, 2 gaps -> 2,1
			"tspan tag is preserved",
			`<tspan font-style="italic">hi</tspan> there you`, 50, 56,
			`<tspan font-style="italic">hi</tspan> ` + rep(2) + "there " + rep(1) + "you",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alignText(tt.text, float64(tt.textLength), float64(tt.targetLength))
			if got != tt.want {
				t.Errorf("alignText() =\n  %q\nwant\n  %q", got, tt.want)
			}
		})
	}
}

type customStringMatcher struct {
	expected string
	T        assert.TestingT
}

func (m *customStringMatcher) Matches(x interface{}) bool {
	s, ok := x.([]byte)
	if !ok {
		return false
	}
	// Custom logic, e.g., checking if string contains a substring
	return assert.Equal(m.T, m.expected, string(s))
}

func (m *customStringMatcher) String() string {
	return "contains " + m.expected
}

func Test_formatAndRenderText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv       func(*gomock.Controller) *canvas.MockCanvasTestify
		y          int
		leftIndent int
		text       string
		want       []string
	}{
		{
			name: "nothing happened",
			canv: func(c *gomock.Controller) *canvas.MockCanvasTestify {
				return canvas.NewMockCanvasTestify(t)
			},
			want: []string{},
		},
		{
			name: "one line without italic",
			canv: func(c *gomock.Controller) *canvas.MockCanvasTestify {
				canvMock := canvas.NewMockCanvasTestify(t)
				canvMock.EXPECT().TextUnescaped(50.0, 100.0, "this is the text without italic")

				return canvMock
			},
			y:    100,
			text: "this is the text without italic",
			want: []string{"this is the text without italic"},
		},
		{
			name: "one linewithout <i>italic</i>",
			canv: func(c *gomock.Controller) *canvas.MockCanvasTestify {
				canvMock := canvas.NewMockCanvasTestify(t)
				canvMock.EXPECT().TextUnescaped(50.0, 100.0, `this is the text with <tspan font-style="italic">Foreign title</tspan>`)

				return canvMock
			},
			y:    100,
			text: "this is the text with <i>Foreign title</i>",
			want: []string{`this is the text with <tspan font-style="italic">Foreign title</tspan>`},
		},
		{
			name: "with italic terminated is broken down to two lines",
			canv: func(c *gomock.Controller) *canvas.MockCanvasTestify {
				canvMock := canvas.NewMockCanvasTestify(t)
				canvMock.EXPECT().TextUnescaped(100.0, 100.0,
					`this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <tspan font-style="italic">Also added a long italic text for breaking down the </tspan>`)

				canvMock.EXPECT().TextUnescaped(100.0, 115.0,
					`<tspan font-style="italic"> text to the new line.</tspan>`)
				return canvMock
			},
			y:          100,
			leftIndent: constant.LAYOUT_INDENT_LENGTH,
			text:       "this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <i>Also added a long italic text for breaking down the text to the new line.</i>",
			want: []string{
				"this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <tspan font-style=\"italic\">Also added a long italic text for breaking down the </tspan>",
				"<tspan font-style=\"italic\"> text to the new line.</tspan>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAndRenderText(context.Background(), tt.canv(ctrl), tt.y, tt.leftIndent, tt.text)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("formatAndRenderText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_renderMusicAndLyric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	initialY := 100
	initialY2 := 100
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv     func(c *gomock.Controller) *canvas.MockCanvas
		y        *int
		metadata repository.HymnData
		want     float64
		wantY    int
	}{
		{
			name: "different music and lyric",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvasMock := canvas.NewMockCanvas(c)
				canvasMock.EXPECT().Text(50, 100, "Syair:")
				canvasMock.EXPECT().TextUnescaped(80.0, 100.0, "Lyric unittest")
				canvasMock.EXPECT().Text(50, 115, "Lagu:")
				canvasMock.EXPECT().TextUnescaped(80.0, 115.0, "Music unittest")

				return canvasMock
			},
			y: &initialY,
			metadata: repository.HymnData{
				Lyric: "Lyric unittest",
				Music: "Music unittest",
			},
			want:  94.65,
			wantY: 115,
		},
		{
			name: "same music and lyric",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvasMock := canvas.NewMockCanvas(c)
				canvasMock.EXPECT().Text(50, 100, "Syair dan lagu: ")
				canvasMock.EXPECT().TextUnescaped(118.0, 100.0,
					`<tspan font-style="italic">unittest rocks!</tspan>`)

				return canvasMock
			},
			y: &initialY2,
			metadata: repository.HymnData{
				Lyric: "<i>unittest rocks!</i>",
				Music: "<i>unittest rocks!</i>",
			},
			want:  utils.CalculateSecondaryLyricWidth("unittest rocks!") + 68,
			wantY: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderMusicAndLyric(context.Background(), tt.canv(ctrl), tt.y, tt.metadata)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("renderMusicAndLyric() = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.wantY, *tt.y) {
				t.Errorf("&y renderMusicAndLyric() = %v, want %v", *tt.y, tt.wantY)
			}
		})
	}
}

func Test_renderCopyright(t *testing.T) {
	newYPtr := func() *int {
		i := 100

		return &i
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv       func(c *gomock.Controller) *canvas.MockCanvas
		y          *int
		leftIndent float64
		wantY      int
		metadata   repository.HymnData
	}{
		{
			name:  "no copyright",
			y:     newYPtr(),
			wantY: 100,
			metadata: repository.HymnData{
				Copyright: sql.NullString{},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				return nil
			},
		},
		{
			name:  "copyright",
			y:     newYPtr(),
			wantY: 115,
			metadata: repository.HymnData{
				Copyright: sql.NullString{
					Valid:  true,
					String: "unittest",
				},
			},
			leftIndent: 100,
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)

				canvMock.EXPECT().Text(684, 100, "© unittest")
				return canvMock
			},
		},
		{
			name:  "copyright offset to new line",
			y:     newYPtr(),
			wantY: 115,
			metadata: repository.HymnData{
				Copyright: sql.NullString{
					Valid:  true,
					String: "this is long copyright",
				},
			},
			leftIndent: 600,
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)

				canvMock.EXPECT().Text(627, 100, "© this is long copyright")
				return canvMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderCopyright(context.Background(), tt.canv(ctrl), tt.y, tt.leftIndent, tt.metadata)

			if !assert.Equal(t, tt.wantY, *tt.y) {
				t.Errorf("&y renderCopyright() = %v, want %v", *tt.y, tt.wantY)
			}
		})
	}
}

func Test_renderReferences(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv     func(c *gomock.Controller) *canvas.MockCanvas
		y        int
		metadata repository.HymnData
	}{
		{
			name: "",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				return nil
			},
			y:        115,
			metadata: repository.HymnData{},
		},
		{
			name: "BE only",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().Text(751, 115, "BE 100")
				return canvMock
			},
			y: 115,
			metadata: repository.HymnData{
				RefBE: sql.NullInt16{
					Valid: true,
					Int16: 100,
				},
			},
		},
		{
			name: "NR only",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().Text(749, 115, "NR 100")
				return canvMock
			},
			y: 115,
			metadata: repository.HymnData{
				RefNR: sql.NullInt16{
					Valid: true,
					Int16: 100,
				},
			},
		},
		{
			name: "Both",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().Text(715, 115, "BE 100, NR 100")
				return canvMock
			},
			y: 115,
			metadata: repository.HymnData{
				RefNR: sql.NullInt16{
					Valid: true,
					Int16: 100,
				},
				RefBE: sql.NullInt16{
					Valid: true,
					Int16: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderReferences(context.Background(), tt.canv(ctrl), tt.y, tt.metadata)
		})
	}
}
