package playback

import (
	"time"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
)

type Step struct {
	Duration      time.Duration `json:"duration"`
	DurationStr   string        `json:"duration_str"`
	LyricPart     int           `json:"lyric_part"`
	MeasureNumber int           `json:"measure_number"`
	TotalBeat     float64       `json:"total_beat"`

	Rect [2]entity.Coordinate `json:"rect"`
}
