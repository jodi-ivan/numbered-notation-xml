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
// TODO: add support for the phrygian mode and mixolydian
// phrygian: kj-144a
// mixolydian: kj-215
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
