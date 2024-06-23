package staff

import (
	"context"
	"encoding/xml"
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

				key := keysig.NewKeySignature(musicxml.KeySignature{
					Fifth: 2, // D Major
				})

				_ = key

				timeSignature := timesig.TimeSignature{IsMixed: false, Signatures: []timesig.Time{timesig.Time{Measure: 1, Beat: 2, BeatType: 4}}}

				renderer := [3]*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX:     50,
						PositionY:     80,
						Note:          5,
						NoteLength:    musicxml.NoteLengthQuarter,
						Beam:          map[int]entity.Beam{},
						MeasureNumber: 1,
					},
					&entity.NoteRenderer{
						PositionX:  50,
						PositionY:  80,
						Note:       7,
						NoteLength: musicxml.NoteLengthEighth,
						Beam: map[int]entity.Beam{
							1: entity.Beam{
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
								Number: 1,
							},
						},
						MeasureNumber: 1,
					},
					&entity.NoteRenderer{
						Articulation: &entity.Articulation{
							BreathMark: &entity.ArticulationTypesBreathMark,
						},
						MeasureNumber: 1,
						Width:         6,
					},
				}

				note := [2]musicxml.Note{
					musicxml.Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "A", Octave: 4},
						Type: musicxml.NoteLengthQuarter,
						Lyric: []musicxml.Lyric{

							musicxml.Lyric{
								Number: 1,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{{Value: "Ha"}},
								Syllabic: musicxml.LyricSyllabicTypeBegin,
							},
						},
					},
					musicxml.Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "C", Octave: 5},
						Type: musicxml.NoteLengthEighth,
						Beam: []*musicxml.NoteBeam{
							&musicxml.NoteBeam{
								Number: 1,
								State:  musicxml.NoteBeamTypeBegin,
							},
						},
						Lyric: []musicxml.Lyric{

							musicxml.Lyric{
								Number: 1,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{{Value: "le"}},
								Syllabic: musicxml.LyricSyllabicTypeMiddle,
							},
						},
						Notations: &musicxml.NoteNotation{
							Slur: []musicxml.NotationSlur{
								musicxml.NotationSlur{
									Type:   musicxml.NoteSlurTypeStart,
									Number: 1,
								},
							},
						},
					},
				}

				measures := musicxml.Measure{
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
					Notes: []musicxml.Note{
						musicxml.Note{
							Rest: &musicxml.Rest{},
							Type: musicxml.NoteLengthHalf,
							MeasureText: []musicxml.MeasureText{
								musicxml.MeasureText{
									Text: "Refrein",
								},
							},
						},
						musicxml.Note{
							Rest: &musicxml.Rest{},
							Type: musicxml.NoteLengthHalf,
						},
						note[0],
						note[1],
					},
					NewLineIndex: -1,
				}
				var emptyCoordinate *entity.Coordinate
				barlineMock.EXPECT().GetRendererLeftBarline(IsEqual(measures, t), int(50), emptyCoordinate).Return(&entity.NoteRenderer{
					Barline: &musicxml.Barline{
						Location: musicxml.BarlineLocationLeft,
						BarStyle: musicxml.BarLineStyleHeavyLight,
						Repeat: &musicxml.BarLineRepeat{
							Direction: musicxml.BarLineRepeatDirectionForward,
						},
					},
				}, &barline.BarlineInfo{})
				numberedMock.EXPECT().GetLengthNote(gomock.Any(), timeSignature, int(1), float64(1)).Return([]numbered.NoteLength{numbered.NoteLength{Type: musicxml.NoteLengthQuarter}})
				numberedMock.EXPECT().GetLengthNote(gomock.Any(), timeSignature, int(1), float64(0.5)).Return([]numbered.NoteLength{numbered.NoteLength{Type: musicxml.NoteLengthEighth}})

				rhythmMock.EXPECT().SetRhythmNotation(IsEqual(renderer[0], t), IsEqual(note[0], t), 5)
				rhythmMock.EXPECT().SetRhythmNotation(IsEqual(renderer[1], t), IsEqual(note[1], t), 7)

				lyricMock.EXPECT().SetLyricRenderer(IsEqual(renderer[0], t), IsEqual(note[0], t)).Return(lyric.VerseInfo{MarginBottom: 25})
				lyricMock.EXPECT().SetLyricRenderer(IsEqual(renderer[1], t), IsEqual(note[1], t)).Return(lyric.VerseInfo{MarginBottom: 25})

				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(IsEqual(renderer[0], t), IsEqual(note[0], t))
				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(IsEqual(renderer[1], t), IsEqual(note[1], t)).Return(&entity.NoteRenderer{
					Articulation: &entity.Articulation{
						BreathMark: &entity.ArticulationTypesBreathMark,
					},
					MeasureNumber: 1,
					Width:         6,
				})

				rhythmMock.EXPECT().AdjustMultiDottedRenderer([]*entity.NoteRenderer{renderer[0], renderer[1], renderer[2]}, int(50), int(80)).Return(int(25), int(80))

				barlineMock.EXPECT().GetRendererRightBarline(IsEqual(measures, t), int(25)).Return(25, &entity.NoteRenderer{Barline: &musicxml.Barline{Location: musicxml.BarlineLocationRight, BarStyle: musicxml.BarLineStyleRegular}})

				renderAlign.EXPECT().RenderWithAlign(gomock.Any(), gomock.Any(), int(80), IsEqual([][]*entity.NoteRenderer{
					[]*entity.NoteRenderer{
						&entity.NoteRenderer{
							Barline: &musicxml.Barline{
								Location: musicxml.BarlineLocationLeft,
								BarStyle: musicxml.BarLineStyleHeavyLight,
								Repeat: &musicxml.BarLineRepeat{
									Direction: musicxml.BarLineRepeatDirectionForward,
								},
							},
						},
						renderer[0],
						renderer[1],
						renderer[2],
						&entity.NoteRenderer{Barline: &musicxml.Barline{Location: musicxml.BarlineLocationRight, BarStyle: musicxml.BarLineStyleRegular}},
					},
				}, t))

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
				notes := []musicxml.Note{
					musicxml.Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{
							Step:   "G",
							Octave: 4,
						},
						Type: musicxml.NoteLengthEighth,
						Beam: []*musicxml.NoteBeam{
							&musicxml.NoteBeam{
								Number: 1,
								State:  musicxml.NoteBeamTypeBegin,
							},
						},
						Lyric: []musicxml.Lyric{
							musicxml.Lyric{
								Number: 1,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "kau,",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
						Dot: []*musicxml.Dot{&musicxml.Dot{}},
					},
					musicxml.Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{
							Step:   "A",
							Octave: 4,
						},
						Type: musicxml.NoteLength16th,
						Beam: []*musicxml.NoteBeam{
							&musicxml.NoteBeam{
								Number: 1,
								State:  musicxml.NoteBeamTypeEnd,
							},
							&musicxml.NoteBeam{
								Number: 2,
								State:  musicxml.NoteBeamTypeBackwardHook,
							},
						},
						Lyric: []musicxml.Lyric{
							musicxml.Lyric{
								Number: 1,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "Yang",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeSingle,
							},
						},
					},
					musicxml.Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{
							Step:   "A",
							Octave: 4,
						},
						Type: musicxml.NoteLengthWhole,
						Notations: &musicxml.NoteNotation{
							Articulation: &musicxml.NotationArticulation{
								BreathMark: &struct {
									Name xml.Name
								}{
									Name: xml.Name{},
								},
							},
						},
						Lyric: []musicxml.Lyric{
							musicxml.Lyric{
								Number: 1,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "duh.",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
					},
				}
				renderer := []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX:  50,
						PositionY:  80,
						Note:       4,
						NoteLength: musicxml.NoteLengthEighth,
						Beam: map[int]entity.Beam{
							1: entity.Beam{
								Number: 1,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
						MeasureNumber: 1,
					},
					&entity.NoteRenderer{
						PositionX:  50,
						PositionY:  80,
						Note:       5,
						NoteLength: musicxml.NoteLength16th,
						Beam: map[int]entity.Beam{
							1: entity.Beam{
								Number: 1,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
							2: entity.Beam{
								Number: 2,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
						MeasureNumber: 1,
						IsNewLine:     true,
					},
					&entity.NoteRenderer{
						PositionX:     50,
						PositionY:     80,
						Note:          5,
						NoteLength:    musicxml.NoteLengthWhole,
						Beam:          map[int]entity.Beam{},
						MeasureNumber: 1,
					},
				}

				measure := musicxml.Measure{
					Number: 1,
					Appendix: []musicxml.Element{
						musicxml.Element{Content: `<pitch><step>G</step><alter>1</alter><octave>4</octave></pitch><duration>3</duration><voice>1</voice><type>eighth</type><dot/><stem>up</stem><beam number="1">begin</beam><lyric number="1"><syllabic>end</syllabic><text>kau,</text></lyric>`},
						musicxml.Element{Content: `<direction-type><words>__layout=br</words></direction-type>`},
						musicxml.Element{Content: `<pitch><step>A</step><octave>4</octave></pitch><duration>1</duration><voice>1</voice><type>16th</type><stem>up</stem><beam number="1">end</beam><beam number="2">backward hook</beam><lyric number="1"><syllabic>single</syllabic><text>Yang</text></lyric>`},
						musicxml.Element{Content: `<pitch><step>A</step><octave>4</octave></pitch><duration>8</duration><voice>1</voice><type>whole</type><notations><articulations><breath-mark/></articulations></notations><lyric number="1"><syllabic>end</syllabic><text>duh.</text></lyric>`},
					},
					Notes:        notes,
					NewLineIndex: 1,
					Barline:      []musicxml.Barline{},
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

				rhythmMock.EXPECT().SetRhythmNotation(IsEqual(renderer[0], t), IsEqual(notes[0], t), int(4))
				rhythmMock.EXPECT().SetRhythmNotation(IsEqual(renderer[1], t), IsEqual(notes[1], t), int(5))
				rhythmMock.EXPECT().SetRhythmNotation(IsEqual(renderer[2], t), IsEqual(notes[2], t), int(5))

				lyricMock.EXPECT().SetLyricRenderer(renderer[0], notes[0])
				lyricMock.EXPECT().SetLyricRenderer(renderer[1], notes[1]).Return(lyric.VerseInfo{MarginBottom: 80})
				lyricMock.EXPECT().SetLyricRenderer(renderer[2], notes[2])

				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(renderer[0], notes[0]).Return(nil)
				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(renderer[1], notes[1]).Return(nil)
				breathpauseMock.EXPECT().SetAndGetBreathPauseRenderer(renderer[2], notes[2]).Return(nil)

				rendererDot := &entity.NoteRenderer{
					PositionY:     80,
					Width:         15,
					IsDotted:      true,
					NoteLength:    musicxml.NoteLengthQuarter,
					Beam:          map[int]entity.Beam{},
					MeasureNumber: 1,
				}
				renderer = append(renderer, rendererDot, rendererDot, rendererDot)
				rhythmMock.EXPECT().AdjustMultiDottedRenderer(IsEqual(renderer, t), int(50), int(80)).Return(int(50), int(80))

				barlineMock.EXPECT().GetRendererRightBarline(IsEqual(measure, t), int(50)).Return(int(50), &entity.NoteRenderer{Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular}})

				// the last param still gomock.Any() cause there is dotted line that needs to be handled for the new line
				renderAlign.EXPECT().RenderWithAlign(gomock.Any(), gomock.Any(), int(80), gomock.Any())

				return si
			},
			wantStaffInfo: StaffInfo{
				Multiline:    true,
				MarginLeft:   65,
				MarginBottom: 0,
				NextLineRenderer: []*entity.NoteRenderer{
					&entity.NoteRenderer{Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular}},
				},
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
