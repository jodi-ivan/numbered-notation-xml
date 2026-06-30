package playback

type Step struct {
	// Duration      time.Duration `json:"duration,omitempty"`
	LyricPart     int `json:"lyric_part,omitempty"`
	MeasureNumber int `json:"measure_number"`

	// Rect [2]entity.Coordinate `json:"rect,omitempty"`
}
