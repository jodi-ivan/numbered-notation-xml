package keysig

import "testing"

func Test_getSpelledNumberedNotation(t *testing.T) {
	type args struct {
		letterNotation string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "C",
			args: args{
				letterNotation: "C",
			},
			want: "c",
		},
		{
			name: "C#",
			args: args{
				letterNotation: "C#",
			},
			want: "cis",
		},
		{
			name: "Db",
			args: args{
				letterNotation: "Db",
			},
			want: "des",
		},
		{
			name: "D",
			args: args{
				letterNotation: "D",
			},
			want: "d",
		},
		{
			name: "D#",
			args: args{
				letterNotation: "D#",
			},
			want: "dis",
		},
		{
			name: "Eb",
			args: args{
				letterNotation: "Eb",
			},
			want: "es",
		},
		{
			name: "E",
			args: args{
				letterNotation: "E",
			},
			want: "e",
		},
		{
			name: "F",
			args: args{
				letterNotation: "F",
			},
			want: "f",
		},
		{
			name: "F#",
			args: args{
				letterNotation: "F#",
			},
			want: "fis",
		},
		{
			name: "Gb",
			args: args{
				letterNotation: "Gb",
			},
			want: "ges",
		},
		{
			name: "G",
			args: args{
				letterNotation: "G",
			},
			want: "g",
		},
		{
			name: "G#",
			args: args{
				letterNotation: "G#",
			},
			want: "gis",
		},
		{
			name: "Ab",
			args: args{
				letterNotation: "Ab",
			},
			want: "as",
		},
		{
			name: "A",
			args: args{
				letterNotation: "A",
			},
			want: "a",
		},
		{
			name: "A#",
			args: args{
				letterNotation: "A#",
			},
			want: "ais",
		},
		{
			name: "Bb",
			args: args{
				letterNotation: "Bb",
			},
			want: "bes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSpelledNumberedNotation(tt.args.letterNotation); got != tt.want {
				t.Errorf("getNumberedNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}
