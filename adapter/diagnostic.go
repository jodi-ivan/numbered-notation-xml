package adapter

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/google/uuid"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
	"github.com/julienschmidt/httprouter"
)

type DiagnosticHTTP struct {
	Usecase   usecase.Usecase
	Interrupt chan os.Signal
	Repo      repository.Repository
}

// SSEvent represents a single server-sent event packet.
type SSEvent struct {
	ID    string
	Event string
	Data  []byte
	Retry int
}

func (dh *DiagnosticHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	target := svg.New(io.Discard)

	canv := canvas.NewCanvasWithDelegator(target, &CanvasDelegator{})

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	no, vars, err := utils.ParseHymnWithVariant(ps.ByName("scope"))
	if err != nil {
		w.Write([]byte("Invalid param" + err.Error()))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		return
	}

	mode := r.FormValue("focus")

	focusMode, err := strconv.ParseBool(mode)
	if mode != "" && err != nil {
		log.Printf("[ServeHTTP] invalid mode: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		return
	}

	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()

	dig := &params.DiagParam{
		VerseSyllMatch:  make(chan map[int]bool, 1),
		VerseDiagnostic: make(chan params.VerseDiagnostic),
		Finish:          make(chan bool, 1),
		Mu:              &sync.RWMutex{},
		MapMtx:          &sync.Map{},
	}

	go func() {
		prm := &params.Param{
			Verse:           3, // hardcoded to trigger the load verse mechanism
			SingleVerseMode: focusMode,
			Diagnostic:      dig,
		}
		rctx := params.NewParamContext(r.Context(), prm)
		v := []string{}
		if vars != "" {
			v = append(v, vars)
		}
		err := dh.Usecase.RenderHymn(rctx, canv, no, v...)
		log.Println("try render", err)
	}()

	for {
		select {

		case <-timeout.C:
			log.Println("timeout")
			body := []string{
				"id: " + uuid.NewString() + "\n",
				"event: close\n",
				"data: stream_finished \n\n",
			}

			w.Write([]byte(strings.Join(body, "")))
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
			return

		case intr := <-dh.Interrupt:
			log.Println("inside the server", intr)
			return
		case <-dig.VerseDiagnostic:

			syncMap := map[int]interface{}{}
			dig.MapMtx.Range(func(key, value any) bool {
				syncMap[key.(int)] = value
				return true
			})

			s, _ := json.Marshal(syncMap)
			body := []string{
				"id: " + uuid.NewString() + "\n",
				"event: data\n",
				"data: " + string(s) + "\n\n",
			}

			w.Write([]byte(strings.Join(body, "")))
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

		}
	}
}
