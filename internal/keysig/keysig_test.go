package keysig

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func TestNewKeySignature(t *testing.T) {
	type args struct {
		key musicxml.KeySignature
	}
	tests := []struct {
		name string
		args args
		want KeySignature
	}{
		{
			name: "empty key mode, expecting to be C major",
			args: args{
				key: musicxml.KeySignature{},
			},
			want: KeySignature{
				Key: "C",
				Mode: Mode{
					rootLettered: "C",
					humanized:    "do = c",
				},
				Humanized: "do = c",
			},
		},
		{
			name: "empty key minor mode, expecting to be A minor",
			args: args{
				key: musicxml.KeySignature{
					Mode: "minor",
				},
			},
			want: KeySignature{
				Key: "A",
				Mode: Mode{
					rootLettered: "A",
					humanized:    "la = a",
					Mode:         KeySignatureModeMinor,
				},
				Humanized: "la = a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKeySignature(tt.args.key); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewKeySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeySignature_GetPitchWithAccidental(t *testing.T) {
	type fields struct {
		Fifth int
	}
	type args struct {
		note musicxml.Note
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "no accidental",
			fields: fields{},
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
			want: "C",
		},
		{
			name:   "no accidental - but the music specified to be so",
			fields: fields{},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "C",
						Octave: 4,
					},
					Accidental: musicxml.NoteAccidentalSharp,
				},
			},
			want: "C#",
		},

		{
			name: "2 sharps - D major",
			fields: fields{
				Fifth: 2,
			},
			args: args{
				note: musicxml.Note{
					Pitch: struct {
						Step   string `xml:"step"`
						Octave int    `xml:"octave"`
					}{
						Step:   "F",
						Octave: 4,
					},
				},
			},
			want: "F#",
		},
		{
			name: "2 flats - Bb Major",
			fields: fields{
				Fifth: -2,
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
			want: "Eb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks := KeySignature{
				Fifth: tt.fields.Fifth,
			}
			if got := ks.GetPitchWithAccidental(tt.args.note); got != tt.want {
				t.Errorf("KeySignature.GetPitchWithAccidental() = %v, want %v", got, tt.want)
			}
		})
	}
}
