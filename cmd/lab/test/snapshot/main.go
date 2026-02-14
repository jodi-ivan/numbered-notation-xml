package main

import (
	"bytes"
	"context"
	"errors"
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

	"github.com/JoshVarga/svgparser"
)

func main() {
	env := os.Environ()

	// Define a string flag
	method := flag.String("snapshot", "gen", "Operation, 'gen' generate golden. 'assert' assert snapshot")
	goldenPath := flag.String("path", "", "Path of the golden files to be generated, only when snapshot=gen")

	// Parse the command-line arguments into the defined flags
	flag.Parse()

	// TODO: [snapshot test] validate here
	xmlFlag := ""
	dbPath := ""

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

		fmt.Printf("Read the musicxml files on %s.\n", xmlFlag)
		fmt.Printf("Generate golden files. Generated on %s. \n\n", *goldenPath)

		err = GenerateGolden(context.Background(), stringRender, *goldenPath)
		if err != nil {
			fmt.Printf("Failed to Generate Snapshot: %s\n", err.Error())
			os.Exit(2)
			return
		}
	case "assert":

		if *goldenPath == "" {
			fmt.Printf("Goldenpath -path cannot be empty")
			os.Exit(2)
			return
		}

		fmt.Printf("Gathering golden files on %s. \n\n", *goldenPath)
		numFiles := 22
		for i := 1; i <= numFiles; i++ {
			fmt.Printf("Asserting kj-%03d....\n", i)
			err = Assert(context.Background(), stringRender, *goldenPath, i)
			if err != nil {
				fmt.Printf("Failed to Assert the golden snapsot: %s\n", err.Error())
				os.Exit(2)
				return
			}
			fmt.Printf("Asserting kj-%03d SUCCESS\n", i)

		}

	}

	os.Exit(0)

}

func getGenElement(stringRenderer *adapter.RenderString, number int) (*svgparser.Element, error) {
	buff := bytes.NewBuffer(nil)
	content, err := stringRenderer.RenderHymn(context.Background(), buff, number)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(content)

	return svgparser.Parse(reader, false)
}

func getGoldenElement(path string, number int) (*svgparser.Element, error) {
	fileName := fmt.Sprintf("%s/goldenfiles/kj-%03d.svg", path, number)
	xmlFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	return svgparser.Parse(xmlFile, false)
}

func breakdownStaffCredit(elmnt *svgparser.Element) (staff []*svgparser.Element, verses *svgparser.Element, credit *svgparser.Element) {
	for i := 4; i < len(elmnt.Children); i++ {
		child := elmnt.Children[i]
		if sty, ok := child.Attributes["style"]; ok && sty == "staff" {
			staff = append(staff, child)
			continue
		}

		if class, ok := child.Attributes["class"]; ok {
			switch class {
			case "verses":
				verses = child
			case "credit":
				credit = child
			}
		}

	}
	return
}

func Assert(ctx context.Context, stringRenderer *adapter.RenderString, path string, number int) error {

	generatedElement, err := getGenElement(stringRenderer, number)
	if err != nil {
		return err
	}

	goldenElement, err := getGoldenElement(path, number)
	if err != nil {
		return err
	}

	cdata := generatedElement.Children[0].Children[0].Content
	expectedFont := []string{
		"Caladea",
		"Old Standard TT",
		"Noto Music",
		"Figtree",
	}

	// FONT VALIDATION
	for _, ef := range expectedFont {
		if !strings.Contains(strings.ToLower(cdata), strings.ToLower(ef)) {
			return fmt.Errorf("font family %s is not found", ef)
		}
	}

	// TITLE
	if !goldenElement.Children[1].Compare(generatedElement.Children[1]) {
		return errors.New("the title element does not match")
	}

	// KEY SIGNATURE
	if !goldenElement.Children[2].Compare(generatedElement.Children[2]) {
		return errors.New("the key signature element does not match")
	}

	// TIME SIGNATURE
	if !goldenElement.Children[3].Compare(generatedElement.Children[3]) {
		return errors.New("the time signature element does not match")
	}

	gs, gv, gc := breakdownStaffCredit(generatedElement)
	gns, gnv, gnc := breakdownStaffCredit(goldenElement)

	// STAFF
	if len(gs) != len(gns) {
		return fmt.Errorf("the staff element count does not match")
	}

	for i, g := range gs {
		if !g.Compare(gns[i]) {
			return fmt.Errorf("the staff element at index %d does not match", i)
		}
	}

	// VERTICES
	if !gv.Compare(gnv) {
		return fmt.Errorf("the vertices element does not match")

	}

	// CREDIT
	if !gc.Compare(gnc) {
		return fmt.Errorf("the credit element does not match")

	}

	return nil

}

var variant = map[int][]string{
	24: []string{"a", "b"},
	30: []string{"a", "b"},
	31: []string{"a", "b"},
	37: []string{"a", "b"},
}

func GenerateGolden(ctx context.Context, stringRenderer *adapter.RenderString, path string) error {
	numFiles := 41

	renderAndSave := func(i int, variants ...string) error {
		buff := bytes.NewBuffer(nil)
		content, err := stringRenderer.RenderHymn(context.Background(), buff, i, variants...)
		if err != nil {
			log.Printf("Problem creating file: %v\n", err)
			return err
		}

		fileName := fmt.Sprintf("%s/goldenfiles/kj-%03d.svg", path, i)
		if len(variants) > 0 {
			fileName = fmt.Sprintf("%s/goldenfiles/kj-%03d%s.svg", path, i, variants[0])
		}
		fmt.Println("Creating golden for", fileName)
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("Problem creating file: %v\n", err)
			return err
		}
		fmt.Fprint(file, content)
		file.Close()

		return nil
	}

	for i := 1; i <= numFiles; i++ {
		vs, ok := variant[i]
		if ok {
			for _, v := range vs {
				err := renderAndSave(i, v)
				if err != nil {
					return err
				}
			}
		} else {
			err := renderAndSave(i)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
