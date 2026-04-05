package splitter_test

import (
	"context"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm/splitter"
	"github.com/stretchr/testify/assert"
)

func TestCleanBeamByNumber(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		notes []*entity.NoteRenderer
		no    int
		want  []splitter.BeamSplitMarker
	}{
		{
			want: []splitter.BeamSplitMarker{},
		},
		{
			notes: []*entity.NoteRenderer{
				{},
			},
			want: []splitter.BeamSplitMarker{},
		},
		{
			no: 1,
			notes: []*entity.NoteRenderer{
				{
					Beam: map[int]entity.Beam{
						1: {
							Number: 1,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				{
					Beam: map[int]entity.Beam{
						1: {
							Number: 1,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
			},
			want: []splitter.BeamSplitMarker{
				{
					StartIndex: 0,
					EndIndex:   1,
				},
			},
		},
		{
			no: 2,
			notes: []*entity.NoteRenderer{
				{
					Beam: map[int]entity.Beam{
						1: {
							Number: 1,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				{
					Beam: map[int]entity.Beam{
						1: {
							Number: 1,
							Type:   musicxml.NoteBeamTypeContinue,
						},
						2: {
							Number: 2,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				{
					Beam: map[int]entity.Beam{
						1: {
							Number: 1,
							Type:   musicxml.NoteBeamTypeEnd,
						},
						2: {
							Number: 2,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				{
					Beam: map[int]entity.Beam{
						1: {
							Number: 1,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				{
					Beam: map[int]entity.Beam{
						1: {
							Number: 1,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
			},
			want: []splitter.BeamSplitMarker{
				{
					StartIndex: 1,
					EndIndex:   2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitter.CleanBeamByNumber(context.Background(), tt.notes, tt.no)

			assert.Equal(t, tt.want, got, "CleanBeamByNumber()")
		})
	}
}
