package verse

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
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

func LoadOtherVerse(notes []*entity.NoteRenderer, metadata *repository.HymnMetadata, startPos int, prevRepeatInfos []*musicxml.RepeatInfo) int {

	verse, ok := metadata.Verse[2] // for now hardcoded two for testing visual
	if !ok {
		return 0
	}

	// TODO: double parsing on the verse.go
	whole := [][]LyricWordVerse{}

	err := json.Unmarshal([]byte(verse.Content.String), &whole)
	if err != nil {
		log.Printf("[LoadOtherVerse] failed to unmarshal for verse, err %s\n", err.Error())
		return 0
	}

	totalSyllable := 0

	flattenSyll := []LyricPartVerse{}
	flattenCombine := map[int]bool{}
	for _, line := range whole {
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

		appendedLyric := lyric.GetMusicxmlLyric(note) // load the lyric on the current music
		if lyricNum == 0 {
			lyricNum = len(appendedLyric) + 1
		}

		txt := flattenSyll[syll].Text
		if lyric.HasPrefix(note) {
			txt = "2." + txt // for now, hardcoded. if the lyric has prefix.
		}

		newLyric := []musicxml.Lyric{
			{
				Text:     ApplyElision(txt, flattenCombine[syll]),
				Syllabic: flattenSyll[syll].Type,
				Number:   lyricNum,
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
						Verse: 2,
						Text: []musicxml.LyricText{
							{}, // add the lyric to them
						},
						Syllabic: musicxml.LyricSyllabicTypeSingle,
						Number:   lyricNum,
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
							Text:     ApplyElision(flattenSyll[syll].Text, flattenCombine[syll]),
							Syllabic: flattenSyll[syll].Type,
							Number:   lyricNum,
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
					Verse: 2,
					Text: []musicxml.LyricText{
						{Value: flattenSyll[syllRepeat].Text}, // add the lyric to them
					},
					Syllabic: flattenSyll[syllRepeat].Type,
					Number:   newLyric[0].Number + 1,
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
