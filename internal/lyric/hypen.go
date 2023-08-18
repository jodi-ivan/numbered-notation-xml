package lyric

import (
	"context"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func CalculateHypen(ctx context.Context, prevLyric, currentLyric *LyricPosition) (location []entity.Coordinate) {

	if prevLyric.Lyrics.Syllabic == musicxml.LyricSyllabicTypeEnd || prevLyric.Lyrics.Syllabic == musicxml.LyricSyllabicTypeSingle {
		return nil
	}

	hypenWidth := CalculateLyricWidth("-")
	/*
		          end position
		               v
			ly    -     ric
			  ^
		 start poition
	*/

	startPosition := prevLyric.Coordinate.X + CalculateLyricWidth(prevLyric.Lyrics.Text)
	endPostion := currentLyric.Coordinate.X
	distance := endPostion - startPosition
	if distance < hypenWidth {
		// TODO: close the gap between these two lyric
		return nil
	}

	// every 20% of layout has 3 hypen
	container := (constant.LAYOUT_WIDTH - (2 - constant.LAYOUT_INDENT_LENGTH)) / 5
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

// TODO: last staff syllable that has no end or has no start
func RenderHypen(ctx context.Context, canv canvas.Canvas, measure []*entity.NoteRenderer) {
	pos := map[int][2]*LyricPosition{}

	hypenLocation := []entity.Coordinate{}
	for _, n := range measure {
		if len(n.Lyric) == 0 {
			continue
		}

		for i, l := range n.Lyric {
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
				if pair[0] != nil {
					hypenLocation = append(CalculateHypen(ctx, pair[0], pair[1]), hypenLocation...)
				}
				pair = [2]*LyricPosition{}

			case musicxml.LyricSyllabicTypeMiddle:
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
					hypenLocation = append(CalculateHypen(ctx, pair[0], pair[1]), hypenLocation...)
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
		canv.Text(int(hl.X), int(hl.Y), "-")
	}
	canv.Gend()
}
