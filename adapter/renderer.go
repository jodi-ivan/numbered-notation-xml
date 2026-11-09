package adapter

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
	"github.com/julienschmidt/httprouter"
)

type CanvasDelegatorHTTP struct {
	w http.ResponseWriter
	r *http.Request
}

func (cdh *CanvasDelegatorHTTP) OnBeforeStartWrite() {
	cdh.w.Header().Set("Content-Type", "image/svg+xml")
	cdh.w.WriteHeader(http.StatusOK)
}

func (cdh *CanvasDelegatorHTTP) OnError(err error) canvas.DelegatorErrorFlowControl {
	if errors.Is(err, repository.ErrHymnNotFound) {
		// metadata is not found
		return canvas.DelegatorErrorFlowControlIgnore

	} else if errors.Is(err, repository.ErrHymnHasMoreThanOneVariant) {
		// Perform the redirect
		cdh.r.URL.Path += "a"
		http.Redirect(cdh.w, cdh.r, cdh.r.URL.RequestURI(), http.StatusSeeOther)
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
	canv := canvas.NewBufferedCanvas(w, delegator)
	raw := ps.ByName("number")

	var variant []string
	num, err := strconv.Atoi(raw)
	if err != nil {
		if len(raw) == 0 {
			log.Printf("[ServeHTTP] invalid number: %v", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid URL"))
			return
		}
		num, err = strconv.Atoi(raw[0 : len(raw)-1])
		if err != nil {
			log.Printf("[ServeHTTP] invalid number: %v", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid URL"))
			return
		}
		variant = []string{string(raw[len(raw)-1])}
	}

	verseRaw := r.FormValue("verse")

	verseNo, err := strconv.Atoi(verseRaw)
	if verseRaw != "" && err != nil {
		log.Printf("[ServeHTTP] invalid verse: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return
	}

	mode := r.FormValue("focus")

	focusMode, err := strconv.ParseBool(mode)
	if mode != "" && err != nil {
		log.Printf("[ServeHTTP] invalid mode: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return
	}

	gregorianStr := r.FormValue("gregorian")
	gregorian := true
	if gregorianStr != "" {
		b, err := strconv.ParseBool(gregorianStr)
		if err == nil {
			gregorian = b
		}
	}

	var renderParam *params.RenderParam
	deviceTypeStr := r.FormValue("device_type")
	deviceTypeOrientation := params.OrientationPortrait
	deviceType := 0
	if deviceTypeStr != "" {
		dt, err := strconv.Atoi(deviceTypeStr)
		if err != nil {
			log.Printf("[ServeHTTP] invalid device_type: %v", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid URL"))
			return
		}
		deviceType = dt
		deviceTypeOrientationStr := r.FormValue("orientation")
		if deviceTypeOrientationStr != "" {
			dt, err := strconv.Atoi(deviceTypeOrientationStr)
			if err != nil {
				log.Printf("[ServeHTTP] invalid orientation: %v", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid URL"))
				return
			}
			deviceTypeOrientation = params.Orientation(dt)
		}

		renderParam = &params.RenderParam{
			DeviceType:  deviceType,
			Orientation: deviceTypeOrientation,
		}
	}

	prm := &params.Param{
		Gregorian:       gregorian,
		Verse:           verseNo,
		SingleVerseMode: focusMode,
		Render:          renderParam,
	}

	err = rh.usecase.RenderHymn(params.NewParamContext(r.Context(), prm), canv, num, variant...)
	if err != nil {
		delegator.OnError(err)
	}

}
