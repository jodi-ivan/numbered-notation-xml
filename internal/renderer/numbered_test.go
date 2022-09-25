package renderer

import (
	"testing"
)

func TestConvertPitchToNumbered(t *testing.T) {
	type args struct {
		ks    KeySignature
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "c",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     2,
					Key:       "d",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     2,
					Key:       "d",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     2,
					Key:       "d",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     2,
					Key:       "d",
					Mode:      KeySignatureModeMajor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     0,
					Key:       "a",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     2,
					Key:       "b",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     2,
					Key:       "b",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     2,
					Key:       "b",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     2,
					Key:       "b",
					Mode:      KeySignatureModeMinor,
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
				ks: KeySignature{
					Fifth:     -1,
					Key:       "f",
					Mode:      KeySignatureModeMajor,
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

func Test_getNextHalfStep(t *testing.T) {
	type args struct {
		pitch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "C",
			args: args{
				pitch: "C",
			},
			want: "C#",
		},
		{
			name: "C#",
			args: args{
				pitch: "C#",
			},
			want: "D",
		},
		{
			name: "D",
			args: args{
				pitch: "D",
			},
			want: "D#",
		},
		{
			name: "D#",
			args: args{
				pitch: "D#",
			},
			want: "E",
		},
		{
			name: "E",
			args: args{
				pitch: "E",
			},
			want: "F",
		},
		{
			name: "B",
			args: args{
				pitch: "B",
			},
			want: "C",
		},
		{
			name: "B#",
			args: args{
				pitch: "B#",
			},
			want: "C#",
		},
		{
			name: "Bx",
			args: args{
				pitch: "Bx",
			},
			want: "D",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNextHalfStep(tt.args.pitch); got != tt.want {
				t.Errorf("getNextHalfStep() = %v, want %v", got, tt.want)
			}
		})
	}
}
