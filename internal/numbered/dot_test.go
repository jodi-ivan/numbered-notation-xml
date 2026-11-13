package numbered

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func TestReplaceDotWithNumbered(t *testing.T) {
	tests := []struct {
		name   string
		dot    *entity.NoteRenderer
		number *entity.NoteRenderer
		want   *entity.NoteRenderer
	}{
		{
			name: "Basic replacement",
			dot: &entity.NoteRenderer{
				IsDotted: true,
			},
			number: &entity.NoteRenderer{
				Note:           1,
				Strikethrough:  true,
				AbsoluteNote:   "C",
				AbsoluteOctave: 4,
			},
			want: &entity.NoteRenderer{
				Note:           1,
				Strikethrough:  true,
				AbsoluteNote:   "C",
				AbsoluteOctave: 4,
				Tie: &entity.Slur{
					Number:       1,
					Type:         musicxml.NoteSlurTypeStop,
					NumberedOnly: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceDotWithNumbered(tt.dot, tt.number)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.number.Tie, &entity.Slur{
				Number:       1,
				Type:         musicxml.NoteSlurTypeStart,
				NumberedOnly: true,
			})
		})
	}
}
