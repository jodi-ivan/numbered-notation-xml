package adapter

import (
	"log"
	"net/http"

	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	customerror "github.com/jodi-ivan/numbered-notation-xml/utils/errors"
	"github.com/jodi-ivan/numbered-notation-xml/utils/webserver"
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
	ctx := r.Context()
	term := r.FormValue("q")
	if term == "" {
		log.Println("[SearchAdapterHTTP] invalid query q ")
		webserver.RenderErrorResponse(w, http.StatusBadRequest, customerror.New("Invalid query"))
		return
	}

	hits, err := sah.Usecase.Lookup(ctx, term)
	if err != nil {
		log.Printf("[SearchAdapterHTTP] failed to lookup, err: %v", err.Error())
		webserver.RenderErrorResponse(w, http.StatusInternalServerError, customerror.NewFromError(err, "Failed to search"))
		return
	}

	highlights := sah.Usecase.Highlight(ctx, hits)

	webserver.RenderSuccessResponse(w, nil, highlights)

}
