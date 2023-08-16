package lyric

import (
	"context"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func CalculateHypen(ctx context.Context, prevLyric, currentLyric LyricPosition) (location []entity.Coordinate) {

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

	// every 20% of layout has 2 hypen
	container := constant.LAYOUT_WIDTH - (2-constant.LAYOUT_INDENT_LENGTH)/5
	if distance < float64(container) {
		hypenPos := (startPosition + (distance - hypenWidth)) / 2
		return []entity.Coordinate{
			entity.Coordinate{
				X: startPosition + hypenPos,
				Y: currentLyric.Coordinate.Y,
			},
		}
	} else {
		result := []entity.Coordinate{}
		totalContainer := math.Floor(distance - (2*hypenWidth)/float64(container))
		totalHypen := (totalContainer * 2) + 1

		for i := float64(0); i < totalHypen; i++ {
			result = append(result, entity.Coordinate{
				X: startPosition + (i * distance),
				Y: currentLyric.Coordinate.Y,
			})
		}

		return result
	}
}
