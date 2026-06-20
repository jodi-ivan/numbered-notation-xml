package main

// #include <stdlib.h>
import "C"
import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"sync"
	"unsafe"

	"github.com/jodi-ivan/numbered-notation-xml/adapter"
	"github.com/jodi-ivan/numbered-notation-xml/internal/renderer"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
	"github.com/jodi-ivan/numbered-notation-xml/utils/storage"
)

var (
	stringAdapter *adapter.RenderString
	engineOnce    sync.Once
	stateMutex    sync.Mutex
)

func GetEngine() *adapter.RenderString {
	engineOnce.Do(func() {
		cfg := config.Config{
			MusicXML: config.MusicXMLConfig{
				Path:       "/home/jodiivan/go/src/github.com/jodi-ivan/numbered-notation-xml/files/scores/musicxml/",
				FilePrefix: "kj",
			},
			SQLite: config.SQLiteConfig{
				DBPath: "/home/jodiivan/go/src/github.com/jodi-ivan/numbered-notation-xml/files/database/kidung-jemaat.db",
			},
		}

		db, err := storage.NewStorage(context.Background(), cfg.SQLite.DBPath)
		if err != nil {
			log.Printf("Failed to connect to storage: %s\n", err.Error())
			return
		}

		repo := repository.New(context.Background(), db)
		usecaseMod := usecase.New(cfg, repo, renderer.NewRenderer())
		stringAdapter = adapter.NewRenderString(usecaseMod)

	})
	return stringAdapter
}

// Define a struct matching your configuration parameters
type RenderConfig struct {
	Verse     int  `json:"verse"`
	FocusMode bool `json:"focus_mode"`
}

//export RenderHymnSVG
func RenderHymnSVG(number C.int, variant *C.char, configJson *C.char) *C.char {

	stateMutex.Lock()
	defer stateMutex.Unlock()

	ctx := context.Background()
	e := GetEngine()
	// 1. Convert C types to Go types
	goNumber := int(number)

	goVariant := []string{}
	if variant != nil {
		goVariant = append(goVariant, C.GoString(variant))
	}

	// 2. Parse the optional configuration JSON
	// Set default values first
	config := RenderConfig{
		Verse:     0, // Default verse
		FocusMode: false,
	}

	if configJson != nil {
		jsonStr := C.GoString(configJson)
		if jsonStr != "" {
			_ = json.Unmarshal([]byte(jsonStr), &config)
			param := params.Param{
				SingleVerseMode: config.FocusMode,
				Verse:           config.Verse,
			}

			ctx = params.NewParamContext(ctx, &param)
		}
	}

	// 3. Call your internal layout engine (Mockup logic here)
	// svgOutput := internalRenderEngine(goNumber, goVariant, config)
	buff := bytes.NewBuffer(nil)
	content, err := e.RenderHymn(context.Background(), buff, goNumber, goVariant...)
	if err != nil {
		log.Printf("Problem creating file: %v\n", err)
		return C.CString(err.Error())
	}
	// 4. Return string back to C++
	// CRITICAL: C.CString allocates memory on the C heap.
	// C++ must free this memory to prevent memory leaks!
	return C.CString(content)
}

// export FreeRenderedString
func FreeRenderedString(ptr *C.char) {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
	}
}

// build it by
// go build -buildmode=c-shared -o libhymn_renderer.so cmd/dll/main.go
func main() {
	// This main function is required for c-shared build modes, but won't be run.
}
