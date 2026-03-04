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

	// beamSets was generated by loop of the map, does not guaranteee the order.
	// sorted need so the element printed in svg would be ALWAYS the same order.
	// it does not matter in the user end, since it is specific use coordinate X and Y
	sort.SliceStable(beamSets, func(i, j int) bool {
		one := beamSets[i]
		two := beamSets[j]

		return one.Start.X < two.Start.X

	})
	for _, b := range beamSets {
		if b.End.X-b.Start.X < constant.LOWERCASE_LENGTH {
			diff := constant.LOWERCASE_LENGTH - (b.End.X - b.Start.X)
			b.Start.X -= diff / 2
			b.End.X += diff / 2
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

func splitSingleBeam(ctx context.Context, ts timesig.TimeSignature, notes []*entity.NoteRenderer, segments []beamSplitMarker) {
	for _, segment := range segments {
		totalBreathmark := 0
		for n := segment.StartIndex; n <= segment.EndIndex && n < len(notes); n++ {
			if breathpause.IsBreathMark(notes[n]) {
				totalBreathmark++
			}
		}
		diff := (segment.EndIndex - segment.StartIndex) + 1
		currTs := ts.GetTimesignatureOnMeasure(ctx, notes[segment.StartIndex].MeasureNumber)
		switch diff {
		case 4: // split 2x2
			if !(currTs.BeatType == 8 && totalBreathmark > 0) { //TODO: need more case for handling this.
				notes[segment.StartIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			}
		case 5: // split 3 x 2
			notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
		case 6:
			if currTs.BeatType == 4 {
				// split 2x2x2
				notes[segment.StartIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
				notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[segment.StartIndex+4].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			} else if currTs.BeatType == 8 {
				// split 3x3
				notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			}
		default:
			if diff > 6 && currTs.Beat == 1 && currTs.BeatType == 4 {
				startIndex := segment.StartIndex
				if diff%2 == 1 {
					notes[startIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
					notes[startIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
					startIndex = startIndex + 2
				}

				// split by 2x2
				for i := startIndex; i < len(notes); i += 2 {
					if breathpause.IsBreathMark(notes[i]) {
						i = i - 1
						continue
					}

					notes[i].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
					if i+1 < len(notes) {
						notes[i+1].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
					}
				}

				lastIndex := len(notes) - 1
				if notes[len(notes)-1].Barline != nil {
					lastIndex--
				}

				if beam, ok := notes[lastIndex].Beam[1]; ok {
					if len(notes) > lastIndex-1 {
						prevBeam, ok := notes[lastIndex-1].Beam[1]
						if ok &&
							prevBeam.Type == musicxml.NoteBeamTypeEnd && beam.Type == musicxml.NoteBeamTypeBegin {

							// notes[lastIndex-1].UpdateBeam(1, musicxml.NoteBeamTypeContinue)
							notes[lastIndex].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						}
					}
				}

			}

			if diff > 6 && currTs.BeatType == 8 {
				for i := segment.StartIndex; i < len(notes); i += 2 {
					if breathpause.IsBreathMark(notes[i]) {
						i = i - 1
						continue
					}
					if i+3 < len(notes) {
						notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
					}
				}

			}

		}
	}

}

// REFACTOR: refactor this until all cases are cover.
func splitBeam(ctx context.Context, ts timesig.TimeSignature, notes []*entity.NoteRenderer, segments map[int][]beamSplitMarker) []*entity.NoteRenderer {

	if len(segments[1]) == 0 && len(segments[2]) == 0 {
		return notes
	}

	if len(segments[2]) == 0 {
		splitSingleBeam(ctx, ts, notes, segments[1])
		return notes
	}

	for _, segment := range segments[1] {
		diff := (segment.EndIndex - segment.StartIndex) + 1
		totalSubSegment := 0
		currTs := ts.GetTimesignatureOnMeasure(ctx, notes[segment.StartIndex].MeasureNumber)

		for _, subSegment := range segments[2] {
			if subSegment.EndIndex <= segment.EndIndex {
				totalSubSegment++
			}

		}
		subSegment := (segments[2])[0]
		distance := (subSegment.EndIndex - subSegment.StartIndex) + 1
		startingPoint := (subSegment.StartIndex - segment.StartIndex)

		if diff == 4 {

			if distance == 1 {
				notes[segment.StartIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			} else if distance == 4 {
				notes[subSegment.StartIndex+1].UpdateBeam(2, musicxml.NoteBeamTypeEnd)
				notes[subSegment.StartIndex+2].UpdateBeam(2, musicxml.NoteBeamTypeBegin)
			}
		} else if diff == 8 {
			// currently taylored to kj-026
			notes[segment.StartIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			notes[segment.StartIndex+4].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
			notes[segment.StartIndex+5].UpdateBeam(1, musicxml.NoteBeamTypeBegin)

		} else if diff > 4 {
			if subSegment.EndIndex <= segment.EndIndex {

				if distance == 1 {

					if startingPoint <= 2 {
						// split 2x3
						notes[segment.StartIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
					} else {
						// split 3x2
						notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
					}
				} else if distance == 2 {

					if currTs.BeatType == 8 {
						for i := segment.StartIndex; i <= segment.EndIndex; i += 2 {
							if breathpause.IsBreathMark(notes[i]) {
								i = i - 1
								continue
							}

							if i+3 < len(notes) {
								notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
								notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
							}
						}
						return notes
					}

					if startingPoint <= 1 {
						// split 3x2
						notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)

					} else {

						// first separate 1st segement of sub segment
						notes[subSegment.StartIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[subSegment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)

						// split 2x2 from segment startSegement to subsegement startSegment
						for i := segment.StartIndex; i < subSegment.StartIndex-1; i += 2 {
							notes[i+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
							notes[i+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
						}

						// TODO: split 2x2 the rest after subsegement end to first segment end

					}
				} else if distance == 3 {
					offset := diff - 5
					if startingPoint == offset-0 {
						// split 3x2
						notes[segment.StartIndex+2+offset].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[segment.StartIndex+3+offset].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
					} else if startingPoint == 2-offset {
						// split 2x3
						notes[segment.StartIndex+1+offset].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[segment.StartIndex+2+offset].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
					}
				} else if distance == 4 {
					// split the subsegment to 2x2
					notes[subSegment.StartIndex+1].UpdateBeam(2, musicxml.NoteBeamTypeEnd)
					notes[subSegment.StartIndex+2].UpdateBeam(2, musicxml.NoteBeamTypeBegin)

					if len(notes) >= subSegment.StartIndex+4 {
						// the subsegment is 2x2, hence the parent need 4 length.
						notes[subSegment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[subSegment.StartIndex+4].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
					}
				} else {
					// split 3x2
					notes[subSegment.StartIndex+2].UpdateBeam(2, musicxml.NoteBeamTypeEnd)
					notes[subSegment.StartIndex+3].UpdateBeam(2, musicxml.NoteBeamTypeBegin)
				}
			}
		}
	}

	return notes
}
