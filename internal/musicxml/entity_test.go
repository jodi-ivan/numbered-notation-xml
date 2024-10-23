package musicxml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeasure_Build(t *testing.T) {

	pitch := `
	<pitch>
		<step>G</step>
		<octave>4</octave>
	</pitch>
	<duration>4</duration>
	<voice>1</voice>
	<type>half</type>
	<stem>up</stem>
	<lyric number="1">
		<syllabic>end</syllabic>
		<text>rap</text>
	</lyric>
	`

	direction := `
	<direction-type>
	<words>Refrein</words>
	</direction-type>`
	_ = pitch
	_ = direction
	tests := []struct {
		name        string
		m           func() *Measure
		wantMeasure *Measure
		wantErr     bool
	}{
		{
			name: "everything went fine",
			m: func() *Measure {
				return &Measure{
					Appendix: []Element{
						Element{
							Content: pitch,
						},
						Element{
							Content: direction,
						},
					},
				}
			},
			wantMeasure: &Measure{
				Appendix: []Element{
					Element{
						Content: pitch,
					},
					Element{
						Content: direction,
					},
				},
				Notes: []Note{
					Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "G", Octave: 4},
						Type: NoteLengthHalf,
						Lyric: []Lyric{
							Lyric{
								Number:   1,
								Syllabic: LyricSyllabicTypeEnd,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "rap",
									},
								},
							},
						},
					},
				},
				RightMeasureText: &MeasureText{
					Text: "Refrein",
				},
				NewLineIndex: -1,
			},
		},
		{
			name: "Failed to parse note",
			m: func() *Measure {
				return &Measure{
					Appendix: []Element{
						Element{
							Content: "<pitch>Nope",
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Has measure text between notes",
			m: func() *Measure {
				return &Measure{
					Appendix: []Element{
						Element{
							Content: pitch,
						},
						Element{
							Content: direction,
						},
						Element{
							Content: pitch,
						},
					},
				}
			},
			wantMeasure: &Measure{
				Appendix: []Element{
					Element{
						Content: pitch,
					},
					Element{
						Content: direction,
					},
					Element{
						Content: pitch,
					},
				},
				Notes: []Note{
					Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "G", Octave: 4},
						Type: NoteLengthHalf,
						Lyric: []Lyric{
							Lyric{
								Number:   1,
								Syllabic: LyricSyllabicTypeEnd,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "rap",
									},
								},
							},
						},
					},
					Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "G", Octave: 4},
						Type: NoteLengthHalf,
						Lyric: []Lyric{
							Lyric{
								Number:   1,
								Syllabic: LyricSyllabicTypeEnd,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "rap",
									},
								},
							},
						},
						MeasureText: []MeasureText{
							MeasureText{
								Text: "Refrein",
							},
						},
					},
				},
				NewLineIndex: -1,
			},
		},
		{
			name: "Failed to parse direction",
			m: func() *Measure {
				return &Measure{
					Appendix: []Element{
						Element{
							Content: "<direction-type>Nope",
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "New Line bartext",
			m: func() *Measure {
				return &Measure{
					Appendix: []Element{
						Element{
							Content: pitch,
						},
						Element{
							Content: `<direction-type><words>__layout=br</words></direction-type>`,
						},
						Element{
							Content: pitch,
						},
					},
				}
			},
			wantMeasure: &Measure{
				Appendix: []Element{
					Element{
						Content: pitch,
					},
					Element{
						Content: `<direction-type><words>__layout=br</words></direction-type>`,
					},
					Element{
						Content: pitch,
					},
				},
				Notes: []Note{
					Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "G", Octave: 4},
						Type: NoteLengthHalf,
						Lyric: []Lyric{
							Lyric{
								Number:   1,
								Syllabic: LyricSyllabicTypeEnd,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "rap",
									},
								},
							},
						},
					},
					Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "G", Octave: 4},
						Type: NoteLengthHalf,
						Lyric: []Lyric{
							Lyric{
								Number:   1,
								Syllabic: LyricSyllabicTypeEnd,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "rap",
									},
								},
							},
						},
					},
				},
				NewLineIndex: 1,
			},
		},
		{
			name: "D.C. al Fine",
			m: func() *Measure {
				return &Measure{
					Appendix: []Element{
						Element{
							Content: pitch,
						},
						Element{
							Content: `<direction-type><words>D.C. al Fine</words></direction-type>`,
						},
						Element{
							Content: pitch,
						},
					},
				}
			},
			wantMeasure: &Measure{
				Appendix: []Element{
					Element{
						Content: pitch,
					},
					Element{
						Content: `<direction-type><words>D.C. al Fine</words></direction-type>`,
					},
					Element{
						Content: pitch,
					},
				},
				Notes: []Note{
					Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "G", Octave: 4},
						Type: NoteLengthHalf,
						Lyric: []Lyric{
							Lyric{
								Number:   1,
								Syllabic: LyricSyllabicTypeEnd,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "rap",
									},
								},
							},
						},
					},
					Note{
						Pitch: struct {
							Step   string `xml:"step"`
							Octave int    `xml:"octave"`
						}{Step: "G", Octave: 4},
						Type: NoteLengthHalf,
						Lyric: []Lyric{
							Lyric{
								Number:   1,
								Syllabic: LyricSyllabicTypeEnd,
								Text: []struct {
									Underline int    `xml:"underline,attr"`
									Value     string `xml:",chardata"`
								}{
									{
										Value: "rap",
									},
								},
							},
						},
					},
				},
				NewLineIndex: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.m()
			if err := m.Build(); (err != nil) != tt.wantErr {
				t.Errorf("Measure.Build() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				assert.Equal(t, tt.wantMeasure, m)
			}
		})
	}
}

func TestNoteAccidental_GetAccidental(t *testing.T) {
	tests := []struct {
		name string
		na   NoteAccidental
		want string
	}{
		{
			name: "natural",
			na:   NoteAccidental("natural"),
			want: "",
		},
		{
			name: "sharp",
			na:   NoteAccidental("sharp"),
			want: "#",
		},
		{
			name: "flat",
			na:   NoteAccidental("flat"),
			want: "b",
		},
		{
			name: "double-sharp",
			na:   NoteAccidental("double-sharp"),
			want: "x",
		},
		{
			name: "double-flat",
			na:   NoteAccidental("double-flat"),
			want: "bb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.na.GetAccidental(); got != tt.want {
				t.Errorf("NoteAccidental.GetAccidental() = %v, want %v", got, tt.want)
			}
		})
	}
}
