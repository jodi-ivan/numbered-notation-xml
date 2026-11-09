package verse

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type SyllableTracker struct {
	Matcher     SyllableMatch
	MetaData    *entity.HymnMetaData
	ReffAtStart bool
}

func (st *SyllableTracker) TrackSyllable(ctx context.Context, notes []*entity.NoteRenderer, currentMeasureSyll *int, syllStart int, repeatOffset map[int]int, tracker entity.SyllableCounterTracker, firstMeasure, fineEnd bool, repeatInfo []*musicxml.RepeatInfo) (newTracker entity.SyllableCounterTracker, _startSyllable, marginBottom int) {
	if (tracker.StartRenderOtherNotes) || firstMeasure {
		start := syllStart

		// matcher := st.Matcher
		// if p, ok := params.GetParamFromContext(ctx); ok && p.Diagnostic != nil {
		// 	matcher = diagnostics.GetVerseDiagnostic(matcher)
		// }
		tracker.RepeatInfo = repeatInfo
		tracker.SyllableOffset, marginBottom = st.Matcher.LoadOtherVerse(ctx, notes, st.MetaData, start, repeatOffset, repeatInfo)
		syllStart += *currentMeasureSyll
		tracker.StartRenderOtherNotes = true
	} else {
		*currentMeasureSyll = 0
		syllStart = 0
		tracker.SyllableCount = 0

		if !st.ReffAtStart && fineEnd {
			tracker.StartRenderOtherNotes = true
		}
	}

	return tracker, syllStart, marginBottom

}
