package splitter

import (
	"context"
	"sort"

	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

type quaterSpliiter struct {
}

func (qs *quaterSpliiter) shouldMergeSegments(notes []*entity.NoteRenderer, s, nextInterval BeamSplitMarker) bool {
	hasOneNoteGap := nextInterval.StartIndex-s.EndIndex == 2
	if !hasOneNoteGap {
		return false
	}
	isBreathmark := breathpause.IsBreathMark(notes[s.EndIndex+1])
	isGapBeam := len(notes[s.EndIndex].Beam) > 0 && len(notes[nextInterval.StartIndex].Beam) > 0
	currIntervalHas1Note := s.EndIndex-s.StartIndex == 0
	nextIntervalHas1Note := nextInterval.EndIndex-nextInterval.StartIndex == 0
	eitherHasOneNote := currIntervalHas1Note || nextIntervalHas1Note

	return isBreathmark && isGapBeam && eitherHasOneNote
}

func (qs *quaterSpliiter) Split(ctx context.Context, notes []*entity.NoteRenderer, ts timesig.TimeSignature, segments map[int][]BeamSplitMarker) {
	if len(segments[2]) == 0 {
		qs.SplitSingle(ctx, notes, ts, segments[1], 1)
		return
	}

	interval := Interval(segments[2])
	sort.Sort(interval)

	eigthSegment := Interval(segments[1])
	sort.Sort(eigthSegment)

	between := BeamSplitMarker{StartIndex: -1, EndIndex: -1}
	unprocessedSegment := []BeamSplitMarker{}
	marker := map[int]bool{}
	afterSegment := map[int][]BeamSplitMarker{}
	leftIdx := -1
	rigthIdx := len(eigthSegment) + 1

	// before
	for i, v := range eigthSegment {
		if v.EndIndex < interval[0].StartIndex {
			leftIdx = i
		}

		if v.StartIndex > interval[len(interval)-1].EndIndex {
			unprocessedSegment = append(unprocessedSegment, v)
			maxInterval := interval[len(interval)-1].EndIndex
			if afterSegment[maxInterval] == nil {
				afterSegment[maxInterval] = []BeamSplitMarker{}
			}
			afterSegment[maxInterval] = append(afterSegment[maxInterval], v)
			if i < rigthIdx {
				rigthIdx = i
			}
		}
	}

	outerMostInterval := map[int][2]BeamSplitMarker{}
	topIdx := 0

	// segment := BeamSplitMarker{StartIndex: -1, EndIndex: -1}

	if leftIdx+1 <= len(eigthSegment)-1 {
		segment := eigthSegment[leftIdx+1]

		diff := (segment.EndIndex - segment.StartIndex) + 1
		if diff%2 == 1 && (interval[0].StartIndex-segment.StartIndex)%2 == 1 { // interval will have 2_3 config when it has odd number
			if interval[0].StartIndex > 0 {
				interval[0].StartIndex--
			}
		}

		if leftIdx >= 0 {
			for i, v := range eigthSegment {
				if i <= leftIdx {
					unprocessedSegment = append(unprocessedSegment, v)
				}
			}
		}
	}

	for is, ss := range interval {

		notes[ss.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
		marker[ss.EndIndex] = true
		offset := 1
		if !notes[ss.StartIndex].IsDotted {
			offset = 0
		}

		if ss.StartIndex-1 >= 0 && breathpause.IsBreathMark(notes[ss.StartIndex-1]) {
			if ss.StartIndex-2 >= 0 && notes[ss.StartIndex-2].IsDotted && len(notes[ss.StartIndex-2].Beam) > 0 {
				notes[ss.StartIndex-1].Beam = map[int]entity.Beam{
					1: entity.Beam{
						Type:   musicxml.NoteBeamTypeContinue,
						Number: 1,
					},
				}
				offset = 2
				notes[ss.StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeContinue)
			}
		}

		ss.StartIndex -= offset
		interval[is] = ss

		notes[ss.StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
		marker[ss.StartIndex] = true

		if between.StartIndex == -1 {
			between.StartIndex = ss.EndIndex + 1
		} else if between.EndIndex == -1 {
			between.EndIndex = ss.StartIndex - 1
			if between.EndIndex-between.StartIndex > 0 {
				// has more than one note, need more processing for splitting
				notes[between.StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
				notes[between.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				unprocessedSegment = append(unprocessedSegment, between)

			} else if between.EndIndex-between.StartIndex == 0 && len(notes[between.StartIndex].Beam) > 0 {
				// just one note, just assign it accordingly
				notes[ss.StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeContinue)
				notes[ss.StartIndex-1].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			}

			between = BeamSplitMarker{StartIndex: -1, EndIndex: -1}
		}

		for topIdx < len(eigthSegment) && ss.StartIndex > eigthSegment[topIdx].EndIndex {
			topIdx++
		}

		if topIdx >= len(eigthSegment) {
			break
		}

		if ss.StartIndex >= eigthSegment[topIdx].StartIndex && ss.EndIndex <= eigthSegment[topIdx].EndIndex {
			if _, exists := outerMostInterval[topIdx]; !exists {
				outerMostInterval[topIdx] = [2]BeamSplitMarker{ss, ss}
			} else {
				outer := outerMostInterval[topIdx]
				outer[1] = ss
				outerMostInterval[topIdx] = outer
			}
		}

	}

	for segmentIdx, ss := range outerMostInterval {
		segment := eigthSegment[segmentIdx]
		minInterval := ss[0]
		maxInterval := ss[1]

		if segment.StartIndex < minInterval.StartIndex {
			if minInterval.StartIndex-segment.StartIndex > 1 {
				notes[segment.StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
				notes[minInterval.StartIndex-1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)

				unprocessedSegment = append(unprocessedSegment, BeamSplitMarker{ // still needed for splitting
					StartIndex: segment.StartIndex,
					EndIndex:   minInterval.StartIndex - 1,
				})
			} else {
				notes[minInterval.StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeContinue)
				notes[segment.StartIndex].UpdateBeam(1, musicxml.NoteBeamTypeBegin)

			}
		}

		if segment.EndIndex > maxInterval.EndIndex {
			mergeable := false
			canCarryOver := (segment.EndIndex-maxInterval.EndIndex == 1 && len(afterSegment[maxInterval.EndIndex]) > 0)
			if canCarryOver {
				mergeable = qs.shouldMergeSegments(notes, segment, afterSegment[maxInterval.EndIndex][0])
			}
			if segment.EndIndex-maxInterval.EndIndex > 1 || mergeable {
				notes[maxInterval.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[maxInterval.EndIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeBegin)

				notes[segment.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				unprocessedSegment = append(unprocessedSegment, BeamSplitMarker{ // still needed for splitting
					StartIndex: maxInterval.EndIndex + 1,
					EndIndex:   segment.EndIndex,
				})
			} else {
				notes[maxInterval.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeContinue)
				notes[segment.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
			}
		}
	}

	qs.SplitSingle(ctx, notes, ts, segments[2], 2)
	qs.SplitSingle(ctx, notes, ts, unprocessedSegment, 1)
}

func (qs *quaterSpliiter) SplitSingle(ctx context.Context, notes []*entity.NoteRenderer, ts timesig.TimeSignature, segments []BeamSplitMarker, beamNo int) {
	interval := Interval(segments)
	sort.Sort(interval)
	skipSplitSegmentIdx := map[int]bool{}
	mergedSegment := []BeamSplitMarker{}

	skipProcess := map[int]bool{}
	for is, s := range interval {
		if skipProcess[is] {
			continue
		}
		if is+1 < len(interval) {
			nextInterval := interval[is+1]

			if qs.shouldMergeSegments(notes, s, nextInterval) { // merge two segements
				notes[s.StartIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
				notes[s.EndIndex+1].Beam = map[int]entity.Beam{
					beamNo: entity.Beam{
						Type:   musicxml.NoteBeamTypeContinue,
						Number: beamNo,
					},
				}
				notes[nextInterval.EndIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)

				// this should be max at 3

				mergedSegment = append(mergedSegment, BeamSplitMarker{
					StartIndex: s.StartIndex,
					EndIndex:   s.EndIndex + 2,
				})
				skipSplitSegmentIdx[len(mergedSegment)-1] = true
				if s.EndIndex+3 < nextInterval.EndIndex {
					mergedSegment = append(mergedSegment, BeamSplitMarker{
						StartIndex: s.EndIndex + 3,
						EndIndex:   nextInterval.EndIndex,
					})
				} else {
					skipProcess[is+1] = true
				}
			} else {
				mergedSegment = append(mergedSegment, s)
			}
		} else {
			hasOneNote := s.EndIndex-s.StartIndex == 0
			isLast2Notes := s.EndIndex+1 == len(notes)-1
			isLastNotesBreathmark := isLast2Notes && breathpause.IsBreathMark(notes[s.EndIndex+1])

			if hasOneNote && isLast2Notes && isLastNotesBreathmark {
				notes[s.StartIndex-1].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
				notes[s.StartIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
				notes[s.EndIndex+1].Beam = map[int]entity.Beam{
					beamNo: entity.Beam{
						Type:   musicxml.NoteBeamTypeEnd,
						Number: beamNo,
					},
				}
				s.EndIndex += 1
			}
			mergedSegment = append(mergedSegment, s)

			// forsome reason, after double segement it send merged
			totalBreathmark := 0
			for i := s.StartIndex; i < s.EndIndex; i++ {
				if breathpause.IsBreathMark(notes[i]) {
					totalBreathmark++
				}
			}

			if totalBreathmark > 0 {
				skipSplitSegmentIdx[len(mergedSegment)-1] = true //&& s.EndIndex-s.StartIndex < 3
			}

		}

	}

	for i, segment := range mergedSegment {
		diff := (segment.EndIndex - segment.StartIndex) + 1
		switch diff {
		case 3:
			if skipSplitSegmentIdx[i] {
				continue
			}
			notes[segment.StartIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+1].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
		case 4: // split 2x2

			notes[segment.StartIndex+1].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+2].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
		case 5:
			notes[segment.StartIndex+1].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+2].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
		case 6:
			// split 2x2x2
			notes[segment.StartIndex+1].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+2].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
			notes[segment.StartIndex+3].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+4].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)

		default:
			if diff > 6 {
				startIndex := segment.StartIndex
				if diff%2 == 1 && !notes[startIndex].IsRest {
					notes[startIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
					notes[startIndex+1].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
					startIndex = startIndex + 2
				}

				for i := startIndex + 1; i < segment.EndIndex; i += 2 {
					notes[i].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)
					if i+1 < len(notes) {
						notes[i+1].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
					}
				}
			}

		}

	}
}
