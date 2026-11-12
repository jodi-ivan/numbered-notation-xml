package lyric_test

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func TestGetMusicxmlLyric(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note *entity.NoteRenderer
		want []musicxml.Lyric
	}{
		{
			name: "No lyric",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{},
			},
			want: []musicxml.Lyric{},
		},
		{
			name: "Single lyric part with single text element",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Underline: 1,
								Value:     "Hello",
							},
						},
						Syllabic: musicxml.LyricSyllabic("syllabic"),
					},
				},
			},
			want: []musicxml.Lyric{
				{
					Number: 1,
					Text: []musicxml.LyricText{
						{
							Underline: 1,
							Value:     "Hello",
						},
					},
					Syllabic: musicxml.LyricSyllabic("syllabic"),
				},
			},
		},
		{
			name: "Single lyric part with multiple text elements",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Underline: 0,
								Value:     "Hello",
							},
							{
								Underline: 1,
								Value:     "World",
							},
						},
						Syllabic: musicxml.LyricSyllabic("syllabic"),
					},
				},
			},
			want: []musicxml.Lyric{
				{
					Number: 1,
					Text: []musicxml.LyricText{
						{
							Underline: 0,
							Value:     "Hello",
						},
						{
							Underline: 1,
							Value:     "World",
						},
					},
					Syllabic: musicxml.LyricSyllabic("syllabic"),
				},
			},
		},
		{
			name: "Multiple lyric parts",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Underline: 0,
								Value:     "First",
							},
						},
						Syllabic: musicxml.LyricSyllabic("syllabic"),
					},
					{
						Text: []entity.Text{
							{
								Underline: 0,
								Value:     "Second",
							},
						},
						Syllabic: musicxml.LyricSyllabic("syllabic"),
					},
				},
			},
			want: []musicxml.Lyric{
				{
					Number: 1,
					Text: []musicxml.LyricText{
						{
							Underline: 0,
							Value:     "First",
						},
					},
					Syllabic: musicxml.LyricSyllabic("syllabic"),
				},
				{
					Number: 2,
					Text: []musicxml.LyricText{
						{
							Underline: 0,
							Value:     "Second",
						},
					},
					Syllabic: musicxml.LyricSyllabic("syllabic"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lyric.GetMusicxmlLyric(tt.note)
			if len(got) != len(tt.want) {
				t.Fatalf("GetMusicxmlLyric() length mismatch, got %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Number != tt.want[i].Number {
					t.Errorf("GetMusicxmlLyric() lyric %d Number mismatch, got %d, want %d", i+1, got[i].Number, tt.want[i].Number)
				}
				if got[i].Syllabic != tt.want[i].Syllabic {
					t.Errorf("GetMusicxmlLyric() lyric %d Syllabic mismatch, got %v, want %v", i+1, got[i].Syllabic, tt.want[i].Syllabic)
				}
				if len(got[i].Text) != len(tt.want[i].Text) {
					t.Errorf("GetMusicxmlLyric() lyric %d Text length mismatch, got %d, want %d", i+1, len(got[i].Text), len(tt.want[i].Text))
				}
				for j := range got[i].Text {
					if got[i].Text[j].Underline != tt.want[i].Text[j].Underline {
						t.Errorf("GetMusicxmlLyric() lyric %d text %d Underline mismatch, got %d, want %d", i+1, j+1, got[i].Text[j].Underline, tt.want[i].Text[j].Underline)
					}
					if got[i].Text[j].Value != tt.want[i].Text[j].Value {
						t.Errorf("GetMusicxmlLyric() lyric %d text %d Value mismatch, got %s, want %s", i+1, j+1, got[i].Text[j].Value, tt.want[i].Text[j].Value)
					}
				}
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note *entity.NoteRenderer
		want bool
	}{
		{
			name: "No lyric",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{},
			},
			want: false,
		},
		{
			name: "Lyric starts with pattern (no digits)",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Value: ". ",
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Lyric starts with pattern (digits, no space)",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Value: "1.",
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Lyric starts with pattern (digits, with space)",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Value: "2. ",
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Lyric does not start with pattern",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Value: "Hello",
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Lyric starts with digits but no dot",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Value: "123",
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Multiple lyrics, one matches",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Value: "Hello",
							},
						},
					},
					{
						Text: []entity.Text{
							{
								Value: "3.",
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Multiple lyrics, none match",
			note: &entity.NoteRenderer{
				Lyric: []entity.Lyric{
					{
						Text: []entity.Text{
							{
								Value: "A",
							},
						},
					},
					{
						Text: []entity.Text{
							{
								Value: "B",
							},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lyric.HasPrefix(tt.note)
			if got != tt.want {
				t.Errorf("HasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
