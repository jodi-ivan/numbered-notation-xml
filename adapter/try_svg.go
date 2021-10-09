package adapter

import (
	"fmt"
	"net/http"
	"strings"

	svg "github.com/ajstarks/svgo"

	"github.com/julienschmidt/httprouter"
)

const lowecaseLength = 15
const upperCaseLegth = 30

type NoteTies struct {
	StartIndex int
	EndIndex   int
}
type NoteGoup struct {
	Notes  []Note
	Ties   []NoteTies
	Dotted []int
}
type Note struct {
	Numbered  int64
	Octave    int
	NumDotted int
	Length    int
}

type Lyric struct {
	Text      string
	IsEndword bool
}

type Row struct {
	NoteGroups []NoteGoup
	Lyrics     map[int][]Lyric
}

type Position struct {
	X int
	Y int
}

type GroupPosition struct {
	Start Position
	End   Position
}

type RowGroupLocator struct {
	NoteGroupLoc GroupPosition
	LyricLoc     GroupPosition
}

type TrySVG struct {
}

func (ts *TrySVG) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	aRow := Row{
		NoteGroups: []NoteGoup{
			NoteGoup{
				Notes: []Note{
					Note{
						Numbered: 3,
					},
					Note{
						Numbered: 6,
					},
					Note{
						Numbered:  7,
						NumDotted: 2,
					},
					Note{
						Numbered: 3,
					},
					Note{
						Numbered: 6,
					},
					Note{
						Numbered: 7,
					},
				},
				Ties: []NoteTies{
					NoteTies{
						StartIndex: 0,
						EndIndex:   4,
					},
					NoteTies{
						StartIndex: 0,
						EndIndex:   5,
					},
				},
			},
			NoteGoup{
				Notes: []Note{
					Note{
						Numbered: 5,
					},
					Note{
						Numbered: 6,
					},
					Note{
						Numbered: 7,
					},
					Note{
						Numbered: 6,
					},
				},
				Ties: []NoteTies{
					NoteTies{
						StartIndex: 0,
						EndIndex:   1,
					},
				},
			},
			NoteGoup{
				Notes: []Note{
					Note{
						Numbered: 3,
					},
				},
			},
			NoteGoup{
				Notes: []Note{
					Note{
						Numbered: 1,
					},
					Note{
						Numbered: 2,
					},
				},
			},
		},
		Lyrics: map[int][]Lyric{
			1: []Lyric{
				Lyric{
					Text:      "ha",
					IsEndword: false,
				},
				Lyric{
					Text:      "le",
					IsEndword: false,
				},
				Lyric{
					Text:      "lu",
					IsEndword: false,
				},
				Lyric{
					Text:      "ya",
					IsEndword: true,
				},
			},
		},
	}

	_ = aRow
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(200)
	s := svg.New(w)
	s.Start(500, 500)

	verse := 1

	lengthLyric := make([]int64, len(aRow.NoteGroups))
	lengthNote := make([]int64, len(aRow.NoteGroups))

	for i, d := range aRow.NoteGroups {
		currLyricLength := len(aRow.Lyrics[verse][i].Text)
		lengthLyric[i] = int64(currLyricLength)

		currNoteLength := len(d.Notes) + len(d.Dotted)
		if len(d.Notes) > 1 {
			// count the space
			currNoteLength = currNoteLength + (len(d.Notes) - 1)
		}

		if len(d.Dotted) > 1 {
			// count the space
			currNoteLength = currNoteLength + (len(d.Dotted) - 1)
		}
		lengthNote[i] = int64(currNoteLength)

		fmt.Println("Note", i, "lyric", currLyricLength, "notes", currNoteLength)
	}

	initialX := 100
	initialY := 100
	offset := 0

	for i, d := range aRow.NoteGroups {
		initialX = initialX + offset

		if len(d.Notes) == 1 {
			s.Text(initialX, initialY, fmt.Sprintf("%d %s", d.Notes[0].Numbered, strings.Repeat(". ", d.Notes[0].NumDotted)))
			initialX += 9 + (8 * d.Notes[0].NumDotted)
		} else {
			locX := map[int]int{}
			for j, g := range d.Notes {
				s.Text(initialX+(lowecaseLength*j), initialY, fmt.Sprintf("%d %s", g.Numbered, strings.Repeat(". ", g.NumDotted)))
				initialX += (8 * g.NumDotted)
				locX[j] = initialX + (lowecaseLength * j) + (8 * g.NumDotted)
			}
			for i, tie := range d.Ties {
				path := fmt.Sprintf("M%f,%d,C%f,%d %f,%d %f,%d", float64(locX[tie.StartIndex])+2.5, initialY+5+(i*5), float64(locX[tie.StartIndex])+2.5, initialY+15+(i*5), float64(locX[tie.EndIndex])+3.5, initialY+15+(i*5), float64(locX[tie.EndIndex])+3.5, initialY+5+(i*5))
				s.Path(path, `fill-opacity="0"`, `stroke="#000"`, `stroke-width="1"`)
			}
		}

		offset = int(lengthNote[i]) * 9
		if offset < int(lengthLyric[i]) {
			offset = int(lengthLyric[i]) * 9
		}
	}

	// for i, d := range aRow.Lyrics[1] {
	// 	text := d.Text
	// 	if !d.IsEndword {
	// 		text = fmt.Sprintf("%s -", text)
	// 	}
	// 	if i == 0 {
	// 		s.Text(100+int(lengthLyric[i]), 125, text)
	// 	} else {
	// 		s.Text(100+(lowecaseLength*i*len(aRow.Lyrics[1][i].Text)), 125, text)
	// 	}
	// }

	// } else {
	// 	locX := map[int]int{}
	// 	for j, g := range d.Notes {
	// 		// s.Text(100+(lowecaseLength*j), 100, fmt.Sprint(g.Numbered))
	// 		locX[j] = 100 + (lowecaseLength * j)
	// 	}
	// 	for _, tie := range d.Ties {
	// 		path := fmt.Sprintf("M%f,%d,C%f,%d %f,%d %f,%d", float64(locX[tie.StartIndex])+2.5, 105, float64(locX[tie.StartIndex])+2.5, 110, float64(locX[tie.EndIndex])+3.5, 110, float64(locX[tie.EndIndex])+3.5, 105)
	// 		s.Path(path, `fill-opacity="0"`, `stroke="#000"`, `stroke-width="1"`)
	// 	}
	// }

	// Text(x int, y int, t string, s ...string)

	// for i, d := range aRow.NoteGroups {
	// 	if i == 0 {
	// 		if len(d.Notes) == 1 {
	// 			s.Text(100, 100, fmt.Sprint(d.Notes[0].Numbered))
	// 		} else {
	// 			locX := map[int]int{}
	// 			for j, g := range d.Notes {
	// 				s.Text(100+(lowecaseLength*j), 100, fmt.Sprint(g.Numbered))
	// 				locX[j] = 100 + (lowecaseLength * j)
	// 			}
	// 			// for _, tie := range d.Ties {
	// 			// 	// path := fmt.Sprintf("M100,200,C100,300 400,300 400,200")
	// 			// 	path := fmt.Sprintf("M%f,%d,C%f,%d %f,%d %f,%d", float64(locX[tie.StartIndex])+2.5, 105, float64(locX[tie.StartIndex])+2.5, 110, float64(locX[tie.EndIndex])+3.5, 110, float64(locX[tie.EndIndex])+3.5, 105)
	// 			// 	s.Path(path, `fill-opacity="0"`, `stroke="#000"`, `stroke-width="1"`)
	// 			// }
	// 		}
	// 	} else {
	// 		if len(d.Notes) == 1 {
	// 			s.Text(100+(lowecaseLength*i*len(aRow.Lyrics[1][i].Text)), 100, fmt.Sprintf(" %d", d.Notes[0].Numbered))
	// 		} else {
	// 			// locX := map[int]int{}
	// 			for j, g := range d.Notes {
	// 				if j == 0 {
	// 					text := fmt.Sprint(g.Numbered)
	// 					if j == 0 {
	// 						text = fmt.Sprintf(" %s", text)
	// 					}
	// 					s.Text(100+(lowecaseLength*i), 100, text)
	// 				} else {
	// 					s.Text(115+(j*(lowecaseLength*i)), 100, fmt.Sprint(g.Numbered))
	// 				}
	// 			}

	// 			// for _, tie := range d.Ties {
	// 			// 	// path := fmt.Sprintf("M100,200,C100,300 400,300 400,200")
	// 			// 	path := fmt.Sprintf("M%f,%d,C%f,%d %f,%d %f,%d", float64(locX[tie.StartIndex])+2.5, 105, float64(locX[tie.StartIndex])+2.5, 110, float64(locX[tie.EndIndex])+3.5, 110, float64(locX[tie.EndIndex])+3.5, 105)
	// 			// 	s.Path(path, `fill-opacity="0"`, `stroke="#000"`, `stroke-width="1"`)
	// 			// }
	// 		}
	// 	}
	// }

	// for i, d := range aRow.Lyrics[1] {
	// 	text := d.Text
	// 	if !d.IsEndword {
	// 		text = fmt.Sprintf("%s -", text)
	// 	}
	// 	if i == 0 {
	// 		s.Text(100, 125, text)
	// 	} else {
	// 		s.Text(100+(lowecaseLength*i*len(aRow.Lyrics[1][i].Text)), 125, text)
	// 	}
	// }
	s.End()
}
