package splitter

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func Test_splitTuplet(t *testing.T) {
	// Helper function to create a dummy NoteRenderer
	createNote := func(tuplet *musicxml.Tuplet) *entity.NoteRenderer {
		return &entity.NoteRenderer{
			Tuplet: tuplet,
		}
	}

	// Helper function to create a dummy segment
	createSegment := func(start, end int) BeamSplitMarker {
		return BeamSplitMarker{StartIndex: start, EndIndex: end}
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		notes    []*entity.NoteRenderer
		segments []BeamSplitMarker
		want     []BeamSplitMarker
	}{
		{
			name: "No Tuplets",
			notes: []*entity.NoteRenderer{
				createNote(nil),
				createNote(nil),
			},
			segments: []BeamSplitMarker{
				createSegment(0, 10),
				createSegment(11, 20),
			},
			want: []BeamSplitMarker{
				createSegment(0, 10),
				createSegment(11, 20),
			},
		},
		{
			name:     "Empty Notes",
			notes:    []*entity.NoteRenderer{},
			segments: []BeamSplitMarker{createSegment(0, 10)},
			want:     []BeamSplitMarker{createSegment(0, 10)},
		},
		{
			name: "Empty Segments",
			notes: []*entity.NoteRenderer{
				createNote(nil),
			},
			segments: []BeamSplitMarker{},
			want:     []BeamSplitMarker{},
		},
		{
			name: "TupletStart_SegmentBefore",
			notes: []*entity.NoteRenderer{
				createNote(nil),
				createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStart}), // i=1
			},
			segments: []BeamSplitMarker{
				createSegment(0, 5), // segment.StartIndex = 0
				createSegment(6, 10),
			},
			want: []BeamSplitMarker{
				createSegment(0, 0),
			},
		},
		{
			name: "TupletStart_SegmentAfter",
			notes: []*entity.NoteRenderer{
				createNote(nil),
				createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStart}), // i=1
			},
			segments: []BeamSplitMarker{
				createSegment(2, 5), // segment.StartIndex = 2
				createSegment(6, 10),
			},
			want: []BeamSplitMarker{
				createSegment(2, 5),
				createSegment(6, 10),
				createSegment(0, 0),
			},
		},
		{
			name: "TupletStop_Basic",
			notes: []*entity.NoteRenderer{
				createNote(nil),
				createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStart}), // i=1 (Start)
				createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStop}),  // i=2 (Stop)
				createNote(nil), // Added a fourth note to prevent out-of-bounds access at notes[EndIndex+1]
			},
			segments: []BeamSplitMarker{
				createSegment(0, 10), // segment.StartIndex = 0, segment.EndIndex = 10
			},
			// TupletStop logic:
			// 1. result appends marker (-1, 2)
			// 2. result appends marker (-1, 2) again
			// 3. result appends (marker.EndIndex + 1, segment.EndIndex) = (3, 10)
			want: []BeamSplitMarker{
				createSegment(0, 0),

				createSegment(1, 2),
				createSegment(1, 2),
				createSegment(3, 10),
			},
		},
		// {
		// 	name: "TupletStop_NoPrecedingSegment",
		// 	notes: []*entity.NoteRenderer{
		// 		createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStart}), // i=1 (Start)
		// 		createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStop}),  // i=2 (Stop)
		// 	},
		// 	segments: []BeamSplitMarker{
		// 		createSegment(1, 10), // segment.StartIndex = 1
		// 	},
		// 	// TupletStop logic:
		// 	// 1. result appends marker (-1, 2)
		// 	// 2. result appends marker (-1, 2) again
		// 	// 3. result appends (3, 10)
		// 	want: []BeamSplitMarker{
		// 		createSegment(-1, 2),
		// 		createSegment(-1, 2),
		// 		createSegment(3, 10),
		// 	},
		// },
		// {
		// 	name: "TupletStop_PrecedingSegmentExists",
		// 	notes: []*entity.NoteRenderer{
		// 		createNote(nil),
		// 		createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStart}), // i=1 (Start)
		// 		createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStop}),  // i=2 (Stop)
		// 	},
		// 	segments: []BeamSplitMarker{
		// 		createSegment(0, 10), // segment.StartIndex = 0
		// 	},
		// 	// TupletStop logic:
		// 	// 1. result appends marker (-1, 2)
		// 	// 2. result appends marker (-1, 2) again
		// 	// 3. result appends (3, 10)
		// 	want: []BeamSplitMarker{
		// 		createSegment(-1, 2),
		// 		createSegment(-1, 2),
		// 		createSegment(3, 10),
		// 	},
		// },
		{
			name: "TupletStop_NoFollowingSegment",
			notes: []*entity.NoteRenderer{
				createNote(nil),
				createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStart}), // i=1 (Start)
				createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStop}),  // i=2 (Stop)
			},
			segments: []BeamSplitMarker{
				createSegment(0, 2), // segment.EndIndex = 2
			},
			// TupletStop logic:
			// 1. result appends marker (-1, 2)
			// 2. result appends marker (-1, 2) again
			// 3. segment.EndIndex (2) > marker.EndIndex (2) is false. No final segment added.
			want: []BeamSplitMarker{

				createSegment(0, 0),
				createSegment(1, 2),
				createSegment(1, 2),
			},
		},
		{
			name: "TupletStop_ExactMatch",
			notes: []*entity.NoteRenderer{
				createNote(nil),
				createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStart}), // i=1 (Start)
				createNote(&musicxml.Tuplet{Type: musicxml.TupletTypeStop}),  // i=2 (Stop)
			},
			segments: []BeamSplitMarker{
				createSegment(0, 2), // segment.EndIndex = 2
			},
			// TupletStop logic:
			// 1. result appends marker (-1, 2)
			// 2. result appends marker (-1, 2) again
			// 3. segment.EndIndex (2) > marker.EndIndex (2) is false. No final segment added.
			want: []BeamSplitMarker{
				createSegment(0, 0),

				createSegment(1, 2),
				createSegment(1, 2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tu := &tuplet{}
			got := tu.Split(tt.notes, tt.segments)

			assert.Equal(t, tt.want, got)
		})
	}
}
