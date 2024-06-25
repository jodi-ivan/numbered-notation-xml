package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
)

type WordBreakdown struct {
	Word      string `json:"word"`
	Breakdown []string
}

type Line []WordBreakdown

func main() {

	result := []Line{}

	verses := `
		Kerubim dan serafim
		memuliakan Yang Trisuci; 
		para rasul dan nabi, 
		martir yang berjubah putih 
		G'reja yang kudus, esa, 
		kepadaMu menyembah. 
	`

	lines := strings.Split(verses, "\n")

	for _, l := range lines {
		line := []WordBreakdown{}
		words := strings.Fields(l)
		if len(words) == 0 {
			continue
		}
		for _, w := range words {
			syllable := lyric.SplitSyllable(w)
			line = append(line, WordBreakdown{
				Word:      w,
				Breakdown: syllable,
			})
		}
		result = append(result, line)
	}

	raw, _ := json.MarshalIndent(result, "", "    ")
	log.Println(string(raw))
}
