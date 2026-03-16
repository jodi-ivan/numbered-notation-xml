package splitter

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func CleanBeamByNumber(ctx context.Context, notes []*entity.NoteRenderer, no int) []BeamSplitMarker {

	switches := map[int]beamMarker{}

	markers := make([]BeamSplitMarker, 0)

	var prev *entity.NoteRenderer

	for noteIdx, note := range notes {
		note.IndexPosition = noteIdx

		if len(note.Beam) == 0 { // stopping the beam
			if noteIdx == 0 {
				prev = note
				continue
			} else {

				t, ok := switches[no]
				if !ok {
					prev = note
					continue
				}

				prev.Beam[no] = entity.Beam{Number: no, Type: musicxml.NoteBeamTypeEnd}
				markers = append(markers, BeamSplitMarker{StartIndex: t.NoteBeginIndex, EndIndex: noteIdx - 1})
				delete(switches, no)

			}
		}

		if t, ok := switches[no]; !ok {
			if _, hasBeam := note.Beam[no]; !hasBeam {
				prev = note
				continue
			}
			newBeam := map[int]entity.Beam{}

			for k, v := range note.Beam {
				newBeam[k] = v
			}

			switches[no] = beamMarker{NoteBeamType: musicxml.NoteBeamTypeBegin, NoteBeginIndex: noteIdx}
			newBeam[no] = entity.Beam{Number: no, Type: musicxml.NoteBeamTypeBegin}
			note.Beam = newBeam
		} else {

			if prev == nil {
				continue
			}

			if _, hasBeam := note.Beam[no]; hasBeam {
				newBeam := map[int]entity.Beam{}

				for k, v := range note.Beam {
					newBeam[k] = v
				}

				switches[no] = beamMarker{
					NoteBeamType:   musicxml.NoteBeamTypeContinue,
					NoteBeginIndex: switches[no].NoteBeginIndex,
				}

				newBeam[no] = entity.Beam{
					Number: no,
					Type:   musicxml.NoteBeamTypeContinue,
				}
				note.Beam = newBeam
				prev = note
				continue
			}

			if t.NoteBeamType == musicxml.NoteBeamTypeBegin || t.NoteBeamType == musicxml.NoteBeamTypeContinue {

				if _, ok := prev.Beam[no]; !ok {
					prev = note
					continue
				}

				prev.Beam[no] = entity.Beam{
					Number: no,
					Type:   musicxml.NoteBeamTypeEnd,
				}

				delete(switches, no)

				markers = append(markers, BeamSplitMarker{
					StartIndex: t.NoteBeginIndex,
					EndIndex:   noteIdx - 1,
				})

			}

		}
		prev = note

	}

	if prev != nil && len(prev.Beam) > 0 {
		additional, ok := prev.Beam[no]

		if ok {
			if additional.Type != musicxml.NoteBeamTypeEnd {
				newBeam := prev.Beam

				newBeam[no] = entity.Beam{
					Type:   musicxml.NoteBeamTypeEnd,
					Number: no,
				}

				prev.Beam = newBeam

				if t, ok := switches[no]; ok {

					markers = append(markers, BeamSplitMarker{
						StartIndex: t.NoteBeginIndex,
						EndIndex:   prev.IndexPosition,
					})
				}

			} else {
				if _, ok := switches[no]; !ok {
					newBeam := prev.Beam
					newBeam[no] = entity.Beam{
						Type:   musicxml.NoteBeamTypeBackwardHook,
						Number: no,
					}
					prev.Beam = newBeam
				}
				markers = append(markers, BeamSplitMarker{
					StartIndex: prev.IndexPosition,
					EndIndex:   prev.IndexPosition,
				})
			}

		}
	}

	return markers
}
