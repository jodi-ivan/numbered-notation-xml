package verse

import (
	"encoding/json"
	"log"
	"slices"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
)

// start position
// check if the lyric ever has prefix 1.
// if it does go we start whenever the 1. location.
// else it starts from start

func getFirstPos(notes []*entity.NoteRenderer) int {
	firstNoteRefrein := slices.ContainsFunc(notes[0].MeasureText, func(mt musicxml.MeasureText) bool {
		return mt.Text == "Refrein"
	})

	verseNumPos := -1

	for i, v := range notes {
		hasVerseNum := slices.ContainsFunc(v.Lyric, func(l entity.Lyric) bool {
			return l.Syllabic == musicxml.LyricSyllabicTypeBegin && strings.HasPrefix(l.Text[0].Value, "1.")
		})

		if hasVerseNum {
			verseNumPos = i
			break
		}
	}

	if verseNumPos >= 0 && firstNoteRefrein {
		return verseNumPos
	}

	return 0
}

func LoadOtherVerse(notes []*entity.NoteRenderer, metadata *repository.HymnMetadata, startPos int, repeatInfo *musicxml.RepeatInfo) int {

	verse, ok := metadata.Verse[2] // for know hardcoded two for testing visual
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

	for _, line := range whole {
		for _, syll := range line {
			totalSyllable += len(syll.Breakdown)
			flattenSyll = append(flattenSyll, syll.Breakdown...)
		}
	}

	syll := startPos

	if syll >= 15 {
		log.Println("adding", syll)
		syll = 30 + (syll - 15)
		log.Println("added", syll)
		log.Println("")
	}
	marginBottom := 0

	li := lyric.NewLyric()
	for i := 0; i < len(notes) && syll < len(flattenSyll); i++ {

		note := notes[i]
		if len(note.Lyric) == 0 {
			continue
		}

		appendedLyric := lyric.GetMusicxmlLyric(note) // load the lyric on the current music

		log.Println(syll, flattenSyll[syll].Text)
		newLyric := []musicxml.Lyric{
			{
				Text: []musicxml.LyricText{
					{Value: flattenSyll[syll].Text}, // add the lyric to them
				},
				Syllabic: flattenSyll[syll].Type,
				Number:   len(appendedLyric) + 1,
			},
		}

		if repeatInfo != nil && syll <= repeatInfo.SyllCntEnd*2 {

			syllRepeat := repeatInfo.OffsetStart + (syll % repeatInfo.OffsetStart)
			newLyric = append(newLyric, musicxml.Lyric{
				Verse: 2,
				Text: []musicxml.LyricText{
					{Value: flattenSyll[syllRepeat].Text}, // add the lyric to them
				},
				Syllabic: flattenSyll[syllRepeat].Type,
				Number:   len(appendedLyric) + 2,
			})

		}
		appendedLyric = append(appendedLyric, newLyric...)
		info := li.SetLyricRenderer(note, appendedLyric) // the one that calculate the margin, padding and width of the notes.
		if marginBottom < info.MarginBottom {
			marginBottom = info.MarginBottom
		}

		if syll == 14 {
			log.Println("adding", syll)
			syll = 30 + (syll - 15)
			log.Println("added", syll)
			log.Println("")
		}
		syll++

	}

	return marginBottom

}
