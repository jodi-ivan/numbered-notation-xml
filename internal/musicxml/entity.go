package musicxml

import (
	"encoding/xml"
	"log"
	"strings"
)

type TextAlignment string

const (
	TextAlignmentRight TextAlignment = "right"
	TextAlignmentLeft  TextAlignment = "left"
)

type Element struct {
	// FIXME: mechanism to check wheter it is notes / direction
	// currently by checking pitch | rest -> notes
	// checking direction-type -> direction
	Content string `xml:",innerxml"`
}

func (e *Element) ParseAsNote() (Note, error) {
	wrapped := `<note>`
	wrapped += e.Content
	wrapped += `</note>`

	result := Note{}

	err := xml.Unmarshal([]byte(wrapped), &result)
	return result, err
}

func (e *Element) ParseAsDirection() (*Direction, error) {
	wrapped := `<direction>`
	wrapped += e.Content
	wrapped += `</direction>`

	result := &Direction{}

	err := xml.Unmarshal([]byte(wrapped), &result)
	return result, err
}

type MusicXML struct {
	XMLName xml.Name `xml:"score-partwise"`
	Credit  []Credit `xml:"credit"`
	Part    Part     `xml:"part"`
	Work    Work     `xml:"work"`
}

type CreditType string

const (
	CreditTypeTitle CreditType = "title"
)

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

type MeasureText struct {
	Text          string
	RelativeY     float64
	TextAlignment TextAlignment
}

type Measure struct {
	Appendix     []Element  `xml:",any"`
	Number       int        `xml:"number,attr"`
	Attribute    *Attribute `xml:"attributes" json:",omitempty"`
	Notes        []Note     `xml:"-" json:",omitempty"`
	Barline      []Barline  `xml:"barline" json:",omitempty"`
	Print        *Print     `xml:"print" json:",omitempty"`
	NewLineIndex int        `xml:"-"`

	// FIXME: one centralized place for the measured text
	RightMeasureText *MeasureText
}

func (m *Measure) Build() error {
	m.NewLineIndex = -1
	var measureText *MeasureText
	for i, elmnt := range m.Appendix {
		cleanedContent := strings.TrimSpace(elmnt.Content)
		if strings.HasPrefix(cleanedContent, "\u003cpitch\u003e") ||
			strings.Contains(cleanedContent, "\u003crest/\u003e") {
			n, err := elmnt.ParseAsNote()
			if err != nil {
				log.Println("error parsing note, err:", err.Error(), "\n\n", elmnt.Content)
				return err
			}

			if measureText != nil {
				if n.MeasureText == nil {
					n.MeasureText = []MeasureText{}
				}

				n.MeasureText = append(n.MeasureText, MeasureText{
					Text:      measureText.Text,
					RelativeY: measureText.RelativeY,
				})

				measureText = nil
			}
			m.Notes = append(m.Notes, n)
		} else if strings.HasPrefix(cleanedContent, "\u003cdirection-type\u003e") {
			d, err := elmnt.ParseAsDirection()
			if err != nil {
				return err
			}
			if d.DirectionType.Word.Value == "__layout=br" {
				m.NewLineIndex = i
			} else if d.DirectionType.Word.Value == "D.C. al Fine" {
				continue
			} else {

				measureText = &MeasureText{
					Text:      d.DirectionType.Word.Value,
					RelativeY: d.DirectionType.Word.RelativeY,
				}
			}
		}
	}

	if measureText != nil {
		m.RightMeasureText = &MeasureText{
			Text:      measureText.Text,
			RelativeY: measureText.RelativeY,
		}
	}

	return nil
}

type PrintNewSystemType string

const (
	PrintNewSystemTypeYes = "yes"
	PrintNewSystemTypeNo  = "no"
)

type Print struct {
	NewSystem PrintNewSystemType `xml:"new-system,attr"`
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

	// used for internal used only
	// indicates that the same note with the same slur number
	// as stop and start at the same time
	NoteSlurTypeHop NoteSlurType = "hop"
)

