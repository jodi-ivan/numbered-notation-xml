package lyric

import (
	"reflect"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func Test_lyricInteractor_SetLyricRenderer(t *testing.T) {
	type args struct {
		noteRenderer *entity.NoteRenderer
		note         musicxml.Note
	}
	tests := []struct {
		name string
		args args
		want VerseInfo

		// since the note renderer is changing, we need to assert the changes here
		wantWidth          int
		wantTakenFromLyric bool
		wantLyric          []entity.Lyric
	}{
		{
			name: "no lyric in the renderer",
			args: args{
				noteRenderer: &entity.NoteRenderer{},
				note:         musicxml.Note{},
			},
			want:      VerseInfo{},
			wantWidth: 15,
		},
		{
			name: "lyric less than the note",
			args: args{
				noteRenderer: &entity.NoteRenderer{},
				note: musicxml.Note{
					Lyric: []musicxml.Lyric{
						musicxml.Lyric{
							Number: 1,
							Text: []struct {
								Underline int    `xml:"underline,attr"`
								Value     string `xml:",chardata"`
							}{
								{
									Value: "a",
								},
							},
							Syllabic: musicxml.LyricSyllabicTypeBegin,
						},
					},
				},
			},
			want:      VerseInfo{},
			wantWidth: 15,
			wantLyric: []entity.Lyric{
				entity.Lyric{
					Text: []entity.Text{
						entity.Text{
							Value: "a",
						},
					},
					Syllabic: musicxml.LyricSyllabicTypeBegin,
				},
			},
		},
		{
			name: "lyric more than the note",
			args: args{
				noteRenderer: &entity.NoteRenderer{},
				note: musicxml.Note{
					Lyric: []musicxml.Lyric{
						musicxml.Lyric{
							Number: 1,
							Text: []struct {
								Underline int    `xml:"underline,attr"`
								Value     string `xml:",chardata"`
							}{
								{
									Value: "Yang",
								},
							},
							Syllabic: musicxml.LyricSyllabicTypeSingle,
						},
					},
				},
			},
			want:               VerseInfo{},
			wantWidth:          34,
			wantTakenFromLyric: true,
			wantLyric: []entity.Lyric{
				entity.Lyric{
					Text: []entity.Text{
						entity.Text{
							Value: "Yang",
						},
					},
					Syllabic: musicxml.LyricSyllabicTypeSingle,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := &lyricInteractor{}
			if got := li.SetLyricRenderer(tt.args.noteRenderer, tt.args.note); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lyricInteractor.SetLyricRenderer() = %v, want %v", got, tt.want)
			}

			assert.Equal(t, tt.wantWidth, tt.args.noteRenderer.Width, "note width after preparation")
			assert.Equal(t, tt.wantTakenFromLyric, tt.args.noteRenderer.IsLengthTakenFromLyric, "the width taken from lyric")
			assert.Equal(t, tt.wantLyric, tt.args.noteRenderer.Lyric, "the renderer lyric")
		})
	}
}

func Test_lyricInteractor_CalculateMarginLeft(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "no margin left",
			args: args{
				txt: "Ha",
			},
			want: 0,
		},
		{
			name: "margin left",
			args: args{
				txt: "1. Ha",
			},
			want: -16.58,
		},
		{
			name: "margin left",
			args: args{
				txt: "15. Be",
			},
			want: -24.19,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := &lyricInteractor{}

			if got := li.CalculateMarginLeft(tt.args.txt); got != tt.want {
				t.Errorf("lyricInteractor.CalculateMarginLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}
