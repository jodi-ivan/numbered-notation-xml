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

func LoadOtherVerse(notes []*entity.NoteRenderer, metadata *repository.HymnMetadata, startPos int) int {

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

	marginBottom := 0

	li := lyric.NewLyric()
	for i := 0; i < len(notes) && syll < len(flattenSyll); i++ {
		note := notes[i]
		if len(note.Lyric) == 0 {
			continue
		}

		newLyric := lyric.GetMusicxmlLyric(note) // load the lyric on the current music
		txt := flattenSyll[syll].Text
		if lyric.HasPrefix(note) {
			txt = "2." + txt // for now, hardcoded. if the lyric has prefix.
		}
		newLyric = append(newLyric, musicxml.Lyric{
			Text: []musicxml.LyricText{
				{Value: txt}, // add the lyric to them
			},
			Syllabic: flattenSyll[syll].Type,
			Number:   len(newLyric) + 1,
		})

		info := li.SetLyricRenderer(note, newLyric) // the one that calculate the margin, padding and width of the notes.
		if marginBottom < info.MarginBottom {
			marginBottom = info.MarginBottom
		}

		syll++

	}

	return marginBottom

}
