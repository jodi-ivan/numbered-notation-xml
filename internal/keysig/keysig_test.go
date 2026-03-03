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

func TestKeySignature_String(t *testing.T) {

	t.Run("Getter", func(t *testing.T) {
		ks := NewKeySignature(musicxml.KeySignature{
			Fifth: 2,
		})
		expect := "do = d"
		if got := ks.String(); got != expect {
			t.Errorf("KeySignature.String() = %v, want %v", got, expect)
		}
	})
}

func TestKeySignature_GetBasedPitch(t *testing.T) {
	tests := []struct {
		name string
		ks   KeySignature
		want string
	}{
		{
			name: "C major",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 0}),
			want: "C",
		},
		{
			name: "D major",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 2}),
			want: "D",
		},
		{
			name: "minor",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 0, Mode: "minor"}),
			want: "A",
		},
		{
			name: "minor",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 2, Mode: "minor"}),
			want: "B",
		},
		{
			name: "dorian",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 0, Mode: "dorian"}),
			want: "D",
		},
		{
			name: "dorian",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 2, Mode: "dorian"}),
			want: "E",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ks.GetBasedPitch(); got != tt.want {
				t.Errorf("KeySignature.GetBasedPitch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeySignature_GetLetteredKeySignature(t *testing.T) {
	tests := []struct {
		name string
		ks   KeySignature
		want string
	}{
		{
			name: "C major",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 0}),
			want: "C",
		},
		{
			name: "D major Circle of Fifth 2",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 2}),
			want: "D",
		},
		{
			name: "minor",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 0, Mode: "minor"}),
			want: "A",
		},
		{
			name: "minor Circle of Fifth 2",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 2, Mode: "minor"}),
			want: "B",
		},
		{
			name: "dorian",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 0, Mode: "dorian"}),
			want: "D",
		},
		{
			name: "dorian Circle of Fifth 2",
			ks:   NewKeySignature(musicxml.KeySignature{Fifth: 2, Mode: "dorian"}),
			want: "E",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ks.GetLetteredKeySignature(); got != tt.want {
				t.Errorf("KeySignature.GetLetteredKeySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeySignature_BuildScale(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		key  musicxml.KeySignature
		want []string
	}{
		// --- MAJOR (Ionian) ---
		{name: "Bb major (-2)", key: musicxml.KeySignature{Fifth: -2, Mode: "major"}, want: []string{"Bb", "C", "D", "Eb", "F", "G", "A"}},
		{name: "F major (-1)", key: musicxml.KeySignature{Fifth: -1, Mode: "major"}, want: []string{"F", "G", "A", "Bb", "C", "D", "E"}},
		{name: "C major (0)", key: musicxml.KeySignature{Fifth: 0, Mode: "major"}, want: []string{"C", "D", "E", "F", "G", "A", "B"}},
		{name: "G major (1)", key: musicxml.KeySignature{Fifth: 1, Mode: "major"}, want: []string{"G", "A", "B", "C", "D", "E", "F#"}},
		{name: "D major (2)", key: musicxml.KeySignature{Fifth: 2, Mode: "major"}, want: []string{"D", "E", "F#", "G", "A", "B", "C#"}},

		// --- DORIAN ---
		{name: "C dorian (-2)", key: musicxml.KeySignature{Fifth: -2, Mode: "dorian"}, want: []string{"C", "D", "Eb", "F", "G", "A", "Bb"}},
		{name: "G dorian (-1)", key: musicxml.KeySignature{Fifth: -1, Mode: "dorian"}, want: []string{"G", "A", "Bb", "C", "D", "E", "F"}},
		{name: "D dorian (0)", key: musicxml.KeySignature{Fifth: 0, Mode: "dorian"}, want: []string{"D", "E", "F", "G", "A", "B", "C"}},
		{name: "A dorian (1)", key: musicxml.KeySignature{Fifth: 1, Mode: "dorian"}, want: []string{"A", "B", "C", "D", "E", "F#", "G"}},
		{name: "E dorian (2)", key: musicxml.KeySignature{Fifth: 2, Mode: "dorian"}, want: []string{"E", "F#", "G", "A", "B", "C#", "D"}},

		// --- MIXOLYDIAN ---
		{name: "F mixolydian (-2)", key: musicxml.KeySignature{Fifth: -2, Mode: "mixolydian"}, want: []string{"F", "G", "A", "Bb", "C", "D", "Eb"}},
		{name: "C mixolydian (-1)", key: musicxml.KeySignature{Fifth: -1, Mode: "mixolydian"}, want: []string{"C", "D", "E", "F", "G", "A", "Bb"}},
		{name: "G mixolydian (0)", key: musicxml.KeySignature{Fifth: 0, Mode: "mixolydian"}, want: []string{"G", "A", "B", "C", "D", "E", "F"}},
		{name: "D mixolydian (1)", key: musicxml.KeySignature{Fifth: 1, Mode: "mixolydian"}, want: []string{"D", "E", "F#", "G", "A", "B", "C"}},
		{name: "A mixolydian (2)", key: musicxml.KeySignature{Fifth: 2, Mode: "mixolydian"}, want: []string{"A", "B", "C#", "D", "E", "F#", "G"}},

		// --- MINOR (Aeolian) ---
		{name: "G minor (-2)", key: musicxml.KeySignature{Fifth: -2, Mode: "minor"}, want: []string{"G", "A", "Bb", "C", "D", "Eb", "F"}},
		{name: "D minor (-1)", key: musicxml.KeySignature{Fifth: -1, Mode: "minor"}, want: []string{"D", "E", "F", "G", "A", "Bb", "C"}},
		{name: "A minor (0)", key: musicxml.KeySignature{Fifth: 0, Mode: "minor"}, want: []string{"A", "B", "C", "D", "E", "F", "G"}},
		{name: "E minor (1)", key: musicxml.KeySignature{Fifth: 1, Mode: "minor"}, want: []string{"E", "F#", "G", "A", "B", "C", "D"}},
		{name: "B minor (2)", key: musicxml.KeySignature{Fifth: 2, Mode: "minor"}, want: []string{"B", "C#", "D", "E", "F#", "G", "A"}},

		// // --- PHRYGIAN ---
		// {name: "D phrygian (-2)", key: musicxml.KeySignature{Fifth: -2, Mode: "phrygian"}, want: []string{"D", "Eb", "F", "G", "A", "Bb", "C"}},
		// {name: "A phrygian (-1)", key: musicxml.KeySignature{Fifth: -1, Mode: "phrygian"}, want: []string{"A", "Bb", "C", "D", "E", "F", "G"}},
		// {name: "E phrygian (0)", key: musicxml.KeySignature{Fifth: 0, Mode: "phrygian"}, want: []string{"E", "F", "G", "A", "B", "C", "D"}},
		// {name: "B phrygian (1)", key: musicxml.KeySignature{Fifth: 1, Mode: "phrygian"}, want: []string{"B", "C", "D", "E", "F#", "G", "A"}},
		// {name: "F# phrygian (2)", key: musicxml.KeySignature{Fifth: 2, Mode: "phrygian"}, want: []string{"F#", "G", "A", "B", "C#", "D", "E"}},

		// --- LYDIAN ---
		// {name: "Eb lydian (-2)", key: musicxml.KeySignature{Fifth: -2, Mode: "lydian"}, want: []string{"Eb", "F", "G", "A", "Bb", "C", "D"}},
		// {name: "Bb lydian (-1)", key: musicxml.KeySignature{Fifth: -1, Mode: "lydian"}, want: []string{"Bb", "C", "D", "E", "F", "G", "A"}},
		// {name: "F lydian (0)", key: musicxml.KeySignature{Fifth: 0, Mode: "lydian"}, want: []string{"F", "G", "A", "B", "C", "D", "E"}},
		// {name: "C lydian (1)", key: musicxml.KeySignature{Fifth: 1, Mode: "lydian"}, want: []string{"C", "D", "E", "F#", "G", "A", "B"}},
		// {name: "G lydian (2)", key: musicxml.KeySignature{Fifth: 2, Mode: "lydian"}, want: []string{"G", "A", "B", "C#", "D", "E", "F#"}},

		// // --- LOCRIAN ---
		// {name: "A locrian (-2)", key: musicxml.KeySignature{Fifth: -2, Mode: "locrian"}, want: []string{"A", "Bb", "C", "D", "Eb", "F", "G"}},
		// {name: "E locrian (-1)", key: musicxml.KeySignature{Fifth: -1, Mode: "locrian"}, want: []string{"E", "F", "G", "A", "Bb", "C", "D"}},
		// {name: "B locrian (0)", key: musicxml.KeySignature{Fifth: 0, Mode: "locrian"}, want: []string{"B", "C", "D", "E", "F", "G", "A"}},
		// {name: "F# locrian (1)", key: musicxml.KeySignature{Fifth: 1, Mode: "locrian"}, want: []string{"F#", "G", "A", "B", "C", "D", "E"}},
		// {name: "C# locrian (2)", key: musicxml.KeySignature{Fifth: 2, Mode: "locrian"}, want: []string{"C#", "D", "E", "F#", "G", "A", "B"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks := NewKeySignature(tt.key)
			got := ks.BuildScale()
			assert.Equal(t, tt.want, got)
		})
	}
}
