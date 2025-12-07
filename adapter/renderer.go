package adapter

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/julienschmidt/httprouter"
)

type CanvasDelegatorHTTP struct {
	w http.ResponseWriter
	r *http.Request
}

func (cdh *CanvasDelegatorHTTP) OnBeforeStartWrite() {
	cdh.w.WriteHeader(http.StatusOK)
	cdh.w.Header().Set("Content-Type", "image/svg+xml")
}

func (cdh *CanvasDelegatorHTTP) OnError(err error) canvas.DelegatorErrorFlowControl {
	if errors.Is(err, repository.ErrHymnNotFound) {
		// metadata is not found
		return canvas.DelegatorErrorFlowControlIgnore

	} else if errors.Is(err, repository.ErrHymnHasMoreThanOneVariant) {
		// Perform the redirect
		http.Redirect(cdh.w, cdh.r, cdh.r.URL.Path+"a", http.StatusSeeOther)

		return canvas.DelegatorErrorFlowControlStop

	}
	cdh.w.WriteHeader(http.StatusInternalServerError)
	cdh.w.Write([]byte(err.Error()))

	return canvas.DelegatorErrorFlowControlStop
}

func New(u usecase.Usecase) *RenderHTTP {
	return &RenderHTTP{
		usecase: u,
	}
}

type RenderHTTP struct {
	usecase usecase.Usecase
}

func (rh *RenderHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	delegator := &CanvasDelegatorHTTP{w: w, r: r}
	canv := canvas.NewCanvasWithDelegator(svg.New(w), delegator)
	raw := ps.ByName("number")

	var variant []string
	num, err := strconv.Atoi(raw)
	if err != nil {
		num, err = strconv.Atoi(raw[0 : len(raw)-1])
		if err != nil {
			log.Printf("invalid number: %v", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid URL"))
			return
		}
		variant = []string{string(raw[len(raw)-1])}
	}

	err = rh.usecase.RenderHymn(r.Context(), canv, num, variant...)
	if err != nil {
		delegator.OnError(err)
	}

}
