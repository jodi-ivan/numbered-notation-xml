package main

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve/v2"
)

func main() {

	indexPath := "files/index/kj.bleve"

	// Load the existing index
	index, err := bleve.Open(indexPath)
	if err != nil {
		// Handle error (e.g., path incorrect, index corrupted)
	}
	defer index.Close()

	// Search the created index
	query := bleve.NewMatchQuery("sembah")
	query.Fuzziness = 1
	query.Prefix = 1
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"hymn_no", "title", "verses.no", "verses.content", "verses.start", "verses.end", "variant"}
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Explain = true
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	result := map[int][]int{}
	for _, hit := range searchResult.Hits {

		// starts := hit.Fields["verses.start"].([]interface{})
		// ends := hit.Fields["verses.end"].([]interface{})

		// for i, s := range starts {
		// 	start := int(s.(float64))
		// 	end := int(ends[i].(float64))
		// }

		// hit.Fields contains your stored values
		// hit.Locations contains the match positions per field
		log.Println("===================", hit.Fields["hymn_no"], hit.Fields["variant"], "|", hit.Fields["verses.start"], hit.Fields["verses.end"], hit.Fields["verses.no"])
		// get match positions in "content"
		if locations, ok := hit.Locations["content"]; ok {
			for term, termLocations := range locations {
				for _, loc := range termLocations {

					fmt.Printf("term %q matched at byte offset %d\n", term, loc.Start)
					hymnNo := hit.Fields["hymn_no"].(float64)
					if result[int(hymnNo)] == nil {
						result[int(hymnNo)] = []int{}
					}

					result[int(hymnNo)] = append(result[int(hymnNo)], int(loc.Start))

					log.Println("document ID", hit.ID)

				}
			}
		}

		log.Println(result)
	}

}
