package keysig

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
)

type Key struct {
	Key       string
	Mode      Mode
	Humanized string
	Fifth     int
	Measure   int
}

type KeySignature struct {
	IsMixed    bool
	Signatures []Key
}

// TODO: support more than once key known (processed) signature changes. currently only support once or all keys
func NewKeySignature(ctx context.Context, measures []musicxml.Measure) KeySignature {
	signatures := []Key{}

	for _, measure := range measures {
		if measure.Attribute != nil && measure.Attribute.Key != nil {
			key := NewKey(measure.Attribute.Key)
			key.Measure = measure.Number
			signatures = append(signatures, key)

		}
	}

	return KeySignature{
		IsMixed:    len(signatures) > 1,
		Signatures: signatures,
	}

}

func (ks *KeySignature) GetKeyOnMeasure(ctx context.Context, measure int) Key {
	if len(ks.Signatures) == 1 {
		return ks.Signatures[0]
	}

	// get the time
	currentTime := ks.Signatures[0]

	counter := 0
	var prev Key
	prev = ks.Signatures[0]
	found := true

	for currentTime.Measure <= measure && counter < len(ks.Signatures) {
		prev = currentTime
		currentTime = ks.Signatures[counter]
		counter++

		if currentTime.Measure <= measure && counter == len(ks.Signatures) {
			found = false
		}

	}

	if !found {
		return ks.Signatures[len(ks.Signatures)-1]
	}

	return prev
}

func NewKey(key *musicxml.KeySignature) Key {
	keyMode := key.Mode
	if keyMode == "" {
		keyMode = "major"
	}
	fifths := key.Fifth

	mode := NewMode(keyMode)

	return Key{
		Fifth:     fifths,
		Key:       mode.GetRoot(fifths),
		Mode:      mode,
		Humanized: mode.GetHumanized(fifths),
	}
}

func (ks *Key) String() string {
	return ks.Humanized
}

func (ks *Key) GetBasedPitch() string {
	result := ks.Mode.GetRoot(ks.Fifth)

	return result
}

func (ks *Key) GetAccidentals() []string {
	return accidentalsSet[ks.Fifth]
}
func (ks Key) GetPitchWithAccidental(note musicxml.Note) string {
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

func (ks Key) GetLetteredKeySignature() string {
	return ks.Mode.GetRoot(ks.Fifth)
}

func (ks Key) BuildScale() []string {

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
