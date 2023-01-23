package renderer

import (
	"context"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func Test_cleanAdditionalBeams(t *testing.T) {
	type args struct {
		notes []*NoteRenderer
	}
	tests := []struct {
		name string
		args args
		want []*NoteRenderer
	}{
		{
			name: "no additional beams",
			args: args{
				notes: []*NoteRenderer{
					&NoteRenderer{
						Note: 4,
						Beam: map[int]Beam{
							1: Beam{
								Number: 1,
								Type:   musicxml.NoteBeamTypeBegin,
							},
						},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							1: Beam{
								Number: 1,
								Type:   musicxml.NoteBeamTypeEnd,
							},
						},
					},
				},
			},
			want: []*NoteRenderer{
				&NoteRenderer{
					Note: 4,
					Beam: map[int]Beam{
						1: Beam{
							Number: 1,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						1: Beam{
							Number: 1,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
			},
		},
		{
			name: "no 2 additional beams",
			args: args{
				notes: []*NoteRenderer{
					&NoteRenderer{
						Note: 4,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeamTypeBegin,
							},
						},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeamTypeEnd,
							},
						},
					},
					&NoteRenderer{
						Note: 2,
						Beam: map[int]Beam{},
					},
				},
			},
			want: []*NoteRenderer{
				&NoteRenderer{
					Note: 4,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
				&NoteRenderer{
					Note: 2,
					Beam: map[int]Beam{},
				},
			},
		},
		{
			name: "last notes additional was supposed to be assigned to end",
			args: args{
				notes: []*NoteRenderer{
					&NoteRenderer{
						Note: 4,
						Beam: map[int]Beam{},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 2,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
				},
			},
			want: []*NoteRenderer{
				&NoteRenderer{
					Note: 4,
					Beam: map[int]Beam{},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 2,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
			},
		},
		{
			name: "ended and start again",
			args: args{
				notes: []*NoteRenderer{
					&NoteRenderer{
						Note: 4,
						Beam: map[int]Beam{},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 2,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 1,
						Beam: map[int]Beam{},
					},
					&NoteRenderer{
						Note: 2,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
				},
			},
			want: []*NoteRenderer{
				&NoteRenderer{
					Note: 4,
					Beam: map[int]Beam{},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 2,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
				&NoteRenderer{
					Note: 1,
					Beam: map[int]Beam{},
				},
				&NoteRenderer{
					Note: 2,
					Beam: map[int]Beam{
						0: Beam{
							Number: 0,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						0: Beam{
							Number: 0,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
			},
		},
		{
			name: "ended and start again with separator has different beam ",
			args: args{
				notes: []*NoteRenderer{
					&NoteRenderer{
						Note: 4,
						Beam: map[int]Beam{},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 7,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 1,
						Beam: map[int]Beam{
							1: Beam{
								Number: 1,
								Type:   musicxml.NoteBeamTypeBegin,
							},
						},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							1: Beam{
								Number: 1,
								Type:   musicxml.NoteBeamTypeEnd,
							},
						},
					},
					&NoteRenderer{
						Note: 2,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 6,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
				},
			},
			want: []*NoteRenderer{
				&NoteRenderer{
					Note: 4,
					Beam: map[int]Beam{},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 7,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
				&NoteRenderer{
					Note: 1,
					Beam: map[int]Beam{
						1: Beam{
							Number: 1,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						1: Beam{
							Number: 1,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
				&NoteRenderer{
					Note: 2,
					Beam: map[int]Beam{
						0: Beam{
							Number: 0,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 6,
					Beam: map[int]Beam{
						0: Beam{
							Number: 0,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
			},
		},
		{
			name: "ended and start again with separator has different beam and multilayered",
			args: args{
				notes: []*NoteRenderer{
					&NoteRenderer{
						Note: 4,
						Beam: map[int]Beam{},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 7,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
						},
					},
					&NoteRenderer{
						Note: 1,
						Beam: map[int]Beam{
							1: Beam{
								Number: 1,
								Type:   musicxml.NoteBeamTypeBegin,
							},
						},
					},
					&NoteRenderer{
						Note: 3,
						Beam: map[int]Beam{
							1: Beam{
								Number: 1,
								Type:   musicxml.NoteBeamTypeEnd,
							},
							2: Beam{
								Number: 2,
								Type:   musicxml.NoteBeamTypeBegin,
							},
						},
					},
					&NoteRenderer{
						Note: 2,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
							2: Beam{
								Number: 2,
								Type:   musicxml.NoteBeamTypeContinue,
							},
						},
					},
					&NoteRenderer{
						Note: 6,
						Beam: map[int]Beam{
							0: Beam{
								Number: 0,
								Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
							},
							2: Beam{
								Number: 2,
								Type:   musicxml.NoteBeamTypeEnd,
							},
						},
					},
				},
			},
			want: []*NoteRenderer{
				&NoteRenderer{
					Note: 4,
					Beam: map[int]Beam{},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 7,
					Beam: map[int]Beam{
						INDEX_BEAM_ADDITIONAL: Beam{
							Number: INDEX_BEAM_ADDITIONAL,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
				&NoteRenderer{
					Note: 1,
					Beam: map[int]Beam{
						1: Beam{
							Number: 1,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 3,
					Beam: map[int]Beam{
						1: Beam{
							Number: 1,
							Type:   musicxml.NoteBeamTypeEnd,
						},
						2: Beam{
							Number: 2,
							Type:   musicxml.NoteBeamTypeBegin,
						},
					},
				},
				&NoteRenderer{
					Note: 2,
					Beam: map[int]Beam{
						0: Beam{
							Number: 0,
							Type:   musicxml.NoteBeamTypeBegin,
						},
						2: Beam{
							Number: 2,
							Type:   musicxml.NoteBeamTypeContinue,
						},
					},
				},
				&NoteRenderer{
					Note: 6,
					Beam: map[int]Beam{
						0: Beam{
							Number: 0,
							Type:   musicxml.NoteBeamTypeEnd,
						},
						2: Beam{
							Number: 2,
							Type:   musicxml.NoteBeamTypeEnd,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanAdditionalBeams(context.Background(), tt.args.notes)

			assert.Equal(t, tt.want, got)

		})
	}
}
