package moveabledo

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
)

func ConvertPitchToNumbered(ks keysig.Key, pitch string) (int, bool) {

	// Build the diatonic scale of the key
	scale := ks.BuildScale() // []string length 7

	// Extract letter only (ignore accidental)
	pitchLetter := string(pitch[0]) // "C" from "C#"

	// Find degree by letter comparison
	for i, scalePitch := range scale {

		scaleLetter := string(scalePitch[0])

		if scaleLetter == pitchLetter {

			// Found correct diatonic degree
			degree := i + 1

			// Now check accidental match
			if utils.IsPitchEqual(scalePitch, pitch) {
				return degree, false
			}

			// Same letter but different accidental → altered
			return degree, true
		}
	}

	// Should never happen unless invalid pitch
	return 0, false
}

func GetNumberedNotation(ks keysig.Key, note musicxml.Note) (numberedNote int, octave int, strikethrough bool) {
	if note.Rest != nil {
		return 0, 0, false
	}
	octave = GetOctave(ks, note)

	pitch := ks.GetPitchWithAccidental(note)
	numberedNote, strikethrough = ConvertPitchToNumbered(ks, pitch)

	return numberedNote, octave, strikethrough
}

func GetOctave(ks keysig.Key, note musicxml.Note) int {
	pitch := ks.GetPitchWithAccidental(note)
	if ks.Mode.Mode == keysig.KeySignatureModeMajor && ks.Fifth == 0 { // C major
		return note.Pitch.Octave - 4
	}

	offset := 0

	if (ks.Fifth == -2 && ks.Mode.Mode == keysig.KeySignatureModeMajor) || ks.Fifth == 3 { // Bb Major or A major or F#min
		offset = +1
	}

	do := ks.GetBasedPitch()
	direction := utils.ComparePitch(do, pitch)
	if direction == 1 { // below do
		return note.Pitch.Octave - 4 - 1 + offset
	}

	return note.Pitch.Octave - 4 + offset

}
