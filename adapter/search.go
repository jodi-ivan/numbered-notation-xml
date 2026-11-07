package adapter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/julienschmidt/httprouter"
)

func NewSearchAdapter(u usecase.Usecase) *SearchAdapterHTTP {
	return &SearchAdapterHTTP{
		Usecase: u,
	}
}

type SearchAdapterHTTP struct {
	Usecase usecase.Usecase
}

func (sah *SearchAdapterHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	term := r.FormValue("q")
	if term == "" {
		log.Println("[SearchAdapterHTTP] invalid query q ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Query can't be empty"))
		return
	}

	hits, err := sah.Usecase.Lookup(r.Context(), term)
	if err != nil {
		log.Printf("[SearchAdapterHTTP] failed to lookup, err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Query can't be empty"))
		return
	}

	raw, _ := json.Marshal(hits)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(raw)

}
