package renderer

import (
	"context"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func RenderBezier(set []SlurBezier, canvas *svg.SVG) {
	canvas.Group("class='slurties'")
	// TODO: check ties across measure bar
	for _, s := range set {

		slurResult := SlurBezier{
			Start: CoordinateWithOctave{
				Coordinate: Coordinate{
					X: s.Start.X + 5,
					Y: s.Start.Y + 5,
				},
				Octave: s.Start.Octave,
			},
			End: CoordinateWithOctave{
				Coordinate: Coordinate{
					X: s.End.X + 5,
					Y: s.End.Y + 5,
				},
				Octave: s.End.Octave,
			},
		}

		if slurResult.Start.Octave < 0 {
			slurResult.Start = CoordinateWithOctave{
				Coordinate: Coordinate{
					X: slurResult.Start.X + 3,
					Y: slurResult.Start.Y + 3,
				},
			}
		}

		if slurResult.End.Octave < 0 {

			slurResult.End = CoordinateWithOctave{
				Coordinate: Coordinate{
					X: slurResult.End.X - 3,
					Y: slurResult.End.Y + 3,
				},
			}
		}

		pull := CoordinateWithOctave{
			Coordinate: Coordinate{
				X: slurResult.Start.X + ((slurResult.End.X - slurResult.Start.X) / 2),
				Y: slurResult.Start.Y + 7.5,
			},
		}
		slurResult.Pull = pull

		canvas.Qbez(
			int(math.Round(slurResult.Start.X)),
			int(math.Round(slurResult.Start.Y)),
			int(math.Round(pull.X)),
			int(math.Round(pull.Y)),
			int(math.Round(slurResult.End.X)),
			int(math.Round(slurResult.End.Y)),
			"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.5",
		)
	}
	canvas.Gend()
}

func RenderSlurAndBeam(ctx context.Context, canvas *svg.SVG, notes []*NoteRenderer, measureNo int) {
	slurs := map[int]SlurBezier{}
	slurSets := []SlurBezier{}

	// [ ] 6/8 time signature
	beams := map[int]BeamLine{}
	beamSets := []BeamLine{}

	ties := map[int]SlurBezier{}
	tiesSet := []SlurBezier{}

	beamSegments := map[int][]beamSplitMarker{}

	var cleanedNote []*NoteRenderer

	cleanedNote, beamSegments[1] = cleanBeamByNumber(ctx, notes, 1)
	cleanedNote, beamSegments[2] = cleanBeamByNumber(ctx, cleanedNote, 2)

	cleanedNote = splitBeam(ctx, cleanedNote, beamSegments)

	for _, note := range cleanedNote {

		for _, s := range note.Slur {
			if s.Type == musicxml.NoteSlurTypeStart {
				slurs[s.Number] = SlurBezier{
					Start: CoordinateWithOctave{
						Coordinate: Coordinate{
							X: float64(note.PositionX),
							Y: float64(note.PositionY),
						},
						Octave: note.Octave,
					},
				}
			} else if s.Type == musicxml.NoteSlurTypeStop {
				temp := slurs[s.Number]
				temp.End = CoordinateWithOctave{
					Coordinate: Coordinate{
						X: float64(note.PositionX),
						Y: float64(note.PositionY),
					},
					Octave: note.Octave,
				}
				slurs[s.Number] = temp

				slurSets = append(slurSets, slurs[s.Number])
				delete(slurs, s.Number)
			}
		}

		// TODO: team beam only support 2 notes grouping
		for _, b := range note.Beam {
			positionY := float64(note.PositionY - 20 + ((b.Number) * 3))

			switch b.Type {
			case musicxml.NoteBeamTypeBegin:
				beams[b.Number] = BeamLine{
					Start: Coordinate{
						X: float64(note.PositionX),
						Y: positionY,
					},
				}
			case musicxml.NoteBeamTypeEnd:

				beam := beams[b.Number]

				if beam.Start.X == 0 {
					beams[b.Number] = BeamLine{
						Start: Coordinate{
							X: float64(note.PositionX),
							Y: positionY,
						},
						End: Coordinate{
							X: float64(note.PositionX) + 8,
							Y: positionY,
						},
					}

				} else {
					beam.End = Coordinate{
						X: float64(note.PositionX) + 8,
						Y: beam.Start.Y,
					}
					beams[b.Number] = beam
				}

				beamSets = append(beamSets, beams[b.Number])

				delete(beams, b.Number)

			}
		}

		if note.Tie != nil {
			if note.Tie.Type == musicxml.NoteSlurTypeStart {
				ties[note.Note] = SlurBezier{
					Start: CoordinateWithOctave{
						Coordinate: Coordinate{
							X: float64(note.PositionX),
							Y: float64(note.PositionY),
						},
						Octave: note.Octave,
					},
				}
			} else if note.Tie.Type == musicxml.NoteSlurTypeStop {
				temp := ties[note.Note]
				temp.End = CoordinateWithOctave{
					Coordinate: Coordinate{
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

	RenderBezier(slurSets, canvas)
	RenderBezier(tiesSet, canvas)

	canvas.Group("class='beam'")
	for _, b := range beamSets {
		canvas.Line(
			int(math.Round(b.Start.X)),
			int(math.Round(b.Start.Y)),
			int(math.Round(b.End.X)),
			int(math.Round(b.End.Y)),
			"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.2",
		)
	}
	canvas.Gend()
}

func cleanBeamByNumber(ctx context.Context, notes []*NoteRenderer, beamNumber int) ([]*NoteRenderer, []beamSplitMarker) {

	switches := map[int]beamMarker{}

	markers := make([]beamSplitMarker, 0)

	var prev *NoteRenderer

	for indexNote, note := range notes {

		if indexNote == len(notes)-1 { // ignore last note or no beam
			continue
		}

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

				prev.Beam[beamNumber] = Beam{
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
			newBeam := map[int]Beam{}

			for k, v := range note.Beam {
				newBeam[k] = v
			}

			switches[beamNumber] = beamMarker{
				NoteBeamType:   musicxml.NoteBeamTypeBegin,
				NoteBeginIndex: indexNote,
			}

			newBeam[beamNumber] = Beam{
				Number: beamNumber,
				Type:   musicxml.NoteBeamTypeBegin,
			}
			note.Beam = newBeam
		} else {

			if prev == nil {
				continue
			}

			if _, hasBeam := note.Beam[beamNumber]; hasBeam {
				newBeam := map[int]Beam{}

				for k, v := range note.Beam {
					newBeam[k] = v
				}

				switches[beamNumber] = beamMarker{
					NoteBeamType:   musicxml.NoteBeamTypeContinue,
					NoteBeginIndex: switches[beamNumber].NoteBeginIndex,
				}

				newBeam[beamNumber] = Beam{
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

				prev.Beam[beamNumber] = Beam{
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

				newBeam[beamNumber] = Beam{
					Type:   musicxml.NoteBeamTypeEnd,
					Number: beamNumber,
				}

				prev.Beam = newBeam

				if t, ok := switches[beamNumber]; ok {

					markers = append(markers, beamSplitMarker{
						StartIndex: t.NoteBeginIndex,
						EndIndex:   prev.indexPosition,
					})
				}

			} else {
				if _, ok := switches[beamNumber]; !ok {
					newBeam := prev.Beam
					newBeam[beamNumber] = Beam{
						Type:   musicxml.NoteBeamTypeBackwardHook,
						Number: beamNumber,
					}
					prev.Beam = newBeam
				}
				markers = append(markers, beamSplitMarker{
					StartIndex: prev.indexPosition,
					EndIndex:   prev.indexPosition,
				})
			}

		}
	}

	return notes, markers
}

// TODO: 7 and 8 notes
func splitBeam(ctx context.Context, notes []*NoteRenderer, segments map[int][]beamSplitMarker) []*NoteRenderer {

	measure, _ := ctx.Value(contextKeyMeasure).(int)

	_ = measure
	if len(segments[1]) == 0 && len(segments[2]) == 0 {
		return notes
	}

	if len(segments[2]) == 0 {
		for _, segment := range segments[1] {
			diff := (segment.EndIndex - segment.StartIndex) + 1
			switch diff {
			case 4: // split 2x2
				notes[segment.StartIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			case 5: // split 3 x 2
				notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			case 6: // spilt 3 x 3
				notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
				notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
			}
		}

		return notes
	}

	for _, segment := range segments[1] {
		diff := (segment.EndIndex - segment.StartIndex) + 1
		totalSubSegment := 0

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
			// 4 beat split to 2x4 : if there is no sub segment
			// 8 beat split to 5x3
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
					if startingPoint <= 1 {
						// split 3x2
						notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[segment.StartIndex+3].UpdateBeam(1, musicxml.NoteBeamTypeBegin)

					} else {
						// split 2x3
						notes[segment.StartIndex+1].UpdateBeam(1, musicxml.NoteBeamTypeEnd)
						notes[segment.StartIndex+2].UpdateBeam(1, musicxml.NoteBeamTypeBegin)
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
				} else {
					// split 3x2
					notes[subSegment.StartIndex+2].UpdateBeam(2, musicxml.NoteBeamTypeEnd)
					notes[subSegment.StartIndex+3].UpdateBeam(2, musicxml.NoteBeamTypeBegin)
				}
			}

		} else {

		}
	}

	return notes
}
