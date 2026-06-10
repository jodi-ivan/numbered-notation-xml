package verse

import (
	"context"
	"fmt"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

func IsVowel(char rune) bool {
	return utils.Contains([]string{"a", "i", "u", "e", "o"}, strings.ToLower(string(char))) >= 0
}

func ApplyElision(syllText string, combine bool) []musicxml.LyricText {
	start, end := -1, -1

	if !combine {
		return []musicxml.LyricText{
			{Value: syllText},
		}
	}

	for ir, r := range syllText {
		if IsVowel(r) || (r == 'h' && start != -1) {
			if start == -1 {
				start = ir
			} else {
				end = ir
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
func LoadOtherVerse(ctx context.Context, notes []*entity.NoteRenderer, metadata *entity.HymnMetaData, startPos int, prevRepeatInfos []*musicxml.RepeatInfo) int {
	prm, _ := params.GetParamFromContext(ctx)

	if prm.Verse == 0 {
		return 0
	}
	targetVerse := 2
	if prm.Verse > 1 {
		targetVerse = prm.Verse
	}

	if targetVerse > 2 {
		// load previous verse
		LoadVerse(ctx, targetVerse-1, true, notes, metadata, startPos, prevRepeatInfos)
	}

	return LoadVerse(ctx, targetVerse, false, notes, metadata, startPos, prevRepeatInfos)

}

func LoadVerse(ctx context.Context, targetVerse int, clear bool, notes []*entity.NoteRenderer, metadata *entity.HymnMetaData, startPos int, prevRepeatInfos []*musicxml.RepeatInfo) int {

	prm, _ := params.GetParamFromContext(ctx)

	verse, ok := metadata.ParsedVerse[targetVerse]
	if !ok {
		return 0
	}

	totalSyllable := 0

	flattenSyll := []entity.LyricPartVerse{}
	flattenCombine := map[int]bool{}
	for _, line := range verse {
		for _, syll := range line {
			flattenSyll = append(flattenSyll, syll.Breakdown...)
			for i, comb := range syll.Breakdown {
				flattenCombine[totalSyllable+i] = comb.Combine
			}
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
		if len(note.Lyric) == 0 {
			continue
		}

		appendedLyric := []musicxml.Lyric{}
		if !clear {
			appendedLyric = lyric.GetMusicxmlLyric(note) // load the lyric on the current music
		}
		if lyricNum == 0 {
			lyricNum = len(appendedLyric) + 1
		}

		txt := flattenSyll[syll].Text
		hasPrefix := lyric.HasPrefix(note)
		if hasPrefix || syll == 0 {
			txt = fmt.Sprintf("%d. %s", targetVerse, txt)
			if len(appendedLyric) > 0 && !hasPrefix {
				lastLyric := appendedLyric[0]
				lastLyric.Text[0].Value = fmt.Sprintf("%d. %s", targetVerse-1, lastLyric.Text[0].Value)
				appendedLyric[0] = lastLyric
			}
		}

		verseIndicator := 2
		if targetVerse != 2 && prm.Verse == verseIndicator-1 {
			verseIndicator = 0

		}

		newLyric := []musicxml.Lyric{
			{
				Text:     ApplyElision(txt, flattenCombine[syll]),
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
							Number: lyricNum, Verse: 2,
							Text:     ApplyElision(flattenSyll[syll].Text, flattenCombine[syll]),
							Syllabic: flattenSyll[syll].Type,
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
					Text:     ApplyElision(flattenSyll[syllRepeat].Text, flattenCombine[syllRepeat]),
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

	return marginBottom

}
