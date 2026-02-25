package keysig

import (
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
)

type KeySignature struct {
	Key       string
	Mode      Mode
	Humanized string
	Fifth     int
	RootPitch string
}

// TODO: add support for key signature changes
// key signature changes on kj-144a, kj-391
// TODO: add support for the phrygian
// phrygian: kj-144a
func NewKeySignature(key musicxml.KeySignature) KeySignature {
	keyMode := key.Mode
	if keyMode == "" {
		keyMode = "major"
	}
	fifths := key.Fifth

	mode := NewMode(keyMode)

	return KeySignature{
		Fifth:     fifths,
		Key:       mode.GetRoot(fifths),
		Mode:      mode,
		Humanized: mode.GetHumanized(fifths),
	}
}

func (ks *KeySignature) String() string {
	return ks.Humanized
}

func (ks *KeySignature) GetBasedPitch() string {
	result := ks.Mode.GetRoot(ks.Fifth)

	return result
}

func (ks KeySignature) GetPitchWithAccidental(note musicxml.Note) string {
	pitch := note.Pitch.Step
	var accidental musicxml.NoteAccidental

	if utils.Contains(accidentalsSet[ks.Fifth], pitch) >= 0 {
		if ks.Fifth > 0 {
			accidental = musicxml.NoteAccidentalSharp
		} else if ks.Fifth < 0 {
			accidental = musicxml.NoteAccidentalFlat
		}
	}

	if note.Accidental != "" {
		accidental = note.Accidental
	}
	pitch = fmt.Sprintf("%s%s", pitch, accidental.GetAccidental())

	return pitch
}

func (ks KeySignature) GetLetteredKeySignature() string {
	return ks.Mode.GetRoot(ks.Fifth)
}

func (ks KeySignature) BuildScale() []string {

	letters := []string{"C", "D", "E", "F", "G", "A", "B"}

	tonic := modeRoot[ks.Mode.Mode.String()][ks.Fifth]
	tonicLetter := string(tonic[0])

	var scaleLetters []string
	startIndex := 0
	for i, l := range letters {
		if l == tonicLetter {
			startIndex = i
			break
		}
	}

	for i := 0; i < 7; i++ {
		scaleLetters = append(scaleLetters, letters[(startIndex+i)%7])
	}

	accidentals := accidentalsSet[ks.Fifth]
	isSharpKey := ks.Fifth > 0
	isFlatKey := ks.Fifth < 0

	var scale []string

	for _, letter := range scaleLetters {

		altered := false

		for _, accLetter := range accidentals {
			if letter == accLetter {
				if isSharpKey {
					scale = append(scale, letter+"#")
				} else if isFlatKey {
					scale = append(scale, letter+"b")
				}
				altered = true
				break
			}
		}

		if !altered {
			scale = append(scale, letter)
		}
	}

	return scale
}
