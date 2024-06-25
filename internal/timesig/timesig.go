package timesig

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type Time struct {
	Measure   int
	Beat      int
	BeatType  int
	notated   string
	Humanized string
}

func (t *Time) GetNotated() string {
	if t.notated != "" {
		return t.notated
	}

	n := fmt.Sprintf("%d/%d", t.Beat, t.BeatType)
	t.notated = n
	return n
}

func (t *Time) String() string {
	return t.GetNotated()
}

// calculateNoteLength in beat
func (t *Time) calculateNoteLength(ctx context.Context, note musicxml.Note) float64 {
	// cases

	baseLength := map[musicxml.NoteLength]float64{
		musicxml.NoteLengthQuarter: 1,
		musicxml.NoteLengthHalf:    2,
		musicxml.NoteLengthWhole:   4,
		musicxml.NoteLengthEighth:  0.5,
		musicxml.NoteLength16th:    0.25,
	}

	ratio := 1 // beat-type 4

	if t.BeatType == 8 {
		ratio = 2 // beat type 8
	}

	base := baseLength[note.Type] * float64(ratio)

	return base + (base * (1 - math.Pow(0.5, float64(len(note.Dot)))))
}

func (t Time) GetNoteLength(ctx context.Context, note musicxml.Note) float64 {
	return t.calculateNoteLength(ctx, note)
}

type TimeSignature struct {
	IsMixed    bool
	Signatures []Time
	humanized  string
}

func (ts *TimeSignature) GetHumanized() string {
	if ts.humanized != "" {
		return ts.humanized
	}

	if !ts.IsMixed {
		ts.humanized = fmt.Sprintf("%d ketuk", ts.Signatures[0].Beat)
		return ts.humanized
	}
	beat := map[int]bool{}
	combined := []string{}
	for _, v := range ts.Signatures {
		if !beat[v.Beat] {
			beat[v.Beat] = true
			combined = append(combined, fmt.Sprintf("%d", v.Beat))
		}
	}

	ts.humanized = strings.Join(combined, " dan ") + " ketuk"

	return ts.humanized
}

func (ts *TimeSignature) GetTimesignatureOnMeasure(ctx context.Context, measure int) Time {
	if len(ts.Signatures) == 1 {
		return ts.Signatures[0]
	}

	// get the time
	currentTime := ts.Signatures[0]

	counter := 0
	var prev Time
	prev = ts.Signatures[0]
	found := true

	for currentTime.Measure <= measure && counter < len(ts.Signatures) {
		prev = currentTime
		currentTime = ts.Signatures[counter]
		counter++

		if currentTime.Measure <= measure && counter == len(ts.Signatures) {
			found = false
		}

	}

	if !found {
		return ts.Signatures[len(ts.Signatures)-1]
	}

	return prev
}

func (ts *TimeSignature) GetNoteLength(ctx context.Context, measure int, note musicxml.Note) float64 {

	return (ts.GetTimesignatureOnMeasure(ctx, measure)).GetNoteLength(ctx, note)
}

func NewTimeSignatures(ctx context.Context, measures []musicxml.Measure) TimeSignature {

	times := []Time{}

	various := map[string]bool{}

	for _, measure := range measures {
		if measure.Attribute != nil && measure.Attribute.Time != nil {
			key := fmt.Sprintf("%d/%d", measure.Attribute.Time.Beats, measure.Attribute.Time.BeatType)
			various[key] = true
			times = append(times, Time{
				Measure:  measure.Number,
				Beat:     measure.Attribute.Time.Beats,
				BeatType: measure.Attribute.Time.BeatType,
			})
		}
	}

	return TimeSignature{
		IsMixed:    len(various) > 1,
		Signatures: times,
	}
}
