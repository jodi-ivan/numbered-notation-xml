package gregorian

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/stretchr/testify/assert"
	ubermock "go.uber.org/mock/gomock"
)

type testifyMatcher struct {
	expected interface{}
	t        *testing.T
}

func (m testifyMatcher) Matches(x interface{}) bool {

	return assert.Equal(m.t, m.expected, x)
}

func (m testifyMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

// Helper function to create a new matcher
func IsEqual(expected interface{}, t *testing.T) gomock.Matcher {
	return testifyMatcher{expected: expected, t: t}
}

func Test_isDottedNote(t *testing.T) {

	uberCtrl := ubermock.NewController(t)
	defer uberCtrl.Finish()
	tests := []struct {
		name    string
		notes   []*entity.NoteRenderer
		notePos int
		ts      func(c *ubermock.Controller) *timesig.MockITimeSig
		want    bool
	}{
		{
			name: "Simple Dotted Note",
			notes: []*entity.NoteRenderer{
				{NoteLength: musicxml.NoteLengthQuarter, NoteValue: 1.0, MeasureNumber: 1},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: []*musicxml.Dot{{}}}}).Return(1.5)
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})
				return res
			},
			want: false,
		},
		{
			name: "Dotted Half Note",
			notes: []*entity.NoteRenderer{
				{NoteLength: musicxml.NoteLengthHalf, NoteValue: 2.0, MeasureNumber: 1},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthHalf, Dot: []*musicxml.Dot{{}}}}).Return(3.0)
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})
				return res
			},
			want: false,
		},
		{
			name: "Dotted Half with Following Dotted Note",
			notes: []*entity.NoteRenderer{
				{NoteLength: musicxml.NoteLengthHalf, NoteValue: 2.0, MeasureNumber: 1},
				{NoteLength: musicxml.NoteLengthQuarter, NoteValue: 1.0, MeasureNumber: 1},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthHalf, Dot: []*musicxml.Dot{{}}}}).Return(3.0)
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})

				return res
			},
			want: false,
		},
		{
			name: "No Dot or Half",
			notes: []*entity.NoteRenderer{
				{NoteLength: musicxml.NoteLengthQuarter, NoteValue: 1.0, MeasureNumber: 1},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: []*musicxml.Dot{{}}}}).Return(1.0)
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})
				return res
			},
			want: true,
		},
		{
			name: "Merged but not in compound time",
			notes: []*entity.NoteRenderer{
				{NoteLength: musicxml.NoteLengthHalf, NoteValue: 2.0, MeasureNumber: 1},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthHalf, Dot: []*musicxml.Dot{{}}}}).Return(3.0)

				// Simulate merged logic being present but time signature doesn't trigger dotted flag
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})
				return res
			},
			want: false, // Assuming logic for merged, non-dotted, non-half is false
		},
		{
			name: "Merged in Compound Time Exception",
			notes: []*entity.NoteRenderer{
				{NoteLength: musicxml.NoteLengthHalf, NoteValue: 2.0, MeasureNumber: 1},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthHalf, Dot: []*musicxml.Dot{{}}}}).Return(3.0)

				// Simulate complex logic where merger in compound time might trigger dotted note check
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 6})
				return res
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := isDottedNote(tt.notes, tt.notePos, tt.ts(uberCtrl))
			if got != tt.want {
				t.Errorf("isDottedNote(pos=%d) = %v, want %v", tt.notePos, got, tt.want)
			}
		})
	}
}

