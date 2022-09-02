package adapter

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/julienschmidt/httprouter"
)

type ReadFile struct {
}

func (rf *ReadFile) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filelocation := "files/scores/kj-001.musicxml"

	xmlFile, err := os.Open(filelocation)
	if err != nil {
		log.Println("failed to read file: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer xmlFile.Close()

	content, _ := ioutil.ReadAll(xmlFile)

	var music musicxml.MusicXML

	err = xml.Unmarshal(content, &music)
	if err != nil {
		log.Println("failed to parse xml: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	toDebug, _ := json.MarshalIndent(music, "", "    ")
	log.Println(string(toDebug))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(toDebug)

}
