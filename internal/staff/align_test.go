package staff

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_renderStaffAlign_RenderWithAlign(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		y            int
		noteRenderer [][]*entity.NoteRenderer
	}
	tests := []struct {
		name       string
		canv       func(ctrl *gomock.Controller) *canvas.MockCanvas
		interactor func(ctrl *gomock.Controller) *renderStaffAlign
		args       args
	}{
		{
			name: "default",
			args: args{
				y: 195,
				noteRenderer: [][]*entity.NoteRenderer{
					[]*entity.NoteRenderer{
						&entity.NoteRenderer{
							MeasureNumber: 1,
							PositionX:     50,
							PositionY:     195,
							Note:          2,
							NoteLength:    musicxml.NoteLengthEighth,
							Width:         17,
							Beam: map[int]entity.Beam{
								1: entity.Beam{
									Number: 1,
									Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
								},
							},
							Lyric: []entity.Lyric{
								entity.Lyric{
									Text: []entity.Text{
										entity.Text{
											Value: "Jo",
										},
									},
									Syllabic: musicxml.LyricSyllabicTypeBegin,
								},
							},
							IsLengthTakenFromLyric: true,
						},
						&entity.NoteRenderer{
							PositionX:  67,
							PositionY:  195,
							Note:       3,
							NoteLength: musicxml.NoteLengthEighth,
							Width:      34,
							Beam: map[int]entity.Beam{
								1: entity.Beam{
									Number: 1,
									Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
								},
							},
							Lyric: []entity.Lyric{
								entity.Lyric{
									Text: []entity.Text{
										entity.Text{
											Value: "dy",
										},
									},
									Syllabic: musicxml.LyricSyllabicTypeEnd,
								},
							},
							MeasureNumber:          1,
							IndexPosition:          1,
							IsLengthTakenFromLyric: true,
						},
						&entity.NoteRenderer{
							PositionX:  101,
							PositionY:  195,
							Note:       4,
							NoteLength: musicxml.NoteLengthEighth,
							Width:      34,
							Lyric: []entity.Lyric{
								entity.Lyric{
									Text: []entity.Text{
										entity.Text{
											Value: "Lum",
										},
									},
									Syllabic: musicxml.LyricSyllabicTypeBegin,
								},
							},
							Beam: map[int]entity.Beam{
								1: entity.Beam{
									Number: 1,
									Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
								},
							},
							MeasureNumber:          1,
							IsLengthTakenFromLyric: true,
							IndexPosition:          2,
						},
						&entity.NoteRenderer{
							PositionX:  135,
							PositionY:  195,
							Note:       5,
							NoteLength: musicxml.NoteLengthEighth,
							Width:      29,
							Beam: map[int]entity.Beam{
								1: entity.Beam{
									Number: 1,
									Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
								},
							},
							Lyric: []entity.Lyric{
								entity.Lyric{
									Text: []entity.Text{
										entity.Text{
											Value: "ban",
										},
									},
									Syllabic: musicxml.LyricSyllabicTypeMiddle,
								},
							},
							IsLengthTakenFromLyric: true,
							MeasureNumber:          1,
							IndexPosition:          3,
						},
						&entity.NoteRenderer{
							PositionX:  110,
							PositionY:  195,
							Note:       3,
							NoteLength: musicxml.NoteLengthQuarter,
							Width:      17,
							Beam:       map[int]entity.Beam{},
							Lyric: []entity.Lyric{
								entity.Lyric{
									Text: []entity.Text{
										entity.Text{
											Value: "to",
										},
									},
									Syllabic: musicxml.LyricSyllabicTypeMiddle,
								},
							},
							IsLengthTakenFromLyric: true,
							MeasureNumber:          1,
							IndexPosition:          4,
						},
						&entity.NoteRenderer{
							PositionX:  181,
							PositionY:  195,
							Note:       1,
							NoteLength: musicxml.NoteLengthEighth,
							Width:      19,
							Beam: map[int]entity.Beam{
								1: entity.Beam{
									Number: 1,
									Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
								},
							},
							Lyric: []entity.Lyric{
								entity.Lyric{
									Text: []entity.Text{
										entity.Text{
											Value: "ru",
										},
									},
									Syllabic: musicxml.LyricSyllabicTypeMiddle,
								},
							},
							IsLengthTakenFromLyric: true,
							MeasureNumber:          1,
							IndexPosition:          5,
						},
						&entity.NoteRenderer{
							PositionX:  200,
							PositionY:  195,
							Note:       2,
							NoteLength: musicxml.NoteLengthEighth,
							Width:      35,
							Beam: map[int]entity.Beam{
								1: entity.Beam{
									Number: 1,
									Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
								},
							},
							Tie: &entity.Slur{
								Number: 2,
								Type:   musicxml.NoteSlurTypeStart,
							},
							Lyric: []entity.Lyric{
								entity.Lyric{
									Text: []entity.Text{
										entity.Text{
											Value: "an",
										},
									},
									Syllabic: musicxml.LyricSyllabicTypeEnd,
								},
							},
							MeasureNumber:          1,
							IsLengthTakenFromLyric: true,
							IndexPosition:          6,
						},
						&entity.NoteRenderer{
							PositionX: 235,
							PositionY: 195,
							Barline: &musicxml.Barline{
								BarStyle: musicxml.BarLineStyleRegular,
							},
							MeasureNumber: 1,
						},
						&entity.NoteRenderer{
							PositionX:  250,
							PositionY:  195,
							Note:       2,
							NoteLength: musicxml.NoteLengthWhole,
							Width:      15,
							Beam:       map[int]entity.Beam{},
							Tie: &entity.Slur{
								Number: 2,
								Type:   musicxml.NoteSlurTypeStop,
							},
							MeasureNumber: 2,
						},
						&entity.NoteRenderer{
							PositionX:     270,
							PositionY:     195,
							IsDotted:      true,
							NoteLength:    musicxml.NoteLengthQuarter,
							Width:         15,
							Beam:          map[int]entity.Beam{},
							MeasureNumber: 2,
							IndexPosition: 1,
						},
						&entity.NoteRenderer{
							PositionX:     290,
							PositionY:     195,
							IsDotted:      true,
							NoteLength:    musicxml.NoteLengthQuarter,
							Width:         15,
							Beam:          map[int]entity.Beam{},
							MeasureNumber: 2,
							IndexPosition: 2,
						},
						&entity.NoteRenderer{
							PositionX:     310,
							PositionY:     195,
							IsDotted:      true,
							NoteLength:    musicxml.NoteLengthQuarter,
							Width:         15,
							Beam:          map[int]entity.Beam{},
							MeasureNumber: 2,
							IndexPosition: 3,
						},
						&entity.NoteRenderer{
							PositionX:  325,
							PositionY:  195,
							NoteLength: musicxml.NoteLengthQuarter,
							Barline: &musicxml.Barline{
								Location: musicxml.BarlineLocationRight,
								BarStyle: musicxml.BarLineStyleLightHeavy,
							},
							Beam:          map[int]entity.Beam{},
							MeasureNumber: 2,
						},
					},
				},
			},
			canv: func(ctrl *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(ctrl)
				canv.EXPECT().Group("staff")
				canv.EXPECT().Group("measure-align")
				canv.EXPECT().Group("class='note'", "style='font-family:Old Standard TT;font-weight:500'")

				canv.EXPECT().Text(int(50), int(195), "2")
				canv.EXPECT().Text(int(50), int(220), "Jo")

				canv.EXPECT().Text(int(97), int(195), "3")
				canv.EXPECT().Text(int(97), int(220), "dy")

				canv.EXPECT().Text(int(162), int(195), "4")
				canv.EXPECT().Text(int(162), int(220), "Lum")

				canv.EXPECT().Text(int(227), int(195), "5")
				canv.EXPECT().Text(int(227), int(220), "ban")

				canv.EXPECT().Text(int(232), int(195), "3")
				canv.EXPECT().Text(int(232), int(220), "to")

				canv.EXPECT().Text(int(334), int(195), "1")
				canv.EXPECT().Text(int(334), int(220), "ru")

				canv.EXPECT().Text(int(384), int(195), "2")
				canv.EXPECT().Text(int(384), int(220), "an")
				canv.EXPECT().Text(int(495), int(195), "2")
				canv.EXPECT().Text(int(546), int(195), ".")
				canv.EXPECT().Text(int(597), int(195), ".")
				canv.EXPECT().Text(int(648), int(195), ".")

				canv.EXPECT().Group("class='lyric'", "style='font-family:Caladea'")
				// canv.EXPECT().Group("class='staff-text'")

				canv.EXPECT().Gend().Times(4)
				return canv
			},
			interactor: func(ctrl *gomock.Controller) *renderStaffAlign {
				barlineMock := barline.NewMockBarline(ctrl)
				numberedMock := numbered.NewMockNumbered(ctrl)
				rhythmMock := rhythm.NewMockRhythm(ctrl)
				lyricMock := lyric.NewMockLyric(ctrl)
				res := &renderStaffAlign{
					Barline:  barlineMock,
					Numbered: numberedMock,
					Rhythm:   rhythmMock,
					Lyric:    lyricMock,
				}

				barlineMock.EXPECT().RenderBarline(gomock.Any(), gomock.Any(), musicxml.Barline{
					BarStyle: musicxml.BarLineStyleRegular,
				}, entity.Coordinate{
					X: 450,
					Y: 195,
				})
				barlineMock.EXPECT().RenderBarline(gomock.Any(), gomock.Any(), musicxml.Barline{
					Location: musicxml.BarlineLocationRight,
					BarStyle: musicxml.BarLineStyleLightHeavy,
				}, entity.Coordinate{
					X: 670,
					Y: 195,
				})

				numberedMock.EXPECT().RenderOctave(gomock.Any(), gomock.Any(), gomock.Any())
				rhythmMock.EXPECT().RenderBeam(gomock.Any(), gomock.Any(), gomock.Any())
				rhythmMock.EXPECT().RenderSlurTies(gomock.Any(), gomock.Any(), IsEqual([]*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX:  384,
						PositionY:  195,
						Note:       2,
						NoteLength: musicxml.NoteLengthEighth,
						Width:      35,
						Beam: map[int]entity.Beam{
							1: entity.Beam{
								Number: 1,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
						Tie: &entity.Slur{
							Number: 2,
							Type:   musicxml.NoteSlurTypeStart,
						},
						Lyric: []entity.Lyric{
							entity.Lyric{
								Text: []entity.Text{
									entity.Text{
										Value: "an",
									},
								},
								Syllabic: musicxml.LyricSyllabicTypeEnd,
							},
						},
						MeasureNumber:          1,
						IsLengthTakenFromLyric: true,
						IndexPosition:          6,
					},
					&entity.NoteRenderer{
						PositionX:  495,
						PositionY:  195,
						Note:       2,
						NoteLength: musicxml.NoteLengthWhole,
						Width:      15,
						Beam:       map[int]entity.Beam{},
						Tie: &entity.Slur{
							Number: 2,
							Type:   musicxml.NoteSlurTypeStop,
						},
						MeasureNumber: 2,
					},
				}, t), float64(670))
				lyricMock.EXPECT().RenderHypen(gomock.Any(), gomock.Any(), gomock.Any())
				lyricMock.EXPECT().CalculateMarginLeft(gomock.Any()).Return(float64(0))

				return res

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			(tt.interactor(ctrl)).RenderWithAlign(context.Background(), tt.canv(ctrl), tt.args.y, tt.args.noteRenderer)
		})
	}
}
