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
		"C":   []string{"B#", "Dbb"},
		"Cb":  []string{"B", "Ax"},
		"C#":  []string{"Db", "Bx"},
		"Cx":  []string{"B", "Ax"},
		"Cbb": []string{"B", "Ax"},
		"D":   []string{"Cx", "Ebb"},
		"Db":  []string{"C#"},
		"D#":  []string{"Eb", "Fbb"},
		"Dbb": []string{"C", "B#"},
		"Dx":  []string{"E", "Fb"},
		"E":   []string{"Dx", "Fb"},
		"Eb":  []string{"D#", "Fbb"},
		"E#":  []string{"F", "Gbb"},
		"Ebb": []string{"D", "Cx"},
		"Ex":  []string{"F#", "Gb"},
		"F":   []string{"E#", "Gbb"},
		"Fb":  []string{"E", "Dx"},
		"F#":  []string{"Gb", "Ex"},
		"Fbb": []string{"Eb", "D#"},
		"Fx":  []string{"G", "Abb"},
		"G":   []string{"Abb", "Fx"},
		"Gb":  []string{"F#", "Ex"},
		"G#":  []string{"Ab"},
		"Gbb": []string{"F", "E#"},
		"Gx":  []string{"A", "Bbb"},
		"A":   []string{"Gx", "Bbb"},
		"Ab":  []string{"G#"},
		"A#":  []string{"Bb", "Cbb"},
		"Abb": []string{"G", "Fx"},
		"Ax":  []string{"B", "Cb"},
		"B":   []string{"Cb", "Ax"},
		"Bb":  []string{"A#", "Cbb"},
		"B#":  []string{"C", "Dbb"},
		"Bbb": []string{"A", "Gx"},
		"Bx":  []string{"C#", "Db"},
	}
	result := Contains(pitches[one], two) >= 0

	return result
}

// compare pitch
// returns
//    0 -> both pitch are equal
//    1 -> 1st param > 2nd param
//   -1 -> 1st param < 2nd param
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
