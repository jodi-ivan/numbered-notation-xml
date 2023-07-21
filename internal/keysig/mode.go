package keysig

import (
	"fmt"
)

type Mode struct {
	rootLettered string
	humanized    string
	Mode         KeySignatureMode
}

func NewMode(mode string) Mode {
	modeMapper := map[string]KeySignatureMode{
		"major":  KeySignatureModeMajor,
		"minor":  KeySignatureModeMinor,
		"dorian": KeySignatureModeDorian,
	}
	currentMode := modeMapper[mode]
	return Mode{
		Mode: currentMode,
	}
}

// GetRoot get the root of the mode based on the fifth
func (m *Mode) GetRoot(fifth int) string {
	if m.rootLettered != "" {
		return m.rootLettered
	}
	root := modeRoot[m.Mode.String()][fifth]
	m.rootLettered = root
	return root
}

func (m *Mode) GetHumanized(fifth int) string {
	if m.humanized != "" {
		return m.humanized
	}

	// <the mode> =  <the root key>
	// for example D Major would be
	// do = D
	humanized := fmt.Sprintf("%s = %s", m.Mode.GetNumberedRoot(), getSpelledNumberedNotation(m.GetRoot(fifth)))
	m.humanized = humanized

	return humanized
}

func (m *Mode) GetScaleSteps() []float64 {
	return modeSteps[m.Mode.String()]
}
