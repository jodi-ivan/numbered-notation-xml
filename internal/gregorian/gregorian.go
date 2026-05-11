package gregorian

import (
	"cmp"
	"context"
	"fmt"
	"math"
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

// DEPRECATED:
func GetYpos(lines [5]int, space int, octave int, pitch rune) float64 {
	noteOrder := []rune{'C', 'D', 'E', 'F', 'G', 'A', 'B'}

	diatonicIndex := func(p rune, oct int) int {
		for i, n := range noteOrder {
			if n == unicode.ToUpper(p) {
				return oct*7 + i
			}
		}
		return -1
	}

	refIndex := diatonicIndex('F', 5) // lines[0] = F5
	noteIndex := diatonicIndex(pitch, octave)

	stepsBelow := refIndex - noteIndex

	return float64(lines[0]) + float64(stepsBelow)*(float64(space)/2)
}

func GetGroupSlueTies(notes []*entity.NoteRenderer, lines [5]int) []SlurTieGroup {
	groupBeamSlurTies := []SlurTieGroup{}

	var tiesTracking *SlurTieGroup
	slurTracking := map[int]SlurTieGroup{}

	for _, note := range notes {
		yPos := 0.0
		direction := 0
		if note.AbsoluteNote != "" && note.AbsoluteOctave > 0 {
			yPos = GetYpos(lines, STAFF_SPACE_WIDTH, note.AbsoluteOctave, rune(note.AbsoluteNote[0]))
			direction = cmp.Compare(int(yPos), lines[2])
			if direction == 0 {
				direction = -1
			}
		}

		if note.Tie != nil {
			if tiesTracking == nil && note.Tie.Type == musicxml.NoteSlurTypeStart && !note.Tie.NumberedOnly {
				tiesTracking = &SlurTieGroup{
					MaxY:       yPos,
					MinY:       yPos,
					Ties:       note.Tie,
					NoteMember: []string{note.UUID},
					Start:      entity.NewCoordinate(float64(note.PositionX), yPos),
				}
				tiesTracking.NoteMember = append(tiesTracking.NoteMember, note.UUID)
				tiesTracking.AccumulativeDirection += direction
			}

			if tiesTracking != nil && note.Tie.Type == musicxml.NoteSlurTypeStop {
				tiesTracking.NoteMember = append(tiesTracking.NoteMember, note.UUID)
				tiesTracking.AccumulativeDirection += direction

				tiesTracking.End = entity.NewCoordinate(float64(note.PositionX), yPos)
				tiesTracking.NoteMember = append(tiesTracking.NoteMember, note.UUID)

				groupBeamSlurTies = append(groupBeamSlurTies, *tiesTracking)
				tiesTracking = nil
			}
		}

		for sid, slur := range note.Slur {
			_, ok := slurTracking[sid]

			if !ok {
				if yPos == 0 {
					slurTracking[sid] = SlurTieGroup{
						MaxY: float64(lines[0]),
						MinY: float64(lines[4]),
					}
				} else {
					slurTracking[sid] = SlurTieGroup{
						MaxY: yPos,
						MinY: yPos,
					}
				}
			}

			temp := slurTracking[sid]
			pos := entity.NewCoordinate(float64(note.PositionX), yPos)

			switch slur.Type {
			case musicxml.NoteSlurTypeStop, musicxml.NoteSlurTypeHop:
				temp.End = pos
			case musicxml.NoteSlurTypeStart:
				temp.Start = pos
				temp.Slur = &slur
			}

			temp.NoteMember = append(temp.NoteMember, note.UUID)
			temp.AccumulativeDirection += direction
			if yPos != 0 {
				temp.MaxY = math.Max(temp.MaxY, yPos)
				temp.MinY = math.Min(temp.MinY, yPos)
			}

			slurTracking[sid] = temp

			if slur.Type == musicxml.NoteSlurTypeStop || slur.Type == musicxml.NoteSlurTypeHop {
				groupBeamSlurTies = append(groupBeamSlurTies, temp)
				delete(slurTracking, sid)
			}

			if slur.Type == musicxml.NoteSlurTypeHop {
				temp := SlurTieGroup{
					NoteMember:            []string{note.UUID},
					Start:                 entity.NewCoordinate(float64(note.PositionX), yPos),
					AccumulativeDirection: temp.AccumulativeDirection,
					Slur:                  &slur,
				}

				if yPos == 0 {
					temp.MaxY = float64(lines[0])
					temp.MinY = float64(lines[4])
				} else {
					temp.MaxY = yPos
					temp.MinY = yPos
				}
				slurTracking[sid] = temp

			}
		}

		if len(note.Slur) > 0 {
			continue
		}

		for i, v := range slurTracking {
			v.AccumulativeDirection += direction
			if yPos != 0 {
				v.MaxY = math.Max(v.MaxY, yPos)
				v.MinY = math.Min(v.MinY, yPos)
			}
			v.NoteMember = append(v.NoteMember, note.UUID)
			slurTracking[i] = v
		}

		if tiesTracking != nil && yPos != 0 {
			tiesTracking.MaxY = math.Max(tiesTracking.MaxY, yPos)
			tiesTracking.MinY = math.Min(tiesTracking.MinY, yPos)
		}

	}

	if tiesTracking != nil {
		groupBeamSlurTies = append(groupBeamSlurTies, *tiesTracking)

	}
	for _, v := range slurTracking {
		groupBeamSlurTies = append(groupBeamSlurTies, v)
	}
	return groupBeamSlurTies
}

func RenderStaffLine(ctx context.Context, staffPos, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) VMargin {
	canv.Group(`class="gregorian"`, "style='font-family:mozart11'")

	lineStaff := NewLineStaff(timeSignature, keySignature)
	lineStaff.Render(canv, y, notes[0].MeasureNumber, staffPos == 0)
	lines := lineStaff.GetLines()
	margin := VMargin{
		Top:           entity.NewCoordinate(0, float64(lines[0])),
		Bottom:        entity.NewCoordinate(0, float64(lines[4])),
		DefaultTop:    lines[0],
		DefaultBottom: lines[4],
	}

	groupBeam := [][]CoordinateWithNoteLength{{}}

	canv.Group(`class="notes"`, `style="font-size:2em"`)
	currentMeasure := 0

	groupBeamSlurTies := GetGroupSlueTies(notes, lines)

	for i, note := range notes {

		if currentMeasure != note.MeasureNumber {
			currentMeasure = note.MeasureNumber
			if i != 0 {
				canv.Gend()
			}
			canv.Group(`class="measure"`, fmt.Sprintf(`number="%d"`, currentMeasure))

		}
		if note.IsAdditional {
			continue
		}

		if breathpause.IsBreathMark(note) {

			xPos := float64(note.PositionX)
			if note.PositionX-notes[i-1].PositionX <= numbered.MIN_DISTANCE_BREATH {
				xPos += (numbered.AVERAGE_CHARACTER_WIDTH + constant.LOWERCASE_LENGTH) / 3
			}

			canv.TextUnescaped(xPos, float64(lines[0])-STAFF_SPACE_WIDTH, "&#xF0E2;", `style="font-size:1.3em"`)
			if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
				groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
			}
			continue
		}
		if note.IsRest {
			canv.TextUnescaped(float64(note.PositionX), float64(lines[2]), restHex[note.NoteLength])
			if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
				groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
			}
			continue
		}

		if note.Barline != nil {
			barline.RenderGregorian(canv, note.Barline, i == len(notes)-1, lines, entity.NewCoordinate(float64(note.PositionX), float64(lines[4])))
			continue
		}

		if note.AbsoluteNote == "" {
			continue
		}

		var noteMargin VMargin
		pairs := []SlurTieGroup{}
		noteMargin, groupBeam, pairs = RenderNote(ctx, canv, lines, groupBeam, groupBeamSlurTies, i, notes, timeSignature, keySignature)
		margin.Merge(noteMargin)

		groupBeamSlurTies = append(groupBeamSlurTies, pairs...)

	}
	canv.Gend()

	for _, gr := range groupBeam {
		if len(gr) == 0 {
			continue
		}
		gMargin := RenderGroupBeam(canv, gr, lines, groupBeamSlurTies)
		margin.Merge(gMargin)
	}

	canv.Gend()

	st := RenderSlurTies(canv, lineStaff, groupBeam, groupBeamSlurTies)
	margin.Merge(st)

	canv.Gend()

	return margin
}

func GetLeftIndentWithTimeSignature(key keysig.Key, timeSig timesig.TimeSignature) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return constant.LAYOUT_INDENT_LENGTH + CLEF_WIDTH + (timesig.GREGORIAN_WIDTH * len(timeSig.UniqueSign)) + (PADDING_WIDTH*(3+(len(timeSig.UniqueSign)-1)) + keySigWith)
}

// DEPRECATED: use the staffLine object instead
func GetLeftIndent(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	offset := 0
	if key.Start && key.Prev != nil {
		offset = (len(key.Prev.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH) + PADDING_WIDTH
	}
	return constant.LAYOUT_INDENT_LENGTH + CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith + offset
}

func GetLeftMargin(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith
}
