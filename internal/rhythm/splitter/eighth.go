package splitter

import (
	"context"
	"sort"

	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

type eighthSplitter struct {
}

func (es *eighthSplitter) Split(ctx context.Context, notes []*entity.NoteRenderer, ts timesig.TimeSignature, segments map[int][]BeamSplitMarker) {
	if len(segments[2]) == 0 {
		es.SplitSingle(ctx, notes, ts, segments[1], 1)
		return
	}

	for _, segment := range segments[1] {
		unprocessedSegment := []BeamSplitMarker{}

		// diff := (segment.EndIndex - segment.StartIndex) + 1
		interval := Interval(segments[2])
		sort.Sort(interval)

		es.SplitSingle(ctx, notes, ts, segments[2], 2)

		before := BeamSplitMarker{
			StartIndex: segment.StartIndex,
			EndIndex:   interval[0].StartIndex - 2,
		}

		if before.EndIndex > 0 {
			unprocessedSegment = append(unprocessedSegment, before)

		}

		for is, ss := range interval {
			for i := ss.StartIndex; i <= ss.EndIndex; i++ {
				notes[i].UpdateBeam(1, musicxml.NoteBeamTypeContinue)
				if notes[i].IsDotted {
					interval[is].EndIndex++
					if is+1 <= len(interval)-1 {
						interval[is+1].StartIndex--
					}
				}

			}

			notes[ss.StartIndex-1].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			notes[ss.EndIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
		}

		after := BeamSplitMarker{
			StartIndex: interval[len(interval)-1].EndIndex + 1,
			EndIndex:   segment.EndIndex,
		}

		if after.EndIndex-after.StartIndex > 1 {
			unprocessedSegment = append(unprocessedSegment, after)
		} else {
			notes[after.StartIndex-1].UpdateBeam(1, musicxml.NoteBeamTypeContinue)
			notes[after.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)

		}

		for _, up := range unprocessedSegment {
			if up.EndIndex >= 0 && up.StartIndex >= 0 && up.EndIndex > up.StartIndex {
				notes[up.StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
				notes[up.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				es.SplitSingle(ctx, notes, ts, unprocessedSegment, 1)
			}

		}

	}

}
func (es *eighthSplitter) SplitSingle(ctx context.Context, notes []*entity.NoteRenderer, ts timesig.TimeSignature, segments []BeamSplitMarker, beamNo int) {
	interval := Interval(segments)
	sort.Sort(interval)
	skipSplitSegmentIdx := map[int]bool{}
	skipProcess := map[int]bool{}
	mergedSegment := []BeamSplitMarker{}

	for is, s := range interval {
		if skipProcess[is] {
			continue
		}

		if is+1 < len(interval) {
			nextInterval := interval[is+1]

			hasOneNoteGap := nextInterval.StartIndex-s.EndIndex == 2
			isBreathmark := breathpause.IsBreathMark(notes[s.EndIndex+1])
			isGapBeam := len(notes[s.EndIndex].Beam) > 0 && len(notes[nextInterval.StartIndex].Beam) > 0
			currIntervalLT2Note := s.EndIndex-s.StartIndex < 2
			nextIntervalLT2Note := nextInterval.EndIndex-nextInterval.StartIndex < 2
			eitherOneLT2Note := currIntervalLT2Note || nextIntervalLT2Note

			if hasOneNoteGap && isBreathmark && isGapBeam && eitherOneLT2Note {
				notes[s.StartIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
				notes[s.EndIndex+1].Beam = map[int]entity.Beam{
					1: entity.Beam{
						Type:   musicxml.NoteBeamTypeContinue,
						Number: 1,
					},
				}
				notes[nextInterval.EndIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)

				mergedSegment = append(mergedSegment, BeamSplitMarker{
					StartIndex: s.StartIndex,
					EndIndex:   nextInterval.EndIndex,
				})
				skipSplitSegmentIdx[len(mergedSegment)-1] = true
				skipProcess[is+1] = true
			} else {
				mergedSegment = append(mergedSegment, s)
			}
		} else {
			mergedSegment = append(mergedSegment, s)
		}
	}
	for _, segment := range mergedSegment {

		diff := (segment.EndIndex - segment.StartIndex) + 1
		for i := segment.StartIndex + 1; i < segment.EndIndex; i++ {
			notes[i].UpdateBeam(beamNo, musicxml.NoteBeamTypeContinue)
		}
		switch diff {
		case 5:
			notes[segment.StartIndex+2].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+3].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
			if segment.EndIndex+1 == len(notes)-1 && breathpause.IsBreathMark(notes[segment.EndIndex+1]) {
				notes[segment.EndIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeContinue)
				notes[segment.EndIndex+1].Beam = map[int]entity.Beam{
					1: entity.Beam{
						Type:   musicxml.NoteBeamTypeEnd,
						Number: 1,
					},
				}
			}
		case 6, 7:
			notes[segment.StartIndex+2].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+3].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
		default:

			if diff > 7 {
				for i := segment.StartIndex; i < segment.EndIndex; i += 3 {
					if i+3 < segment.EndIndex {
						notes[i+2].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
						notes[i+3].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)

					}
				}

			}

		}
	}
}
