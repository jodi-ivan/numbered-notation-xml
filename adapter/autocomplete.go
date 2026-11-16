package adapter

import (
	"log"
	"net/http"

	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	customerror "github.com/jodi-ivan/numbered-notation-xml/utils/errors"
	"github.com/jodi-ivan/numbered-notation-xml/utils/webserver"
	"github.com/julienschmidt/httprouter"
)

func AutocompleteAdapter(u usecase.Usecase) *AutocompleteAdapterHTTP {
	return &AutocompleteAdapterHTTP{
		Usecase: u,
	}
}

type AutocompleteAdapterHTTP struct {
	Usecase usecase.Usecase
}

func (aah *AutocompleteAdapterHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	term := r.FormValue("q")
	if term == "" {
		log.Println("[AutocompleteAdapterHTTP] invalid query q ")
		webserver.RenderErrorResponse(w, http.StatusBadRequest, customerror.New("Invalid query"))
		return
	}

	if len(term) < 3 {
		log.Printf("[AutocompleteAdapterHTTP] query term too short: %s", term)
		webserver.RenderErrorResponse(w, http.StatusBadRequest, customerror.New("Query must be at least 3 characters long"))
		return
	}

	result, err := aah.Usecase.Autocomplete(ctx, term)
	if err != nil {
		log.Printf("[AutocompleteAdapterHTTP] failed to autocomplete, err: %v", err.Error())
		webserver.RenderErrorResponse(w, http.StatusInternalServerError, customerror.NewFromError(err, "Failed to search"))
		return
	}

	webserver.RenderSuccessResponse(w, nil, result)

}
