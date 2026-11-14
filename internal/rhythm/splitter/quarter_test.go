package splitter

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/stretchr/testify/assert"
)

func Test_quaterSpliiter_shouldMergeSegments(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		notes        []*entity.NoteRenderer
		s            BeamSplitMarker
		nextInterval BeamSplitMarker
		want         bool
	}{

		// // 1. Basic success case: All conditions met
		{
			name: "Merge - conditions met",
			notes: []*entity.NoteRenderer{
				// Mock notes for index 0 and 2
				{Beam: map[int]entity.Beam{1: {}}},
				{Beam: map[int]entity.Beam{1: {}}, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
				{Beam: map[int]entity.Beam{1: {}}},
			},
			s:            BeamSplitMarker{StartIndex: 0, EndIndex: 0}, // Interval 1: 1 note
			nextInterval: BeamSplitMarker{StartIndex: 2, EndIndex: 2}, // Interval 2: 1 note
			want:         true,
		},
		// 2. Failure: No one note gap
		{
			name: "No Merge - gap too small",
			notes: []*entity.NoteRenderer{
				{Beam: map[int]entity.Beam{1: {}}},
				{Beam: map[int]entity.Beam{1: {}}, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
				{Beam: map[int]entity.Beam{1: {}}},
			},
			s:            BeamSplitMarker{StartIndex: 0, EndIndex: 1}, // Gap size 1
			nextInterval: BeamSplitMarker{StartIndex: 2, EndIndex: 2}, // Gap size 2
			want:         false,
		},
		// // 3. Failure: Not a breath mark
		{
			name: "No Merge - not breath mark",
			notes: []*entity.NoteRenderer{
				{Beam: map[int]entity.Beam{1: {}}},
				{Beam: map[int]entity.Beam{1: {}}},
				{Beam: map[int]entity.Beam{1: {}}},
			},
			s:            BeamSplitMarker{StartIndex: 0, EndIndex: 0},
			nextInterval: BeamSplitMarker{StartIndex: 2, EndIndex: 2},
			want:         false,
		},
		// // 4. Failure: No beams on first note
		{
			name: "No Merge - first note no beam",
			notes: []*entity.NoteRenderer{
				{Beam: map[int]entity.Beam{}}, // No beam
				{Beam: map[int]entity.Beam{1: {}}, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
				{Beam: map[int]entity.Beam{1: {}}},
			},
			s:            BeamSplitMarker{StartIndex: 0, EndIndex: 0},
			nextInterval: BeamSplitMarker{StartIndex: 2, EndIndex: 2},
			want:         false,
		},
		// 5. Failure: No beams on second note
		{
			name: "No Merge - next note no beam",
			notes: []*entity.NoteRenderer{
				{Beam: map[int]entity.Beam{1: {}}},
				{Beam: map[int]entity.Beam{1: {}}, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
				{Beam: map[int]entity.Beam{}}, // No beam
			},
			s:            BeamSplitMarker{StartIndex: 0, EndIndex: 0},
			nextInterval: BeamSplitMarker{StartIndex: 2, EndIndex: 2},
			want:         false,
		},
		// 6. Failure: Neither interval has only one note
		{
			name: "No Merge - neither interval has 1 note",
			notes: []*entity.NoteRenderer{
				{Beam: map[int]entity.Beam{1: {}}},
				{Beam: map[int]entity.Beam{1: {}}, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
				{Beam: map[int]entity.Beam{1: {}}}, // 2 notes
			},
			s:            BeamSplitMarker{StartIndex: 0, EndIndex: 1}, // 2 notes
			nextInterval: BeamSplitMarker{StartIndex: 2, EndIndex: 3}, // 2 notes
			want:         false,
		},
		// 7. Success: Only one interval is 1 note, and all other conditions hold
		{
			name: "Merge - only first interval is 1 note",
			notes: []*entity.NoteRenderer{
				{Beam: map[int]entity.Beam{1: {}}},
				{Beam: map[int]entity.Beam{1: {}}, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
				{Beam: map[int]entity.Beam{1: {}}},
			},
			s:            BeamSplitMarker{StartIndex: 0, EndIndex: 0}, // 1 note
			nextInterval: BeamSplitMarker{StartIndex: 2, EndIndex: 3}, // 2 notes
			want:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var qs quaterSpliiter
			got := qs.shouldMergeSegments(tt.notes, tt.s, tt.nextInterval)
			assert.Equal(t, tt.want, got)
		})
	}
}
