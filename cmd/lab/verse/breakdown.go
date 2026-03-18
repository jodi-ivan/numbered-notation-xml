package verse

import (
	"database/sql"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
)

type WordBreakdown struct {
	Word      string `json:"word"`
	Breakdown []string
}

// WordResult holds per-word processing output.
type WordResult struct {
	Breakdown string
	InDB      bool
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
	parts[elisionIndex] = "_" + parts[elisionIndex]
	return strings.Join(parts, "-")
}

// fallbackBreakdown is your own algorithm for words not in DB.
// Replace the body with your real implementation.
func fallbackBreakdown(word string) string {
	return strings.Join(SplitSyllable(word), "-")
}

// ProcessSentence looks up each word in the DB, applies elision if marked,
// falls back to your algorithm for unknown words, and preserves casing/trailing symbols.
func ProcessSentence(db *sqlx.DB, sentence string) ([]string, map[string]string, error) {
	notInDB := map[string]string{}
	result := []string{}

	lines := strings.Split(sentence, "\n")

	for _, line := range lines {
		tokens := strings.Fields(line)

		if len(tokens) == 0 {
			continue
		}

		var breakdownParts []string
		for _, token := range tokens {

			// 1. detect elision marker
			raw, elided := isElided(token)

			// 2. split trailing punctuation/quotes
			core, trail := splitTrailing(raw)

			pretrail := ""
			if core[0] == '"' || core[0] == '\'' {
				pretrail = string(core[0])
				core = core[1:]
			} else if unicode.IsNumber(rune(core[0])) {
				breakdownParts = append(breakdownParts, core+trail)
				continue

			}

			// 3. preserve original casing for output; query in lowercase
			lookup := strings.ToLower(core)

			// 4. query DB
			bd, elisionIdx, found, err := queryWord(db, lookup, elided)
			if err != nil {
				return nil, nil, err
			}
			var wordBreakdown string

			if found {
				if elided && elisionIdx.Valid {
					wordBreakdown = applyElision(bd, int(elisionIdx.Int64))
				} else {
					wordBreakdown = bd
				}
			} else {
				// // not in DB — use fallback algorithm
				// notInDB = append(notInDB, core) // original casing, no trail
				wordBreakdown = fallbackBreakdown(lookup)
				if elided {
					wordBreakdown = "_" + wordBreakdown
					lookup = "_" + lookup
				}
				notInDB[lookup] = wordBreakdown

			}

			res := syncCasing(core, wordBreakdown)

			// 5. re-attach trailing symbols
			breakdownParts = append(breakdownParts, pretrail+res+trail)
		}
		result = append(result, strings.Join(breakdownParts, " "))
	}

	return result, notInDB, nil

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

		// If the original character was uppercase, make this one uppercase
		if origIdx < len(origRunes) && unicode.IsUpper(origRunes[origIdx]) {
			result.WriteRune(unicode.ToUpper(char))
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
