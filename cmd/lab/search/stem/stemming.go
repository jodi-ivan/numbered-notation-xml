package stem

import (
	"github.com/RadhiFadlillah/go-sastrawi"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

const Name = "indonesian_stem"

type filter struct {
	stemmer sastrawi.Stemmer
}

func (f *filter) Filter(tokens analysis.TokenStream) analysis.TokenStream {
	result := tokens[:0]

	for _, token := range tokens {
		token.Term = []byte(f.stemmer.Stem(string(token.Term)))
		result = append(result, token)

	}

	return result
}

func newFilter(_ map[string]interface{}, _ *registry.Cache) (analysis.TokenFilter, error) {
	return &filter{
		stemmer: sastrawi.NewStemmer(sastrawi.DefaultDictionary()),
	}, nil
}

func init() {
	if err := registry.RegisterTokenFilter(Name, newFilter); err != nil {
		panic(err)
	}
}
