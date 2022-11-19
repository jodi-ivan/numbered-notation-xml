package musicxml

import "encoding/xml"

type MusicXML struct {
	XMLName xml.Name `xml:"score-partwise"`
	Credit  Credit   `xml:"credit"`
	Part    Part     `xml:"part"`
	Work    Work     `xml:"work"`
}

type CreditType string

type Credit struct {
	Type  CreditType `xml:"credit-type"`
	Words string     `xml:"credit-words"`
}

type Work struct {
	Title string `xml:"work-title"`
}

type locationAttr struct {
	RelativeX float32 `xml:"relative-x,attr"`
	RelativeY float32 `xml:"relative-y,attr"`
	DefaultX  float32 `xml:"default-x,attr"`
	DefaultY  float32 `xml:"default-y,attr"`
}

type Part struct {
	ID       string    `xml:"id,attr"`
	Measures []Measure `xml:"measure"`
}

type Measure struct {
	Number    int        `xml:"number,attr"`
	Attribute *Attribute `xml:"attributes" json:",omitempty"`
	Notes     []Note     `xml:"note"`
	Direction []struct {
		Placement string `xml:"placement,attr"`
		Type      struct {
			Words struct {
				locationAttr
				Word string `xml:",chardata"`
			} `xml:"words"`
		} `xml:"direction-type"`
	} `xml:"direction" json:",omitempty"`
	Barline *Barline `xml:"barline" json:",omitempty"`
}

type KeySignature struct {
	Fifth int    `xml:"fifths"`
	Mode  string `xml:"mode"`
}

type Attribute struct {
	Key  KeySignature `xml:"key"`
	Time *struct {
		Beats    int `xml:"beats"`
		BeatType int `xml:"beat-type"`
	} `xml:"time"`
}

type NoteLength string

const (
	NoteLength256th   NoteLength = "256th"
	NoteLength128th   NoteLength = "128th"
	NoteLength64th    NoteLength = "64th"
	NoteLength32nd    NoteLength = "32nd"
	NoteLength16th    NoteLength = "16th"
	NoteLengthEighth  NoteLength = "eighth"
	NoteLengthQuarter NoteLength = "quarter"
	NoteLengthHalf    NoteLength = "half"
	NoteLengthWhole   NoteLength = "whole"
	NoteLengthBreve   NoteLength = "breve"
	NoteLengthLong    NoteLength = "long"
)

type NoteAccidental string

const (
	NoteAccidentalNatural     NoteAccidental = "natural"
	NoteAccidentalSharp       NoteAccidental = "sharp"
	NoteAccidentalFlat        NoteAccidental = "flat"
	NoteAccidentalDoubleSharp NoteAccidental = "double-sharp"
	NoteAccidentalDoubleFlat  NoteAccidental = "double-flat"
)

type LyricSyllabic string

const (
	LyricSyllabicTypeBegin  LyricSyllabic = "begin"
	LyricSyllabicTypeMiddle LyricSyllabic = "middle"
	LyricSyllabicTypeEnd    LyricSyllabic = "end"
	LyricSyllabicTypeSingle LyricSyllabic = "single"
)

type NoteSlurType string

const (
	NoteSlurTypeStart NoteSlurType = "start"
	NoteSlurTypeStop  NoteSlurType = "stop"
)

type NoteBeamType string

const (
	NoteBeamTypeBegin        NoteBeamType = "begin"
	NoteBeamTypeContinue     NoteBeamType = "continue"
	NoteBeamTypeEnd          NoteBeamType = "end"
	NoteBeamTypeForwardHook  NoteBeamType = "forward hook"
	NoteBeamTypeBackwardHook NoteBeamType = "backward hook"
)

func (na NoteAccidental) GetAccidental() string {
	sign := map[string]string{
		"natural":      "",
		"sharp":        "#",
		"flat":         "b",
		"double-sharp": "x",
		"double-flat":  "bb",
	}

	return sign[string(na)]
}

type Dot struct {
	Name xml.Name `xml:"dot"`
}

type Rest struct {
	Name xml.Name `xml:"rest"`
}

type Tie struct {
	Name xml.Name     `xml:"tied"`
	Type NoteSlurType `xml:"type,attr"`
}

type Note struct {
	Pitch struct {
		Step   string `xml:"step"`
		Octave int    `xml:"octave"`
	} `xml:"pitch"`
	Type       NoteLength     `xml:"type"`
	Beam       []*NoteBeam    `xml:"beam" json:",omitempty"`
	Notations  *NoteNotation  `xml:"notations" json:",omitempty"`
	Lyric      []Lyric        `xml:"lyric"`
	Accidental NoteAccidental `xml:"accidental"`
	Dot        []*Dot         `xml:"dot"`
	Rest       *Rest          `xml:"rest"`
}

type NoteBeam struct {
	Number int          `xml:"number,attr"`
	State  NoteBeamType `xml:",chardata"`
}

type NoteNotation struct {
	Slur         []NotationSlur        `xml:"slur" json:",omitempty"`
	Tied         *Tie                  `xml:"tied" json:",omitempty"`
	Articulation *NotationArticulation `xml:"articulations" json:",omitempty"`
}

type NotationArticulation struct {
	BreathMark *struct {
		Name xml.Name
	} `xml:"breath-mark"`
}
type NotationSlur struct {
	Type   NoteSlurType `xml:"type,attr"`
	Number int          `xml:"number,attr"`
}

type Lyric struct {
	locationAttr
	Number int `xml:"number,attr"`
	Text   struct {
		Underline int    `xml:"underline,attr"`
		Value     string `xml:",chardata"`
	} `xml:"text"`
	Syllabic LyricSyllabic `xml:"syllabic"`
}

type Barline struct {
	Location string `xml:"location,attr"`
	BarStyle string `xml:"bar-style"`
}
