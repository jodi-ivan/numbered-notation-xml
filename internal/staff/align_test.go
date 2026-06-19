package staff

import (
	"context"
	"slices"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/gregorian"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/text"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/toping"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_alignJustify(t *testing.T) {
	count := func() *int {
		c := 10

		return &c
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		measure      []*entity.NoteRenderer
		y            int
		addedSpace   float64
		count        *int
		measureIndex int
		lastMeasure  bool
		wantMeasures []*entity.NoteRenderer
		wantCount    int
	}{
		{
			name:         "First measure",
			y:            100,
			addedSpace:   8,
			count:        count(),
			wantCount:    14,
			measureIndex: 0,
			measure: []*entity.NoteRenderer{
				{PositionX: 50},
				{PositionX: 55, IsDotted: true}, {PositionX: 60, IsDotted: true}, {PositionX: 65, IsDotted: true},
				{PositionX: 70, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
			},
			wantMeasures: []*entity.NoteRenderer{
				{PositionX: 50, PositionY: 100},
				{PositionX: 79, PositionY: 100, IsDotted: true}, {PositionX: 108, PositionY: 100, IsDotted: true}, {PositionX: 137, PositionY: 100, IsDotted: true},
				{PositionX: 169, PositionY: 100, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
			},
		},
		{
			name:         "Last measure",
			y:            100,
			addedSpace:   8,
			count:        count(),
			wantCount:    14,
			measureIndex: 1,
			lastMeasure:  true,
			measure: []*entity.NoteRenderer{
				{PositionX: 50},
				{PositionX: 55, IsDotted: true}, {PositionX: 60, IsDotted: true}, {PositionX: 65, IsDotted: true},
				{PositionX: 70, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
			},
			wantMeasures: []*entity.NoteRenderer{
				{PositionX: 130, PositionY: 100},
				{PositionX: 115, PositionY: 100, IsDotted: true}, {PositionX: 100, PositionY: 100, IsDotted: true}, {PositionX: 85, PositionY: 100, IsDotted: true},
				{PositionX: 70, PositionY: 100, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alignJustify(tt.measure, tt.y, tt.addedSpace, tt.count, tt.measureIndex, tt.lastMeasure)

			assert.Equal(t, tt.wantMeasures, tt.measure, "alignJustify")
			assert.Equal(t, tt.wantCount, *tt.count, "Count aligned")
		})
	}
}

func Test_renderStaffAlign_getAddedSpace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rightAlignOffset := func() *int {
		rao := 0

		return &rao
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		lastNote   *entity.NoteRenderer
		lyricMock  func(c *gomock.Controller) *lyric.MockLyric
		totalNotes int
		want       float64
		want2      int

		rightAlignOffset      *int
		wantRrightAlignOffset int
	}{
		{
			name:             "last note is barline",
			rightAlignOffset: rightAlignOffset(),
			lastNote: &entity.NoteRenderer{
				PositionX: 650,
				Barline: &musicxml.Barline{
					BarStyle: musicxml.BarLineStyleHeavyLight,
				},
			},
			totalNotes: 8,

			want:  11.5375,
			want2: 750,

			wantRrightAlignOffset: 0,
		},
		{
			name:             "last note is lyric",
			rightAlignOffset: rightAlignOffset(),
			lastNote: &entity.NoteRenderer{
				PositionX: 650,
				Lyric: []entity.Lyric{
					{Text: []entity.Text{{Value: "unit"}}},
					{Text: []entity.Text{{Value: "testing"}}},
				},
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateOverallWidth([]entity.Lyric{
					{Text: []entity.Text{{Value: "unit"}}},
					{Text: []entity.Text{{Value: "testing"}}},
				}).Return(10.0)
				return li
			},
			totalNotes: 8,
			want:       11.25,
			want2:      750,

			wantRrightAlignOffset: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rsa renderStaffAlign
			if tt.lyricMock != nil {
				rsa.Lyric = tt.lyricMock(ctrl)
			}
			got, got2 := rsa.getAddedSpace(tt.lastNote, tt.rightAlignOffset, tt.totalNotes)

			assert.Equal(t, tt.want, got, "renderStaffAlign_getAddedSpace --> added")
			assert.Equal(t, tt.want2, got2, "renderStaffAlign_getAddedSpace --> lastPos")
			assert.Equal(t, tt.wantRrightAlignOffset, *tt.rightAlignOffset, "renderStaffAlign_getAddedSpace --> &(rightAlignOffset)")

		})
	}
}

func Test_renderStaffAlign_RenderWithAlign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	measures := []musicxml.Measure{
		{
			Number: 1,
			Attribute: &musicxml.Attribute{
				Key: &musicxml.KeySignature{
					Fifth: 2, // D major
				},
				Time: &struct {
					Beats    int `xml:"beats"`
					BeatType int `xml:"beat-type"`
				}{
					Beats:    4,
					BeatType: 4,
				},
			},
		},
	}

	expectMeasures := map[int][]*entity.NoteRenderer{
		1: {
			{MeasureNumber: 1, PositionX: 50, PositionY: 100},
			{MeasureNumber: 1, PositionX: 106, PositionY: 100, IsDotted: true}, {MeasureNumber: 1, PositionX: 162, PositionY: 100, IsDotted: true}, {MeasureNumber: 1, PositionX: 218, PositionY: 100, IsDotted: true},
			{MeasureNumber: 1, PositionX: 277, PositionY: 100, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
			{MeasureNumber: 1, PositionX: 340, PositionY: 100, Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular}},
		},
		2: {
			{MeasureNumber: 2, PositionX: 398, PositionY: 100},
			{MeasureNumber: 2, PositionX: 455, PositionY: 100, IsDotted: true}, {MeasureNumber: 2, PositionX: 512, PositionY: 100, IsDotted: true}, {MeasureNumber: 2, PositionX: 569, PositionY: 100, IsDotted: true},
			{MeasureNumber: 2, PositionX: 626, PositionY: 100, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
			{MeasureNumber: 2, PositionX: 754, PositionY: 100, Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleHeavyLight}},
		},
	}

	ks := keysig.NewKeySignature(context.TODO(), measures)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv         func(c *gomock.Controller) *canvas.MockCanvasTestify
		numberedMock func(c *gomock.Controller) *numbered.MockNumbered
		rhythmMock   func(c *gomock.Controller) *rhythm.MockRhythm
		lyricMock    func(c *gomock.Controller) *lyric.MockLyric
		y            int
		ts           timesig.TimeSignature
		noteRenderer [][]*entity.NoteRenderer
	}{
		{
			name: "default",
			noteRenderer: [][]*entity.NoteRenderer{
				{
					{MeasureNumber: 1, PositionX: 50},
					{MeasureNumber: 1, PositionX: 55, IsDotted: true}, {MeasureNumber: 1, PositionX: 60, IsDotted: true}, {MeasureNumber: 1, PositionX: 65, IsDotted: true},
					{MeasureNumber: 1, PositionX: 70, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
					{MeasureNumber: 1, PositionX: 75, Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular}},
				},
				{
					{MeasureNumber: 2, PositionX: 80},
					{MeasureNumber: 2, PositionX: 85, IsDotted: true}, {MeasureNumber: 2, PositionX: 90, IsDotted: true}, {MeasureNumber: 2, PositionX: 95, IsDotted: true},
					{MeasureNumber: 2, PositionX: 100, Articulation: &entity.Articulation{BreathMark: &entity.ArticulationTypesBreathMark}},
					{MeasureNumber: 2, PositionX: 105, Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleHeavyLight}},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvasTestify {
				canv := canvas.NewMockCanvasTestify(t)
				canv.EXPECT().Group([]string{"class='staff'"})
				canv.EXPECT().Group([]string{`class="numbered"`})
				canv.EXPECT().Group([]string{"class='staff-text'"})
				canv.EXPECT().Group([]string{"class=\"gregorian\"", "style='font-family:mozart11'"})
				canv.EXPECT().Group([]string{"class='measure-align'", "number='1'"})
				canv.EXPECT().Group([]string{"class='measure-align'", "number='2'"})
				canv.EXPECT().Group([]string{"class='note'", "style='font-family:Old Standard TT;font-weight:500'"}).Times(2)
				canv.EXPECT().Gend().Times(9)
				return canv
			},
			ts: timesig.NewTimeSignatures(context.Background(), measures),
			y:  100,
			numberedMock: func(c *gomock.Controller) *numbered.MockNumbered {
				mock := numbered.NewMockNumbered(c)
				mock.EXPECT().RenderNote(gomock.Any(), gomock.Any(), gomock.Any(), 165, 0).Times(2)

				return mock
			},
			rhythmMock: func(c *gomock.Controller) *rhythm.MockRhythm {
				mock := rhythm.NewMockRhythm(c)
				ts := timesig.NewTimeSignatures(context.Background(), measures)
				mock.EXPECT().RenderBeam(gomock.Any(), 165, gomock.Any(), ts, gomock.Any()).Times(2)
				mock.EXPECT().RenderSlurTies(gomock.Any(), 165, gomock.Any(),
					testifyMatcher{expected: []*entity.NoteRenderer{}, t: t},
					float64(750),
				)
				mock.EXPECT().Split(gomock.Any(), ts, gomock.Any())
				mock.EXPECT().Split(gomock.Any(), ts, gomock.Any())

				return mock
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				mock := lyric.NewMockLyric(c)
				mock.EXPECT().RenderLyrics(gomock.Any(), 165, gomock.Any(), gomock.Any()).Times(2) //&testifyMatcher{t: t, expected: expectMeasures[1]})
				mock.EXPECT().RenderHypen(gomock.Any(), 165, 0, gomock.Any(), slices.Concat(expectMeasures[1], expectMeasures[2]))
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var rsa renderStaffAlign
			rsa.Numbered = tt.numberedMock(ctrl)
			rsa.Rhythm = tt.rhythmMock(ctrl)
			rsa.Lyric = tt.lyricMock(ctrl)
			gregMock := gregorian.NewMockGregorian(ctrl)
			gregMock.EXPECT().RenderStaffLine(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			rsa.Gregorian = gregMock
			canv := tt.canv(ctrl)

			topingMock := toping.NewMockToping(t)
			topingMock.EXPECT().RenderRepeatMeasure(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
			topingMock.EXPECT().RenderRepeatMeasure(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
			topingMock.EXPECT().RenderStaffLineDash(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
			topingMock.EXPECT().RenderTuplet(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
			rsa.Toping = topingMock

			mtMock := text.NewMockText(t)
			mtMock.EXPECT().RenderMeasureText(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
			rsa.Text = mtMock
			rsa.RenderWithAlign(context.Background(), canv, 0, tt.y, tt.ts, ks, tt.noteRenderer)
			canv.AssertExpectations(t)
			mtMock.AssertExpectations(t)
		})
	}
}
