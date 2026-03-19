package timesig

import (
	"context"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func TestTime_calculateNoteLength(t *testing.T) {

	// quarter                    =       1 beat
	//      quarter with .        =   1 1/2 beat
	//      quarter with . .      =   1 3/4 beat
	// half                       =       2 beat
	//      half with .           =       3 beat
	//      half with . .         =   3 1/2 beat
	//      half with . . .       =   3 3/4 beat
	// whole                      =       4 beat
	//      whole with .          =       6 beat
	//      whole with . .        =       7 beat
	//      whole with . . .      =   7 1/2 beat
	//      whole with . . . .    =   7 3/4 beat
	// eighth                     =     1/2 beat
	//      eight with .          =     3/4 beat
	// 16th                       =     1/4 beat
	type fields struct {
		BeatType int
	}

	type args struct {
		measure int
		note    musicxml.Note
	}
	tests := []struct {
		name  string
		args  args
		field fields
		want  float64
	}{
		{
			name: "quarter",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthQuarter,
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 1,
		},
		{
			name: "quarter with 1 dot -  1 1/2 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthQuarter,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 1.5,
		},
		{
			name: "quarter with 2 dot -  1 3/4 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthQuarter,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 1.75,
		},
		{
			name: "half",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthHalf,
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 2,
		},
		{
			name: "half with 1 dot -  3 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthHalf,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 3,
		},
		{
			name: "half with 2 dot -  3 1/2 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthHalf,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 3.5,
		},
		{
			name: "half with 3 dot -  3 3/4 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthHalf,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
						&musicxml.Dot{},
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 3.75,
		},
		{
			name: "half",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthWhole,
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 4,
		},
		{
			name: "half with 1 dot - 6 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthWhole,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 6,
		},
		{
			name: "half with 2 dot - 7 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthWhole,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 7,
		},
		{
			name: "half with 3 dot - 7.5 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthWhole,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
						&musicxml.Dot{},
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 7.5,
		},
		{
			name: "half with 4 dot - 7.75 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthWhole,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
						&musicxml.Dot{},
						&musicxml.Dot{},
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 7.75,
		},
		{
			name: "eighth",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthEighth,
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 0.5,
		},
		{
			name: "eighth with 1 dot - 0.75 beat",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLengthEighth,
					Dot: []*musicxml.Dot{
						&musicxml.Dot{},
					},
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 0.75,
		},
		{
			name: "16th",
			args: args{
				note: musicxml.Note{
					Type: musicxml.NoteLength16th,
				},
			},
			field: fields{
				BeatType: 4,
			},
			want: 0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Time{
				BeatType: tt.field.BeatType,
			}
			if got := tr.calculateNoteLength(context.Background(), tt.args.note); got != tt.want {
				t.Errorf("Time.calculateNoteLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeSignature_GetNoteLength(t *testing.T) {
	type fields struct {
		IsMixed    bool
		Signatures []Time
	}
	type args struct {
		measure int
		note    musicxml.Note
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "mixed timesignatures, targeted measure at the front",
			fields: fields{
				IsMixed: true,
				Signatures: []Time{
					Time{
						Measure:  1,
						Beat:     4,
						BeatType: 4,
					},
					Time{
						Measure:  10,
						Beat:     6,
						BeatType: 8,
					},
				},
			},
			args: args{
				measure: 1,
				note: musicxml.Note{
					Type: musicxml.NoteLengthQuarter,
				},
			},
			want: 1,
		},
		{
			name: "mixed timesignatures, targeted measure at the front edge",
			fields: fields{
				IsMixed: true,
				Signatures: []Time{
					Time{
						Measure:  1,
						Beat:     4,
						BeatType: 4,
					},
					Time{
						Measure:  10,
						Beat:     6,
						BeatType: 8,
					},
				},
			},
			args: args{
				measure: 9,
				note: musicxml.Note{
					Type: musicxml.NoteLengthQuarter,
				},
			},
			want: 1,
		},
		{
			name: "mixed timesignatures, targeted measure at the back start",
			fields: fields{
				IsMixed: true,
				Signatures: []Time{
					Time{
						Measure:  1,
						Beat:     4,
						BeatType: 4,
					},
					Time{
						Measure:  10,
						Beat:     6,
						BeatType: 8,
					},
					Time{
						Measure:  18,
						Beat:     4,
						BeatType: 4,
					},
				},
			},
			args: args{
				measure: 10,
				note: musicxml.Note{
					Type: musicxml.NoteLengthQuarter,
				},
			},
			want: 2,
		},
		{
			name: "mixed timesignatures, targeted measure at the back rear start",
			fields: fields{
				IsMixed: true,
				Signatures: []Time{
					Time{
						Measure:  1,
						Beat:     4,
						BeatType: 4,
					},
					Time{
						Measure:  10,
						Beat:     6,
						BeatType: 8,
					},
					Time{
						Measure:  18,
						Beat:     4,
						BeatType: 4,
					},
				},
			},
			args: args{
				measure: 18,
				note: musicxml.Note{
					Type: musicxml.NoteLengthQuarter,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TimeSignature{
				IsMixed:    tt.fields.IsMixed,
				Signatures: tt.fields.Signatures,
			}
			if got := ts.GetNoteLength(context.Background(), tt.args.measure, tt.args.note); got != tt.want {
				t.Errorf("TimeSignature.GetNoteLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTimeSignatures(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		measures []musicxml.Measure
		want     TimeSignature
	}{
		{
			name: "single time signature",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
			},
			want: TimeSignature{
				IsMixed: false,
				Signatures: []Time{
					Time{
						Measure:  1,
						Beat:     4,
						BeatType: 4,
					},
				},
			},
		},
		{
			name: "multi signature",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    1,
							BeatType: 1,
						},
					},
				},
				musicxml.Measure{
					Number: 2,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
			},
			want: TimeSignature{
				IsMixed: true,
				Signatures: []Time{
					Time{
						Measure:  1,
						Beat:     1,
						BeatType: 4,
					},
					Time{
						Measure:  2,
						Beat:     4,
						BeatType: 4,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTimeSignatures(context.Background(), tt.measures)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeSignature_GetTimesignatureOnMeasure(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		measures []musicxml.Measure
		// Named input parameters for target function.
		measure int
		want    Time
	}{
		{
			name: "single",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
			},
			measure: 5,
			want: Time{
				Measure:  1,
				Beat:     4,
				BeatType: 4,
			},
		},
		{
			name: "multi#1",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
				musicxml.Measure{
					Number: 2,
				},
				musicxml.Measure{
					Number: 3,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    6,
							BeatType: 8,
						},
					},
				},
			},
			measure: 1,
			want: Time{
				Measure:  1,
				Beat:     4,
				BeatType: 4,
			},
		},
		{
			name: "multi#2",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
				musicxml.Measure{
					Number: 2,
				},
				musicxml.Measure{
					Number: 3,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    6,
							BeatType: 8,
						},
					},
				},
			},
			measure: 3,
			want: Time{
				Measure:  3,
				Beat:     6,
				BeatType: 8,
			},
		},
		{
			name: "multi#3",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
				musicxml.Measure{
					Number: 2,
				},
				musicxml.Measure{
					Number: 3,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    6,
							BeatType: 8,
						},
					},
				},
			},
			measure: 5,
			want: Time{
				Measure:  3,
				Beat:     6,
				BeatType: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTimeSignatures(context.Background(), tt.measures)
			got := ts.GetTimesignatureOnMeasure(context.Background(), tt.measure)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeSignature_GetHumanized(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		measures []musicxml.Measure
		want     string
	}{
		{
			name: "single",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
			},
			want: "4 ketuk",
		},
		{
			name: "single 8 beat",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    6,
							BeatType: 8,
						},
					},
				},
			},
			want: "6 (2 x 3) ketuk",
		},
		{
			name: "multi",
			measures: []musicxml.Measure{
				musicxml.Measure{
					Number: 1,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
				musicxml.Measure{
					Number: 2,
				},
				musicxml.Measure{
					Number: 3,
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int `xml:"beats"`
							BeatType int `xml:"beat-type"`
						}{
							Beats:    6,
							BeatType: 8,
						},
					},
				},
			},
			want: "4 dan 6 (2 x 3) ketuk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTimeSignatures(context.Background(), tt.measures)
			got := ts.GetHumanized()
			assert.Equal(t, tt.want, got)
		})
	}
}
