package lyric

import (
	"context"

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
	// if distance < hypenWidth {
	// 	// TODO: close the gap between these two lyric
	// 	return nil
	// }

	// every 20% of layout has 2 hypen
	container := constant.LAYOUT_WIDTH - (2-constant.LAYOUT_INDENT_LENGTH)/3
	_ = container
	// if distance < float64(container) {
	offset := (distance / 2) - hypenWidth
	return []entity.Coordinate{
		entity.Coordinate{
			X: startPosition + offset,
			Y: currentLyric.Coordinate.Y,
		},
	}
	// } else {
	// 	result := []entity.Coordinate{}
	// 	totalContainer := math.Floor(distance - (2*hypenWidth)/float64(container))
	// 	totalHypen := (totalContainer * 2) + 1

	// 	for i := float64(0); i < totalHypen; i++ {
	// 		result = append(result, entity.Coordinate{
	// 			X: startPosition + (i * distance),
	// 			Y: currentLyric.Coordinate.Y,
	// 		})
	// 	}

	// 	return result
	// }
}

func RenderHypen(ctx context.Context, canv canvas.Canvas, measure []*entity.NoteRenderer) {
	var curr, prev *LyricPosition
	_, _ = curr, prev
	hypenLocation := []entity.Coordinate{}
	for _, n := range measure {
		if len(n.Lyric) == 0 {
			continue
		}

		for i, l := range n.Lyric {
			switch l.Syllabic {
			case musicxml.LyricSyllabicTypeBegin:
				prev = &LyricPosition{
					Coordinate: entity.Coordinate{
						X: float64(n.PositionX),
						Y: float64(n.PositionY) + 25 + float64(i*20),
					},
					Lyrics: l,
				}
			case musicxml.LyricSyllabicTypeEnd:
				curr = &LyricPosition{
					Coordinate: entity.Coordinate{
						X: float64(n.PositionX),
						Y: float64(n.PositionY) + 25 + float64(i*20),
					},
					Lyrics: l,
				}
				// do the calculation here
				if prev != nil {
					hypenLocation = append(CalculateHypen(ctx, prev, curr), hypenLocation...)
				}
				prev, curr = nil, nil
			case musicxml.LyricSyllabicTypeMiddle:
				if prev == nil {
					prev = &LyricPosition{
						Coordinate: entity.Coordinate{
							X: float64(n.PositionX),
							Y: float64(n.PositionY) + 25 + float64(i*20),
						},
						Lyrics: l,
					}
					continue
				}
				if curr == nil {
					curr = &LyricPosition{
						Coordinate: entity.Coordinate{
							X: float64(n.PositionX),
							Y: float64(n.PositionY) + 25 + float64(i*20),
						},
						Lyrics: l,
					}
					hypenLocation = append(CalculateHypen(ctx, prev, curr), hypenLocation...)
					prev, curr = curr, nil
				}
			}
		}

	}

	for _, hl := range hypenLocation {
		canv.Text(int(hl.X), int(hl.Y), "-")
	}
}
