package renderer

import (
	"context"
	"testing"

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
			if got := SplitLines(context.Background(), tt.args.part); !assert.Equal(t, tt.want, got) {
				t.Errorf("SplitLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
