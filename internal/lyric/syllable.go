package lyric

import (
	"strings"
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
)

// FIXME: number of syllable is equal to number of vowel. it will be false, if word kacau should be splitted [ka][cau]
func SplitSyllable(word string) []string {
	vowels := []string{"a", "i", "u", "e", "o", "A", "I", "U", "E", "O"}

	vowelIndex := []int{}

	for i, char := range word {
		if index := utils.Contains(vowels, string(char)); index >= 0 {
			vowelIndex = append(vowelIndex, i)
		}
	}

	result := []string{}
	for _, v := range vowelIndex {
		syllable := string(word[v])

		// previous char
		if v > 0 && v <= len(word) {
			prevChar := string(word[v-1])
			if !IsVowel(prevChar) {
				syllable = prevChar + syllable
				lowered := strings.ToLower(prevChar)
				if lowered == "y" || lowered == "g" { // case ya or ga
					if v-2 >= 0 {
						// get the prev of the prev
						prevOfPrev := string(word[v-2])
						if strings.ToLower(prevOfPrev) == "n" { // case nya or nga
							syllable = prevOfPrev + syllable
						}
					}
				}
			}
		}

		// next char

		if v+1 < len(word) {
			nextChar := string(word[v+1])
			if !unicode.IsLetter([]rune(nextChar)[0]) {
				syllable = syllable + nextChar
			} else if !IsVowel(nextChar) {
				if strings.ToLower(nextChar) == "n" {
					if v+2 < len(word) {
						nextToNext := string(word[v+2])
						if strings.ToLower(nextToNext) == "y" || strings.ToLower(nextToNext) == "g" {
							// syllable = syllable + nextChar
							if v+2 == len(word)-1 {
								syllable = syllable + nextChar + nextToNext
							} else if v+3 < len(word) {
								next3Char := string(word[v+3])
								if !IsVowel(next3Char) {
									syllable = syllable + nextChar + nextToNext
								}
								if !unicode.IsLetter([]rune(next3Char)[0]) {
									syllable = syllable + next3Char
								}
							}
						} else {
							syllable = syllable + nextChar
							if v+3 < len(word) {
								next3Char := string(word[v+3])
								if !unicode.IsLetter([]rune(next3Char)[0]) {
									syllable = syllable + next3Char
								}
							}
						}
					} else {
						syllable = syllable + nextChar
						if v+3 < len(word) {
							nextToNext := string(word[v+3])
							if !unicode.IsLetter([]rune(nextToNext)[0]) {
								syllable = syllable + nextToNext
							}
						}

					}
				} else {

					if v+2 < len(word) {
						nextToNext := string(word[v+2])
						if !IsVowel((nextToNext)) {
							syllable = syllable + nextChar
						}

						if !unicode.IsLetter([]rune(nextToNext)[0]) {
							syllable = syllable + nextToNext
						}
					} else {
						syllable = syllable + nextChar

					}
				}
			}
		}

		result = append(result, syllable)

	}

	return result
}

func IsVowel(char string) bool {
	return utils.Contains([]string{"a", "i", "u", "e", "o"}, strings.ToLower(char)) >= 0
}
