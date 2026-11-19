package stem

import (
	"regexp"
	"strings"

	"github.com/RadhiFadlillah/go-sastrawi"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

// memujimu
// puji vs muji
// anak-anakmu is registered as anakanakmu
// porakporanda
// b'rilah
// t'rus 268 & 26 failed.
// ombang-ambing
const Name = "indonesian_stem"

var (
	// Built once and shared. Read-only after init, so safe for concurrent use.
	// Do not call dict.Add/Remove after this point.
	dict    = sastrawi.DefaultDictionary()
	stemmer = sastrawi.NewStemmer(dict)

	// go-sastrawi cannot strip "ter-" when the root starts with consonant + "e"
	// (tersembunyi, tersenyum, tertekan, terdengar, ...): its ter-CP regex
	// uses [^e] where it should only reject "er".
	rxBrokenTer = regexp.MustCompile(`^ter[^aiueor]e`)
)

func stem(word string) string {
	word = strings.ToLower(word)

	// A dictionary root is final. This also protects words like "terang".
	if dict.Contains(word) {
		return word
	}

	if s := stemmer.Stem(word); s != word {
		return s
	}

	// The library gave up. For the known-broken shape, strip "ter" ourselves
	// and let the library handle whatever is left (suffixes like -lah, -nya).
	if rxBrokenTer.MatchString(word) {
		if s := stemmer.Stem(word[3:]); dict.Contains(s) {
			return s
		}
	}

	return word
}

type filter struct {
	stemmer sastrawi.Stemmer
}

func (f *filter) Filter(tokens analysis.TokenStream) analysis.TokenStream {
	result := tokens[:0]

	for _, token := range tokens {
		token.Term = []byte(stem(string(token.Term)))
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
