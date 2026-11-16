package entity

type HitPosition struct {
	Pos   int
	Start int
	End   int
}

type TitleSuggestion struct {
	Score    float64
	FieldHit map[string][]HitPosition
	Title    string
	HymnNo   int
	Variant  string
}

type SearchMatch struct {
	TitleSuggestion
	Content       string
	ContentVerse  string
	VersePosNo    []float64
	VersePosStart []float64
	VersePosEnd   []float64
}
