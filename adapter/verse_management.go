package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/julienschmidt/httprouter"
)

type VerseManagement struct {
	VerseRepo repository.Repository
}

type Input struct {
	Style   int                      `json:"style"`
	Col     int                      `json:"col"`
	Row     int                      `json:"row"`
	Content [][]lyric.LyricWordVerse `json:"content"`
}

func (vm *VerseManagement) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	hymn, err := strconv.Atoi(ps.ByName("hymn"))
	if err != nil {
		log.Printf("invalid hymn number: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL hymn"))
		return
	}

	verse, err := strconv.Atoi(ps.ByName("verse"))
	if err != nil {
		log.Printf("invalid verse number: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL verse"))
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	input := &Input{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error input: :%s", err.Error())))
		return
	}

	stringify, err := json.Marshal(input.Content)
	if err != nil {
		log.Printf("cannot stringfy content: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot stringfy content"))
		return
	}
	id, err := vm.VerseRepo.InsertVerse(ctx, hymn, verse, input.Style, input.Col, input.Row, string(stringify))
	if err != nil {
		log.Printf("Failed to insert the verse: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert verse"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("sucess insert with ID: %d", id)))

}
