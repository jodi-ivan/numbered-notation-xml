package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

// add the lyric here
var wholeLyric string = `
`

func main() {

	result := []string{}

	lines := strings.Split(strings.TrimSpace(wholeLyric), "\n")

	for _, line := range lines {
		lineResult := []string{}
		words := strings.Fields(strings.TrimSpace(line))

		for _, word := range words {
			breakdown := []string{}
			syllables := lyric.SplitSyllable(word)

			breakdownFmt := `{"word":"%s", "breakdown":[%s]}`

			if len(syllables) == 1 {
				info := fmt.Sprintf(`{"text": "%s","type":"%s"}`, syllables[0], musicxml.LyricSyllabicTypeSingle)
				breakdown = append(breakdown, info)
			}
			for i, syll := range syllables {
				var t musicxml.LyricSyllabic
				if i == len(syllables)-1 {
					t = musicxml.LyricSyllabicTypeEnd
				} else if i == 0 {
					t = musicxml.LyricSyllabicTypeBegin
				} else {
					t = musicxml.LyricSyllabicTypeMiddle
				}
				info := fmt.Sprintf(`{"text": "%s","type":"%s"}`, syll, t)
				breakdown = append(breakdown, info)
			}

			lineResult = append(lineResult, fmt.Sprintf(breakdownFmt, word, strings.Join(breakdown, ",")))

		}

		result = append(result, "["+strings.Join(lineResult, ",")+"]")
	}

	log.Println("[", strings.Join(result, ","), "]")

}
