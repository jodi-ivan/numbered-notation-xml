package splitter

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func splitTuplet(notes []*entity.NoteRenderer, segments []BeamSplitMarker) []BeamSplitMarker {
	marker := BeamSplitMarker{-1, -1}

	result := []BeamSplitMarker{}
	for i, n := range notes {
		if n.Tuplet == nil {
			continue
		}
		var segment BeamSplitMarker
		for _, v := range segments {
			if i >= v.StartIndex {
				segment = v
				break
			}
			result = append(result, v)
		}

		switch n.Tuplet.Type {
		case musicxml.TupletTypeStart:
			marker.StartIndex = i
			if segment.StartIndex < i {
				result = append(result, BeamSplitMarker{
					StartIndex: segment.StartIndex,
					EndIndex:   marker.StartIndex - 1,
				})
			}
		case musicxml.TupletTypeStop:
			marker.EndIndex = i
			result = append(result, marker)

			if segment.StartIndex < marker.StartIndex && segment.StartIndex > 0 {
				notes[marker.StartIndex-1].UpdateBeamWithLock(1, musicxml.NoteBeamTypeEnd)
			}
			notes[marker.StartIndex].UpdateBeamWithLock(1, musicxml.NoteBeamTypeBegin)
			notes[marker.EndIndex].UpdateBeamWithLock(1, musicxml.NoteBeamTypeEnd)
			if marker.EndIndex < segment.EndIndex {
				notes[marker.EndIndex+1].UpdateBeamWithLock(1, musicxml.NoteBeamTypeBegin)

			}
			result = append(result, marker)
			if segment.EndIndex > marker.EndIndex {
				result = append(result, BeamSplitMarker{
					StartIndex: marker.EndIndex + 1,
					EndIndex:   segment.EndIndex,
				})
			}

			for s := marker.StartIndex; s < marker.EndIndex; s++ {
				notes[s].UpdateBeamWithLock(1, musicxml.NoteBeamTypeContinue)
			}
			marker = BeamSplitMarker{-1, -1}
		}

	}
	if len(result) == 0 {
		return segments
	}
	return result
}
