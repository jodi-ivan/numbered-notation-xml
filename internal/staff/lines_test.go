package staff

import (
	"context"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func TestSplitLines(t *testing.T) {
	type args struct {
		part musicxml.Part
	}
	tests := []struct {
		name string
		args args
		want [][]musicxml.Measure
	}{
		{
			name: "everything went fine",
			args: args{
				part: musicxml.Part{
					Measures: []musicxml.Measure{
						musicxml.Measure{
							Number: 1,
						},
						musicxml.Measure{
							Number: 2,
						},
						musicxml.Measure{
							Number: 3,
						},
						musicxml.Measure{
							Number: 4,
							Print: &musicxml.Print{
								NewSystem: musicxml.PrintNewSystemTypeYes,
							},
						},
						musicxml.Measure{
							Number: 5,
						},
						musicxml.Measure{
							Number: 6,
						},
						musicxml.Measure{
							Number: 7,
						},
						musicxml.Measure{
							Number: 8,
							Print: &musicxml.Print{
								NewSystem: musicxml.PrintNewSystemTypeYes,
							},
						},
						musicxml.Measure{
							Number: 9,
						},
						musicxml.Measure{
							Number: 10,
						},
					},
				},
			},
			want: [][]musicxml.Measure{
				[]musicxml.Measure{
					musicxml.Measure{
						Number: 1,
					},
					musicxml.Measure{
						Number: 2,
					},
					musicxml.Measure{
						Number: 3,
					},
				},
				[]musicxml.Measure{
					musicxml.Measure{
						Number: 4,
						Print: &musicxml.Print{
							NewSystem: musicxml.PrintNewSystemTypeYes,
						},
					},
					musicxml.Measure{
						Number: 5,
					},
					musicxml.Measure{
						Number: 6,
					},
					musicxml.Measure{
						Number: 7,
					},
				},
				[]musicxml.Measure{
					musicxml.Measure{
						Number: 8,
						Print: &musicxml.Print{
							NewSystem: musicxml.PrintNewSystemTypeYes,
						},
					},
					musicxml.Measure{
						Number: 9,
					},
					musicxml.Measure{
						Number: 10,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := staffInteractor{}
			if got := si.SplitLines(context.Background(), tt.args.part); !assert.Equal(t, tt.want, got) {
				t.Errorf("SplitLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessPreviousLines(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		prevNotes      []*entity.NoteRenderer
		yPos           int
		wantTotalNotes int
		wantStaffInfo  StaffInfo
	}{
		{
			name: "no new line",
			prevNotes: []*entity.NoteRenderer{
				&entity.NoteRenderer{},
				&entity.NoteRenderer{},
				&entity.NoteRenderer{},
			},
			wantTotalNotes: 3,
			wantStaffInfo: StaffInfo{
				MarginLeft: constant.LAYOUT_INDENT_LENGTH,
			},
		},
		{
			name: "has line",
			prevNotes: []*entity.NoteRenderer{
				&entity.NoteRenderer{},
				&entity.NoteRenderer{
					IsNewLine: true,
				},
				&entity.NoteRenderer{},
			},
			wantTotalNotes: 3,
			wantStaffInfo: StaffInfo{
				Multiline: true,
				NextLineRenderer: []*entity.NoteRenderer{
					&entity.NoteRenderer{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := ProcessPreviousLines(tt.prevNotes, tt.yPos)
			if len(got) == tt.wantTotalNotes {
				t.Errorf("ProcessPreviousLines() total notes = %v, want %v", got, tt.wantTotalNotes)
			}
			assert.Equal(t, tt.wantStaffInfo.MarginBottom, got2.MarginBottom, "Margin bottom")
			assert.Equal(t, tt.wantStaffInfo.MarginLeft, got2.MarginLeft, "Margin left")
			assert.Equal(t, len(tt.wantStaffInfo.NextLineRenderer), len(got2.NextLineRenderer), "staff info next line")

		})
	}
}
