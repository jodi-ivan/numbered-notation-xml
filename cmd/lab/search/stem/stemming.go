package stem

import (
	"log"

	"github.com/RadhiFadlillah/go-sastrawi"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

const Name = "indonesian_stem"

type filter struct {
	stemmer   sastrawi.Stemmer
	stopwords sastrawi.Dictionary
}

func (f *filter) Filter(tokens analysis.TokenStream) analysis.TokenStream {
	result := tokens[:0]

	for _, token := range tokens {
		word := string(token.Term)

		if f.stopwords.Contains(word) {
			continue
		}

		token.Term = []byte(f.stemmer.Stem(word))
		result = append(result, token)
	}

	return result
}

func newFilter(_ map[string]interface{}, _ *registry.Cache) (analysis.TokenFilter, error) {
	return &filter{
		stemmer:   sastrawi.NewStemmer(sastrawi.DefaultDictionary()),
		stopwords: sastrawi.DefaultStopword(),
	}, nil
}

func init() {
	log.Println("CALLED")
	if err := registry.RegisterTokenFilter(Name, newFilter); err != nil {
		panic(err)
	}
}
