package rhythm

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (ri *rhythmInteractor) RenderBezier(set []SlurBezier, canv canvas.Canvas) {
	if len(set) == 0 {
		return
	}
	canv.Group("class='slurties'")
	for _, s := range set {

		slurResult := SlurBezier{
			SlurTieType: s.SlurTieType,
			Start: CoordinateWithOctave{
				Coordinate: entity.Coordinate{
					X: s.Start.X + 5,
					Y: s.Start.Y + 5,
				},
				Octave: s.Start.Octave,
			},
			End: CoordinateWithOctave{
				Coordinate: entity.Coordinate{
					X: s.End.X + 5,
					Y: s.End.Y + 5,
				},
				Octave: s.End.Octave,
			},
			LineType: s.LineType,
			Pull:     s.Pull,
		}

		offset := float64(2.25)

		if slurResult.Start.Octave < 0 {
			slurResult.Start = CoordinateWithOctave{
				Coordinate: entity.Coordinate{
					X: slurResult.Start.X + offset,
					Y: slurResult.Start.Y + offset,
				},
			}
		}

		if slurResult.End.Octave < 0 {

			slurResult.End = CoordinateWithOctave{
				Coordinate: entity.Coordinate{
					X: slurResult.End.X - offset,
					Y: slurResult.End.Y + offset,
				},
			}
		}

		pullY := slurResult.Start.Y

		block := ((slurResult.End.X - slurResult.Start.X) / constant.UPPERCASE_LENGTH) // * 2
		if block < 2 {
			pullY += 5
		} else if block < 5 {
			pullY += 12
		} else {
			pullY += 15 //long distance ties, need more height
		}

		pull := CoordinateWithOctave{
			Coordinate: entity.Coordinate{
				X: slurResult.Start.X + ((slurResult.End.X - slurResult.Start.X) / 2),
				Y: pullY + 0.3,
			},
		}

		slurResult.Pull = pull
		lineType := "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1"
		if slurResult.LineType == musicxml.NoteSlurLineTypeDashed {
			// Calculate approximate curve length
			curveLen := quadBezierLength(slurResult.Start.Coordinate, pull.Coordinate, slurResult.End.Coordinate, 30)

			dash := 3.5
			patternCount := math.Floor(curveLen / (dash * 2))
			gap := (curveLen / patternCount) - dash

			lineType += fmt.Sprintf(";stroke-dasharray:%.1f %.1f;", dash, gap)
			lineType += "stroke-dashoffset:" + fmt.Sprintf("%f", dash/2) + ";"
		}
		canv.Qbez(
			int(math.Round(slurResult.Start.X)),
			int(math.Round(slurResult.Start.Y)),
			int(math.Round(pull.X)),
			int(math.Ceil(pull.Y)),
			int(math.Round(slurResult.End.X)),
			int(math.Round(slurResult.End.Y)),
			lineType,
		)

	}
	canv.Gend()
}

func quadBezierLength(p0, p1, p2 entity.Coordinate, steps int) float64 {
	var length float64
	prev := p0

	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)

		x := math.Pow(1-t, 2)*p0.X +
			2*(1-t)*t*p1.X +
			math.Pow(t, 2)*p2.X

		y := math.Pow(1-t, 2)*p0.Y +
			2*(1-t)*t*p1.Y +
			math.Pow(t, 2)*p2.Y

		dx := x - prev.X
		dy := y - prev.Y

		length += math.Hypot(dx, dy)
		prev = entity.Coordinate{X: x, Y: y}
	}

	return length
}

