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

func NewUnifiedSearchAdapter(u usecase.Usecase) *SearchAdapterHTTP {
	return &SearchAdapterHTTP{
		Usecase: u,
		unified: true,
	}
}

type SearchAdapterHTTP struct {
	unified bool
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

	if sah.unified && len(term) < 3 {
		log.Printf("[SearchAdapterHTTP] query term too short: %s", term)
		webserver.RenderErrorResponse(w, http.StatusBadRequest, customerror.New("Query must be at least 3 characters long"))
		return
	}

	var data interface{}

	hits, err := sah.Usecase.Lookup(ctx, term)
	if err != nil {
		log.Printf("[SearchAdapterHTTP] failed to lookup, err: %v", err.Error())
		webserver.RenderErrorResponse(w, http.StatusInternalServerError, customerror.NewFromError(err, "Failed to search"))
		return
	}
	highlights := sah.Usecase.Highlight(ctx, hits)
	data = highlights

	if sah.unified {
		autocomplete, err := sah.Usecase.Autocomplete(ctx, term)
		if err != nil {
			log.Printf("[SearchAdapterHTTP] failed to autocomplete, err: %v", err.Error())
			webserver.RenderErrorResponse(w, http.StatusInternalServerError, customerror.NewFromError(err, "Failed to search"))
			return
		}

		data = map[string]interface{}{
			"autocomplete": autocomplete,
			"search":       highlights,
		}
	}

	webserver.RenderSuccessResponse(w, nil, data)

}
