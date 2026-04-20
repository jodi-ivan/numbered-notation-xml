package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/verse"
	"github.com/jodi-ivan/numbered-notation-xml/utils/storage"
)

type WordBreakdown struct {
	Word      string `json:"word"`
	Breakdown []string
}

// Result holds the full processing output for a sentence.
type Result struct {
	Breakdown string   // e.g. "ma-(lai)-kat se-la-lu"
	NotInDB   []string // clean words that were not found in DB
}

// WordResult holds per-word processing output.
type WordResult struct {
	Breakdown string
	InDB      bool
}

type Line []WordBreakdown

func main() {
	dbPath := `/home/jodiivan/go/src/github.com/jodi-ivan/numbered-notation-xml/files/database/kidung-jemaat.db`
	db, err := storage.NewStorage(context.Background(), dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to storage: %s", err.Error())
		return
	}

	defer db.Close()

	verses := `
	_malaikat _kekuatan _Daud _Tahu
	`

	lines := strings.Split(verses, "\n")

	for _, l := range lines {
		res, _ := ProcessSentence(db, l)
		raw, _ := json.MarshalIndent(res, "", "    ")
		log.Println(string(raw))
	}

}

// splitTrailing splits a word into its alpha core and trailing symbols.
// Trailing = punctuation (.,!?;:) and quotes (" ' ")
func splitTrailing(word string) (core, trail string) {
	i := len(word)
	for i > 0 {
		r := rune(word[i-1])
		if unicode.IsPunct(r) || r == '"' {
			i--
		} else {
			break
		}
	}
	return word[:i], word[i:]
}

// isElided returns true and strips the underscore prefix.
func isElided(word string) (string, bool) {
	if strings.HasPrefix(word, "_") {
		return word[1:], true
	}
	return word, false
}

// applyElision wraps the elision_index-th syllable (1-based) in parentheses.
// e.g. breakdown="ma-lai-kat", elisionIndex=1 → "ma-(lai)-kat"
func applyElision(breakdown string, elisionIndex int) string {
	parts := strings.Split(breakdown, "-")
	if elisionIndex+1 > len(parts) {
		return breakdown
	}

	log.Println(elisionIndex, parts)
	parts[elisionIndex] = "_" + parts[elisionIndex]
	return strings.Join(parts, "-")
}

// fallbackBreakdown is your own algorithm for words not in DB.
// Replace the body with your real implementation.
func fallbackBreakdown(word string) string {
	return strings.Join(verse.SplitSyllable(word), "-")
}

// ProcessSentence looks up each word in the DB, applies elision if marked,
// falls back to your algorithm for unknown words, and preserves casing/trailing symbols.
func ProcessSentence(db *sqlx.DB, sentence string) (Result, error) {
	tokens := strings.Fields(sentence)
	var breakdownParts []string
	var notInDB []string

	for _, token := range tokens {
		// 1. detect elision marker
		raw, elided := isElided(token)

		// 2. split trailing punctuation/quotes
		core, trail := splitTrailing(raw)

		pretrail := ""
		if core[0] == '"' {
			core = core[1:]
			pretrail = ""
		}

		// 3. preserve original casing for output; query in lowercase
		lookup := strings.ToLower(core)

		// 4. query DB
		bd, elisionIdx, found, err := queryWord(db, lookup, elided)
		if err != nil {
			return Result{}, err
		}

		var wordBreakdown string
		if found {
			if elided && elisionIdx.Valid {
				wordBreakdown = applyElision(bd, int(elisionIdx.Int64))
			} else {
				wordBreakdown = bd
			}
		} else {
			// not in DB — use fallback algorithm
			wordBreakdown = fallbackBreakdown(lookup)
			if elided {
				wordBreakdown = "_" + wordBreakdown
				lookup = "_" + lookup
			}
			notInDB = append(notInDB, lookup) // original casing, no trail
		}

		res := syncCasing(core, wordBreakdown)

		log.Println(core, wordBreakdown)

		// 5. re-attach trailing symbols
		breakdownParts = append(breakdownParts, pretrail+res+trail)
	}

	return Result{
		Breakdown: strings.Join(breakdownParts, " "),
		NotInDB:   notInDB,
	}, nil
}

func syncCasing(original, hyphenated string) string {
	var result strings.Builder
	origRunes := []rune(original)
	origIdx := 0

	prefixes := map[int]string{}

	for _, char := range hyphenated {
		if char == '-' {
			result.WriteRune('-')
			continue
		}

		if char == '_' {
			prefixes[origIdx] = "_"
			continue
		}

		if p, ok := prefixes[origIdx]; ok {
			result.WriteRune(rune(p[0]))
		}
		// log.Println("char should be", string(origRunes[origIdx]))
		// If the original character was uppercase, make this one uppercase
		if origIdx < len(origRunes) && unicode.IsUpper(origRunes[origIdx]) {
			result.WriteRune(unicode.ToUpper(origRunes[origIdx]))
		} else {
			result.WriteRune(unicode.ToLower(char))
		}
		origIdx++
	}

	return result.String()
}

// queryWord fetches breakdown and elision_index from DB.
// When elided=true it looks for the row WHERE elision_index IS NOT NULL,
// otherwise WHERE elision_index IS NULL.
func queryWord(db *sqlx.DB, word string, elided bool) (breakdown string, elisionIdx sql.NullInt64, found bool, err error) {
	var query string
	if elided {
		query = `
			SELECT breakdown, elision_index
			FROM syllable_breakdown
			WHERE whole = ? AND elision_index IS NOT NULL
			LIMIT 1`
	} else {
		query = `
			SELECT breakdown, elision_index
			FROM syllable_breakdown
			WHERE whole = ? AND elision_index IS NULL
			LIMIT 1`
	}
	row := db.QueryRow(query, strings.ToLower(word))
	err = row.Scan(&breakdown, &elisionIdx)
	if err == sql.ErrNoRows {
		return "", sql.NullInt64{}, false, nil
	}
	if err != nil {
		return "", sql.NullInt64{}, false, err
	}
	return breakdown, elisionIdx, true, nil
}
