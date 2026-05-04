package staff

import (
	"context"
	"fmt"
	reflect "reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/stretchr/testify/assert"
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

func TestNewStaff(t *testing.T) {
	t.Run("new staff", func(t *testing.T) {
		got := NewStaff()
		interactor := got.(*staffInteractor)

		v := reflect.ValueOf(*interactor)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			assert.NotNil(t, v.Field(i).Interface(), fmt.Sprintf("Field name: %s", typeOfS.Field(i).Name))
		}
	})
}

func Test_staffInteractor_RenderStaff(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	kj001 := []musicxml.Measure{
		{
			Number: 1,
			Attribute: &musicxml.Attribute{
				Key: &musicxml.KeySignature{
					Fifth: -1, // F major
				},
				Time: &struct {
					Beats    int `xml:"beats"`
					BeatType int `xml:"beat-type"`
				}{
					Beats:    4,
					BeatType: 4,
				},
			},
			Appendix: []musicxml.Element{
				{
					Content: `<direction-type><words>Refrein</words></direction-type>`,
				},
				{
					Content: `<pitch><step>A</step><octave>4</octave></pitch><duration>2</duration><type>quarter</type><lyric number="1"><syllabic>begin</syllabic><text>Ha</text></lyric>`,
				},
			},
		},
	}

	kj075 := []musicxml.Measure{
		{
			Number: 2,
			Attribute: &musicxml.Attribute{
				Key: &musicxml.KeySignature{
					Fifth: 0, // a minor
					Mode:  "minor",
				},
				Time: &struct {
					Beats    int `xml:"beats"`
					BeatType int `xml:"beat-type"`
				}{
					Beats:    4,
					BeatType: 4,
				},
			},
			Appendix: []musicxml.Element{
				{
					Content: `<pitch><step>F</step><octave>5</octave></pitch><duration>4</duration><tie type="start"/><type>half</type><notations><tied type="start"/><slur type="start" number="1"/></notations><lyric number="1"><syllabic>begin</syllabic><text>Da</text></lyric>`,
				},
				{
					Content: `<pitch><step>F</step><octave>5</octave></pitch><duration>1</duration><tie type="stop"/><type>eighth</type><notations><tied type="stop"/></notations>`,
				},
			},
		},
	}
	signatures := map[string]struct {
		ts timesig.TimeSignature
		ks keysig.KeySignature
	}{
		"kj001": {
			ts: timesig.NewTimeSignatures(context.Background(), kj001),
			ks: keysig.NewKeySignature(context.Background(), kj001),
		},
		"kj075": {
			ts: timesig.NewTimeSignatures(context.Background(), kj075),
			ks: keysig.NewKeySignature(context.Background(), kj075),
		},
	}

	getNote := func(appx []musicxml.Element, i int) musicxml.Note {
		n, _ := appx[i].ParseAsNote()
		if i > 0 {
			d, err := appx[i-1].ParseAsDirection()
			if err == nil && len(d.DirectionType) > 0 {
				if d.DirectionType[0].Word.Value != "" {
					n.MeasureText = []musicxml.MeasureText{{Text: d.DirectionType[0].Word.Value}}
				}
			}
		}
		return n
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		x               int
		y               int
		isLastStaff     bool
		keySignature    keysig.KeySignature
		timeSignature   timesig.TimeSignature
		measures        []musicxml.Measure
		prevNotes       []*entity.NoteRenderer
		numberedMock    func(c *gomock.Controller) *numbered.MockNumbered
		rhythmMock      func(c *gomock.Controller) *rhythm.MockRhythm
		lyricMock       func(c *gomock.Controller) *lyric.MockLyric
		breathpauseMock func(c *gomock.Controller) *breathpause.MockBreathPause
		barlineMock     func(c *gomock.Controller) *barline.MockBarline
		renderAlignMock func(c *gomock.Controller) *MockRenderStaffWithAlign
		want            StaffInfo
	}{
		{
			name:          "kj-001-first-staff",
			x:             50,
			y:             95,
			keySignature:  signatures["kj001"].ks,
			timeSignature: signatures["kj001"].ts,
			measures:      kj001,
			want: StaffInfo{
				MarginBottom:     15,
				NextLineRenderer: []*entity.NoteRenderer{},
			},
			numberedMock: func(c *gomock.Controller) *numbered.MockNumbered {
				nm := numbered.NewMockNumbered(c)
				nm.EXPECT().GetLengthNote(gomock.Any(), signatures["kj001"].ts, 1, float64(1)).Return([]numbered.NoteLength{{Type: musicxml.NoteLengthQuarter}})
				nm.EXPECT().RendererFromAdditional(testifyMatcher{t: t, expected: getNote(kj001[0].Appendix, 1)}, gomock.Any(), []numbered.NoteLength{{Type: musicxml.NoteLengthQuarter}})
				return nm
			},
			rhythmMock: func(c *gomock.Controller) *rhythm.MockRhythm {
				rm := rhythm.NewMockRhythm(c)
				rm.EXPECT().SetRhythmNotation(gomock.Any(), gomock.Any(), 3)
				rm.EXPECT().AdjustMultiDottedRenderer(gomock.Any(), 50, 95, signatures["kj001"].ks).Return(50, 95)
				return rm
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().SetLyricRenderer(gomock.Any(), gomock.Any())

				return li
			},
			breathpauseMock: func(c *gomock.Controller) *breathpause.MockBreathPause {
				bp := breathpause.NewMockBreathPause(c)
				bp.EXPECT().SetAndGetBreathPauseRenderer(gomock.Any(), testifyMatcher{t: t, expected: getNote(kj001[0].Appendix, 1)})
				return bp
			},
			barlineMock: func(c *gomock.Controller) *barline.MockBarline {
				bm := barline.NewMockBarline(c)
				bm.EXPECT().GetRendererRightBarline(gomock.Any(), 50).Return(50, &entity.NoteRenderer{PositionX: 90, Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular}})
				return bm
			},
			renderAlignMock: func(c *gomock.Controller) *MockRenderStaffWithAlign {
				ra := NewMockRenderStaffWithAlign(c)
				ra.EXPECT().RenderWithAlign(gomock.Any(), gomock.Any(), 0, 110, signatures["kj001"].ts, signatures["kj001"].ks, gomock.Any())
				return ra
			},
		},
		{
			name:          "kj-075-first-staff",
			x:             50,
			y:             95,
			keySignature:  signatures["kj075"].ks,
			timeSignature: signatures["kj075"].ts,
			measures:      kj075,
			want: StaffInfo{
				MarginBottom:     0,
				NextLineRenderer: []*entity.NoteRenderer{},
			},
			numberedMock: func(c *gomock.Controller) *numbered.MockNumbered {
				nm := numbered.NewMockNumbered(c)
				nm.EXPECT().GetLengthNote(gomock.Any(), signatures["kj075"].ts, 2, 2.5).Return([]numbered.NoteLength{{Type: musicxml.NoteLengthQuarter}})
				ts := signatures["kj075"].ts.GetTimesignatureOnMeasure(context.Background(), 2)
				nm.EXPECT().SplitNote(gomock.Any(), 2.5, ts, musicxml.NoteLengthHalf, musicxml.NoteLengthEighth).
					Return([]numbered.NoteLength{
						{Type: musicxml.NoteLengthQuarter},
						{Type: musicxml.NoteLengthQuarter, IsDotted: true},
						{Type: musicxml.NoteLengthEighth, IsDotted: true},
					})
				expectedNote := getNote(kj075[0].Appendix, 0)
				if expectedNote.Notations != nil {
					expectedNote.Notations.Tied = nil
				}
				nm.EXPECT().RendererFromAdditional(testifyMatcher{t: t, expected: expectedNote}, gomock.Any(), []numbered.NoteLength{
					{Type: musicxml.NoteLengthQuarter},
					{Type: musicxml.NoteLengthQuarter, IsDotted: true},
					{Type: musicxml.NoteLengthEighth, IsDotted: true},
				}).Return([]*entity.NoteRenderer{
					{NoteLength: musicxml.NoteLengthQuarter},
					{NoteLength: musicxml.NoteLengthQuarter, IsDotted: true},
					{NoteLength: musicxml.NoteLengthHalf, IsDotted: true},
				})
				return nm
			},
			rhythmMock: func(c *gomock.Controller) *rhythm.MockRhythm {
				rm := rhythm.NewMockRhythm(c)
				rm.EXPECT().SetRhythmNotation(gomock.Any(), gomock.Any(), 6)
				rm.EXPECT().AdjustMultiDottedRenderer(gomock.Any(), 50, 95, signatures["kj075"].ks).Return(50, 95)
				return rm
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().SetLyricRenderer(gomock.Any(), gomock.Any())

				return li
			},
			breathpauseMock: func(c *gomock.Controller) *breathpause.MockBreathPause {
				bp := breathpause.NewMockBreathPause(c)
				bp.EXPECT().SetAndGetBreathPauseRenderer(gomock.Any(), gomock.Any())
				return bp
			},
			barlineMock: func(c *gomock.Controller) *barline.MockBarline {
				bm := barline.NewMockBarline(c)
				bm.EXPECT().GetRendererRightBarline(gomock.Any(), 50).Return(50, &entity.NoteRenderer{PositionX: 90, Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular}})
				return bm
			},
			renderAlignMock: func(c *gomock.Controller) *MockRenderStaffWithAlign {
				ra := NewMockRenderStaffWithAlign(c)
				ra.EXPECT().RenderWithAlign(gomock.Any(), gomock.Any(), 0, 95, signatures["kj075"].ts, signatures["kj075"].ks, gomock.Any())
				return ra
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var si staffInteractor
			si.Numbered = tt.numberedMock(ctrl)
			si.Rhythm = tt.rhythmMock(ctrl)
			si.Lyric = tt.lyricMock(ctrl)
			si.BreathPause = tt.breathpauseMock(ctrl)
			si.Barline = tt.barlineMock(ctrl)
			si.RenderAlign = tt.renderAlignMock(ctrl)

			got := si.RenderStaff(context.Background(), nil, tt.x, tt.y, 0, tt.isLastStaff, tt.keySignature, tt.timeSignature, tt.measures, tt.prevNotes...)

			assert.Equal(t, tt.want, got, "StaffInfo assert")

		})
	}
}
