package keysig

type KeySignatureMode int

const (
	KeySignatureModeMajor    KeySignatureMode = 0
	KeySignatureModeMinor    KeySignatureMode = 1
	KeySignatureModeDorian   KeySignatureMode = 2
	KeySignatureModePhrygian KeySignatureMode = 3
)

func (ksm KeySignatureMode) String() string {
	return []string{"major", "minor", "dorian", "phrygian"}[int(ksm)]
}

func (ksm KeySignatureMode) GetNumberedRoot() string {
	return []string{"do", "la", "re", "mi"}[int(ksm)]
}

var accidentalsSet = map[int][]string{
	7:  []string{"C", "D", "E", "F", "G", "A", "B"},
	6:  []string{"F", "G", "A", "C", "D", "E"},
	5:  []string{"C", "D", "F", "G", "A"},
	4:  []string{"F", "G", "C", "D"},
	3:  []string{"C", "F", "G"},
	2:  []string{"F", "C"},
	1:  []string{"F"},
	0:  []string{},
	-1: []string{"B"},
	-2: []string{"B", "E"},
	-3: []string{"E", "A", "B"},
	-4: []string{"A", "B", "D", "E"},
	-5: []string{"D", "E", "G", "A", "B"},
	-6: []string{"G", "A", "B", "C", "D", "E"},
	-7: []string{"G", "A", "B", "C", "D", "E"},
}

var modeRoot = map[string]map[int]string{
	"major": map[int]string{
		7:  "C#",
		6:  "F#",
		5:  "B",
		4:  "E",
		3:  "A",
		2:  "D",
		1:  "G",
		0:  "C",
		-1: "F",
		-2: "Bb",
		-3: "Eb",
		-4: "Ab",
		-5: "Db",
		-6: "Gb",
		-7: "Cb",
	},
	"minor": map[int]string{
		7:  "A#",
		6:  "D#",
		5:  "G#",
		4:  "C#",
		3:  "F#",
		2:  "B",
		1:  "E",
		0:  "A",
		-1: "D",
		-2: "G",
		-3: "C",
		-4: "F",
		-5: "Bb",
		-6: "Eb",
		-7: "Ab",
	},
	"dorian": map[int]string{
		// FIXME: find the 6th and the 7th of the circle of  fifth
		5:  "C#",
		4:  "F#",
		3:  "B",
		2:  "E",
		1:  "A",
		0:  "D",
		-1: "G",
		-2: "C",
		-3: "F",
		-4: "Bb",
		-5: "Eb",
	},
	"phrygian": map[int]string{
		-5: "F",
		-4: "C",
		-3: "G",
		-2: "D",
		-1: "A",
		0:  "E",
		1:  "B",
		2:  "F#",
		3:  "C#",
		4:  "G#",
		5:  "D#",
		6:  "A#",
	},
}

var modeSteps = map[string][]float64{
	"major": []float64{ //
		1,   // do -> re
		1,   // re -> mi
		0.5, // mi -> fa
		1,   // fa -> sol
		1,   // sol -> la
		1,   // la -> si (ti)
		0.5, // si -> do
	},
	"minor": []float64{ //
		1,   // do -> re
		0.5, // re -> mi
		1,   // mi -> fa
		1,   // fa -> sol
		0.5, // sol -> la
		1,   // la -> si (ti)
		1,   // si -> do
	},
	"dorian": []float64{
		1,   // do -> re
		0.5, // re-> mi
		1,   // mi -> fa
		1,   // fa -> sol
		1,   // sol -> la
		0.5, // la -> si
		1,   // si -> do
	},
	"phrygian": []float64{
		0.5, // do -> re (The defining half step)
		1,   // re -> mi
		1,   // mi -> fa
		1,   // fa -> sol
		0.5, // sol -> la
		1,   // la -> si (ti)
		1,   // si -> do
	},
}