func (ri *rhythmInteractor) RenderSlurTies(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer, maxXPosition float64) {
	slurs := map[int]SlurBezier{}
	slurSets := []SlurBezier{}

	ties := map[int]SlurBezier{}
	tiesSet := []SlurBezier{}
	for _, note := range notes {
		for _, s := range note.Slur {
			if s.Type == musicxml.NoteSlurTypeStop || s.Type == musicxml.NoteSlurTypeHop {
				temp := slurs[s.Number]
				temp.End = CoordinateWithOctave{
					Coordinate: entity.Coordinate{
						X: float64(note.PositionX - 2),
						Y: float64(note.PositionY),
					},
					Octave: note.Octave,
				}

				if temp.Start.X == 0 && temp.Start.Y == 0 {
					temp.Start = CoordinateWithOctave{
						Coordinate: entity.Coordinate{
							X: float64(note.PositionX - constant.UPPERCASE_LENGTH),
							Y: float64(note.PositionY),
						},
						Octave: 0,
					}
				}
				slurs[s.Number] = temp

				slurSets = append(slurSets, slurs[s.Number])

				delete(slurs, s.Number)

			}

			if s.Type == musicxml.NoteSlurTypeStart || s.Type == musicxml.NoteSlurTypeHop {
				slurs[s.Number] = SlurBezier{
					SlurTieType: SlurTieTypeSlur,
					Start: CoordinateWithOctave{
						Coordinate: entity.Coordinate{
							X: float64(note.PositionX + 2),
							Y: float64(note.PositionY),
						},
						Octave: note.Octave,
					},
					LineType: s.LineType,
				}
			}

		}

		if note.Tie != nil {
			if note.Tie.Type == musicxml.NoteSlurTypeStart {
				ties[note.Note] = SlurBezier{
					SlurTieType: SlurTieTypeTie,
					Start: CoordinateWithOctave{
						Coordinate: entity.Coordinate{
							X: float64(note.PositionX),
							Y: float64(note.PositionY),
						},
						Octave: note.Octave,
					},
					LineType: note.Tie.LineType,
				}
			} else if note.Tie.Type == musicxml.NoteSlurTypeStop {
				temp := ties[note.Note]
				temp.End = CoordinateWithOctave{
					Coordinate: entity.Coordinate{
						X: float64(note.PositionX),
						Y: float64(note.PositionY),
					},
					Octave: note.Octave,
				}
				ties[note.Note] = temp

				tiesSet = append(tiesSet, ties[note.Note])
				delete(slurs, note.Note)
			}
		}

	}

	if len(slurs) > 0 { // there is start, but no end
		for _, slur := range slurs {
			temp := slur
			if temp.End.Coordinate.X == 0 && temp.End.Coordinate.Y == 0 {
				temp.End = CoordinateWithOctave{
					Coordinate: entity.Coordinate{
						X: float64(maxXPosition - 5),
						Y: float64(temp.Start.Y),
					},
					Octave: 0,
				}
			}
			slurSets = append(slurSets, temp)
		}
	}

	ri.RenderBezier(slurSets, canv)
	ri.RenderBezier(tiesSet, canv)

}

func (ri *rhythmInteractor) RenderBeam(ctx context.Context, canv canvas.Canvas, ts timesig.TimeSignature, notes []*entity.NoteRenderer) {

	beams := map[int]BeamLine{}
	beamSets := []BeamLine{}

	beamSegments := map[int][]beamSplitMarker{}

	var cleanedNote []*entity.NoteRenderer

	cleanedNote, beamSegments[1] = cleanBeamByNumber(ctx, notes, 1)
	cleanedNote, beamSegments[2] = cleanBeamByNumber(ctx, cleanedNote, 2)
	// TODO: more than 16th beam support check

	cleanedNote = splitBeam(ctx, ts, cleanedNote, beamSegments)

	for _, note := range cleanedNote {

		for _, b := range note.Beam {
			positionY := float64(note.PositionY - 22 + ((b.Number) * 3))

			switch b.Type {
			case musicxml.NoteBeamTypeBegin:
				beams[b.Number] = BeamLine{
					Start: entity.Coordinate{
						X: float64(note.PositionX),
						Y: positionY,
					},
					Number: b.Number,
				}
			case musicxml.NoteBeamTypeEnd:

				beam := beams[b.Number]

				if beam.Start.X == 0 {
					beams[b.Number] = BeamLine{
						Start: entity.Coordinate{
							X: float64(note.PositionX),
							Y: positionY,
						},
						End: entity.Coordinate{
							X: float64(note.PositionX) + 8,
							Y: positionY,
						},
						Number: b.Number,
					}

				} else {
					beam.End = entity.Coordinate{
						X: float64(note.PositionX) + 8,
						Y: beam.Start.Y,
					}
					beams[b.Number] = beam
				}

				beamSets = append(beamSets, beams[b.Number])

				delete(beams, b.Number)
			}
		}

	}

	if len(beamSets) == 0 {
		return
	}

	canv.Group("class='beam'")

	sort.SliceStable(beamSets, func(i, j int) bool {
		one := beamSets[i]
		two := beamSets[j]

		return one.Number > two.Number
	})

	m := map[[2]float64]bool{}

	for _, b := range beamSets {
		if b.End.X-b.Start.X < constant.LOWERCASE_LENGTH {
			diff := constant.LOWERCASE_LENGTH - (b.End.X - b.Start.X)
			b.Start.X -= diff / 2
			b.End.X += diff / 2
		}

		if b.Number == 2 {
			m[[2]float64{b.Start.X, b.End.X}] = true
		}

		if b.Number == 1 && m[[2]float64{b.Start.X, b.End.X}] && b.Start.X == constant.LAYOUT_INDENT_LENGTH {
			b.Start.X -= (constant.UPPERCASE_LENGTH / 2)
		}
		canv.Line(
			int(math.Round(b.Start.X)),
			int(math.Round(b.Start.Y)),
			int(math.Round(b.End.X)),
			int(math.Round(b.End.Y)),
			"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.2",
		)
	}
	canv.Gend()
}

