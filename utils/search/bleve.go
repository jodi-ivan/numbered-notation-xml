package search

import (
	"context"
	"log"
	"sort"
	"strings"

	core "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/token/ngram"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/search"
	_ "github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/stem"
	"github.com/jodi-ivan/numbered-notation-xml/svc/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
)

func init() {
	indexMapping := core.NewIndexMapping()
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
}

type Provider interface {
	Close() error
	Search(ctx context.Context, req *core.SearchRequest) (*core.SearchResult, error)
}
type index struct {
	bleveIndex core.Index
	open       bool
}

func (i index) Close() error {
	if !i.open {
		return nil
	}
	return i.bleveIndex.Close()
}

func (i index) Search(ctx context.Context, req *core.SearchRequest) (*core.SearchResult, error) {
	req.Highlight = core.NewHighlight()
	req.Explain = true
	searchResult, err := i.bleveIndex.Search(req)

	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

func BuildQuery(ctx context.Context, term string) *core.SearchRequest {
	// should these move to usecase?
	q := core.NewMatchQuery(term)
	q.SetField("content")

	q2 := core.NewMatchQuery(term)
	q2.SetField("verses.content")

	q3 := core.NewMatchQuery(term)
	q3.SetField("title")

	combined1 := core.NewDisjunctionQuery(q, q2, q3)
	searchRequest := core.NewSearchRequest(combined1)
	searchRequest.IncludeLocations = true
	searchRequest.Fields = []string{"hymn_no", "content", "title", "verses.no", "verses.content", "verses.start", "verses.end", "variant", "original_title"}
	searchRequest.Highlight = core.NewHighlight()
	searchRequest.Explain = true

	return searchRequest
}

func NewSearhProvider(conf config.Config) (Provider, error) {
	idx, err := core.Open(conf.Bleve.IndexPath)
	if err != nil {
		return nil, err
	}
	result := &index{
		bleveIndex: idx,
		open:       true,
	}
	return result, nil
}

func NewSearchMatch(hit *search.DocumentMatch) entity.SearchMatch {
	m := entity.SearchMatch{
		Score:         hit.Score,
		FieldHit:      map[string][]entity.HitPosition{},
		Title:         asString(hit.Fields["title"]),
		HymnNo:        asInt(hit.Fields["hymn_no"]),
		Variant:       asString(hit.Fields["variant"]),
		Content:       asString(hit.Fields["content"]),
		ContentVerse:  asString(hit.Fields["verses.content"]),
		VersePosNo:    asFloats(hit.Fields["verses.no"]),
		VersePosStart: asFloats(hit.Fields["verses.start"]),
		VersePosEnd:   asFloats(hit.Fields["verses.end"]),
	}

	// hit.Locations: field -> term -> []*Location
	for field, terms := range hit.Locations {
		for _, locs := range terms {
			for _, l := range locs {
				m.FieldHit[field] = append(m.FieldHit[field], entity.HitPosition{
					Pos:   int(l.Pos),
					Start: int(l.Start),
					End:   int(l.End),
				})
			}
		}
	}
	for _, ps := range m.FieldHit {
		ps := ps
		sort.Slice(ps, func(i, j int) bool { return ps[i].Start < ps[j].Start })
	}
	return m
}

func asString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case []interface{}:
		parts := make([]string, 0, len(t))
		for _, e := range t {
			if s, ok := e.(string); ok {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

func asInt(v interface{}) int {
	if f, ok := v.(float64); ok {
		return int(f)
	}
	return 0
}

func asFloats(v interface{}) []float64 {
	switch t := v.(type) {
	case float64:
		return []float64{t}
	case []interface{}:
		out := make([]float64, 0, len(t))
		for _, e := range t {
			if f, ok := e.(float64); ok {
				out = append(out, f)
			}
		}
		return out
	}
	return nil
}
