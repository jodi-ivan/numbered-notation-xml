package adapter

import (
	"net/http"
	"strings"
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/document"
	webutils "github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	customerror "github.com/jodi-ivan/numbered-notation-xml/utils/errors"
	"github.com/jodi-ivan/numbered-notation-xml/utils/search"
	"github.com/jodi-ivan/numbered-notation-xml/utils/webserver"

	"github.com/julienschmidt/httprouter"
)

type KeywordsAdapter struct {
	Repo   repository.Repository
	Search search.Provider
}

func (ka *KeywordsAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	n, _, err := webutils.ParseHymnWithVariant(ps.ByName("number"))
	if err != nil {
		webserver.RenderErrorResponse(w, http.StatusBadRequest, customerror.NewFromError(err, "Invalid request"))
		return
	}
	doc, err := document.BuildDocument(ka.Repo, n)
	if err != nil {
		webserver.RenderErrorResponse(w, http.StatusInternalServerError, customerror.NewFromError(err, "Failed to get the document"))
		return
	}

	allText := doc.Title + " " + doc.Content
	for _, verse := range doc.Verses {
		allText += " " + verse.Content
	}

	words := strings.Fields(strings.ToLower(allText))
	uniqueWords := make(map[string]bool)
	var cleanedWords []string

	for _, word := range words {
		// Clean punctuation and numbers
		cleanedWord := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || r == '\'' {
				return r
			}
			return -1 // Drop non-letters
		}, word)

		if cleanedWord != "" {
			if !uniqueWords[cleanedWord] {
				uniqueWords[cleanedWord] = true
				cleanedWords = append(cleanedWords, cleanedWord)
			}

			for _, v := range ka.Search.Analyzer(cleanedWord) {
				if !uniqueWords[v] {
					cleanedWords = append(cleanedWords, v)
					uniqueWords[v] = true

				}
			}

		}
	}

	webserver.RenderSuccessResponse(w, nil, cleanedWords)
}
