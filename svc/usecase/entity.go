package usecase

type verseRange struct {
	no    int
	start int
	end   int
}

type Span struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type HighlightText struct {
	Snippet   string `json:"snippet"`
	Spans     []Span `json:"spans"`
	SpanVerse []Span `json:"span_verses,omitempty"`
}

type HighlightSearchResult struct {
	Score      float64               `json:"score,omitempty"`
	Title      HighlightText         `json:"title"`
	Verse      map[int]HighlightText `json:"verse,omitempty"`
	HymnNo     int                   `json:"hymnno,omitempty"`
	Variant    string                `json:"variant,omitempty"`
	FullTitle  string                `json:"full_title"`
	MoreVerses int                   `json:"moreverses,omitempty"`
	Tier       HitTier               `json:"tier,omitempty"`

	TotalHitVerse int `json:"total_verse,omitempty"`
}
