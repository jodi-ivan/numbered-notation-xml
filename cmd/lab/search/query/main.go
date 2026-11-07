package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/v2/analysis/lang/id"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/token/ngram"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	_ "github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/stem"
)

type VerseDoc struct {
	No      int    `json:"no"`
	Content string `json:"content"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
}

func toFloat64Slice(v interface{}) []float64 {
	switch val := v.(type) {
	case []interface{}:
		result := make([]float64, len(val))
		for i, item := range val {
			result[i] = item.(float64)
		}
		return result
	case float64:
		return []float64{val} // wrap scalar in slice
	}
	return nil
}

func main() {

	indexMapping := bleve.NewIndexMapping()
	var err error
	if err := indexMapping.AddCustomAnalyzer("indonesian", map[string]interface{}{
		"type":          custom.Name,
		"tokenizer":     unicode.Name,
		"token_filters": []string{lowercase.Name, "indonesian_stem"},
	}); err != nil {
		panic(err)
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

	// db, err := storage.NewStorage(context.Background(), "files/database/kidung-jemaat.db")
	// if err != nil {
	// 	log.Fatalf("Failed to connect to storage: %s", err.Error())
	// 	return
	// }
	// repo := repository.New(context.Background(), db)

	indexPath := "files/index/kj.bleve"

	// Load the existing index
	index, err := bleve.Open(indexPath)
	if err != nil {
		panic(err)
		// Handle error (e.g., path incorrect, index corrupted)
	}
	defer index.Close()

	/*
		fmt.Println("1. Check which analyzer each field path resolves to")
		m := index.Mapping()
		for _, p := range []string{"content", "verses.content", "title"} {
			fmt.Printf("%s -> %q\n", p, m.AnalyzerNameForPath(p))
		}
		fmt.Println("-------------------------")

		fmt.Println("2. Check what the analyzer does to the query word")

		name := m.AnalyzerNameForPath("verses.content")
		a := m.AnalyzerNamed(name)
		for _, w := range []string{"menahan", "tahan"} {
			for _, t := range a.Analyze([]byte(w)) {
				fmt.Printf("%s -> %s\n", w, string(t.Term))
			}
		}
	*/

	q := bleve.NewMatchQuery("nyanyi")
	q.SetBoost(3)
	q.SetField("content")

	q2 := bleve.NewMatchQuery("nyanyi")
	q2.SetBoost(1)
	q2.SetField("verses.content")

	q3 := bleve.NewMatchQuery("nyanyi")
	q3.SetBoost(5)
	q3.SetField("title")

	combined1 := bleve.NewDisjunctionQuery(q, q2, q3)

	searchRequest := bleve.NewSearchRequest(combined1)
	searchRequest.IncludeLocations = true
	searchRequest.Fields = []string{"hymn_no", "content", "title", "verses.no", "verses.content", "verses.start", "verses.end", "variant", "original_title"}
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Explain = true
	searchResult, err := index.Search(searchRequest)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("taken:", searchResult.Took)

	// result := map[int][]int{}
	for _, hit := range searchResult.Hits {

		matchedVerseNos := map[int][]string{}

		if locations, ok := hit.Locations["content"]; ok {
			raw, _ := json.Marshal(locations)
			log.Println(string(raw))

			rawStarts := toFloat64Slice(hit.Fields["verses.start"])
			rawEnds := toFloat64Slice(hit.Fields["verses.end"])
			rawNos := toFloat64Slice(hit.Fields["verses.no"])
			log.Println("Score", hit.Fields["hymn_no"].(float64), ": ", hit.Score)
			for term, termLocations := range locations {
				for _, loc := range termLocations {
					// find which verse this offset falls in
					for i := range rawStarts {
						start := rawStarts[i]
						end := rawEnds[i]
						if loc.Start >= uint64(start) && loc.Start < uint64(end) {
							if matchedVerseNos[int(rawNos[i])] == nil {
								matchedVerseNos[int(rawNos[i])] = []string{}
							}
							matchedVerseNos[int(rawNos[i])] = append(matchedVerseNos[int(rawNos[i])], term)
							break
						}
					}
				}
			}
			var verseNos []int
			for no := range matchedVerseNos {
				verseNos = append(verseNos, no)
			}
			sort.Ints(verseNos)

			for _, no := range verseNos {
				fmt.Printf("Found: %s. hymn %.0f — matched verses: %v\n", matchedVerseNos[no], hit.Fields["hymn_no"].(float64), no)
			}
			fmt.Println("")
		}

		// for _, hit := range searchResult.Hits {
		// 	for field, terms := range hit.Locations {
		// 		for term := range terms {
		// 			fmt.Printf("hymn %.0f matched field=%s term=%q\n", hit.Fields["hymn_no"].(float64), field, term)
		// 		}
		// 	}
		// }

		if locations, ok := hit.Locations["title"]; ok {
			for term, termLocations := range locations {
				for _, loc := range termLocations {
					// find which verse this offset falls in
					fmt.Printf("[%.0f] Title term %q matched at byte offset %d: %s\n", hit.Fields["hymn_no"].(float64), term, loc.Start, hit.Fields["title"])

				}
			}
		}

	}

	// q1 := bleve.NewWildcardQuery("*dalam*")
	// q1.SetField("title")
	// q1.SetBoost(3)

	// q2 := bleve.NewWildcardQuery("*dalam*")
	// q2.SetField("original_title")

	// combined := bleve.NewDisjunctionQuery(q1, q2)

	// searchRequest1 := bleve.NewSearchRequest(combined)
	// searchRequest1.IncludeLocations = true
	// searchRequest1.Fields = []string{"hymn_no", "content", "title", "verses.no", "verses.content", "verses.start", "verses.end", "variant"}
	// searchRequest1.Highlight = bleve.NewHighlight()
	// searchRequest1.Explain = true
	// searchResult1, err := index.Search(searchRequest1)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("taken:", searchResult1.Took)

	// // result := map[int][]int{}
	// for _, hit := range searchResult1.Hits {

	// 	log.Println("===================", hit.Fields["hymn_no"], hit.Fields["variant"])

	// 	log.Println("Score", hit.Fields["hymn_no"].(float64), ": ", hit.Score)

	// 	if locations, ok := hit.Locations["original_title"]; ok {
	// 		for term, termLocations := range locations {
	// 			for _, loc := range termLocations {
	// 				// find which verse this offset falls in
	// 				fmt.Printf("Original title term %q matched at byte offset %d: %s\n", term, loc.Start, hit.Fields["original_title"])
	// 			}
	// 		}
	// 	}

	// 	if locations, ok := hit.Locations["title"]; ok {
	// 		for term, termLocations := range locations {
	// 			for _, loc := range termLocations {
	// 				// find which verse this offset falls in
	// 				fmt.Printf("Title term %q matched at byte offset %d: %s\n", term, loc.Start, hit.Fields["title"])
	// 			}
	// 		}
	// 	}

	// 	fmt.Println("")

	// }

}
