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

var definedSplitter = map[string][2]string{
	"selamalamanya": {"selama", "lamanya"},
	"anakanak":      {"anak", "anak"},
	"porakporanda":  {"porak", "poranda"},
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
	for _, t := range result {
		term := string(t.Term)
		if _, repeatExist := definedSplitter[strings.Repeat(term, 2)]; repeatExist {
			term = strings.Repeat(term, 2)
		}

		splitted, ok := definedSplitter[term]
		if !ok {
			return result
		}

		if splitted[0] == splitted[1] {
			at := &analysis.Token{
				Term:     []byte(splitted[0] + "-" + splitted[0]),
				Position: t.Position,
				Start:    t.Start,
				End:      t.Start + len(splitted[1])*2,
				Type:     t.Type,
			}
			result = append(result, at)
			return result
		}
		for _, s := range splitted {
			at := &analysis.Token{
				Term:     []byte(s),
				Position: t.Position,
				Start:    t.Start,
				End:      t.End,
				Type:     t.Type,
			}
			result = append(result, at)
		}

	}

	return result
}

type termMapFilter struct{ m map[string]string }

func (f *termMapFilter) Filter(in analysis.TokenStream) analysis.TokenStream {
	for _, t := range in {
		if r, ok := f.m[string(t.Term)]; ok {
			t.Term = []byte(r) // Start/End untouched
		}
	}
	return in
}

func init() {

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

	registry.RegisterTokenFilter("id_spelling", func(cfg map[string]interface{}, _ *registry.Cache) (analysis.TokenFilter, error) {
		return &termMapFilter{m: map[string]string{
			"cerlang":   "cemerlang",
			"guruh":     "gemuruh",
			"putra":     "putera",
			"putri":     "puteri",
			"kurnia":    "karunia",
			"'ngajak":   "mengajak",
			"selamanya": "selamalamanya",
		}}, nil
	})
}
