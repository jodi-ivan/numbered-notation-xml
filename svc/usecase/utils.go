package usecase

import (
	"sort"

	"github.com/jodi-ivan/numbered-notation-xml/svc/entity"
)

func toSpans(ps []entity.HitPosition) []Span {
	resutl := make([]Span, 0, len(ps))
	for _, p := range ps {
		resutl = append(resutl, Span{Start: p.Start, End: p.End})
	}
	return resutl
}

func normalize(spans []Span, size int) []Span {
	valid := []Span{}
	for _, s := range spans {
		if s.Start >= 0 && s.End <= size && s.Start < s.End {
			valid = append(valid, s)
		}
	}
	sort.Slice(valid, func(i, j int) bool { return valid[i].Start < valid[j].Start })

	merged := valid[:0]
	for _, s := range valid {
		if n := len(merged); n > 0 && s.Start <= merged[n-1].End {
			if s.End > merged[n-1].End {
				merged[n-1].End = s.End
			}
			continue
		}
		merged = append(merged, s)
	}
	return merged
}
