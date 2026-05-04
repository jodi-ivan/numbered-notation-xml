package musicxml

import (
	"log"
	"strings"
)

type RepeatInfoType string

const (
	RepeatInfoTypeOpening RepeatInfoType = "opening"
	RepeatInfoTypeMiddle  RepeatInfoType = "middle"
	RepeatInfoTypeClosing RepeatInfoType = "closing"
	RepeatInfoTypeBoth    RepeatInfoType = "both"
)

type RepeatInfo struct {
	Type                 RepeatInfoType
	SyllableCount        int
	StartPosition        int
	OffsetSyllable       int
	StartIndex           int
	SectionSyllableCount int
}

type Measure struct {
	Appendix     []Element    `xml:",any"`
	Number       int          `xml:"number,attr"`
	Attribute    *Attribute   `xml:"attributes" json:",omitempty"`
	Notes        []Note       `xml:"-" json:",omitempty"`
	Barline      []Barline    `xml:"barline" json:",omitempty"`
	Print        *Print       `xml:"print" json:",omitempty"`
	NewLineIndex map[int]bool `xml:"-"`

	DirectionDashes map[int]map[int]DirectionDashesType `xml:"-"`

	// FIXME: one centralized place for the measured text
	RightMeasureText *MeasureText   `xml:"-"`
	PrefixHeader     map[int]string `xml:"-"`
	RepeatInfo       *RepeatInfo
}

func (m *Measure) Build() error {
	m.NewLineIndex = map[int]bool{}
	var measureText []MeasureText
	foundDirectionType := 0
	for i, elmnt := range m.Appendix {
		cleanedContent := strings.TrimSpace(elmnt.Content)
		if strings.HasPrefix(cleanedContent, "\u003cpitch\u003e") ||
			strings.Contains(cleanedContent, "\u003crest /\u003e") ||
			strings.Contains(cleanedContent, "\u003crest/\u003e") {
			n, err := elmnt.ParseAsNote()
			if err != nil {
				log.Println("error parsing note, err:", err.Error(), "\n\n", elmnt.Content)
				return err
			}

			n.IndexPosition = i //+ m.StartIndex

			if len(measureText) > 0 {
				if n.MeasureText == nil {
					n.MeasureText = []MeasureText{}
				}

				n.MeasureText = append(n.MeasureText, measureText...)
				measureText = []MeasureText{}
			}
			m.Notes = append(m.Notes, n)
		} else if strings.HasPrefix(cleanedContent, "\u003cdirection-type\u003e") {
			d, err := elmnt.ParseAsDirection()
			if err != nil {
				return err
			}
			if len(d.DirectionType) == 0 {
				continue
			}
			initalDirection := d.DirectionType[0]

			if initalDirection.Word.Value == "__layout=br" {
				m.NewLineIndex[i-foundDirectionType] = true
				foundDirectionType++
			} else if initalDirection.Word.Value == "D.C. al Fine" {
				continue
			} else {
				measureText = append(measureText, MeasureText{
					Text:      initalDirection.Word.Value,
					RelativeY: initalDirection.Word.RelativeY,
				})
			}

			if initalDirection.Rehearshal != nil {
				if m.PrefixHeader == nil {
					m.PrefixHeader = map[int]string{}
				}
				m.PrefixHeader[i-foundDirectionType] = initalDirection.Rehearshal.Value
			}

			if len(d.DirectionType) == 2 && d.DirectionType[1].Dashes != nil || (len(d.DirectionType) == 1 && d.DirectionType[0].Dashes != nil) {
				pos := 1
				if d.DirectionType[0].Dashes != nil {
					pos = 0
				}
				direction := d.DirectionType[pos].Dashes

				if _, ok := m.DirectionDashes[i]; !ok {
					if m.DirectionDashes == nil {
						m.DirectionDashes = map[int]map[int]DirectionDashesType{}
					}
					m.DirectionDashes[i] = map[int]DirectionDashesType{}
				}

				m.DirectionDashes[i][direction.Number] = direction.Type
			}
		}
	}

	if len(measureText) > 0 {
		m.RightMeasureText = &measureText[0]
	}

	return nil
}
