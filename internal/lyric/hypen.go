package lyric

import (
	"context"
	"fmt"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (li *lyricInteractor) CalculateHypen(ctx context.Context, prevLyric, currentLyric *LyricPosition) (location []entity.Coordinate) {

	if prevLyric.Lyrics.Syllabic == musicxml.LyricSyllabicTypeEnd || prevLyric.Lyrics.Syllabic == musicxml.LyricSyllabicTypeSingle {
		return nil
	}

	hypenWidth := CalculateLyricWidth("-")

	lyricText := entity.LyricVal(prevLyric.Lyrics.Text).String()

	startPosition := prevLyric.Coordinate.X + CalculateLyricWidth(lyricText)
	endPostion := currentLyric.Coordinate.X
	distance := endPostion - startPosition
	if distance < hypenWidth {
		// TODO: close the gap between these two lyric
		return nil
	}

	// every 20% of layout has 3 hypen
	container := (constant.LAYOUT_WIDTH - (2 - constant.LAYOUT_INDENT_LENGTH)) / 6
	if distance < float64(container) {
		offset := (distance / 2) - hypenWidth
		if offset < 0 {
			offset = 0
		}
		return []entity.Coordinate{
			entity.Coordinate{
				X: startPosition + offset,
				Y: currentLyric.Coordinate.Y,
			},
		}
	} else {
		result := []entity.Coordinate{}
		totalContainer := math.Floor((distance - (2 * hypenWidth)) / float64(container))
		totalHypen := (totalContainer * 3)

		for i := float64(0); i < totalHypen; i++ {
			result = append(result, entity.Coordinate{
				X: startPosition + (i * (distance / totalHypen)),
				Y: currentLyric.Coordinate.Y,
			})
		}

		return result
	}
}

// measure is the notes for the whole staff
func (li *lyricInteractor) RenderHypen(ctx context.Context, canv canvas.Canvas, measure []*entity.NoteRenderer) {
	pos := map[int][2]*LyricPosition{}

	// for tracking the pair of begin to end
	hs := NewHypenStack()

	hypenLocation := []entity.Coordinate{}
	for notePos, n := range measure {

		// DONE: new line inside the measure
		// TODO: multi line lyrics
		if pos[0][0] != nil && (notePos == (len(measure) - 1)) && !hs.IsEmpty() {
			endNotePos := entity.Coordinate{
				X: float64(n.PositionX),
				Y: float64(n.PositionY) + 25,
			}

			hypenEnd := li.CalculateHypen(ctx, pos[0][0], &LyricPosition{
				Coordinate: endNotePos,
			})
			//add at the end of the lines
			if len(hypenEnd) != 1 {
				hypenEnd = append(hypenEnd, endNotePos)
			}

			hypenLocation = append(hypenEnd, hypenLocation...)
			continue
		}

		if len(n.Lyric) == 0 {
			continue
		}

		for i, l := range n.Lyric {
			hs.Process(ctx, l.Syllabic)
			if len(pos[i]) == 0 {
				pos[i] = [2]*LyricPosition{}
			}
			pair := pos[i]
			switch l.Syllabic {
			case musicxml.LyricSyllabicTypeBegin:
				// prev
				pair[0] = &LyricPosition{
					Coordinate: entity.Coordinate{
						X: float64(n.PositionX),
						Y: float64(n.PositionY) + 25 + float64(i*20),
					},
					Lyrics: l,
				}
			case musicxml.LyricSyllabicTypeEnd:
				// curr
				pair[1] = &LyricPosition{
					Coordinate: entity.Coordinate{
						X: float64(n.PositionX),
						Y: float64(n.PositionY) + 25 + float64(i*20),
					},
					Lyrics: l,
				}
				// do the calculation here
				start := pair[0]
				if start == nil {
					start = &LyricPosition{
						Coordinate: entity.Coordinate{
							X: float64(measure[0].PositionX),
							Y: float64(n.PositionY) + 25 + float64(i*20),
						},
					}
				}
				hypenLocation = append(li.CalculateHypen(ctx, start, pair[1]), hypenLocation...)
				pair = [2]*LyricPosition{}

			case musicxml.LyricSyllabicTypeMiddle:
				//TODO: no start hypen
				if pair[0] == nil {
					pair[0] = &LyricPosition{
						Coordinate: entity.Coordinate{
							X: float64(n.PositionX),
							Y: float64(n.PositionY) + 25 + float64(i*20),
						},
						Lyrics: l,
					}
					continue
				}
				if pair[1] == nil {
					pair[1] = &LyricPosition{
						Coordinate: entity.Coordinate{
							X: float64(n.PositionX),
							Y: float64(n.PositionY) + 25 + float64(i*20),
						},
						Lyrics: l,
					}
					hypenLocation = append(li.CalculateHypen(ctx, pair[0], pair[1]), hypenLocation...)
					pair = [2]*LyricPosition{
						pair[1],
						nil,
					}
				}
			}
			pos[i] = pair
		}

	}
	canv.Group("hyphens")
	for _, hl := range hypenLocation {
		fmt.Fprintf(canv.Writer(), `<text x="%.4f" y="%.0f">-</text>`, hl.X, hl.Y)
	}
	canv.Gend()
}
