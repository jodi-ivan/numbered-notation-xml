package usecase

import (
	"context"
	"sort"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/svc/entity"
)

var maxVersesShown = 4
var maxSnippetBytes = 125
var ellipsis = "..."

type HitTier int

const (
	TierTitle HitTier = iota
	TierFirstVerse
	TierOtherVerses
	TierUnknown
)

func TierOf(h entity.SearchMatch) HitTier {
	switch {
	case len(h.FieldHit["title"]) > 0:
		return TierTitle
	case len(h.FieldHit["content"]) > 0:
		return TierFirstVerse
	case len(h.FieldHit["verses.content"]) > 0:
		return TierOtherVerses
	}
	return TierUnknown
}

func (u *interactor) rankSearchResult(ctx context.Context, hits []entity.SearchMatch) {
	sort.SliceStable(hits, func(i, j int) bool {
		a, b := hits[i], hits[j]
		if ta, tb := TierOf(a), TierOf(b); ta != tb {
			return ta < tb
		}
		return a.HymnNo < b.HymnNo
	})
}

func (u *interactor) Highlight(ctx context.Context, hits []entity.SearchMatch) []HighlightSearchResult {
	result := make([]HighlightSearchResult, 0, len(hits))
	for _, h := range hits {
		result = append(result, highlightOne(h))
	}
	return result
}

func highlightOne(m entity.SearchMatch) HighlightSearchResult {
	res := HighlightSearchResult{
		HymnNo:  m.HymnNo,
		Variant: m.Variant,
		Title: HighlightText{
			Snippet: m.Title,
			Spans:   normalize(toSpans(m.FieldHit["title"]), len(m.Title)),
		},
		FullTitle: m.Title,
		Verse:     map[int]HighlightText{},
	}

	verses := versesOf(m)
	byNo := make(map[int]verseRange, len(verses))
	for _, v := range verses {
		byNo[v.no] = v
	}

	spansByVerse := map[int][]Span{}
	for _, p := range m.FieldHit["content"] {
		v, ok := verseAt(verses, p.Start)
		if !ok {
			continue
		}
		end := p.End
		if end > v.end {
			end = v.end
		}
		spansByVerse[v.no] = append(spansByVerse[v.no], Span{Start: p.Start - v.start, End: end - v.start})
	}

	nos := make([]int, 0, len(spansByVerse))
	for no := range spansByVerse {
		nos = append(nos, no)
	}
	sort.Ints(nos)

	shown := nos
	if len(nos) > maxVersesShown {
		shown = nos[:maxVersesShown]
		res.MoreVerses = len(nos) - maxVersesShown
	}

	for _, no := range shown {
		text, ok := byNo[no].text(m.Content)
		if !ok {
			continue
		}
		snippet, spans := windowAround(text, normalize(spansByVerse[no], len(text)), maxSnippetBytes)
		res.Verse[no] = HighlightText{Snippet: snippet, Spans: spans}
	}

	res.Tier = TierOf(m)

	return res

}

func verseAt(vs []verseRange, offset int) (verseRange, bool) {
	for _, v := range vs {
		if offset >= v.start && offset < v.end {
			return v, true
		}
	}
	return verseRange{}, false
}

func versesOf(m entity.SearchMatch) []verseRange {
	n := len(m.VersePosNo)
	if len(m.VersePosStart) < n {
		n = len(m.VersePosStart)
	}
	if len(m.VersePosEnd) < n {
		n = len(m.VersePosEnd)
	}
	out := make([]verseRange, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, verseRange{
			no:    int(m.VersePosNo[i]),
			start: int(m.VersePosStart[i]),
			end:   int(m.VersePosEnd[i]),
		})
	}
	return out
}

func (v verseRange) text(content string) (string, bool) {
	if v.start < 0 || v.start > len(content) {
		return "", false
	}
	end := v.end
	if end > len(content) {
		end = len(content)
	}
	if end <= v.start {
		return "", false
	}
	return strings.TrimRight(content[v.start:end], " \n\t"), true
}

func windowAround(text string, spans []Span, limit int) (string, []Span) {
	if len(text) <= limit || len(spans) == 0 {
		return text, spans
	}

	first := spans[0]
	lo := first.Start - limit/3 // some left context
	if lo < 0 {
		lo = 0
	}
	hi := lo + limit
	if hi > len(text) {
		hi = len(text)
		lo = hi - limit
		if lo < 0 {
			lo = 0
		}
	}
	if hi < first.End {
		hi = first.End
	}

	if lo > 0 {
		if i := strings.LastIndexAny(text[:lo], " \n\t"); i >= 0 {
			lo = i + 1
		} else {
			lo = 0
		}
	}
	if hi < len(text) {
		if i := strings.IndexAny(text[hi:], " \n\t"); i >= 0 {
			hi += i
		} else {
			hi = len(text)
		}
	}

	prefix, suffix := "", ""
	if lo > 0 {
		prefix = ellipsis
	}
	if hi < len(text) {
		suffix = ellipsis
	}

	shift := len(prefix) - lo // ellipsis is multi-byte, so spans move by its length
	out := make([]Span, 0, len(spans))
	for _, s := range spans {
		if s.Start >= lo && s.End <= hi {
			out = append(out, Span{Start: s.Start + shift, End: s.End + shift})
		}
	}
	return prefix + text[lo:hi] + suffix, out
}
