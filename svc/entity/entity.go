package entity

type HitPosition struct {
	Pos   int
	Start int
	End   int
}

type SearchMatch struct {
	Score         float64
	FieldHit      map[string][]HitPosition
	Title         string
	HymnNo        int
	Variant       string
	Content       string
	ContentVerse  string
	VersePosNo    []float64
	VersePosStart []float64
	VersePosEnd   []float64
}
