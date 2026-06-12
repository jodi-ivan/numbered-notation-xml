package verse

import (
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func GetLineSyllableBound(metadata *entity.HymnMetaData) []int {
	result := []int{}
	verse, ok := metadata.ParsedVerse[2]
	if !ok {
		return result
	}

	syllLine := make([]int, len(verse))
	for i, line := range verse {
		totalSyllablePerLine := 0
		for _, syll := range line {
			totalSyllablePerLine += len(syll.Breakdown)
		}
		syllLine[i] = totalSyllablePerLine
		if i > 0 {
			syllLine[i] += syllLine[i-1]
		}
	}

	return syllLine
}

func BuildContent(music musicxml.MusicXML, metadata *entity.HymnMetaData) [][]entity.LyricWordVerse {

	lMapper := map[int][]entity.LyricWordVerse{}

	prevTotalLyric := -1
	wordVerses := map[int]entity.LyricWordVerse{}
	for _, measure := range music.Part.Measures {
		measure.Build()
		for _, note := range measure.Notes {
			if len(note.Lyric) == 0 {
				continue
			}

			if prevTotalLyric != len(note.Lyric) && prevTotalLyric != -1 {
				lMapper[1] = append(lMapper[1], lMapper[2]...)
				lMapper[2] = []entity.LyricWordVerse{}
			}

			for _, l := range note.Lyric {

				syl := ""

				part := entity.LyricPartVerse{
					Text: syl,
					Type: l.Syllabic,
				}

				if len(l.Text) > 1 {
					part.Combine = true
					part.Breakdown = []entity.LyricStylePart{}
				}
				for _, s := range l.Text {
					syl += s.Value
					part.Text += s.Value

					part.Breakdown = append(part.Breakdown, entity.LyricStylePart{
						Text:      s.Value,
						Underline: s.Underline != 0,
					})
				}

				// li := lyric.NewLyric()
				if unicode.IsDigit(rune(syl[0])) {
					syl = syl[2:]
				}
				wordVerse := wordVerses[l.Number]
				if len(wordVerse.Breakdown) == 0 {
					wordVerse.Breakdown = []entity.LyricPartVerse{}
				}
				wordVerse.Word += syl
				wordVerse.Breakdown = append(wordVerse.Breakdown, part)

				wordVerses[l.Number] = wordVerse

				if l.Syllabic == musicxml.LyricSyllabicTypeEnd || l.Syllabic == musicxml.LyricSyllabicTypeSingle {
					lMapper[l.Number] = append(lMapper[l.Number], wordVerse)
					wordVerses[l.Number] = entity.LyricWordVerse{
						Breakdown: []entity.LyricPartVerse{},
					}

				}

			}

			prevTotalLyric = len(note.Lyric)
		}
	}

	for part := 2; part <= 4; part++ {
		if len(lMapper[part]) > 0 {
			lMapper[1] = append(lMapper[1], lMapper[part]...)
		}
	}

	return SplitLyricsIntoLines(lMapper[1], GetLineSyllableBound(metadata))

}

func SplitLyricsIntoLines(words []entity.LyricWordVerse, maxSyllables []int) [][]entity.LyricWordVerse {
	result := make([][]entity.LyricWordVerse, len(maxSyllables))

	wordIdx := 0  // index into words[]
	consumed := 0 // syllables consumed across all lines so far

	for lineIdx, ceiling := range maxSyllables {
		var lineWords []entity.LyricWordVerse

		for wordIdx < len(words) {
			w := words[wordIdx]
			syllableCount := len(w.Breakdown)

			// Stop if adding this word exceeds the line's ceiling
			if consumed+syllableCount > ceiling {
				break
			}

			lineWords = append(lineWords, w)
			consumed += syllableCount
			wordIdx++
		}

		result[lineIdx] = lineWords
	}

	return result
}
