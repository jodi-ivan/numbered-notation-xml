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

type Attribute struct {
	Key struct {
		Fifth int    `xml:"fifths"`
		Mode  string `xml:"mode"`
	} `xml:"key"`
	Time struct {
		Beats    int `xml:"beats"`
		BeatType int `xml:"beat-type"`
	} `xml:"time"`
}

type NoteType string

type Note struct {
	Pitch struct {
		Step   string `xml:"step"`
		Octave int    `xml:"octave"`
	} `xml:"pitch"`
	Type      NoteType      `xml:"type"`
	Beam      *NoteBeam     `xml:"beam" json:",omitempty"`
	Notations *NoteNotation `xml:"notations" json:",omitempty"`
	Lyric     []Lyric       `xml:"lyric"`
}

type NoteBeam struct {
	Number int    `xml:"number,attr"`
	State  string `xml:",chardata"`
}

type NoteNotation struct {
	Slur         *NotationSlur         `xml:"slur" json:",omitempty"`
	Articulation *NotationArticulation `xml:"articulations" json:",omitempty"`
}

type NotationArticulation struct {
	BreathMark *struct {
		Name xml.Name
	} `xml:"breath-mark"`
}
type NotationSlur struct {
	Type   string `xml:"type,attr"`
	Number int    `xml:"number,attr"`
}

type Lyric struct {
	locationAttr
	Number int `xml:"number,attr"`
	Text   struct {
		Underline int    `xml:"underline,attr"`
		Value     string `xml:",chardata"`
	} `xml:"text"`
	Syllabic string `xml:"syllabic"`
}

type Barline struct {
	Location string `xml:"location,attr"`
	BarStyle string `xml:"bar-style"`
}
