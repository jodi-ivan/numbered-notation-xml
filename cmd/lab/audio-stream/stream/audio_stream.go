package stream

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type AudioStream struct {
	Sig chan os.Signal
}

func (dh *AudioStream) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	fileLocation := `/home/jodiivan/go/src/github.com/jodi-ivan/numbered-notation-xml/files/audio/kj-001.mp3`

	file, err := os.Open(fileLocation)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	defer file.Close()

	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "audio/mpeg")

	// Create a buffer for streaming data chunks
	buf := make([]byte, 4096) // 4KB chunks

	for {

		select {
		case <-dh.Sig:
			log.Println("os signal interrupt")
			return

		default:

		}

		// Read a chunk from the MP3 file
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				// Loop the audio by seeking back to the start
				return
			}
			// Stop streaming if the connection drops or error occurs
			break
		}

		// Write the chunk to the HTTP response
		w.Write(buf[:n])

		// Flush the buffer to send data to the client immediately
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

	}

}
