package timesignature

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type Time struct {
	Beat      int
	BeatType  int
	Notated   string
	Humanized string
}

type TimeSignature struct {
	IsMixed        bool
	Signatures     map[string]Time
	MeasureAddress map[string][]int64
}

func NewTimeSignatures(ctx context.Context, measures []musicxml.Measure) TimeSignature {
	return TimeSignature{}
}
