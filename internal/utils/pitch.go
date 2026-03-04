package utils

func IsPitchEqual(one, two string) bool {
	if one == two {
		return true
	}
	pitches := map[string][]string{
		"C": []string{"B#", "Dbb"}, "C#": []string{"Db", "Bx"},
		"Cb": []string{"B", "Ax"}, "Cx": []string{"Ax", "D"},
		"Cbb": []string{"B", "Ax"},

		"D": []string{"Cx", "Ebb"}, "Dbb": []string{"C", "B#"},
		"Db": []string{"C#"}, "Dx": []string{"E", "Fb"},
		"D#": []string{"Eb", "Fbb"},

		"E": []string{"Dx", "Fb"}, "Ebb": []string{"D", "Cx"},
		"Eb": []string{"D#", "Fbb"}, "Ex": []string{"F#", "Gb"},
		"E#": []string{"F", "Gbb"},

		"F": []string{"E#", "Gbb"}, "Fbb": []string{"Eb", "D#"},
		"Fb": []string{"E", "Dx"}, "Fx": []string{"G", "Abb"},
		"F#": []string{"Gb", "Ex"},

		"G": []string{"Abb", "Fx"}, "Gbb": []string{"F", "E#"},
		"Gb": []string{"F#", "Ex"}, "Gx": []string{"A", "Bbb"},
		"G#": []string{"Ab"},

		"A": []string{"Gx", "Bbb"}, "Abb": []string{"G", "Fx"},
		"Ab": []string{"G#"}, "Ax": []string{"B", "Cb"},
		"A#": []string{"Bb", "Cbb"},

		"B": []string{"Cb", "Ax"}, "Bbb": []string{"A", "Gx"},
		"Bb": []string{"A#", "Cbb"}, "Bx": []string{"C#", "Db"},
		"B#": []string{"C", "Dbb"},
	}
	result := Contains(pitches[one], two) >= 0

	return result
}

// pitchSemitone returns the semitone value (0-11) for a given pitch string.
// Base note semitones: C=0, D=2, E=4, F=5, G=7, A=9, B=11
// Accidentals: bb=-2, b=-1, none=0, #=+1, x=+2
func pitchSemitone(pitch string) int {
	if len(pitch) == 0 {
		return -1
	}

	baseSemitones := map[byte]int{
		'C': 0, 'D': 2, 'E': 4, 'F': 5, 'G': 7, 'A': 9, 'B': 11,
	}

	base, ok := baseSemitones[pitch[0]]
	if !ok {
		return -1
	}

	accidental := 0
	if len(pitch) > 1 {
		suffix := pitch[1:]
		switch suffix {
		case "bb":
			accidental = -2
		case "b":
			accidental = -1
		case "#":
			accidental = 1
		case "x":
			accidental = 2
		}
	}

	// Wrap around with mod 12 to handle edge cases like B# = C (12 -> 0)
	return ((base+accidental)%12 + 12) % 12
}

// ComparePitch compares two pitches assumed to be in the same octave (C4 -> B4).
// Returns:
//
//	 0 -> both pitches are equal
//	 1 -> 1st param > 2nd param
//	-1 -> 1st param < 2nd param
func ComparePitch(one, two string) int {
	if IsPitchEqual(one, two) {
		return 0
	}

	oneSemitone := pitchSemitone(one)
	twoSemitone := pitchSemitone(two)

	if oneSemitone > twoSemitone {
		return 1
	}

	return -1
}