type NoteBeamType string

const (
	NoteBeamTypeBegin        NoteBeamType = "begin"
	NoteBeamTypeContinue     NoteBeamType = "continue"
	NoteBeamTypeEnd          NoteBeamType = "end"
	NoteBeamTypeForwardHook  NoteBeamType = "forward hook"
	NoteBeamTypeBackwardHook NoteBeamType = "backward hook"

	// additional beam. used for internal used for rendering the numbered notes
	// location 0
	NoteBeam_INTERNAL_TypeAdditional NoteBeamType = "additional"
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

type Direction struct {
	Placement     string        `xml:"placement,attr"`
	DirectionType DirectionType `xml:"direction-type"`
}

type DirectionType struct {
	Word struct {
		Value     string  `xml:",chardata"`
		RelativeY float64 `xml:"relative-y,attr"`
	} `xml:"words"`
}

type Bool string

const (
	BoolYes Bool = "yes"
	BoolNo  Bool = "no"
)

type ChardataInt struct {
	Value int `xml:",chardata"`
}

type TimeModification struct {
	Name        xml.Name    `xml:"time-modification"`
	ActualNotes ChardataInt `xml:"actual-notes"`
	NormalNotes ChardataInt `xml:"normal-notes"`
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

	TimeModification *TimeModification `xml:"time-modification"`

	MeasureText []MeasureText `xml:"-"`
}

type NoteBeam struct {
	Number int          `xml:"number,attr"`
	State  NoteBeamType `xml:",chardata"`
}

type TupletType string

const (
	TupletTypeStart TupletType = "start"
	TupletTypeStop  TupletType = "stop"
)

type Tuplet struct {
	Type   TupletType `xml:"type,attr"`
	Braket Bool       `xml:"braket,attr"`
}

type NoteNotation struct {
	Slur         []NotationSlur        `xml:"slur" json:",omitempty"`
	Tied         *Tie                  `xml:"tied" json:",omitempty"`
	Articulation *NotationArticulation `xml:"articulations" json:",omitempty"`

	// TODO: stacking tuplets / tuplet inside the tuplet
	Tuplet *Tuplet `xml:"tuplet"`
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
	Text   []struct {
		Underline int    `xml:"underline,attr"`
		Value     string `xml:",chardata"`
	} `xml:"text"`
	Syllabic LyricSyllabic `xml:"syllabic"`
}

type BarLineStyle string

const (
	BarLineStyleRegular    BarLineStyle = "regular"
	BarLineStyleLightHeavy BarLineStyle = "light-heavy"
	BarLineStyleLightLight BarLineStyle = "light-light"
	BarLineStyleHeavyHeavy BarLineStyle = "heavy-heavy"
	BarLineStyleHeavyLight BarLineStyle = "heavy-light"
)

type BarLineRepeatDirection string

const (
	BarLineRepeatDirectionBackward BarLineRepeatDirection = "backward"
	BarLineRepeatDirectionForward  BarLineRepeatDirection = "forward"
)

type BarlineLocation string

const (
	BarlineLocationLeft  = "left"
	BarlineLocationRight = "right"
)

type Barline struct {
	Location BarlineLocation `xml:"location,attr"`
	BarStyle BarLineStyle    `xml:"bar-style"`
	Repeat   *BarLineRepeat  `xml:"repeat"`
	Ending   *BarlineEnding  `xml:"ending"`
}

type BarLineRepeat struct {
	Direction BarLineRepeatDirection `xml:"direction,attr"`
}

type BarlineEndingType string

const (
	BarlineEndingTypeStart       = "start"
	BarlineEndingTypeStop        = "stop"
	BarlineEndingTypeDiscontinue = "discontinue"
)

type BarlineEnding struct {
	Name   xml.Name          `xml:"ending"`
	Number string            `xml:"number,attr"`
	Type   BarlineEndingType `xml:"type,attr"`
}
