package lyric

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

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

func Test_lyricInteractor_SplitLyricPrefix(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note        *entity.NoteRenderer
		y           int
		part        int
		leftBarline *entity.NoteRenderer
		want        []LyricPosition
	}{
		{
			name: "Bait prefix",
			note: &entity.NoteRenderer{
				PositionX: 85,
				Lyric: []entity.Lyric{
					entity.Lyric{
						Text: []entity.Text{
							entity.Text{
								Value: "bait1:unitest",
							},
						},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
					},
				},
			},
			want: []LyricPosition{
				{
					Coordinate: entity.NewCoordinate(34.7, 25),
					Lyrics: entity.Lyric{
						Text: []entity.Text{
							{
								Value:  "bait 1: ",
								Bold:   true,
								Italic: true,
							},
						},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
					},
				},
				{
					Coordinate: entity.NewCoordinate(85, 25),
					Lyrics: entity.Lyric{
						Text: []entity.Text{
							{
								Value: "unitest",
							},
						},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
					},
				},
			},
		},
		{
			name: "Verse dot prefix",
			note: &entity.NoteRenderer{
				PositionX: 60,
				Lyric: []entity.Lyric{
					entity.Lyric{
						Text: []entity.Text{
							entity.Text{
								Value: "2.unitest",
							},
						},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
					},
				},
			},
			want: []LyricPosition{
				{
					Coordinate: entity.NewCoordinate(45.85, 25),
					Lyrics: entity.Lyric{
						Text: []entity.Text{
							{
								Value: "2.",
							},
						},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
					},
				},
				{
					Coordinate: entity.NewCoordinate(60, 25),
					Lyrics: entity.Lyric{
						Text: []entity.Text{
							{
								Value: "unitest",
							},
						},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
					},
				},
			},
		},
		{
			name: "No prefix to be splitted",
			note: &entity.NoteRenderer{
				PositionX: 60,
				Lyric: []entity.Lyric{
					entity.Lyric{
						Text: []entity.Text{
							entity.Text{
								Value: "unitest",
							},
						},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
					},
				},
			},
			// want: []LyricPosition{
			// 	{
			// 		Coordinate: entity.NewCoordinate(60, 25),
			// 		Lyrics: entity.Lyric{
			// 			Text: []entity.Text{
			// 				{
			// 					Value: "unittest",
			// 				},
			// 			},
			// 			Syllabic: musicxml.LyricSyllabicTypeSingle,
			// 		},
			// 	},
			// },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var li lyricInteractor
			got := li.SplitLyricPrefix(tt.note, tt.y, tt.part, tt.leftBarline)
			assert.Equal(t, tt.want, got)

		})
	}
}
