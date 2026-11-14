package splitter

import (
	"context"
	"fmt"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
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

func Test_beamSplitter_Split(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		eighth    func(c *gomock.Controller) *MockBeamSplitter
		quarter   func(c *gomock.Controller) *MockBeamSplitter
		tuplet    func(c *gomock.Controller) *MockTupletSplitter
		sanitizer func(c *gomock.Controller) *MockSanitizer
		ts        func(c *gomock.Controller) *timesig.MockITimeSig
		notes     []*entity.NoteRenderer
	}{
		{
			name: "Standard Quarter Split",
			quarter: func(c *gomock.Controller) *MockBeamSplitter {
				mock := NewMockBeamSplitter(c)
				// Mock Split call for quarter
				mock.EXPECT().
					Split(gomock.Any(), gomock.Any(), gomock.Any(), testifyMatcher{t: t, expected: map[int][]BeamSplitMarker{1: {{StartIndex: 0, EndIndex: 4}}, 2: {{StartIndex: 0, EndIndex: 2}, {StartIndex: 3, EndIndex: 4}}}}).
					Return()
				return mock
			},
			sanitizer: func(c *gomock.Controller) *MockSanitizer {
				mock := NewMockSanitizer(c)
				// Mock CleanBeamByNumber calls
				mock.EXPECT().
					CleanBeamByNumber(gomock.Any(), gomock.Any(), 1).
					Return([]BeamSplitMarker{{StartIndex: 0, EndIndex: 4}})
				mock.EXPECT().
					CleanBeamByNumber(gomock.Any(), gomock.Any(), 2).
					Return([]BeamSplitMarker{{StartIndex: 0, EndIndex: 2}, {StartIndex: 3, EndIndex: 4}})
				return mock
			},
			tuplet: func(c *gomock.Controller) *MockTupletSplitter {
				mock := NewMockTupletSplitter(c)
				// Mock Split call for tuplet
				mock.EXPECT().
					Split(gomock.Any(), testifyMatcher{t: t, expected: []BeamSplitMarker{{StartIndex: 0, EndIndex: 4}}}).
					Return([]BeamSplitMarker{{StartIndex: 0, EndIndex: 4}})
				return mock
			},
			ts: func(c *gomock.Controller) *timesig.MockITimeSig {
				mock := timesig.NewMockITimeSig(c)
				// Mock GetTimesignatureOnMeasure to return a time signature with BeatType 4
				mock.EXPECT().
					GetTimesignatureOnMeasure(gomock.Any(), gomock.Any()).
					Return(timesig.Time{BeatType: 4})
				return mock
			},
			notes: []*entity.NoteRenderer{
				{MeasureNumber: 1},
			},
		},
		{
			name: "Standard Eighth Split",
			eighth: func(c *gomock.Controller) *MockBeamSplitter {
				mock := NewMockBeamSplitter(c)
				// Mock Split call for quarter
				mock.EXPECT().
					Split(gomock.Any(), gomock.Any(), gomock.Any(), testifyMatcher{t: t, expected: map[int][]BeamSplitMarker{1: {{StartIndex: 0, EndIndex: 4}}, 2: {{StartIndex: 0, EndIndex: 2}, {StartIndex: 3, EndIndex: 4}}}}).
					Return()
				return mock
			},
			sanitizer: func(c *gomock.Controller) *MockSanitizer {
				mock := NewMockSanitizer(c)
				// Mock CleanBeamByNumber calls
				mock.EXPECT().
					CleanBeamByNumber(gomock.Any(), gomock.Any(), 1).
					Return([]BeamSplitMarker{{StartIndex: 0, EndIndex: 4}})
				mock.EXPECT().
					CleanBeamByNumber(gomock.Any(), gomock.Any(), 2).
					Return([]BeamSplitMarker{{StartIndex: 0, EndIndex: 2}, {StartIndex: 3, EndIndex: 4}})
				return mock
			},
			tuplet: func(c *gomock.Controller) *MockTupletSplitter {
				mock := NewMockTupletSplitter(c)
				// Mock Split call for tuplet
				mock.EXPECT().
					Split(gomock.Any(), testifyMatcher{t: t, expected: []BeamSplitMarker{{StartIndex: 0, EndIndex: 4}}}).
					Return([]BeamSplitMarker{{StartIndex: 0, EndIndex: 4}})
				return mock
			},
			ts: func(c *gomock.Controller) *timesig.MockITimeSig {
				mock := timesig.NewMockITimeSig(c)
				// Mock GetTimesignatureOnMeasure to return a time signature with BeatType 4
				mock.EXPECT().
					GetTimesignatureOnMeasure(gomock.Any(), gomock.Any()).
					Return(timesig.Time{BeatType: 8})
				return mock
			},
			notes: []*entity.NoteRenderer{
				{MeasureNumber: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := beamSplitter{
				tuplet:    tt.tuplet(ctrl),
				sanitizer: tt.sanitizer(ctrl),
			}
			//
			if tt.eighth != nil {
				bs.eighth = tt.eighth(ctrl)
			}
			if tt.quarter != nil {
				bs.quarter = tt.quarter(ctrl)
			}

			bs.Split(context.Background(), tt.notes, tt.ts(ctrl), nil)
		})
	}
}

func Test_beamSplitter_SplitSingle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		quarter  func(c *gomock.Controller) *MockBeamSplitter
		ts       func(c *gomock.Controller) *timesig.MockITimeSig
		notes    []*entity.NoteRenderer
		segments []BeamSplitMarker
		beamNo   int
		eighth   func(c *gomock.Controller) *MockBeamSplitter
	}{
		{
			name: "Standard Quarter Split Single",
			quarter: func(c *gomock.Controller) *MockBeamSplitter {
				mock := NewMockBeamSplitter(c)
				// Mock SplitSingle call for quarter
				mock.EXPECT().
					SplitSingle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return()
				return mock
			},
			ts: func(c *gomock.Controller) *timesig.MockITimeSig {
				mock := timesig.NewMockITimeSig(c)
				// Mock GetTimesignatureOnMeasure to return a time signature with BeatType 4
				mock.EXPECT().
					GetTimesignatureOnMeasure(gomock.Any(), gomock.Any()).
					Return(timesig.Time{BeatType: 4})
				return mock
			},
			notes: []*entity.NoteRenderer{
				{MeasureNumber: 1},
			},
			segments: []BeamSplitMarker{{StartIndex: 0, EndIndex: 4}},
			beamNo:   1,
		},
		{
			name: "Standard Eighth Split Single",
			eighth: func(c *gomock.Controller) *MockBeamSplitter {
				mock := NewMockBeamSplitter(c)
				// Mock SplitSingle call for eighth
				mock.EXPECT().
					SplitSingle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return()
				return mock
			},
			ts: func(c *gomock.Controller) *timesig.MockITimeSig {
				mock := timesig.NewMockITimeSig(c)
				// Mock GetTimesignatureOnMeasure to return a time signature with BeatType 8
				mock.EXPECT().
					GetTimesignatureOnMeasure(gomock.Any(), gomock.Any()).
					Return(timesig.Time{BeatType: 8})
				return mock
			},
			notes: []*entity.NoteRenderer{
				{MeasureNumber: 1},
			},
			segments: []BeamSplitMarker{{StartIndex: 0, EndIndex: 4}},
			beamNo:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := beamSplitter{
				// eighth:  tt.eighth(ctrl),
				// quarter: tt.quarter(ctrl),
				// tuplet and sanitizer are not used in SplitSingle logic, but we must initialize them
				tuplet:    &tuplet{},
				sanitizer: &sanitizer{},
			}

			if tt.eighth != nil {
				bs.eighth = tt.eighth(ctrl)
			}
			if tt.quarter != nil {
				bs.quarter = tt.quarter(ctrl)
			}

			bs.SplitSingle(context.Background(), tt.notes, tt.ts(ctrl), tt.segments, tt.beamNo)
		})
	}
}
