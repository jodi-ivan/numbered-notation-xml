package adapter

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/verse"
	"github.com/julienschmidt/httprouter"
)

type WordBreakdown struct {
	Word      string `json:"word"`
	Breakdown []string
}

type Line []WordBreakdown

type LyricParser struct{}

func (lp *LyricParser) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	result := []Line{}
	input := strings.ReplaceAll(strings.TrimSpace(string(b)), "\\t", "")
	lines := strings.Split(input, "\\n")

	for _, l := range lines {
		line := []WordBreakdown{}
		words := strings.Fields(l)
		if len(words) == 0 {
			continue
		}
		for _, w := range words {
			syllable := verse.SplitSyllable(w)
			line = append(line, WordBreakdown{
				Word:      w,
				Breakdown: syllable,
			})
		}
		result = append(result, line)
	}

	raw, _ := json.MarshalIndent(result, "", "    ")
	w.WriteHeader(http.StatusOK)
	w.Write(raw)
}

type LyricParserV2 struct {
	Db *sqlx.DB
}

func (lpv2 LyricParserV2) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	input := strings.ReplaceAll(strings.TrimSpace(string(b)), "\\t", "")
	breakdown, notindb, err := verse.ProcessSentence(lpv2.Db, input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	data := map[string]interface{}{
		"breakdown": breakdown,
		"generated": notindb,
	}

	raw, _ := json.MarshalIndent(data, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(raw)

}
