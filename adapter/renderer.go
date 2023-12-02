package adapter

import (
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
}

func (cdh *CanvasDelegatorHTTP) OnBeforeStartWrite() {
	cdh.w.WriteHeader(http.StatusOK)
	cdh.w.Header().Set("Content-Type", "image/svg+xml")
}

func (cdh *CanvasDelegatorHTTP) OnError(err error) canvas.DelegatorErrorFlowControl {
	if err == repository.ErrHymnNotFound {
		// metadata is not found
		return canvas.DelegatorErrorFlowControlIgnore
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
	canv := canvas.NewCanvasWithDelegator(svg.New(w), &CanvasDelegatorHTTP{w: w})

	num, err := strconv.Atoi(ps.ByName("number"))
	if err != nil {
		log.Printf("invalid number: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return
	}

	rh.usecase.RenderHymn(r.Context(), canv, num)

}
