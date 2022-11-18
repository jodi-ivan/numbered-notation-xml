package keysig

import (
	"fmt"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
)

func getSpelledNumberedNotation(letterNotation string) string {
	if strings.HasSuffix(letterNotation, "#") {
		return fmt.Sprintf("%sis", strings.ToLower(string(letterNotation[0])))
	}

	if strings.HasSuffix(letterNotation, "b") {
		if utils.Contains([]string{"A", "E"}, string(letterNotation[0])) >= 0 {
			return fmt.Sprintf("%ss", strings.ToLower(string(letterNotation[0])))
		}

		return fmt.Sprintf("%ses", strings.ToLower(string(letterNotation[0])))
	}

	return strings.ToLower(letterNotation)

}