func cleanBeamByNumber(ctx context.Context, notes []*entity.NoteRenderer, beamNumber int) ([]*entity.NoteRenderer, []beamSplitMarker) {

	switches := map[int]beamMarker{}

	markers := make([]beamSplitMarker, 0)

	var prev *entity.NoteRenderer

	for indexNote, note := range notes {
		note.IndexPosition = indexNote

		if len(note.Beam) == 0 { // stopping the beam
			if indexNote == 0 {
				prev = note
				continue
			} else {

				t, ok := switches[beamNumber]
				if !ok {
					prev = note
					continue
				}

				prev.Beam[beamNumber] = entity.Beam{
					Number: beamNumber,
					Type:   musicxml.NoteBeamTypeEnd,
				}

				markers = append(markers, beamSplitMarker{
					StartIndex: t.NoteBeginIndex,
					EndIndex:   indexNote - 1,
				})

				delete(switches, beamNumber)

			}
		}

		if t, ok := switches[beamNumber]; !ok {

			if _, hasBeam := note.Beam[beamNumber]; !hasBeam {
				prev = note
				continue
			}
			newBeam := map[int]entity.Beam{}

			for k, v := range note.Beam {
				newBeam[k] = v
			}

			switches[beamNumber] = beamMarker{
				NoteBeamType:   musicxml.NoteBeamTypeBegin,
				NoteBeginIndex: indexNote,
			}

			newBeam[beamNumber] = entity.Beam{
				Number: beamNumber,
				Type:   musicxml.NoteBeamTypeBegin,
			}
			note.Beam = newBeam
		} else {

			if prev == nil {
				continue
			}

			if _, hasBeam := note.Beam[beamNumber]; hasBeam {
				newBeam := map[int]entity.Beam{}

				for k, v := range note.Beam {
					newBeam[k] = v
				}

				switches[beamNumber] = beamMarker{
					NoteBeamType:   musicxml.NoteBeamTypeContinue,
					NoteBeginIndex: switches[beamNumber].NoteBeginIndex,
				}

				newBeam[beamNumber] = entity.Beam{
					Number: beamNumber,
					Type:   musicxml.NoteBeamTypeContinue,
				}
				note.Beam = newBeam
				prev = note
				continue
			}

			if t.NoteBeamType == musicxml.NoteBeamTypeBegin || t.NoteBeamType == musicxml.NoteBeamTypeContinue {

				if _, ok := prev.Beam[beamNumber]; !ok {
					prev = note
					continue
				}

				prev.Beam[beamNumber] = entity.Beam{
					Number: beamNumber,
					Type:   musicxml.NoteBeamTypeEnd,
				}

				delete(switches, beamNumber)

				markers = append(markers, beamSplitMarker{
					StartIndex: t.NoteBeginIndex,
					EndIndex:   indexNote - 1,
				})

			}

		}
		prev = note

	}

	if prev != nil && len(prev.Beam) > 0 {
		additional, ok := prev.Beam[beamNumber]

		if ok {
			if additional.Type != musicxml.NoteBeamTypeEnd {
				newBeam := prev.Beam

				newBeam[beamNumber] = entity.Beam{
					Type:   musicxml.NoteBeamTypeEnd,
					Number: beamNumber,
				}

				prev.Beam = newBeam

				if t, ok := switches[beamNumber]; ok {

					markers = append(markers, beamSplitMarker{
						StartIndex: t.NoteBeginIndex,
						EndIndex:   prev.IndexPosition,
					})
				}

			} else {
				if _, ok := switches[beamNumber]; !ok {
					newBeam := prev.Beam
					newBeam[beamNumber] = entity.Beam{
						Type:   musicxml.NoteBeamTypeBackwardHook,
						Number: beamNumber,
					}
					prev.Beam = newBeam
				}
				markers = append(markers, beamSplitMarker{
					StartIndex: prev.IndexPosition,
					EndIndex:   prev.IndexPosition,
				})
			}

		}
	}

	return notes, markers
}

