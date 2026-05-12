package gregorian

import (
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
)

func GetGroupSlueTies(notes []*entity.NoteRenderer, staffLine lines.LineStaff) []SlurTieGroup {
	groupBeamSlurTies := []SlurTieGroup{}

	var tiesTracking *SlurTieGroup
	slurTracking := map[int]SlurTieGroup{}

	for _, note := range notes {
		yPos := 0.0
		direction := 0
		if note.AbsoluteNote != "" && note.AbsoluteOctave > 0 {
			yPos = staffLine.GetYPos(rune(note.AbsoluteNote[0]), note.AbsoluteOctave)
			direction = staffLine.GetStemDirection(rune(note.AbsoluteNote[0]), note.AbsoluteOctave)
		}

		if note.Tie != nil {
			if tiesTracking == nil && note.Tie.Type == musicxml.NoteSlurTypeStart && !note.Tie.NumberedOnly {
				tiesTracking = &SlurTieGroup{
					MaxY: yPos, MinY: yPos,
					Ties:       note.Tie,
					NoteMember: []string{note.UUID},
					Start:      entity.NewCoordinate(float64(note.PositionX), yPos),
				}
				tiesTracking.NoteMember = append(tiesTracking.NoteMember, note.UUID)
				tiesTracking.AccumulativeDirection += direction
			}

			if tiesTracking != nil && note.Tie.Type == musicxml.NoteSlurTypeStop {
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
						MaxY: float64(staffLine.GetTopLine()),
						MinY: float64(staffLine.GetBottomLine()),
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
					temp.MaxY = float64(staffLine.GetTopLine())
					temp.MinY = float64(staffLine.GetBottomLine())
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
