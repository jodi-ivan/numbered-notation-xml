package utils

import (
	"fmt"
	"strings"
)

func GetNextHalfStep(pitch string) string {
	step := []string{"C", "D", "E", "F", "G", "A", "B"}

	wholeStep := func(p string) string {
		index := Contains(step, p)
		if index == len(step)-1 {
			return step[0]
		}
		return step[index+1]
	}

	if strings.HasSuffix(pitch, "bb") { // double flat
		return fmt.Sprintf("%sb", string(pitch[0]))
	}

	if strings.HasSuffix(pitch, "b") { // flat, just remove the flat
		return string(pitch[0])
	}

	if strings.HasSuffix(pitch, "x") { // double sharp
		if string(pitch[0]) == "B" {
			return GetNextHalfStep("C#")
		}

		if string(pitch[0]) == "E" {
			return GetNextHalfStep("F#")
		}

		newPitch := wholeStep(string(pitch[0]))
		if newPitch == "B" || newPitch == "E" {
			return wholeStep(newPitch)
		}

		return fmt.Sprintf("%s#", newPitch)
	}

	if strings.HasSuffix(pitch, "#") {
		if string(pitch[0]) == "B" || string(pitch[0]) == "E" {
			return fmt.Sprintf("%s#", wholeStep(string(pitch[0])))
		}

		return wholeStep(string(pitch[0]))
	}

	if string(pitch[0]) == "B" || string(pitch[0]) == "E" {
		return wholeStep(string(pitch[0]))
	}

	return fmt.Sprintf("%s#", string(pitch[0]))

}

func IsPitchEqual(one, two string) bool {
	if one == two {
		return true
	}
	pitches := map[string][]string{
		"C": []string{"B#", "Dbb"}, "C#": []string{"Db", "Bx"},
		"Cb": []string{"B", "Ax"}, "Cx": []string{"B", "Ax"},
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

// compare pitch
// returns
//
//	 0 -> both pitch are equal
//	 1 -> 1st param > 2nd param
//	-1 -> 1st param < 2nd param
//
// Assume both pitch on the same octave between C4 -> B4
func ComparePitch(one, two string) int {
	if IsPitchEqual(one, two) {
		return 0
	}

	pitches := []string{"C", "D", "E", "F", "G", "A", "B"}

	if string(one[0]) != string(two[0]) {

		oneIndex := Contains(pitches, string(one[0]))
		twoIndex := Contains(pitches, string(two[0]))

		if oneIndex > twoIndex {
			return 1
		}

		return -1
	}
	if len(two) > 1 {
		return 1
	}

	return -1

}
