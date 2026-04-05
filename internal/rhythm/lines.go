package rhythm

import (
	"context"
	"fmt"
	"math"
	"sort"

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
		hasTies := false
		if note.Tie != nil {
			hasTies = true
			switch note.Tie.Type {
			case musicxml.NoteSlurTypeStart:
				ties[note.Note] = SlurBezier{
					SlurTieType: SlurTieTypeTie,
					Start: CoordinateWithOctave{
						Coordinate: entity.NewCoordinate(float64(note.PositionX), float64(note.PositionY)),
						Octave:     note.Octave,
					},
					LineType: note.Tie.LineType,
				}
			case musicxml.NoteSlurTypeStop:
				temp := ties[note.Note]
				temp.End = CoordinateWithOctave{
					Coordinate: entity.NewCoordinate(float64(note.PositionX), float64(note.PositionY)),
					Octave:     note.Octave,
				}
				ties[note.Note] = temp

				tiesSet = append(tiesSet, ties[note.Note])
				delete(ties, note.Note)
			}
		}

		for _, s := range note.Slur {
			yPos := note.PositionY
			if hasTies {
				yPos += 3
			}
			if s.Type == musicxml.NoteSlurTypeStop || s.Type == musicxml.NoteSlurTypeHop {
				temp := slurs[s.Number]
				temp.End = CoordinateWithOctave{
					Coordinate: entity.NewCoordinate(float64(note.PositionX-2), float64(yPos)),
					Octave:     note.Octave,
				}

				if temp.Start.X == 0 && temp.Start.Y == 0 {
					temp.Start = CoordinateWithOctave{
						Coordinate: entity.NewCoordinate(float64(note.PositionX-constant.UPPERCASE_LENGTH), float64(yPos)),
						Octave:     0,
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
						Coordinate: entity.NewCoordinate(float64(note.PositionX+2), float64(yPos)),
						Octave:     note.Octave,
					},
					LineType: s.LineType,
				}
			}

		}

	}

	if len(slurs) > 0 { // there is start, but no end
		for _, slur := range slurs {
			temp := slur
			if temp.End.Coordinate.X == 0 && temp.End.Coordinate.Y == 0 {
				temp.End = CoordinateWithOctave{
					Coordinate: entity.NewCoordinate(float64(maxXPosition-5), float64(temp.Start.Y)),
					Octave:     0,
				}
			}
			slurSets = append(slurSets, temp)
		}
	}

	ri.RenderBezier(tiesSet, canv)
	ri.RenderBezier(slurSets, canv)

}

func (ri *rhythmInteractor) RenderBeam(ctx context.Context, canv canvas.Canvas, ts timesig.TimeSignature, notes []*entity.NoteRenderer) {

	beams := map[int]BeamLine{}
	beamSets := []BeamLine{}

	ri.BeamSplitter.Split(ctx, notes, ts, nil)
	for _, note := range notes {

		for _, b := range note.Beam {
			positionY := float64(note.PositionY - 22 + ((b.Number) * 3))

			switch b.Type {
			case musicxml.NoteBeamTypeBegin:
				beams[b.Number] = BeamLine{
					Start:  entity.NewCoordinate(float64(note.PositionX), positionY),
					Number: b.Number,
				}
			case musicxml.NoteBeamTypeEnd:
				beam := beams[b.Number]

				if beam.Start.X == 0 {
					beams[b.Number] = BeamLine{
						Start:  entity.NewCoordinate(float64(note.PositionX), positionY),
						End:    entity.NewCoordinate(float64(note.PositionX)+8, positionY),
						Number: b.Number,
					}
				} else {
					beam.End = entity.NewCoordinate(float64(note.PositionX)+8, beam.Start.Y)
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
