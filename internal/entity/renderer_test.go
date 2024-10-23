package entity

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func TestLyricVal_String(t *testing.T) {
	tests := []struct {
		name string
		lv   LyricVal
		want string
	}{
		{
			name: "single",
			lv: LyricVal([]Text{
				Text{
					Value: "single",
				},
			}),
			want: "single",
		},
		{
			name: "multiple",
			lv: LyricVal([]Text{
				Text{
					Value: "one",
				},
				Text{
					Value: "two",
				},
			}),
			want: "onetwo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lv.String(); !assert.Equal(t, tt.want, got) {
				t.Errorf("LyricVal.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoteRenderer_UpdateBeam(t *testing.T) {

	t.Run("UpdateBeam", func(t *testing.T) {
		nr := NoteRenderer{
			Beam: map[int]Beam{},
		}
		nr.UpdateBeam(1, musicxml.NoteBeamTypeBegin)
		expect := map[int]Beam{1: Beam{1, musicxml.NoteBeamTypeBegin}}
		assert.Equal(t, expect, nr.Beam)
	})

}
