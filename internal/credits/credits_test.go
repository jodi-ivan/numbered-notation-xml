package credits

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
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

func TestCalculateLyric(t *testing.T) {
	type args struct {
		text   string
		italic bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "with italic",
			args: args{
				italic: true,
				text:   "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKMNOPQRSTUVWXYZ1234567890-.!;:-/",
			},
			want: 361.79999999999984,
		},
		{
			name: "without italic",
			args: args{
				italic: false,
				text:   "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKMNOPQRSTUVWXYZ1234567890-.!;:-/",
			},
			want: 535,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateLyric(tt.args.text, tt.args.italic); got != tt.want {
				t.Errorf("CalculateLyric() = %v, want %v", got, tt.want)
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
			lenLines: []int{125},
		},
		{
			name: "with italic terminated in the middle sentence, no new line",
			args: args{
				text:       "this is a simple text <i>with italic</i> added",
				leftIndent: constant.LAYOUT_INDENT_LENGTH,
			},
			lines:    []string{"this is a simple text <tspan font-style=\"italic\">with italic</tspan> added"},
			lenLines: []int{213},
		},
		{
			name: "with italic terminated in the end sentence, no new line",
			args: args{
				text:       "this is a simple text <i>with italic</i>",
				leftIndent: constant.LAYOUT_INDENT_LENGTH,
			},
			lines:    []string{"this is a simple text <tspan font-style=\"italic\">with italic</tspan>"},
			lenLines: []int{172},
		},
		{
			name: "with italic terminated is broken down to two lines",
			args: args{
				text:       "this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <i>Also added a long italic text for breaking down the text to the new line.</i>",
				leftIndent: constant.LAYOUT_INDENT_LENGTH,
			},
			lines: []string{
				"this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <tspan font-style=\"italic\">Also",
				"added a long italic text for breaking down the text to the new line.</tspan>",
			},
			lenLines: []int{625, 281},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := creditsInteractor{}
			lines, lenLines := ci.autoWrapText(tt.args.text, tt.args.leftIndent)
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

func Test_creditsInteractor_RenderCredits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		y        int
		metadata repository.HymnData
	}
	tests := []struct {
		name     string
		args     args
		initCanv func(ctrl *gomock.Controller) *canvas.MockCanvas
	}{
		{
			name: "music and lyric same, no BE, no NR",
			args: args{
				y: 150,
				metadata: repository.HymnData{
					Lyric: "this is from unittest",
					Music: "this is from unittest",
				},
			},
			initCanv: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(ctrl)
				writer := canvas.NewMockWriter(ctrl)
				res.EXPECT().Writer().Return(writer)

				res.EXPECT().Group("class='credit'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
				res.EXPECT().Text(50, 150, "Syair dan lagu :")
				writer.EXPECT().Write([]byte("<text x=\"118\" y=\"150\">this is from unittest</text>"))
				res.EXPECT().Gend()
				return res
			},
		},
		{
			name: "music and lyric different, with BE and with NR. no italic",
			args: args{
				y: 150,
				metadata: repository.HymnData{
					Lyric: "this is lyric from unittest",
					Music: "this is music from unittest",
					RefNR: sql.NullInt16{
						Valid: true,
						Int16: 1,
					},
					RefBE: sql.NullInt16{
						Valid: true,
						Int16: 1,
					},
				},
			},
			initCanv: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(ctrl)
				writer := canvas.NewMockWriter(ctrl)
				res.EXPECT().Writer().Return(writer).Times(2)

				res.EXPECT().Group("class='credit'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
				res.EXPECT().Gend()
				res.EXPECT().Text(50, 150, "Syair: ")
				writer.EXPECT().Write([]byte("<text x=\"80\" y=\"150\">this is lyric from unittest</text>"))
				writer.EXPECT().Write([]byte("<text x=\"50\" y=\"165\">Lagu: this is music from unittest</text>"))
				res.EXPECT().Text(644, 165, "BE 1, NR 1")

				return res
			},
		},
		{
			name: "music and lyric different,  italic break down to multiple lines",
			args: args{
				y: 150,
				metadata: repository.HymnData{
					Lyric: "this is a very long text, this intentionally added with a lot of text just for satisfy requirement. <i>Also added a long italic text for breaking down the text to the new line.</i>",
					Music: "this is music from unittest",
				},
			},
			initCanv: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(ctrl)
				writer := canvas.NewMockWriter(ctrl)
				res.EXPECT().Writer().Return(writer).Times(3)

				res.EXPECT().Group("class='credit'", `style="font-size:60%;font-family:'Figtree';font-weight:600"`)
				res.EXPECT().Gend()
				res.EXPECT().Text(50, 150, "Syair: ")
				line1 := `this&#160;&#160;&#160;is&#160;&#160;&#160;a&#160;&#160;&#160;very&#160;&#160;&#160;long&#160;&#160;&#160;text,&#160;&#160;&#160;this&#160;&#160;&#160;intentionally&#160;&#160;&#160;added&#160;&#160;&#160;with&#160;&#160;&#160;a&#160;&#160;&#160;lot&#160;&#160;&#160;of&#160;&#160;&#160;text&#160;&#160;&#160;just&#160;&#160;&#160;for&#160;&#160;&#160;satisfy&#160;&#160;&#160;requirement.&#160;&#160;&#160;<tspan font-style="italic">Also&#160;&#160;&#160;added</tspan>`
				line2 := `<tspan font-style="italic">a long italic text for breaking down the text to the new line.</tspan>`

				writer.EXPECT().Write([]byte(fmt.Sprintf(`<text x="80" y="150">%s</text>`, line1)))
				writer.EXPECT().Write([]byte(fmt.Sprintf(`<text x="80" y="165">%s</text>`, line2)))
				writer.EXPECT().Write([]byte("<text x=\"50\" y=\"180\">Lagu: this is music from unittest</text>"))

				return res
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := creditsInteractor{}
			ci.RenderCredits(context.Background(), tt.initCanv(ctrl), tt.args.y, tt.args.metadata)
		})
	}
}
