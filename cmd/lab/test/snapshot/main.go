package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/adapter"
	"github.com/jodi-ivan/numbered-notation-xml/internal/renderer"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
	"github.com/jodi-ivan/numbered-notation-xml/utils/storage"
)

// TODO: the snapshot testing
/*
	- generate golden file
	- generate the svg
	- sampling
	- assert
*/

var defaultxml = "~/go/src/github.com/jodi-ivan/numbered-notation-xml/files/scores/musicxml/"
var defaultdb = "~/go/src/github.com/jodi-ivan/numbered-notation-xml/files/database/kidung-jemaat.db"

func main() {
	env := os.Environ()

	// Define a string flag
	method := flag.String("snapshot", "gen", "Operation, 'gen' generate golden. 'assert' assert snapshot")
	goldenPath := flag.String("path", "", "Path of the golden files to be generated, only when snapshot=gen")

	// Parse the command-line arguments into the defined flags
	flag.Parse()

	log.Println("method:", *method)
	xmlFlag := defaultxml
	dbPath := defaultdb

	for _, e := range env {
		keyval := strings.Split(e, "=")
		switch keyval[0] {
		case "KJ_TEST_MUSICXML_PATH":
			xmlFlag = keyval[1]
		case "KJ_TEST_DB_PATH":
			dbPath = keyval[1]
		}
	}

	// Parse the flags from the command line

	cfg := config.Config{
		MusicXML: config.MusicXMLConfig{
			Path:       xmlFlag,
			FilePrefix: "kj",
		},
		SQLite: config.SQLiteConfig{
			DBPath: dbPath,
		},
	}

	db, err := storage.NewStorage(context.Background(), cfg.SQLite.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to storage: %s", err.Error())
		return
	}

	defer db.Close()

	repo := repository.New(context.Background(), db)
	usecaseMod := usecase.New(cfg, repo, renderer.NewRenderer())
	stringRender := adapter.NewRenderString(usecaseMod)

	switch *method {
	case "gen":
		if *goldenPath == "" {
			fmt.Printf("Goldenpath -path cannot be empty")
			os.Exit(2)
			return
		}

		fmt.Printf("Read the musicxml files on %s.\n", defaultxml)
		fmt.Printf("Generate golden files. Generated on %s. \n\n", *goldenPath)

		err = GenerateGolden(context.Background(), stringRender, *goldenPath)
		if err != nil {
			fmt.Printf("Failed to Generate Snapshot: %s\n", err.Error())
			os.Exit(2)
			return
		}

	}

	os.Exit(0)

}

func GenerateGolden(ctx context.Context, stringRenderer *adapter.RenderString, path string) error {
	numFiles := 22
	for i := 1; i <= numFiles; i++ {
		buff := bytes.NewBuffer(nil)
		content, err := stringRenderer.RenderHymn(context.Background(), buff, i)
		if err != nil {
			log.Fatalf("Problem creating file: %v", err)
			return err
		}

		fileName := fmt.Sprintf("%s/kj-%03d.svg", path, i)
		fmt.Println("Creating golden for", fileName)
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("Problem creating file: %v", err)
			return err
		}
		fmt.Fprintf(file, "%s", content)
		file.Close()
	}
	return nil
}
