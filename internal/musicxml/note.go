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

	MeasureText   []MeasureText `xml:"-"`
	IndexPosition int           `xml:"-"`
}

func (n Note) IsBreathMark() bool {
	return n.Notations != nil &&
		n.Notations.Articulation != nil &&
		n.Notations.Articulation.BreathMark != nil
}

type TextFontWeight string

var (
	TextFontWeightBold   TextFontWeight = "bold"
	TextFontWeightNormal TextFontWeight = "normal"
)

type TextFontStyle string

var (
	TextFontStyleItalic TextFontStyle = "bold"
	TextFontStyleNormal TextFontStyle = "normal"
)

type LyricText struct {
	Underline  int    `xml:"underline,attr"`
	Value      string `xml:",chardata"`
	FontWeight string `xml:"font-weight,attr"`
	FontStyle  string `xml:"font-style,attr"`
}

type Lyric struct {
	locationAttr
	Number   int           `xml:"number,attr"`
	Text     []LyricText   `xml:"text"`
	Syllabic LyricSyllabic `xml:"syllabic"`
}
