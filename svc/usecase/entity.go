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
	Snippet string `json:"snippet"`
	Spans   []Span `json:"spans"`
}

type HighlightSearchResult struct {
	Title      HighlightText         `json:"title"`
	Verse      map[int]HighlightText `json:"verse,omitempty"`
	HymnNo     int                   `json:"hymnno,omitempty"`
	Variant    string                `json:"variant,omitempty"`
	FullTitle  string                `json:"full_title"`
	MoreVerses int                   `json:"moreverses,omitempty"`
	Tier       HitTier               `json:"tier"`
}
