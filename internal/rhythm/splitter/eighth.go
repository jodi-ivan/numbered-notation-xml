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

	rightMostBeforeSpan, leftMostAfterSpan := -1, -1
	topSpan := Interval(segments[1])
	sort.Sort(topSpan)

	bottomSpan := Interval(segments[2])
	sort.Sort(bottomSpan)

	unprocessedSegment := []BeamSplitMarker{}

	if bottomSpan[0].StartIndex > 0 && notes[bottomSpan[0].StartIndex].IsDotted {
		bottomSpan[0].StartIndex--
	}

	for i, top := range topSpan {
		if top.EndIndex < bottomSpan[0].StartIndex {
			rightMostBeforeSpan = i
			unprocessedSegment = append(unprocessedSegment, top)
		} else if top.StartIndex > bottomSpan[len(bottomSpan)-1].EndIndex {
			if leftMostAfterSpan == -1 {
				leftMostAfterSpan = i
			}
			unprocessedSegment = append(unprocessedSegment, top)
		}
	}

	leftIndexBeforeInterval := topSpan[0].StartIndex
	if rightMostBeforeSpan != -1 {
		leftIndexBeforeInterval = topSpan[rightMostBeforeSpan+1].StartIndex
	}

	if bottomSpan[0].StartIndex > leftIndexBeforeInterval {
		span := (bottomSpan[0].StartIndex - leftIndexBeforeInterval) % 3
		switch span {
		case 0:
			notes[leftIndexBeforeInterval].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			notes[bottomSpan[0].StartIndex-1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
			notes[bottomSpan[0].StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			unprocessedSegment = append(unprocessedSegment, BeamSplitMarker{StartIndex: leftIndexBeforeInterval, EndIndex: bottomSpan[0].StartIndex - 1})
		}
		//TODO 2 or 1 left on the span
	}

	rightIndexAfterIntervalStartIndex := bottomSpan[len(bottomSpan)-1].EndIndex + 1
	rightIndexAfterIntervalEndIndex := topSpan[len(topSpan)-1].EndIndex
	if leftMostAfterSpan != -1 {
		rightIndexAfterIntervalEndIndex = topSpan[leftMostAfterSpan-1].EndIndex - 1
	}

	diff := (rightIndexAfterIntervalEndIndex - rightIndexAfterIntervalStartIndex + 1) % 3

	notes[rightIndexAfterIntervalStartIndex-1+diff].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
	notes[rightIndexAfterIntervalStartIndex+diff].UpdateBeam(1, musicxml.NoteBeamTypeBegin)

	unprocessedSegment = append(unprocessedSegment, BeamSplitMarker{StartIndex: rightIndexAfterIntervalStartIndex + diff, EndIndex: rightIndexAfterIntervalEndIndex})

	// TODO : in betweeners
	// TODO: merge breathmatk algorithm
	es.SplitSingle(ctx, notes, ts, unprocessedSegment, 1)

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
