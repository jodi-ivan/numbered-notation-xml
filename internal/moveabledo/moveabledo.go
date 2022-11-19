package moveabledo

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
)

func ConvertPitchToNumbered(ks keysig.KeySignature, pitch string) (numbered int, strike bool) {

	var current string

	steps := ks.Mode.GetScaleSteps()
	current = ks.GetLetteredKeySignature()

	counter := 0

	for !(utils.IsPitchEqual(current, pitch)) {
		current = utils.GetNextHalfStep(current)
		counter++
	}

	stepped := 0.0
	increase := 0

	for stepped < float64(counter)/2 {

		stepped += steps[increase]
		increase++
	}

	if stepped > (float64(counter) / 2) {
		return increase, true
	}
	return increase + 1, false

}

func GetNumberedNotation(ks keysig.KeySignature, note musicxml.Note) (numberedNote int, octave int, strikethrough bool) {
	if note.Rest != nil {
		return 0, 0, false
	}
	octave = GetOctave(ks, note)

	pitch := ks.GetPitchWithAccidental(note)
	numberedNote, strikethrough = ConvertPitchToNumbered(ks, pitch)

	return numberedNote, octave, strikethrough
}

func GetOctave(ks keysig.KeySignature, note musicxml.Note) int {
	pitch := ks.GetPitchWithAccidental(note)
	if ks.Mode.Mode == keysig.KeySignatureModeMajor && ks.Fifth == 0 { // C major
		return note.Pitch.Octave - 4
	}

	offset := 0

	if ks.Fifth == -2 && ks.Mode.Mode == keysig.KeySignatureModeMajor {
		offset = +1
	}

	do := ks.GetBasedPitch()
	direction := utils.ComparePitch(do, pitch)
	if direction == 1 { // below do
		return note.Pitch.Octave - 4 - 1 + offset
	}

	return note.Pitch.Octave - 4 + offset

}
