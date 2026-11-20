package stem

import (
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

var redupSuffix = map[string]bool{
	"": true, "ku": true, "mu": true, "nya": true,
	"lah": true, "kah": true, "tah": true, "pun": true,
}

// anakanakmu -> anakmu, anakanak -> anak
func collapseRedup(s string) string {
	r := []rune(s)
	for i := 3; i*2 <= len(r); i++ { // base >= 3 runes
		if string(r[:i]) == string(r[i:2*i]) && redupSuffix[string(r[2*i:])] {
			return string(r[:i]) + string(r[2*i:])
		}
	}
	return s
}

type redupFilter struct{}

func (redupFilter) Filter(in analysis.TokenStream) analysis.TokenStream {
	for _, t := range in {
		t.Term = []byte(collapseRedup(string(t.Term)))
	}
	return in
}

func init() {
	registry.RegisterTokenFilter("id_redup", func(map[string]interface{}, *registry.Cache) (analysis.TokenFilter, error) {
		return redupFilter{}, nil
	})
}