func Test_getAdditionalNotes(t *testing.T) {
	ctrl := ubermock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ts                    func(c *ubermock.Controller) *timesig.MockITimeSig
		staffLines            func(c *ubermock.Controller) *lines.MockLineStaff
		notes                 []*entity.NoteRenderer
		notePos               int
		accumulativeDirection int
		groupBeam             [][]entity.CoordinateWithNoteLength
		wantGroupBeam         [][]entity.CoordinateWithNoteLength
		want                  *entity.NoteRenderer
	}{

		{
			name: "Successful creation of additional note",
			notes: []*entity.NoteRenderer{
				{
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.25, // Between 1.0 (non-dotted) and 1.5 (dotted)
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note1",
				},
				{
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.0,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note2",
					IsDotted:       true,
				},
				{
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.0,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note3",
				},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				// nonDottedValue = 1.0
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: nil}}).Return(1.0)
				// dottedValue = 1.5
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: []*musicxml.Dot{{}}}}).Return(1.5)
				// Check for tied notes (sameNextTwoDottedReplaced)
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})
				return res
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				res.EXPECT().GetYPos('C', 4).Return(50.0)
				return res
			},
			accumulativeDirection: 0,
			groupBeam:             [][]entity.CoordinateWithNoteLength{},
			want: &entity.NoteRenderer{
				AbsoluteNote:       "C",
				AbsoluteAccidental: "",
				AbsoluteOctave:     4,
				PositionX:          0, // notes[0+2].PositionX
				UUID:               "note3",
				NoteLength:         musicxml.NoteLength16th, // remaining = 0.25, noteType[0.25] = Quarter
			},
			wantGroupBeam: [][]entity.CoordinateWithNoteLength{
				{
					{
						Coordinate: entity.NewCoordinate(0, 50),
						NoteLength: musicxml.NoteLength16th,
						NoteID:     "note3",
					},
				},
			},
		},
		{
			name: "Failure: No remaining note value (equal to non-dotted)",
			notes: []*entity.NoteRenderer{
				{
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.0, // Equal to non-dotted value
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note1",
				},
				{
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.0,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note2",
				},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: nil}}).Return(1.0)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: []*musicxml.Dot{{}}}}).Return(1.5)
				// res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})
				return res
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				res.EXPECT().GetYPos('C', 4).Return(50.0)
				return res
			},
			accumulativeDirection: 0,
			groupBeam:             [][]entity.CoordinateWithNoteLength{},
			wantGroupBeam:         [][]entity.CoordinateWithNoteLength{},
			want:                  nil,
		},
		{
			name: "Failure: No tied notes",
			notes: []*entity.NoteRenderer{
				{
					PositionX:      25,
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.25,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note1",
				},
				{
					PositionX:      25,
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.0,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note2",
				},
				{
					PositionX:      25,
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.0,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note3",
				},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: nil}}).Return(1.0)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: []*musicxml.Dot{{}}}}).Return(1.5)
				// Ensure no tied notes condition is met
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})
				return res
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				res.EXPECT().GetYPos('C', 4).Return(50.0)
				return res
			},
			accumulativeDirection: 0,
			groupBeam:             nil,
			want: &entity.NoteRenderer{
				UUID:           "note3",
				PositionX:      25,
				AbsoluteNote:   "C",
				AbsoluteOctave: 4,
				NoteLength:     musicxml.NoteLength16th,
			},
			wantGroupBeam: [][]entity.CoordinateWithNoteLength{
				{
					{
						Coordinate: entity.NewCoordinate(25, 50),
						NoteID:     "note3",
						NoteLength: musicxml.NoteLength16th,
					},
				},
			},
		},
		{
			name: "Failure: Merged (Eighth note + Dotted next note)",
			notes: []*entity.NoteRenderer{
				{
					PositionX:      25,
					NoteLength:     musicxml.NoteLengthEighth,
					NoteValue:      0.75, // Example value for remaining note
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note1",
				},
				{
					PositionX:      25,
					NoteLength:     musicxml.NoteLengthEighth,
					NoteValue:      0.5,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note2",
					IsDotted:       true,
				},
				{
					PositionX:      25,
					NoteLength:     musicxml.NoteLengthEighth,
					NoteValue:      0.5,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note3",
				},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthEighth, Dot: nil}}).Return(0.5)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthEighth, Dot: []*musicxml.Dot{{}}}}).Return(0.75)
				// res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 4})
				return res
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				res.EXPECT().GetYPos('C', 4).Return(50.0)
				return res
			},
			accumulativeDirection: 0,
			groupBeam:             nil,
			want:                  nil,
		},
		{
			name: "Failure: Merged (Quarter note in compound time)",
			notes: []*entity.NoteRenderer{
				{
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.25,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note1",
				},
				{
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.0,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note2",
					IsDotted:       true,
				},
				{
					NoteLength:     musicxml.NoteLengthQuarter,
					NoteValue:      1.0,
					MeasureNumber:  1,
					AbsoluteNote:   "C",
					AbsoluteOctave: 4,
					UUID:           "note3",
				},
			},
			notePos: 0,
			ts: func(c *ubermock.Controller) *timesig.MockITimeSig {
				res := timesig.NewMockITimeSig(c)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: nil}}).Return(1.0)
				res.EXPECT().GetNoteLength(ubermock.Any(), ubermock.Any(), testifyMatcher{t: t, expected: musicxml.Note{Type: musicxml.NoteLengthQuarter, Dot: []*musicxml.Dot{{}}}}).Return(1.5)
				// Simulate compound time (e.g., 6/8)
				res.EXPECT().GetTimesignatureOnMeasure(ubermock.Any(), 1).Return(timesig.Time{BeatType: 6})
				return res
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				res.EXPECT().GetYPos('C', 4).Return(50.0)
				return res
			},
			accumulativeDirection: 0,
			groupBeam:             [][]entity.CoordinateWithNoteLength{},
			want: &entity.NoteRenderer{
				UUID:           "note3",
				AbsoluteNote:   "C",
				Octave:         0,
				AbsoluteOctave: 4,
				NoteLength:     musicxml.NoteLength16th,
			},
			wantGroupBeam: [][]entity.CoordinateWithNoteLength{
				{
					{
						Coordinate: entity.NewCoordinate(0, 50),
						NoteLength: musicxml.NoteLength16th,
						NoteID:     "note3",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAdditionalNotes(context.Background(), tt.ts(ctrl), tt.notes, tt.staffLines(ctrl), tt.notePos, tt.accumulativeDirection, &tt.groupBeam)
			assert.Equal(t, tt.want, got, "additional notes")
			assert.Equal(t, tt.wantGroupBeam, tt.groupBeam, "groupped beam")

		})
	}
}

