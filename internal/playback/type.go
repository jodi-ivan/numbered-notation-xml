package playback

import (
	"time"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
)

type Step struct {
	Duration      time.Duration `json:"duration"`
	LyricPart     int           `json:"lyric_part"`
	MeasureNumber int           `json:"measure_number"`

	Rect [2]entity.Coordinate `json:"rect"`
}
