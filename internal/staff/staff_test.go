package staff

import (
	"context"
	"fmt"
	reflect "reflect"
	"slices"
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
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
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

func Test_staffInteractor_RenderStaff(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		x             int
		y             int
		keySignature  keysig.KeySignature
		timeSignature timesig.TimeSignature
		measures      []musicxml.Measure
		prevNotes     []*entity.NoteRenderer
	}
	tests := []struct {
		name string
		args args
		si   func(ctrl *gomock.Controller) *staffInteractor
		canv func(ctrl *gomock.Controller) *canvas.MockCanvas

		wantStaffInfo StaffInfo
	}{
		{
			name: "case #1",
			args: args{
				x: 50,
				y: 80,
				keySignature: keysig.NewKeySignature(musicxml.KeySignature{
					Fifth: 2, // D Major
				}),
				timeSignature: timesig.TimeSignature{IsMixed: false, Signatures: []timesig.Time{timesig.Time{Measure: 1, Beat: 2, BeatType: 4}}},
				measures: []musicxml.Measure{
					musicxml.Measure{
						Number: 1,
						Appendix: []musicxml.Element{
							musicxml.Element{Content: `<direction-type><words>Refrein</words></direction-type>`},
							musicxml.Element{Content: `<rest/><duration>4</duration><voice>1</voice><type>half</type>`},
							musicxml.Element{Content: `<rest/><duration>4</duration><voice>1</voice><type>half</type>`},
							musicxml.Element{Content: `<pitch><step>A</step><octave>4</octave></pitch><duration>2</duration><voice>1</voice><type>quarter</type><stem>up</stem><lyric number="1"><syllabic>begin</syllabic><text>Ha</text></lyric>`},
							musicxml.Element{Content: `<pitch><step>C</step><octave>5</octave></pitch><duration>1</duration><voice>1</voice><type>eighth</type><stem>down</stem><beam number="1">begin</beam><notations><slur type="start"placement="above"number="1"/></notations><lyric number="1"><syllabic>middle</syllabic><text>le</text></lyric>`},
						},
						Barline: []musicxml.Barline{
							musicxml.Barline{
								Location: musicxml.BarlineLocationLeft,
								BarStyle: musicxml.BarLineStyleHeavyLight,
								Repeat: &musicxml.BarLineRepeat{
									Direction: musicxml.BarLineRepeatDirectionForward,
								},
							},
						},
					},
				},
			},
			canv: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)

				return canv
			},
			si: func(ctrl *gomock.Controller) *staffInteractor {

				barlineMock := barline.NewMockBarline(ctrl)
				lyricMock := lyric.NewMockLyric(ctrl)
				numberedMock := numbered.NewMockNumbered(ctrl)
				rhythmMock := rhythm.NewMockRhythm(ctrl)
				breathpauseMock := breathpause.NewMockBreathPause(ctrl)
				renderAlign := NewMockRenderStaffWithAlign(ctrl)
				timeSignature := timesig.TimeSignature{IsMixed: false, Signatures: []timesig.Time{timesig.Time{Measure: 1, Beat: 2, BeatType: 4}}}

				var emptyCoordinate *entity.Coordinate
				barlineMock.EXPECT().GetRendererLeftBarline(gomock.Any(), int(50), emptyCoordinate).Return(&entity.NoteRenderer{
					Barline: &musicxml.Barline{
						Location: musicxml.BarlineLocationLeft,
						BarStyle: musicxml.BarLineStyleHeavyLight,
						Repeat: &musicxml.BarLineRepeat{
							Direction: musicxml.BarLineRepeatDirectionForward,
						},
					},
				}, &barline.BarlineInfo{})
				numberedMock.EXPECT().GetLengthNote(gomock.Any(), timeSignature, int(1), float64(2)).Return([]numbered.NoteLength{numbered.NoteLength{Type: musicxml.NoteLengthQuarter}}).Times(2)
				numberedMock.EXPECT().GetLengthNote(gomock.Any(), timeSignature, int(1), float64(1)).Return([]numbered.NoteLength{numbered.NoteLength{Type: musicxml.NoteLengthQuarter}})
				numberedMock.EXPECT().GetLengthNote(gomock.Any(), timeSignature, int(1), float64(0.5)).Return([]numbered.NoteLength{numbered.NoteLength{Type: musicxml.NoteLengthEighth}})
				rhythmMock.EXPECT().SetRhythmNotation(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				lyricMock.EXPECT().SetLyricRenderer(gomock.Any(), gomock.Any()).Return(lyric.VerseInfo{MarginBottom: 25}).AnyTimes()
				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				rhythmMock.EXPECT().AdjustMultiDottedRenderer(gomock.Any(), int(50), int(80)).Return(int(25), int(80))
				barlineMock.EXPECT().GetRendererRightBarline(gomock.Any(), int(25)).Return(25, &entity.NoteRenderer{Barline: &musicxml.Barline{Location: musicxml.BarlineLocationRight, BarStyle: musicxml.BarLineStyleRegular}})
				renderAlign.EXPECT().RenderWithAlign(gomock.Any(), gomock.Any(), int(80), timeSignature, gomock.Any())

				si := &staffInteractor{
					Barline:     barlineMock,
					Lyric:       lyricMock,
					Numbered:    numberedMock,
					Rhythm:      rhythmMock,
					BreathPause: breathpauseMock,
					RenderAlign: renderAlign,
				}

				return si
			},
			wantStaffInfo: StaffInfo{
				MarginBottom:     25,
				NextLineRenderer: []*entity.NoteRenderer{},
			},
		},
		{
			name: "case #2",
			args: args{
				x: 50,
				y: 80,
				keySignature: keysig.NewKeySignature(musicxml.KeySignature{
					Fifth: 2, // D Major
				}),
				timeSignature: timesig.TimeSignature{IsMixed: false, Signatures: []timesig.Time{timesig.Time{Measure: 1, Beat: 2, BeatType: 4}}},
				measures: []musicxml.Measure{
					musicxml.Measure{
						Number: 1,
						Appendix: []musicxml.Element{
							musicxml.Element{Content: `<pitch><step>G</step><alter>1</alter><octave>4</octave></pitch><duration>3</duration><voice>1</voice><type>eighth</type><dot/><stem>up</stem><beam number="1">begin</beam><lyric number="1"><syllabic>end</syllabic><text>kau,</text></lyric>`},
							musicxml.Element{Content: `<direction-type><words>__layout=br</words></direction-type>`},
							musicxml.Element{Content: `<pitch><step>A</step><octave>4</octave></pitch><duration>1</duration><voice>1</voice><type>16th</type><stem>up</stem><beam number="1">end</beam><beam number="2">backward hook</beam><lyric number="1"><syllabic>single</syllabic><text>Yang</text></lyric>`},
							musicxml.Element{Content: `<pitch><step>A</step><octave>4</octave></pitch><duration>8</duration><voice>1</voice><type>whole</type><notations><articulations><breath-mark/></articulations></notations><lyric number="1"><syllabic>end</syllabic><text>duh.</text></lyric>`},
						},
						Barline: []musicxml.Barline{},
					},
				},
			},
			canv: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)

				return canv
			},
			si: func(ctrl *gomock.Controller) *staffInteractor {

				barlineMock := barline.NewMockBarline(ctrl)
				lyricMock := lyric.NewMockLyric(ctrl)
				numberedMock := numbered.NewMockNumbered(ctrl)
				rhythmMock := rhythm.NewMockRhythm(ctrl)
				breathpauseMock := breathpause.NewMockBreathPause(ctrl)
				renderAlign := NewMockRenderStaffWithAlign(ctrl)

				key := keysig.NewKeySignature(musicxml.KeySignature{
					Fifth: 2, // D Major
				})

				_ = key

				timeSignature := timesig.TimeSignature{IsMixed: false, Signatures: []timesig.Time{timesig.Time{Measure: 1, Beat: 2, BeatType: 4}}}

				_ = timeSignature

				si := &staffInteractor{
					Barline:     barlineMock,
					Lyric:       lyricMock,
					Numbered:    numberedMock,
					Rhythm:      rhythmMock,
					BreathPause: breathpauseMock,
					RenderAlign: renderAlign,
				}

				// _ = measure

				numberedMock.EXPECT().GetLengthNote(gomock.Any(), timeSignature, int(1), float64(0.75)).Return([]numbered.NoteLength{numbered.NoteLength{Type: musicxml.NoteLengthEighth}})
				numberedMock.EXPECT().GetLengthNote(gomock.Any(), timeSignature, int(1), float64(0.25)).Return([]numbered.NoteLength{numbered.NoteLength{Type: musicxml.NoteLength16th}})
				numberedMock.EXPECT().GetLengthNote(gomock.Any(), timeSignature, int(1), float64(4)).AnyTimes().Return([]numbered.NoteLength{
					numbered.NoteLength{Type: musicxml.NoteLengthQuarter},
					numbered.NoteLength{Type: musicxml.NoteLengthQuarter, IsDotted: true},
					numbered.NoteLength{Type: musicxml.NoteLengthQuarter, IsDotted: true},
					numbered.NoteLength{Type: musicxml.NoteLengthQuarter, IsDotted: true},
				})

				rhythmMock.EXPECT().SetRhythmNotation(gomock.Any(), gomock.Any(), int(4))
				rhythmMock.EXPECT().SetRhythmNotation(gomock.Any(), gomock.Any(), int(5))
				rhythmMock.EXPECT().SetRhythmNotation(gomock.Any(), gomock.Any(), int(5))

				lyricMock.EXPECT().SetLyricRenderer(gomock.Any(), gomock.Any())
				lyricMock.EXPECT().SetLyricRenderer(gomock.Any(), gomock.Any()).Return(lyric.VerseInfo{MarginBottom: 80})
				lyricMock.EXPECT().SetLyricRenderer(gomock.Any(), gomock.Any())

				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(gomock.Any(), gomock.Any()).Return(nil)
				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(gomock.Any(), gomock.Any()).Return(nil)
				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(gomock.Any(), gomock.Any()).Return(nil)

				rhythmMock.EXPECT().AdjustMultiDottedRenderer(gomock.Any(), int(50), int(80)).Return(int(50), int(80))

				barlineMock.EXPECT().GetRendererRightBarline(gomock.Any(), int(50)).Return(int(50), &entity.NoteRenderer{Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular}})

				// the last param still gomock.Any() cause there is dotted line that needs to be handled for the new line
				renderAlign.EXPECT().RenderWithAlign(gomock.Any(), gomock.Any(), int(80), gomock.Any(), gomock.Any())

				return si
			},
			wantStaffInfo: StaffInfo{
				Multiline:    true,
				MarginLeft:   65,
				MarginBottom: 0,
				NextLineRenderer: append([]*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX:     50,
						PositionY:     80,
						Note:          5,
						MeasureNumber: 1,
						Beam:          map[int]entity.Beam{},
						NoteLength:    musicxml.NoteLengthWhole,
					},
				}, append(slices.Repeat([]*entity.NoteRenderer{
					&entity.NoteRenderer{
						IsDotted:      true,
						PositionY:     80,
						Width:         15,
						Beam:          map[int]entity.Beam{},
						MeasureNumber: 1,
						NoteLength:    musicxml.NoteLengthQuarter,
					},
				}, 3), &entity.NoteRenderer{
					Barline: &musicxml.Barline{
						BarStyle: musicxml.BarLineStyleRegular,
					},
				})...),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStaffInfo := tt.si(ctrl).RenderStaff(context.Background(), tt.canv(ctrl), tt.args.x, tt.args.y, tt.args.keySignature, tt.args.timeSignature, tt.args.measures, tt.args.prevNotes...); !assert.Equal(t, tt.wantStaffInfo, gotStaffInfo) {
				t.Errorf("staffInteractor.RenderStaff() = %v, want %v", gotStaffInfo, tt.wantStaffInfo)
			}
		})
	}
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
