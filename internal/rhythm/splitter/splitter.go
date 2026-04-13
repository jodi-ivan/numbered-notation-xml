package splitter

import (
	"context"
	"log"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

type BeamSplitter interface {
	Split(ctx context.Context, notes []*entity.NoteRenderer, ts timesig.TimeSignature, segments map[int][]BeamSplitMarker)
	SplitSingle(ctx context.Context, notes []*entity.NoteRenderer, ts timesig.TimeSignature, segments []BeamSplitMarker, beamNo int)
}

type beamSplitter struct {
	eighth  eighthSplitter
	quarter quaterSpliiter
}

func New() BeamSplitter {
	return &beamSplitter{}
}

func (bs *beamSplitter) Split(ctx context.Context, notes []*entity.NoteRenderer, ts timesig.TimeSignature, _ map[int][]BeamSplitMarker) {
	if len(notes) == 0 {
		return
	}
	measureNumber := notes[0].MeasureNumber

	beamSegments := map[int][]BeamSplitMarker{}

	beamSegments[1] = CleanBeamByNumber(ctx, notes, 1)
	beamSegments[2] = CleanBeamByNumber(ctx, notes, 2)

	currentTimesig := ts.GetTimesignatureOnMeasure(ctx, measureNumber)
	beamSegments[1] = splitTuplet(notes, beamSegments[1])
	switch currentTimesig.BeatType {
	case 4, 2:
		bs.quarter.Split(ctx, notes, ts, beamSegments)
	case 8:
		bs.eighth.Split(ctx, notes, ts, beamSegments)
	default:
		log.Printf("[rhythm] Unable to split, On measure: %d: Unsupported beat-type %d.  \n", measureNumber, currentTimesig.BeatType)
	}
}
func (bs *beamSplitter) SplitSingle(ctx context.Context, notes []*entity.NoteRenderer, ts timesig.TimeSignature, segments []BeamSplitMarker, beamNo int) {
	if len(notes) == 0 {
		return
	}

	measureNumber := notes[0].MeasureNumber

	currentTimesig := ts.GetTimesignatureOnMeasure(ctx, measureNumber)
	switch currentTimesig.BeatType {
	case 4, 2:
		bs.quarter.SplitSingle(ctx, notes, ts, segments, beamNo)
	case 8:
		bs.eighth.SplitSingle(ctx, notes, ts, segments, beamNo)
	default:
		log.Printf("[rhythm] Unable to single split, On measure: %d: Unsupported beat-type %d.  \n", measureNumber, currentTimesig.BeatType)
	}
}
