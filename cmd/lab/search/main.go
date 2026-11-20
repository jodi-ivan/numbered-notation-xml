package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/char/regexp"
	_ "github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/v2/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	bleve_unicode "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/document"
	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/stem"
	_ "github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/stem"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/category"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/storage"
	progressbar "github.com/schollz/progressbar/v3"
)

func dump(m mapping.IndexMapping, analyzer, text string) {
	an := m.AnalyzerNamed(analyzer)
	var terms []string
	for _, t := range an.Analyze([]byte(text)) {
		terms = append(terms, fmt.Sprintf("%q", t.Term))
	}
	fmt.Printf("%-20s %-10s -> %v\n", analyzer, text, terms)
}

func main() {
	log.Println("starting the indexing...")
	db, err := storage.NewStorage(context.Background(), "files/database/kidung-jemaat.db")
	if err != nil {
		log.Fatalf("Failed to connect to storage: %s", err.Error())
		return
	}

	repo := repository.New(context.Background(), db)

	indexPath := "files/index/kj.bleve"
	// Create a new index
	indexMapping := bleve.NewIndexMapping()

	var charFilters []string
	for _, r := range stem.ElisionRules {
		err := indexMapping.AddCustomCharFilter(r.Name, map[string]interface{}{
			"type":    regexp.Name,
			"regexp":  r.Pattern,
			"replace": r.Replace,
		})
		if err != nil {
			log.Fatal(err)

		}
		charFilters = append(charFilters, r.Name)
	}

	// 1. Create a custom token filter for n-grams (3 to 6 characters is usually sweet spot)
	err = indexMapping.AddCustomTokenFilter("my_ngram",
		map[string]interface{}{
			"type": ngram.Name,
			"min":  3,
			"max":  6,
		})
	if err != nil {
		log.Fatal(err)
	}
	// 2. Create an analyzer using that filter
	err = indexMapping.AddCustomAnalyzer("part_word_analyzer",
		map[string]interface{}{
			"type":      "custom",
			"tokenizer": "unicode",
			"token_filters": []string{
				lowercase.Name,
				"my_ngram",
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	if err := indexMapping.AddCustomAnalyzer("indonesian", map[string]interface{}{
		"type":          custom.Name,
		"char_filters":  charFilters,
		"tokenizer":     bleve_unicode.Name,
		"token_filters": []string{lowercase.Name, "id_spelling", "id_redup", "indonesian_stem"},
	}); err != nil {
		panic(err)
	}

	// index-time analyzer for autocomplete
	if err := indexMapping.AddCustomAnalyzer("autocomplete", map[string]interface{}{
		"type":          custom.Name,
		"char_filters":  charFilters, // b'rikan -> berikan first
		"tokenizer":     bleve_unicode.Name,
		"token_filters": []string{lowercase.Name, "id_spelling", "id_redup", "my_ngram"},
	}); err != nil {
		log.Fatal(err)
	}

	docMapping := bleve.NewDocumentMapping()

	contentMapping := bleve.NewTextFieldMapping()
	contentMapping.Store = true
	contentMapping.IncludeTermVectors = true
	contentMapping.Analyzer = "indonesian"

	docMapping.AddFieldMappingsAt("content", contentMapping)

	versesMapping := bleve.NewDocumentMapping()

	noMapping := bleve.NewNumericFieldMapping()
	noMapping.Store = true

	verseContentMapping := bleve.NewTextFieldMapping()
	verseContentMapping.Store = true
	verseContentMapping.Analyzer = "indonesian"
	verseContentMapping.IncludeTermVectors = true

	startMapping := bleve.NewNumericFieldMapping()
	startMapping.Store = true

	endMapping := bleve.NewNumericFieldMapping()
	endMapping.Store = true

	variantMapping := bleve.NewTextFieldMapping()
	variantMapping.Store = true

	hymnMap := bleve.NewNumericFieldMapping()
	hymnMap.Store = true

	titleMap := bleve.NewTextFieldMapping()
	titleMap.Store = true
	titleMap.Analyzer = "indonesian"

	// ngram for autocomplete
	titleNgram := bleve.NewTextFieldMapping()
	titleNgram.Store = true
	titleNgram.IncludeTermVectors = true
	titleNgram.Analyzer = "autocomplete"

	ogMap := bleve.NewTextFieldMapping()
	ogMap.Store = true

	versesMapping.AddFieldMappingsAt("no", noMapping)
	versesMapping.AddFieldMappingsAt("content", verseContentMapping)
	versesMapping.AddFieldMappingsAt("start", startMapping)
	versesMapping.AddFieldMappingsAt("end", endMapping)
	docMapping.AddSubDocumentMapping("verses", versesMapping)
	docMapping.AddFieldMappingsAt("variant", variantMapping)
	docMapping.AddFieldMappingsAt("hymn_no", hymnMap)
	docMapping.AddFieldMappingsAt("title", titleMap)
	docMapping.AddFieldMappingsAt("original_title", ogMap)
	docMapping.AddFieldMappingsAt("title_ngram", titleNgram)

	indexMapping.DefaultMapping = docMapping

	// Add documents
	documents := []*document.Document{}
	log.Println("Building documents....")
	bar := progressbar.Default(478)
	for no := 1; no <= 478; no++ {
		var doc *document.Document
		var err error
		if variant, ok := category.Variants[no]; ok {
			for _, vs := range variant {
				doc, err = document.BuildDocument(repo, no, vs)
				if err != nil {
					log.Fatalf("Failed to build docs: %s", err.Error())
					return
				}
				documents = append(documents, doc)

			}
		} else {
			doc, err = document.BuildDocument(repo, no)
			if err != nil {
				log.Fatalf("Failed to build docs: %s", err.Error())
				return
			}
			documents = append(documents, doc)

		}
		bar.Add(1)
	}

	log.Println("Batch indexing starting....")
	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	indexBar := progressbar.Default(int64(len(documents)))

	// Iterate and index the documents
	batch := index.NewBatch()
	for _, doc := range documents {
		// raw, _ := json.Marshal(doc.Verses)
		// log.Printf("%d. %s \n %s \n %s\n\n", doc.HymnNo, doc.Title, doc.Content, string(raw))
		err := batch.Index(doc.ID, doc)
		if err != nil {
			log.Fatal(err)
		}
		indexBar.Add(1)
	}
	if err := index.Batch(batch); err != nil {
		log.Fatal(err)
	}
	// size := map[string]interface{}
	statsraw := index.StatsMap()
	stats, _ := statsraw["index"].(map[string]interface{})
	log.Printf("success to build the index. Index size %.2f mb", float64((stats["CurOnDiskBytes"].(uint64)))/1e6)
}
