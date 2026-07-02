package params

import "github.com/jodi-ivan/numbered-notation-xml/internal/entity"

type PlaybackParams struct {
	Rect      map[int][][2]entity.Coordinate
	TotalBeat map[int][]float64
}
