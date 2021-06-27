package webserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/jodi-ivan/numbered-notation-xml/utils/errors"
	"github.com/julienschmidt/httprouter"
)

type ResponseWriterWithStatus struct {
	w       http.ResponseWriter
	status  int
	err     *ErrorResponse
	message string
}

func (rwws *ResponseWriterWithStatus) Header() http.Header {
	return rwws.w.Header()
}

func (rwws *ResponseWriterWithStatus) Write(data []byte) (int, error) {
	return rwws.w.Write(data)
}

func (rwws *ResponseWriterWithStatus) WriteHeader(statusCode int) {
	rwws.status = statusCode
	rwws.w.WriteHeader(statusCode)
}

//HTTPAdapter
type HTTPAdapter interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

func commonMiddleware(wg *sync.WaitGroup, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("Panic: ", err)
				debug.PrintStack()
				RenderErrorResponse(w, http.StatusInternalServerError, errors.NewFromError(fmt.Errorf("panic: %v", err), "Something went wrong, please try again"))
			}

		}()
		responseWriter := &ResponseWriterWithStatus{
			w: w,
		}
		t := time.Now()
		wg.Add(1)
		defer wg.Done()

		next(responseWriter, r, ps)

		if responseWriter.err != nil {
			log.Printf("[Webserver][%d ms] %s %s -> %d. %s: %s.\n", time.Since(t).Milliseconds(), strings.ToUpper(r.Method), r.URL.Path, responseWriter.status, responseWriter.err.Source.Pointer, responseWriter.err.Detail)
		} else {
			log.Printf("[Webserver][%d ms] %s %s -> %d. %s.\n", time.Since(t).Milliseconds(), strings.ToUpper(r.Method), r.URL.Path, responseWriter.status, responseWriter.message)
		}
	}
}

// WebServer the server object
type WebServer struct {
	httpServer *http.Server
	listener   net.Listener
	wg         *sync.WaitGroup
	httpRouter *httprouter.Router
}

func (ws *WebServer) Serve(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("[Webserver] failed to listen to port %s. Err: %v\n", port, err.Error())
		return err

	}
	log.Printf("[Webserver] http server running on address: %v\n", port)
	ws.httpServer.Handler = ws.httpRouter
	ws.listener = l
	go ws.httpServer.Serve(ws.listener)
	return nil
}

func (ws *WebServer) Register(method, path string, adapter HTTPAdapter) {
	ws.httpRouter.Handle(method, path, commonMiddleware(ws.wg, adapter.ServeHTTP))
}

func (ws *WebServer) Stop() {
	log.Println("[Webserver] Closing the webserver...")

	ws.wg.Wait()
	err := ws.listener.Close()
	if err != nil {
		log.Printf("[Webserver] failed to close the listener's web server. Err : %s\n", err.Error())
	}

	err = ws.httpServer.Close()
	if err != nil {
		log.Printf("[Webserver] failed to close the web server. Err : %s\n", err.Error())
	}

	log.Println("[Webserver] Webserver closed.")
}

func InitWebserver() (*WebServer, error) {

	return &WebServer{
		httpServer: &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		wg:         &sync.WaitGroup{},
		httpRouter: httprouter.New(),
	}, nil
}
