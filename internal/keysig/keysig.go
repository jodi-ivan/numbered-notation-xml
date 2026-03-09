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
	IsMixed     bool
	MeasureText map[int]string
	Signatures  []Key
}

// TODO: support more than once key known (processed) signature changes. currently only support once or all keys
func NewKeySignature(ctx context.Context, measures []musicxml.Measure) KeySignature {
	signatures := []Key{}
	various := map[int]bool{}
	inidcator := map[int]string{}

	for _, measure := range measures {
		if measure.Attribute != nil && measure.Attribute.Key != nil {
			if various[measure.Attribute.Key.Fifth] {
				continue
			}
			various[measure.Attribute.Key.Fifth] = true

			if len(signatures) == 0 {
				key := NewKey(measure.Attribute.Key)
				key.Measure = measure.Number
				signatures = append(signatures, key)
				continue
			}

			// detect changes
			lastKey := signatures[len(signatures)-1]
			indicatorText, mode := Transtion(lastKey.Fifth, lastKey.Mode.String(), measure.Attribute.Key.Fifth)
			inidcator[measure.Number] = indicatorText
			key := NewKey(&musicxml.KeySignature{Fifth: measure.Attribute.Key.Fifth, Mode: mode.String()})
			key.Measure = measure.Number
			signatures = append(signatures, key)

		}
	}

	return KeySignature{
		IsMixed:     len(signatures) > 1,
		Signatures:  signatures,
		MeasureText: inidcator,
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

// Mode metadata for offsets (semitones from Ionian) and Solfège mapping
var modeInfo = map[string]struct {
	Offset  int
	Solfege []string
}{
	"major":      {0, []string{"do", "re", "mi", "fa", "sol", "la", "ti"}},
	"dorian":     {2, []string{"re", "mi", "fa", "sol", "la", "ti", "do"}},
	"phrygian":   {4, []string{"mi", "fa", "sol", "la", "ti", "do", "re"}},
	"lydian":     {5, []string{"fa", "sol", "la", "ti", "do", "re", "mi"}},
	"mixolydian": {7, []string{"sol", "la", "ti", "do", "re", "mi", "fa"}},
	"minor":      {9, []string{"la", "ti", "do", "re", "mi", "fa", "sol"}},
	"locrian":    {11, []string{"ti", "do", "re", "mi", "fa", "sol", "la"}},
}

func Transtion(oldFifths int, oldMode string, newFifths int) (string, Mode) {
	// 1. Calculate the Tonic Pitch of the old key
	// Pitch = (Fifths * 7 + ModeOffset) mod 12
	oldTonic := (oldFifths*7 + modeInfo[oldMode].Offset) % 12
	if oldTonic < 0 {
		oldTonic += 12
	}

	// 2. Educated Guess: The Tonic remains the same (Parallel Key)
	// Find which mode for 'newFifths' results in the same Tonic
	guessedMode := "major"
	for mode, info := range modeInfo {
		currentTonic := (newFifths*7 + info.Offset) % 12
		if currentTonic < 0 {
			currentTonic += 12
		}
		if currentTonic == oldTonic {
			guessedMode = mode
			break
		}
	}

	// 3. Generate Pivot Notation for the Dominant (5th degree)
	// In hymns, the pivot is almost always the Dominant note (index 4)
	oldPivotName := modeInfo[oldMode].Solfege[4]
	newPivotName := modeInfo[guessedMode].Solfege[4]

	return fmt.Sprintf("%s = %s", oldPivotName, newPivotName), NewMode(guessedMode)
}

func TranstionFromTwoKeySignatures(from, to Key) string {
	indicator := ""
	if from.BuildScale()[0] == to.BuildScale()[0] {
		indicator = fmt.Sprintf("%s = %s", to.Mode.Mode.GetNumberedRoot(), from.Mode.Mode.GetNumberedRoot())
	} else {
		indicator, _ = Transtion(to.Fifth, to.Mode.String(), from.Fifth)
	}

	return indicator
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
