package musicxml

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

	TimeModification *TimeModification `xml:"time-modification"`

	MeasureText []MeasureText `xml:"-"`
}

func (n Note) IsBreathMark() bool {
	return n.Notations != nil &&
		n.Notations.Articulation != nil &&
		n.Notations.Articulation.BreathMark != nil
}