func Test_setMarginTop(t *testing.T) {
	ctrl := ubermock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		notes             []*entity.NoteRenderer
		staffLines        func(c *ubermock.Controller) *lines.MockLineStaff
		notesChangedAfter []*entity.NoteRenderer
	}{
		{
			name: "Note below top line, StemDirection 0 (downwards)",
			notes: []*entity.NoteRenderer{
				{
					Note:               1,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      0,
					MarginTopFromStaff: 0, // Initial value
				},
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				// Mock Y position below top line
				res.EXPECT().GetYPos('C', 4).Return(100.0)
				// Mock top line Y coordinate
				res.EXPECT().GetTopLine().Return(110).Times(2)
				return res
			},
			notesChangedAfter: []*entity.NoteRenderer{
				{
					Note:               1,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      0,
					MarginTopFromStaff: 10, // 110 - 100
				},
			},
		},
		{
			name: "Note below top line, StemDirection 1 (upwards)",
			notes: []*entity.NoteRenderer{
				{
					Note:               1,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      1,
					MarginTopFromStaff: 0, // Initial value
				},
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				// Mock Y position below top line
				res.EXPECT().GetYPos('C', 4).Return(100.0)
				// Mock top line Y coordinate
				res.EXPECT().GetTopLine().Return(110)
				// Mock STAFF_SPACE_WIDTH for calculation: 2.5 * 10 = 25
				// maxY = 100 + floor(1 * 25) = 125.0
				// 125.0 < 110.0 is false, so MarginTopFromStaff remains 0
				return res
			},
			notesChangedAfter: []*entity.NoteRenderer{
				{
					Note:               1,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      1,
					MarginTopFromStaff: 0,
				},
			},
		},
		{
			name: "Note above top line",
			notes: []*entity.NoteRenderer{
				{
					Note:               1,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      0,
					MarginTopFromStaff: 0, // Initial value
				},
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				// Mock Y position above top line
				res.EXPECT().GetYPos('C', 4).Return(120.0)
				// Mock top line Y coordinate
				res.EXPECT().GetTopLine().Return(110)
				return res
			},
			notesChangedAfter: []*entity.NoteRenderer{
				{
					Note:               1,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      0,
					MarginTopFromStaff: 0,
				},
			},
		},
		{
			name: "Note exactly on top line",
			notes: []*entity.NoteRenderer{
				{
					Note:               1,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      0,
					MarginTopFromStaff: 0, // Initial value
				},
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				// Mock Y position exactly on top line
				res.EXPECT().GetYPos('C', 4).Return(110.0)
				// Mock top line Y coordinate
				res.EXPECT().GetTopLine().Return(110)
				return res
			},
			notesChangedAfter: []*entity.NoteRenderer{
				{
					Note:               1,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      0,
					MarginTopFromStaff: 0,
				},
			},
		},
		{
			name: "Note with Note=0 (should be skipped)",
			notes: []*entity.NoteRenderer{
				{
					Note:               0,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      1,
					MarginTopFromStaff: 0, // Initial value
				},
			},
			staffLines: func(c *ubermock.Controller) *lines.MockLineStaff {
				res := lines.NewMockLineStaff(c)
				// Should not be called
				return res
			},
			notesChangedAfter: []*entity.NoteRenderer{
				{
					Note:               0,
					AbsoluteNote:       "C",
					AbsoluteOctave:     4,
					StemDirection:      1,
					MarginTopFromStaff: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// // Create a deep copy of notes to check changes
			// notesCopy := make([]*entity.NoteRenderer, len(tt.notes))
			// for i, n := range tt.notes {
			// 	notesCopy[i] = &entity.NoteRenderer{
			// 		Note:               n.Note,
			// 		AbsoluteNote:       n.AbsoluteNote,
			// 		AbsoluteOctave:     n.AbsoluteOctave,
			// 		StemDirection:      n.StemDirection,
			// 		MarginTopFromStaff: n.MarginTopFromStaff,
			// 	}
			// }

			setMarginTop(tt.notes, tt.staffLines(ctrl))
			assert.Equal(t, tt.notesChangedAfter, tt.notes)
		})
	}
}