func shouldMergeSegments(notes []*entity.NoteRenderer, s, nextInterval beamSplitMarker) bool {
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

func splitSingleBeamQuarter(ctx context.Context, notes []*entity.NoteRenderer, segments []beamSplitMarker, beamNo int) {
	interval := Interval(segments)
	sort.Sort(interval)
	skipSplitSegmentIdx := map[int]bool{}
	mergedSegment := []beamSplitMarker{}

	skipProcess := map[int]bool{}
	for is, s := range interval {
		if skipProcess[is] {
			continue
		}
		if is+1 < len(interval) {
			nextInterval := interval[is+1]

			if shouldMergeSegments(notes, s, nextInterval) { // merge two segements
				notes[s.StartIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeBegin)
				notes[s.EndIndex+1].Beam = map[int]entity.Beam{
					beamNo: entity.Beam{
						Type:   musicxml.NoteBeamTypeContinue,
						Number: beamNo,
					},
				}
				notes[nextInterval.EndIndex].UpdateBeam(beamNo, musicxml.NoteBeamTypeEnd)

				// this should be max at 3

				mergedSegment = append(mergedSegment, beamSplitMarker{
					StartIndex: s.StartIndex,
					EndIndex:   s.EndIndex + 2,
				})
				skipSplitSegmentIdx[len(mergedSegment)-1] = true
				if s.EndIndex+3 < nextInterval.EndIndex {
					mergedSegment = append(mergedSegment, beamSplitMarker{
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
				if diff%2 == 1 {
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
func splitSingleBeam(ctx context.Context, notes []*entity.NoteRenderer, segments []beamSplitMarker, beamNo int) {
	interval := Interval(segments)
	sort.Sort(interval)
	skipSplitSegmentIdx := map[int]bool{}
	skipProcess := map[int]bool{}
	mergedSegment := []beamSplitMarker{}

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

				mergedSegment = append(mergedSegment, beamSplitMarker{
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

func splitBeamQuater(ctx context.Context, notes []*entity.NoteRenderer, segments map[int][]beamSplitMarker) {
	if len(segments[2]) == 0 {
		splitSingleBeamQuarter(ctx, notes, segments[1], 1)
		return
	}

	interval := Interval(segments[2])
	sort.Sort(interval)

	eigthSegment := Interval(segments[1])
	sort.Sort(eigthSegment)

	between := beamSplitMarker{StartIndex: -1, EndIndex: -1}
	unprocessedSegment := []beamSplitMarker{}
	marker := map[int]bool{}
	afterSegment := map[int][]beamSplitMarker{}
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
				afterSegment[maxInterval] = []beamSplitMarker{}
			}
			afterSegment[maxInterval] = append(afterSegment[maxInterval], v)
			if i < rigthIdx {
				rigthIdx = i
			}
		}
	}

	outerMostInterval := map[int][2]beamSplitMarker{}
	topIdx := 0

	// segment := beamSplitMarker{StartIndex: -1, EndIndex: -1}

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
				if i < leftIdx {
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

			between = beamSplitMarker{StartIndex: -1, EndIndex: -1}
		}

		for topIdx < len(eigthSegment) && ss.StartIndex > eigthSegment[topIdx].EndIndex {
			topIdx++
		}

		if topIdx >= len(eigthSegment) {
			break
		}

		if ss.StartIndex >= eigthSegment[topIdx].StartIndex && ss.EndIndex <= eigthSegment[topIdx].EndIndex {
			if _, exists := outerMostInterval[topIdx]; !exists {
				outerMostInterval[topIdx] = [2]beamSplitMarker{ss, ss}
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

				unprocessedSegment = append(unprocessedSegment, beamSplitMarker{ // still needed for splitting
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
				mergeable = shouldMergeSegments(notes, segment, afterSegment[maxInterval.EndIndex][0])
			}
			if segment.EndIndex-maxInterval.EndIndex > 1 || mergeable {
				notes[maxInterval.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[maxInterval.EndIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeBegin)

				notes[segment.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				unprocessedSegment = append(unprocessedSegment, beamSplitMarker{ // still needed for splitting
					StartIndex: maxInterval.EndIndex + 1,
					EndIndex:   segment.EndIndex,
				})
			} else {
				notes[maxInterval.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeContinue)
				notes[segment.EndIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
			}
		}
	}

	splitSingleBeamQuarter(ctx, notes, segments[2], 2)
	splitSingleBeamQuarter(ctx, notes, unprocessedSegment, 1)
}

func splitBeam(ctx context.Context, ts timesig.TimeSignature, notes []*entity.NoteRenderer, segments map[int][]beamSplitMarker) []*entity.NoteRenderer {

	if len(segments[1]) == 0 && len(segments[2]) == 0 {
		return notes
	}

	currTs := ts.GetTimesignatureOnMeasure(ctx, notes[segments[1][0].StartIndex].MeasureNumber)
	if currTs.BeatType == 4 || currTs.BeatType == 2 {
		splitBeamQuater(ctx, notes, segments)
		return notes
	}

	if len(segments[2]) == 0 {
		splitSingleBeam(ctx, notes, segments[1], 1)
		return notes
	}

	for _, segment := range segments[1] {
		unprocessedSegment := []beamSplitMarker{}

		// diff := (segment.EndIndex - segment.StartIndex) + 1
		interval := Interval(segments[2])
		sort.Sort(interval)

		splitSingleBeam(ctx, notes, segments[2], 2)

		before := beamSplitMarker{
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

		after := beamSplitMarker{
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
				splitSingleBeam(ctx, notes, unprocessedSegment, 1)
			}

		}

	}

	return notes
}
