package moveabledo

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func TestConvertPitchToNumbered(t *testing.T) {
	type args struct {
		ks    keysig.KeySignature
		pitch string
	}
	tests := []struct {
		name         string
		args         args
		wantNumbered int
		wantStrike   bool
	}{
		{
			name: "C major input C",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = c",
				},
				pitch: "C",
			},
			wantNumbered: 1,
			wantStrike:   false,
		},
		{
			name: "C major input C#",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = c",
				},
				pitch: "C#",
			},
			wantNumbered: 1,
			wantStrike:   true,
		},
		{
			name: "C major input E",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = c",
				},
				pitch: "E",
			},
			wantNumbered: 3,
			wantStrike:   false,
		},
		{
			name: "C major input F",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = c",
				},
				pitch: "F",
			},
			wantNumbered: 4,
			wantStrike:   false,
		},
		{
			name: "C major input F#",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = c",
				},
				pitch: "F#",
			},
			wantNumbered: 4,
			wantStrike:   true,
		},
		{
			name: "C major input B",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = c",
				},
				pitch: "B",
			},
			wantNumbered: 7,
			wantStrike:   false,
		},
		{
			name: "D major input D",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     2,
					Key:       "d",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = d",
				},
				pitch: "D",
			},
			wantNumbered: 1,
			wantStrike:   false,
		},
		{
			name: "D major input B",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     2,
					Key:       "d",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = d",
				},
				pitch: "B",
			},
			wantNumbered: 6,
			wantStrike:   false,
		},
		{
			name: "D major input F#",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     2,
					Key:       "d",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = d",
				},
				pitch: "F#",
			},
			wantNumbered: 3,
			wantStrike:   false,
		},
		{
			name: "D major input C#",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     2,
					Key:       "d",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = d",
				},
				pitch: "C#",
			},
			wantNumbered: 7,
			wantStrike:   false,
		},
		{
			name: "A minor input A",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = a",
				},
				pitch: "A",
			},
			wantNumbered: 1,
			wantStrike:   false,
		},
		{
			name: "A minor input B",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = a",
				},
				pitch: "B",
			},
			wantNumbered: 2,
			wantStrike:   false,
		},
		{
			name: "A minor input C",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = a",
				},
				pitch: "C",
			},
			wantNumbered: 3,
			wantStrike:   false,
		},
		{
			name: "A minor input C#",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = a",
				},
				pitch: "C#",
			},
			wantNumbered: 3,
			wantStrike:   true,
		},
		{
			name: "A minor input G",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = a",
				},
				pitch: "G",
			},
			wantNumbered: 7,
			wantStrike:   false,
		},

		/*
			 	B minor scale
				fifth 2

				                    C#                     F#
				  |   |     |     |   |   |   |     |    |   |   |   |   |   |     |
				  |   |     |     |   |   |   |     |    |   |   |   |   |   |     |
				  |   |     |     |   |   |   |     |    |   |   |   |   |   |     |
				  |   |     |     | 2 |   |   |     |    | 5 |   |   |   |   |     |
				  +---+     |     +---+   +---+     |    +---+   +---+   +---+     |
				    |       |       |       |       |      |       |       |       |
				    |       |       |       |       |      |       |       |   *   |
				    |   1   |       |   3   |   4   |      |   6   |   7   |   1   |
				    |       |       |       |       |      |       |       |       |
				----+-------+-------+-------+-------+------+-------+-------+-------+
						B                D       E              G       A       B
		*/
		{
			name: "B minor input B",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     2,
					Key:       "b",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = b",
				},
				pitch: "B",
			},
			wantNumbered: 1,
			wantStrike:   false,
		},
		{
			name: "B minor input C#",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     2,
					Key:       "b",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = b",
				},
				pitch: "C#",
			},
			wantNumbered: 2,
			wantStrike:   false,
		},
		{
			name: "B minor input C",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     2,
					Key:       "b",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = b",
				},
				pitch: "C",
			},
			wantNumbered: 1,
			wantStrike:   true,
		},
		{
			name: "B minor input F#",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     2,
					Key:       "b",
					Mode:      keysig.NewMode("minor"),
					Humanized: "la = b",
				},
				pitch: "F#",
			},
			wantNumbered: 5,
			wantStrike:   false,
		},
		{
			name: "F major input Bb",
			args: args{
				ks: keysig.KeySignature{
					Fifth:     -1,
					Key:       "f",
					Mode:      keysig.NewMode("major"),
					Humanized: "do = f",
				},
				pitch: "Bb",
			},
			wantNumbered: 4,
			wantStrike:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNumbered, gotStrike := ConvertPitchToNumbered(tt.args.ks, tt.args.pitch)
			if gotNumbered != tt.wantNumbered {
				t.Errorf("ConvertPitchToNumbered() gotNumbered = %v, want %v", gotNumbered, tt.wantNumbered)
			}
			if gotStrike != tt.wantStrike {
				t.Errorf("ConvertPitchToNumbered() gotStrike = %v, want %v", gotStrike, tt.wantStrike)
			}
		})
	}
}

func TestGetOctave(t *testing.T) {
	type fields struct {
		Key       string
		Humanized string
		Fifth     int
		Mode      keysig.KeySignatureMode
	}
	type args struct {
		note musicxml.Note
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "C major octave 0",
			fields: fields{
				Key:       "c",
				Humanized: "do = c",
				Fifth:     0,
				Mode:      keysig.KeySignatureModeMajor,
			},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "D",
						Octave: 4,
					},
				},
			},
			want: 0,
		},
		{
			name: "C major octave -1",
			fields: fields{
				Key:       "c",
				Humanized: "do = c",
				Fifth:     0,
				Mode:      keysig.KeySignatureModeMajor,
			},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "D",
						Octave: 3,
					},
				},
			},
			want: -1,
		},
		{
			name: "C major octave (+)1",
			fields: fields{
				Key:       "c",
				Humanized: "do = c",
				Fifth:     0,
				Mode:      keysig.KeySignatureModeMajor,
			},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "D",
						Octave: 5,
					},
				},
			},
			want: 1,
		},
		{
			name: "(C4) ti/si on D major",
			fields: fields{
				Key:       "d",
				Humanized: "do = d",
				Fifth:     2,
				Mode:      keysig.KeySignatureModeMajor,
			},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "C",
						Octave: 4,
					},
				},
			},
			want: -1,
		},
		{
			name: "E4 - re - on D major",
			fields: fields{
				Key:       "d",
				Humanized: "do = d",
				Fifth:     2,
				Mode:      keysig.KeySignatureModeMajor,
			},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "E",
						Octave: 4,
					},
				},
			},
			want: 0,
		},
		{
			name: "C5 - si - on D major",
			fields: fields{
				Key:       "d",
				Humanized: "do = d",
				Fifth:     2,
				Mode:      keysig.KeySignatureModeMajor,
			},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "C",
						Octave: 5,
					},
					Accidental: musicxml.NoteAccidentalNatural,
				},
			},
			want: 0,
		},
		{
			name: "D5 - do - on D major",
			fields: fields{
				Key:       "d",
				Humanized: "do = d",
				Fifth:     2,
				Mode:      keysig.KeySignatureModeMajor,
			},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "D",
						Octave: 5,
					},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks := keysig.KeySignature{
				Key:       tt.fields.Key,
				Humanized: tt.fields.Humanized,
				Fifth:     tt.fields.Fifth,
			}
			if got := GetOctave(ks, tt.args.note); got != tt.want {
				t.Errorf("KeySignature.GetOctave() = %v, want %v", got, tt.want)
			}
		})
	}
}
