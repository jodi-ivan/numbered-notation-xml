package adapter

import (
	"fmt"
	"net/http"

	svg "github.com/ajstarks/svgo"

	"github.com/julienschmidt/httprouter"
)

const lowecaseLength = 30
const upperCaseLegth = 30

type NoteGoup struct {
	Notes []Note
}
type Note struct {
	Numbered int64
	Octave   int
	IsDotted bool
	Length   int
}

type Lyric struct {
	Text      string
	IsEndword bool
}

type Row struct {
	NoteGroups []NoteGoup
	Lyrics     map[int][]Lyric
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
	// Text(x int, y int, t string, s ...string)

	for i, d := range aRow.NoteGroups {
		if i == 0 {
			if len(d.Notes) == 1 {
				s.Text(100, 100, fmt.Sprint(d.Notes[0].Numbered))
			} else {
				for j, g := range d.Notes {
					s.Text(100+(lowecaseLength*j), 100, fmt.Sprint(g.Numbered))
				}
			}
		} else {
			if len(d.Notes) == 1 {
				s.Text(100+(lowecaseLength*i*len(aRow.Lyrics[1][i].Text)), 100, fmt.Sprintf(" %d", d.Notes[0].Numbered))
			} else {
				for j, g := range d.Notes {
					if j == 0 {
						text := fmt.Sprint(g.Numbered)
						if j == 0 {
							text = fmt.Sprintf(" %s", text)
						}
						s.Text(100+((j+1)*(lowecaseLength*i*len(aRow.Lyrics[1][i].Text))), 100, text)
					} else {
						s.Text(115+(j*(lowecaseLength*i*len(aRow.Lyrics[1][i].Text))), 100, fmt.Sprint(g.Numbered))

					}
				}
			}
		}
	}

	for i, d := range aRow.Lyrics[1] {
		text := d.Text
		if !d.IsEndword {
			text = fmt.Sprintf("%s -", text)
		}
		if i == 0 {
			s.Text(100, 120, text)
		} else {
			s.Text(100+(lowecaseLength*i*len(aRow.Lyrics[1][i].Text)), 120, text)
		}
	}
	s.End()
}
