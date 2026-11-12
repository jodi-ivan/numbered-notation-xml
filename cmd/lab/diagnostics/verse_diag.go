package diagnostics

import (
	"context"
	"log"
	"time"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/verse"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

var verseDiagnostic *VerseDiagnostic

type VerseDiagnostic struct {
	Matcher verse.SyllableMatch
}

func (vd *VerseDiagnostic) IsVowel(char rune) bool {
	return vd.Matcher.IsVowel(char)
}
func (vd *VerseDiagnostic) ApplyElision(syllText string, combine bool) []musicxml.LyricText {
	return vd.Matcher.ApplyElision(syllText, combine)
}
func (vd *VerseDiagnostic) LoadOtherVerse(ctx context.Context, notes []*entity.NoteRenderer, metadata *entity.HymnMetaData, startPos int, offset map[int]int, prevRepeatInfos []*musicxml.RepeatInfo) (map[int]int, int) {
	all := map[int]int{}
	// singleResult := map[int]bool{}

	timeout := time.NewTimer(10 * time.Second)
	defer timeout.Stop()

	prm, _ := params.GetParamFromContext(ctx)
	go func() {

		for {
			select {
			case data := <-prm.Diagnostic.VerseSyllMatch:
				currentData := data
				prm.Diagnostic.Mu.RLock()
				for k, v := range currentData {
					prm.Diagnostic.MapMtx.Store(k, v)
				}

				prm.Diagnostic.Mu.RUnlock()

				prm.Diagnostic.VerseDiagnostic <- params.VerseDiagnostic{
					SingleMode: currentData,
				}

			case <-timeout.C:
				log.Println("timeout")

				prm.Diagnostic.Finish <- true
				return

			}

		}
	}()

	allOffset := map[int]map[int]int{}
	for i := 2; i <= len(metadata.Verse)+1; i++ {
		newParam := &params.Param{
			Verse:           i,
			SingleVerseMode: prm.SingleVerseMode,
			Diagnostic:      prm.Diagnostic,
		}

		rctx := params.NewParamContext(ctx, newParam)
		allOffset[i], _ = vd.Matcher.LoadOtherVerse(rctx, notes, metadata, startPos, offset, prevRepeatInfos)

		for i, v := range allOffset {
			for _, off := range v {
				all[i] = off
			}
		}

	}
	// prm.Diagnostic.Finish <- true

	return all, 0

}
func (vd *VerseDiagnostic) LoadVerse(ctx context.Context, targetVerse int, clear bool, notes []*entity.NoteRenderer, metadata *entity.HymnMetaData, startPos int, prevRepeatInfos []*musicxml.RepeatInfo) (int, int) {
	return vd.Matcher.LoadVerse(ctx, targetVerse, clear, notes, metadata, startPos, prevRepeatInfos)
}

func GetVerseDiagnostic(matcher verse.SyllableMatch) *VerseDiagnostic {
	if verseDiagnostic == nil {
		verseDiagnostic = &VerseDiagnostic{
			Matcher: matcher,
		}
	}

	return verseDiagnostic
}
