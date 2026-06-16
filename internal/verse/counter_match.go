package verse

import (
	"context"
	"fmt"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

type SyllableMatch interface {
	IsVowel(char rune) bool
	ApplyElision(syllText string, combine bool) []musicxml.LyricText
	LoadOtherVerse(ctx context.Context, notes []*entity.NoteRenderer, metadata *entity.HymnMetaData, startPos int, offset map[int]int, prevRepeatInfos []*musicxml.RepeatInfo) (map[int]int, int)
	LoadVerse(ctx context.Context, targetVerse int, clear bool, notes []*entity.NoteRenderer, metadata *entity.HymnMetaData, startPos int, prevRepeatInfos []*musicxml.RepeatInfo) (int, int)
}

func NewSyllableMatcher() SyllableMatch {
	return &matcher{}
}

type matcher struct {
}

func (m *matcher) IsVowel(char rune) bool {
	return utils.Contains([]string{"a", "i", "u", "e", "o"}, strings.ToLower(string(char))) >= 0
}

func (m *matcher) ApplyElision(syllText string, combine bool) []musicxml.LyricText {
	start, end := -1, -1

	if !combine {
		return []musicxml.LyricText{
			{Value: syllText},
		}
	}

	if defElistion, ok := defaultElision[syllText]; ok {
		start, end = defElistion[0], defElistion[1]
	} else {
		for ir, r := range syllText {
			if m.IsVowel(r) || (r == 'h' && start != -1) {
				if start == -1 {
					start = ir
				} else {
					end = ir
				}
			}
		}

	}

	if end < start {
		return []musicxml.LyricText{
			{Value: syllText},
		}
	}

	partBreakdown := []musicxml.LyricText{}

	if start == 0 {
		partBreakdown = []musicxml.LyricText{
			{Underline: 1, Value: syllText[start : end+1]},
			{Underline: 0, Value: syllText[end+1:]},
		}
	} else if end == len(syllText)-1 {
		partBreakdown = []musicxml.LyricText{
			{Underline: 0, Value: syllText[0:start]},
			{Underline: 1, Value: syllText[start:]},
		}
	} else {
		partBreakdown = []musicxml.LyricText{
			{Underline: 0, Value: syllText[0:start]},
			{Underline: 1, Value: syllText[start : end+1]},
			{Underline: 0, Value: syllText[end+1:]},
		}

	}

	return partBreakdown

}
func (m *matcher) LoadOtherVerse(ctx context.Context, notes []*entity.NoteRenderer, metadata *entity.HymnMetaData, startPos int, offset map[int]int, prevRepeatInfos []*musicxml.RepeatInfo) (map[int]int, int) {
	prm, _ := params.GetParamFromContext(ctx)

	if offset == nil {
		offset = map[int]int{}
	}

	if prm.Verse < 2 {
		return offset, 0
	}
	targetVerse := 2
	if prm.Verse > 1 {
		targetVerse = prm.Verse
	}

	if targetVerse > 2 && !prm.SingleVerseMode {
		// load previous verse
		prvOffset, _ := m.LoadVerse(ctx, targetVerse-1, true, notes, metadata, startPos+offset[targetVerse-1], prevRepeatInfos)
		offset[targetVerse-1] += prvOffset
	}

	targetVerseOffset, margin := m.LoadVerse(ctx, targetVerse, prm.SingleVerseMode, notes, metadata, startPos+offset[targetVerse], prevRepeatInfos)
	offset[targetVerse] += targetVerseOffset
	return offset, margin

}

func (m *matcher) fillableByLyric(n *entity.NoteRenderer) bool {
	return !(n.Barline != nil || breathpause.IsBreathMark(n) || n.IsDotted)
}

func (m *matcher) LoadVerse(ctx context.Context, targetVerse int, clear bool, notes []*entity.NoteRenderer, metadata *entity.HymnMetaData, startPos int, prevRepeatInfos []*musicxml.RepeatInfo) (int, int) {

	prm, _ := params.GetParamFromContext(ctx)
	verse, ok := metadata.ParsedVerse[targetVerse]
	if !ok {
		return 0, 0
	}

	offset := 0

	totalSyllable := 0

	flattenSyll := []entity.LyricPartVerse{}
	flattenCombine := map[int]bool{}
	for _, line := range verse {
		for _, syll := range line {
			if syll.VerseOnly {
				continue
			}
			bd := syll.Breakdown
			for i, comb := range syll.Breakdown {
				bd[i].Load1stVerse = syll.Load1stVerse
				flattenCombine[totalSyllable+i] = comb.Combine
			}
			flattenSyll = append(flattenSyll, bd...)
			totalSyllable += len(syll.Breakdown)
		}
	}

	syll := startPos
	var repeatInfo *musicxml.RepeatInfo

	lastOffset := 0
	if len(prevRepeatInfos) > 0 {
		repeatInfo = prevRepeatInfos[len(prevRepeatInfos)-1]
		lastOffset = repeatInfo.OffsetStart
		if syll >= lastOffset {
			syll = lastOffset*2 + (syll - lastOffset)
		}
	}

	marginBottom := 0
	insert := true
	lyricNum := 0
	firstNote := false
	li := lyric.NewLyric()
	for i := 0; i < len(notes) && syll < len(flattenSyll); i++ {

		note := notes[i]

		if !m.fillableByLyric(note) {
			continue
		}
		if len(note.Lyric) == 0 {
			if flattenSyll[syll].Offset == 1 {
				// fill up the empty notes with current syllable (shift left)
				offset++
			} else {
				continue

			}
		} else if flattenSyll[syll].Offset == 1 {
			offset++
		}
		if flattenSyll[syll].Offset == -1 { // fill the current notes with empty syllable. shift right
			flattenSyll[syll].Offset = 0
			offset--
			appendedLyric := lyric.GetMusicxmlLyric(note)
			if clear {
				appendedLyric = []musicxml.Lyric{}
			}
			lyricVerse := 0
			if len(appendedLyric) > 0 {
				lyricVerse = 2
			}
			appendedLyric = append(appendedLyric, musicxml.Lyric{
				Number:   len(appendedLyric) + 1,
				Syllabic: musicxml.LyricSyllabicTypeMiddle,
				Verse:    lyricVerse,
			})

			li.SetLyricRenderer(note, appendedLyric)
			continue
		}
		load1stVerse := false
		if clear && flattenSyll[syll].Load1stVerse {
			if prm.Verse > 2 && prm.Verse-1 == targetVerse {
				continue
			}

			load1stVerse = prm.Verse > 1
		}
		appendedLyric := []musicxml.Lyric{}
		if !clear || load1stVerse {
			appendedLyric = lyric.GetMusicxmlLyric(note) // load the lyric on the current music
			if load1stVerse {
				for i, al := range appendedLyric {
					al.Verse = 2
					appendedLyric[i] = al
				}
				li.SetLyricRenderer(note, appendedLyric)
				continue
			}
		}
		if lyricNum == 0 {
			lyricNum = len(appendedLyric) + 1
		}

		txt := flattenSyll[syll].Text
		hasPrefix := lyric.HasPrefix(note)
		start1stVerse := !flattenSyll[syll].Load1stVerse && (syll > 0 && flattenSyll[syll-1].Load1stVerse)
		if ((hasPrefix || syll == 0) || start1stVerse) && txt != "" {
			txt = fmt.Sprintf("%d. %s", targetVerse, txt)
			if len(appendedLyric) > 0 && !hasPrefix {
				lastLyric := appendedLyric[0]
				lastLyric.Text[0].Value = fmt.Sprintf("%d. %s", targetVerse-1, lastLyric.Text[0].Value)
				appendedLyric[0] = lastLyric
			}
		}

		verseIndicator := 2
		if (targetVerse != 2 && prm.Verse == verseIndicator-1) || prm.SingleVerseMode {
			verseIndicator = 0

		}

		newLyric := []musicxml.Lyric{
			{
				Text:     m.ApplyElision(txt, flattenCombine[syll]),
				Syllabic: flattenSyll[syll].Type,
				Number:   lyricNum,
				Verse:    verseIndicator,
			},
		}
		insertLastMeasure := syll < lastOffset*2
		if repeatInfo != nil && repeatInfo.BarlineEnding != nil {
			switch repeatInfo.BarlineEnding.Number {
			case "1":
				decendnewLyric := newLyric[0]
				decendnewLyric.Number++
				newLyric[0] = decendnewLyric
				insert = false

			case "2":

				if note.MeasureNumber == repeatInfo.MeasureNumber {

					newLyric = []musicxml.Lyric{{
						Number: lyricNum, Verse: 2,
						Text:     []musicxml.LyricText{{}},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
					}}
					insertLastMeasure = syll < lastOffset*2
					syll = prevRepeatInfos[len(prevRepeatInfos)-2].SyllCntEnd - 1
					insert = true
				} else if note.MeasureNumber > repeatInfo.MeasureNumber && len(prevRepeatInfos) > 1 && prevRepeatInfos[len(prevRepeatInfos)-2].BarlineEnding != nil {
					if !firstNote {
						syll -= ((repeatInfo.OffsetStart - repeatInfo.SyllCntStart) + (prevRepeatInfos[len(prevRepeatInfos)-2].OffsetStart - prevRepeatInfos[len(prevRepeatInfos)-2].SyllCntEnd)) + 1
						firstNote = true
					}

					newLyric = []musicxml.Lyric{
						{
							Number: lyricNum, Verse: 2, Syllabic: flattenSyll[syll].Type,
							Text: m.ApplyElision(flattenSyll[syll].Text, flattenCombine[syll]),
						},
					}
					insert = false

				}

			}

		}

		if insert && repeatInfo != nil && insertLastMeasure {

			syllRepeat := repeatInfo.OffsetStart + (syll % repeatInfo.OffsetStart)
			if syllRepeat < len(flattenSyll) {
				newLyric = append(newLyric, musicxml.Lyric{
					Number: newLyric[0].Number + 1, Verse: 2,
					Text:     m.ApplyElision(flattenSyll[syllRepeat].Text, flattenCombine[syllRepeat]),
					Syllabic: flattenSyll[syllRepeat].Type,
				})
			}

		}
		appendedLyric = append(appendedLyric, newLyric...)
		info := li.SetLyricRenderer(note, appendedLyric) // the one that calculate the margin, padding and width of the notes.
		if marginBottom < info.MarginBottom {
			marginBottom = info.MarginBottom
		}

		if syll == lastOffset-1 {
			syll = lastOffset*2 + (syll - lastOffset)
		}
		syll++

	}

	if prm.Diagnostic != nil {
		res := map[int]bool{targetVerse: syll == len(flattenSyll)}
		prm.Diagnostic.VerseSyllMatch <- res

	}

	return offset, marginBottom

}
