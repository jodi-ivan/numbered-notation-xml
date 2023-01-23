package numbered

import (
	"context"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/stretchr/testify/assert"
)

func TestRenderLengthNote(t *testing.T) {
	type args struct {
		ts         timesig.TimeSignature
		measure    int
		noteLength float64
	}
	tests := []struct {
		name string
		args args
		want []renderedNote
	}{
		{
			name: "quarter .5 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 0.5,
			},
			want: []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthEighth,
				},
			},
		},
		{
			name: "quarter 1 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 1,
			},
			want: []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthQuarter,
				},
			},
		},
		{
			name: "quarter 1.5 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 1.5,
			},
			want: []renderedNote{ /// 1 ^.
				renderedNote{
					Type: musicxml.NoteLengthQuarter,
				},
				renderedNote{
					Type:     musicxml.NoteLengthEighth,
					IsDotted: true,
				},
			},
		},
		{
			name: "quarter 2 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 2,
			},
			want: []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthQuarter,
				},
				renderedNote{
					IsDotted: true,
				},
			},
		},
		{
			name: "quarter 2.5 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 2.5,
			},
			want: []renderedNote{ /// 1 ^.
				renderedNote{
					Type: musicxml.NoteLengthQuarter,
				},
				renderedNote{
					IsDotted: true,
				},
				renderedNote{
					Type:     musicxml.NoteLengthEighth,
					IsDotted: true,
				},
			},
		},
		{
			name: "quarter 3 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 3,
			},
			want: []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthQuarter,
				},
				renderedNote{
					IsDotted: true,
				},
				renderedNote{
					IsDotted: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderLengthNote(context.Background(), tt.args.ts, tt.args.measure, tt.args.noteLength); !assert.Equal(t, tt.want, got) {
				t.Errorf("RenderLengthNote() = %v, want %v", got, tt.want)
			}
		})
	}
}
