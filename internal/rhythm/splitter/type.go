package splitter

import "github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"

type beamMarker struct {
	NoteBeamType   musicxml.NoteBeamType
	NoteBeginIndex int
}

type BeamSplitMarker struct {
	StartIndex int
	EndIndex   int
}

type Interval []BeamSplitMarker

func (c Interval) Len() int           { return len(c) }
func (c Interval) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Interval) Less(i, j int) bool { return c[i].StartIndex < c[j].StartIndex }
