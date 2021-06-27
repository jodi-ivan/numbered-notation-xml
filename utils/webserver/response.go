package webserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jodi-ivan/numbered-notation-xml/utils/errors"
)

type ErrorSource struct {
	Pointer string `json:"pointer"`
}
type ErrorResponse struct {
	Code   int         `json:"code,string"`
	Source ErrorSource `json:"source"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
	Err    error       `json:"-"`
}

func (er *ErrorResponse) Error() string {
	return er.Err.Error()
}

type Pagination struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type InsertSuccessMessage struct {
	Message string `json:"message"`
	ID      int64  `json:"id,string"`
}

type SuccessResponse struct {
	Links *Pagination `json:"links,omitempty"`
	Data  interface{} `json:"data"`
}
type DataWrapper struct {
	ID         int64       `json:"id,string"`
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes"`
}

func NewErrorSource(pointer string) ErrorSource {
	return ErrorSource{
		Pointer: pointer,
	}
}

func RenderErrorResponse(w http.ResponseWriter, code int, err *errors.Error) {
	resp := &ErrorResponse{
		Code: code,
		Source: ErrorSource{
			Pointer: err.GetSource(),
		},
		Title:  err.GetTitle(),
		Detail: err.Error(),
	}
	raw, errMarshal := json.Marshal(resp)
	if errMarshal != nil {
		log.Printf("[Webserver][RenderError] failed to marshal the error object, err: %s", err.Error())
		return
	}
	rw, ok := w.(*ResponseWriterWithStatus)
	if ok {
		rw.err = resp
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(code)
		rw.Write(raw)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(raw)
	}
}

func RenderSuccessInsertResponse(w http.ResponseWriter, newID int64, message string) {
	RenderSuccessResponse(w, nil, InsertSuccessMessage{
		Message: message,
		ID:      newID,
	})
}

func RenderSuccessResponse(w http.ResponseWriter, pagination *Pagination, data interface{}, logging ...string) {
	resp := SuccessResponse{
		Links: pagination,
		Data:  data,
	}

	raw, errMarshal := json.Marshal(resp)
	if errMarshal != nil {
		log.Printf("[Webserver][RenderSuccess] failed to marshal the error object, err: %s", errMarshal.Error())
		return
	}

	rw, ok := w.(*ResponseWriterWithStatus)
	if ok {
		if len(logging) > 0 {
			rw.message = logging[0]
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(raw)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(raw)
	}
}
