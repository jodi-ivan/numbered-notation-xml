package stream

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hcl/audioduration"
	"github.com/julienschmidt/httprouter"
)

type AudioStream struct {
	Sig chan os.Signal
}

func (dh *AudioStream) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	paths := []string{
		`/home/jodiivan/go/src/github.com/jodi-ivan/numbered-notation-xml/files/audio/kj-001-opening.mp3`,
		`/home/jodiivan/go/src/github.com/jodi-ivan/numbered-notation-xml/files/audio/kj-001-core.mp3`,
	}

	openingFile, err := os.Open(paths[0])
	if err != nil {
		log.Println("Failed to open openingFile", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	defer openingFile.Close()

	openingDuration, err := audioduration.Mp3(openingFile)
	if err != nil {
		log.Println("Failed to calculate mp3 opening duration", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	coreFile, err := os.Open(paths[1])
	if err != nil {
		log.Println("Failed to open coreFile", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	defer coreFile.Close()

	repeatRaw := r.FormValue("repeat")

	repeat, err := strconv.Atoi(repeatRaw)
	if repeatRaw != "" && err != nil {
		log.Printf("[ServeHTTP] invalid verse: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return
	}

	if repeat == 1 {
		repeat = 0
	}

	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "audio/mpeg")

	ctx := r.Context()
	currRepeat := 1
	measure := 0
	measureMtx := sync.Mutex{}

	songMeasure := 23.0
	songTempo := 85.0
	songDuration := ((songMeasure * 4) / songTempo) * 60.0
	op := time.Duration(openingDuration * float64(time.Second))
	core := time.Duration(songDuration * float64(time.Second)) // for 1 timesignature no changes
	log.Println("wait for ", op.String())

	files := []*os.File{
		openingFile,
		coreFile,
	}
	go func() {

		<-time.After(op)
		measure = 1
		log.Println("Start: Measure ", measure)
		ticker := time.NewTicker(core / 23)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				measureMtx.Lock()
				measure++
				measureMtx.Unlock()
				log.Println("Measure ", measure)
				if measure == 23 {
					return
				}
			case <-dh.Sig:
				log.Println("os signal interrupt")
				return
			}
		}

	}()

	for i, file := range files {

		buf := make([]byte, 4096) // 4KB chunks

		// Loop core indefinitely
		for {
			select {
			case <-ctx.Done():
				log.Println("context timeout")

				return
			case <-dh.Sig:
				log.Println("os signal interrupt")
				return
			default:
			}

			n, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("Music ", file.Name(), " is finished")
					// if i == 1 && currRepeat == repeat {
					// 	break
					// }
					measureMtx.Lock()
					measure = 0
					measureMtx.Unlock()

					if i == 0 || ((i == 1 && repeat == 0) || (i == 1 && currRepeat == repeat)) {
						break
					} else if currRepeat < repeat {
						currRepeat++
						file.Seek(0, 0)
						continue
					}
				}
				log.Println("Failed to ", err.Error())
				return
			}

			_, werr := w.Write(buf[:n])
			if werr != nil {
				log.Println("Failed to ", werr.Error())
				return
			}

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

		}
	}

	// wg.Wait()

	log.Println("It is done")

}
