package credits

import (
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

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
				"this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <tspan font-style=\"italic\">Also added a long italic text for breaking </tspan>",
				"down the text to the new line.</tspan>",
			},
			lenLines: []int{661, 156},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, lenLines := autoWrapText(tt.args.text, tt.args.leftIndent)
			if !assert.Equal(t, tt.lines, lines) {
				t.Errorf("creditsInteractor.autoWrapText() lines got = %v, want %v", lines, tt.lines)
			}
			if !assert.Equal(t, tt.lenLines, lenLines) {
				t.Errorf("creditsInteractor.autoWrapText() lenLines got = %v, want %v", lenLines, tt.lenLines)
			}
		})
	}
}

func Test_alignText(t *testing.T) {
	type args struct {
		text         string
		textLength   int
		targetLength int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "everything went fine",
			args: args{
				text:         "this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <tspan font-style=\"italic\">Also",
				textLength:   625,
				targetLength: constant.LAYOUT_WIDTH,
			},
			// 5 spaces
			want: "this&#160;&#160;&#160;&#160;&#160;is&#160;&#160;&#160;&#160;&#160;a&#160;&#160;&#160;&#160;&#160;very&#160;&#160;&#160;&#160;&#160;long&#160;&#160;&#160;&#160;&#160;text,&#160;&#160;&#160;&#160;&#160;this&#160;&#160;&#160;&#160;&#160;intentionally&#160;&#160;&#160;&#160;&#160;added&#160;&#160;&#160;&#160;&#160;with&#160;&#160;&#160;&#160;&#160;a&#160;&#160;&#160;&#160;&#160;lot&#160;&#160;&#160;&#160;&#160;of&#160;&#160;&#160;&#160;&#160;text&#160;&#160;&#160;&#160;&#160;just&#160;&#160;&#160;&#160;&#160;for&#160;&#160;&#160;&#160;&#160;satisfy&#160;&#160;&#160;&#160;&#160;requirement.&#160;&#160;&#160;&#160;&#160;<tspan font-style=\"italic\">Also",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alignText(tt.args.text, tt.args.textLength, tt.args.targetLength); !assert.Equal(t, tt.want, got) {
				t.Errorf("alignText() = %v, want %v", got, tt.want)
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
		canv       func(*gomock.Controller) *canvas.MockCanvas
		y          int
		leftIndent int
		text       string
		want       []string
	}{
		{
			name: "nothing happened",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				return canvas.NewMockCanvas(c)
			},
			want: []string{},
		},
		{
			name: "one line without italic",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().TextUnescaped(50.0, 100.0, "this is the text without italic")

				return canvMock
			},
			y:    100,
			text: "this is the text without italic",
			want: []string{"this is the text without italic"},
		},
		{
			name: "one linewithout <i>italic</i>",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().TextUnescaped(50.0, 100.0, `this is the text with <tspan font-style="italic">Foreign title</tspan>`)

				return canvMock
			},
			y:    100,
			text: "this is the text with <i>Foreign title</i>",
			want: []string{`this is the text with <tspan font-style="italic">Foreign title</tspan>`},
		},
		{
			name: "with italic terminated is broken down to two lines",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)
				canvMock.EXPECT().TextUnescaped(100.0, 100.0,
					`this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <tspan font-style="italic">Also added a long italic text for breaking </tspan>`)

				canvMock.EXPECT().TextUnescaped(100.0, 115.0,
					`<tspan font-style="italic"> down the text to the new line.</tspan>`)
				return canvMock
			},
			y:          100,
			leftIndent: constant.LAYOUT_INDENT_LENGTH,
			text:       "this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <i>Also added a long italic text for breaking down the text to the new line.</i>",
			want: []string{
				"this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <tspan font-style=\"italic\">Also added a long italic text for breaking </tspan>",
				"<tspan font-style=\"italic\"> down the text to the new line.</tspan>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAndRenderText(tt.canv(ctrl), tt.y, tt.leftIndent, tt.text)
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
			got := renderMusicAndLyric(tt.canv(ctrl), tt.y, tt.metadata)
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

				canvMock.EXPECT().Text(654, 100, "© unittest")
				return canvMock
			},
		},
		{
			name:  "copyright offset to new line",
			y:     newYPtr(),
			wantY: 130,
			metadata: repository.HymnData{
				Copyright: sql.NullString{
					Valid:  true,
					String: "this is long copyright",
				},
			},
			leftIndent: 600,
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canvMock := canvas.NewMockCanvas(c)

				canvMock.EXPECT().Text(597, 115, "© this is long copyright")
				return canvMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderCopyright(tt.canv(ctrl), tt.y, tt.leftIndent, tt.metadata)

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
				canvMock.EXPECT().Text(671, 115, "BE 100")
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
				canvMock.EXPECT().Text(669, 115, "NR 100")
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
				canvMock.EXPECT().Text(635, 115, "BE 100, NR 100")
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
			renderReferences(tt.canv(ctrl), tt.y, tt.metadata)
		})
	}
}
