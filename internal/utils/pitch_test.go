package utils

import "testing"

func TestIsPitchEqual(t *testing.T) {
	tests := []struct {
		name     string
		one      string
		two      string
		expected bool
	}{
		// Identical pitches
		{"identical C", "C", "C", true},
		{"identical F#", "F#", "F#", true},
		{"identical Bb", "Bb", "Bb", true},

		// Enharmonic equivalents - natural notes
		{"C = B#", "C", "B#", true},
		{"C = Dbb", "C", "Dbb", true},
		{"D = Cx", "D", "Cx", true},
		{"D = Ebb", "D", "Ebb", true},
		{"E = Dx", "E", "Dx", true},
		{"E = Fb", "E", "Fb", true},
		{"F = E#", "F", "E#", true},
		{"F = Gbb", "F", "Gbb", true},
		{"G = Fx", "G", "Fx", true},
		{"G = Abb", "G", "Abb", true},
		{"A = Gx", "A", "Gx", true},
		{"A = Bbb", "A", "Bbb", true},
		{"B = Cb", "B", "Cb", true},
		{"B = Ax", "B", "Ax", true},

		// Enharmonic equivalents - sharps/flats
		{"C# = Db", "C#", "Db", true},
		{"C# = Bx", "C#", "Bx", true},
		{"D# = Eb", "D#", "Eb", true},
		{"Eb = Fbb", "Eb", "Fbb", true},
		{"F# = Gb", "F#", "Gb", true},
		{"F# = Ex", "F#", "Ex", true},
		{"G# = Ab", "G#", "Ab", true},
		{"A# = Bb", "A#", "Bb", true},
		{"Bb = Cbb", "Bb", "Cbb", true},

		// Double accidentals enharmonics
		{"Bx = C#", "Bx", "C#", true},
		{"Bx = Db", "Bx", "Db", true},
		{"Cx = D", "Cx", "D", true},
		{"Dx = E", "Dx", "E", true},
		{"Ex = F#", "Ex", "F#", true},
		{"Fx = G", "Fx", "G", true},
		{"Gx = A", "Gx", "A", true},
		{"Ax = B", "Ax", "B", true},

		// Symmetry: reverse direction
		{"B# = C", "B#", "C", true},
		{"Dbb = C", "Dbb", "C", true},
		{"Gb = F#", "Gb", "F#", true},
		{"Ab = G#", "Ab", "G#", true},
		{"Cb = B", "Cb", "B", true},

		// Non-equal pitches
		{"C != D", "C", "D", false},
		{"C != G", "C", "G", false},
		{"F# != G", "F#", "G", false},
		{"A != B", "A", "B", false},
		{"Bb != B", "Bb", "B", false},
		{"C# != C", "C#", "C", false},
		{"Eb != E", "Eb", "E", false},
		{"Ab != A", "Ab", "A", false},

		// Unknown/empty pitches
		{"empty strings", "", "", true},
		{"one empty", "C", "", false},
		{"two empty", "", "C", false},
		{"unknown pitch", "X", "C", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPitchEqual(tt.one, tt.two)
			if result != tt.expected {
				t.Errorf("IsPitchEqual(%q, %q) = %v, want %v", tt.one, tt.two, result, tt.expected)
			}
		})
	}
}

func TestComparePitch(t *testing.T) {
	tests := []struct {
		name     string
		one      string
		two      string
		expected int
	}{
		// Equal pitches -> 0
		{"identical C", "C", "C", 0},
		{"identical F#", "F#", "F#", 0},
		{"enharmonic C = B#", "C", "B#", 0},
		{"enharmonic C = Dbb", "C", "Dbb", 0},
		{"enharmonic F# = Gb", "F#", "Gb", 0},
		{"enharmonic Eb = D#", "Eb", "D#", 0},
		{"enharmonic G# = Ab", "G#", "Ab", 0},
		{"enharmonic B = Cb", "B", "Cb", 0},

		// Different base notes -> compare by position in C D E F G A B
		{"C < D", "C", "D", -1},
		{"D > C", "D", "C", 1},
		{"C < G", "C", "G", -1},
		{"G > C", "G", "C", 1},
		{"A > F", "A", "F", 1},
		{"F < A", "F", "A", -1},
		{"B > A", "B", "A", 1},
		{"A < B", "A", "B", -1},
		{"E < G", "E", "G", -1},
		{"G > E", "G", "E", 1},

		// Different base notes with accidentals -> still uses base letter
		{"C# vs D -> C# < D", "C#", "D", -1},
		{"Db vs C -> Db > C", "Db", "C", 1},
		{"F# vs G -> F# < G", "F#", "G", -1},
		{"Gb vs F -> Gb > F", "Gb", "F", 1},
		{"Bb vs C -> Bb > C", "Bb", "C", 1},
		{"Ab vs B -> Ab < B", "Ab", "B", -1},

		// Same base letter, one has accidental -> natural vs modified
		{"C vs C# -> C < C#", "C", "C#", -1},
		{"C# vs C -> C# > C", "C#", "C", 1},
		{"G vs Gb -> G > Gb", "G", "Gb", 1},
		{"Gb vs G -> Gb < G", "Gb", "G", -1},
		{"A vs Ab -> A > Ab", "A", "Ab", 1},
		{"Ab vs A -> Ab < A", "Ab", "A", -1},
		{"B vs Bb -> B > Bb", "B", "Bb", 1},
		{"Bb vs B -> Bb < B", "Bb", "B", -1},
		{"D vs D# -> D < D#", "D", "D#", -1},
		{"E vs Eb -> E > Eb", "E", "Eb", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComparePitch(tt.one, tt.two)
			if result != tt.expected {
				t.Errorf("ComparePitch(%q, %q) = %v, want %v", tt.one, tt.two, result, tt.expected)
			}
		})
	}
}

func TestPitchSemitone(t *testing.T) {
	tests := []struct {
		pitch    string
		expected int
	}{
		{"C", 0}, {"C#", 1}, {"Db", 1}, {"Cx", 2},
		{"D", 2}, {"D#", 3}, {"Eb", 3}, {"Dbb", 0},
		{"E", 4}, {"Eb", 3}, {"E#", 5}, {"Ebb", 2},
		{"F", 5}, {"F#", 6}, {"Fb", 4}, {"Fx", 7},
		{"G", 7}, {"G#", 8}, {"Gb", 6}, {"Gx", 9},
		{"A", 9}, {"A#", 10}, {"Ab", 8}, {"Ax", 11},
		{"B", 11}, {"Bb", 10}, {"B#", 0}, {"Bx", 1},
	}

	for _, tt := range tests {
		t.Run(tt.pitch, func(t *testing.T) {
			result := pitchSemitone(tt.pitch)
			if result != tt.expected {
				t.Errorf("pitchSemitone(%q) = %d, want %d", tt.pitch, result, tt.expected)
			}
		})
	}
}
