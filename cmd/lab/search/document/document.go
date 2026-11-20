package document

import (
	"context"
	"strconv"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/category"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
)

type VerseDoc struct {
	No      int    `json:"no"`
	Content string `json:"content"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
}

type Document struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	Content       string           `json:"content"`
	Catergory     []category.Entry `json:"categories,omitempty"`
	HymnNo        int              `json:"hymn_no"`
	Variant       string           `json:"variant,omitempty"`
	Verses        []VerseDoc       `json:"verses"`
	OriginalTitle []string         `json:"original_title"`
	BE            int              `json:"be_num,omitempty"`
	NR            int              `json:"nr_num,omitempty"`
	ForKids       bool             `json:"for_kids"`
	Copyright     string           `json:"copyright,omitempty"`
	MusicCredit   string           `json:"music"`
	Title_ngram   string           `json:"title_ngram"`
}

func GetOriginalTitle(metadata *repository.HymnMetadata) []string {
	result := []string{}
	if strings.HasPrefix(metadata.Lyric, "<i>") {
		index := strings.Index(metadata.Lyric, "</i>")
		lyric := metadata.Lyric[3:index]
		lyric = strings.TrimSpace(lyric)
		lyric = strings.TrimSuffix(lyric, ",")
		result = strings.Split(lyric, "/")
	} else {
		result = []string{metadata.Title}
	}

	for i, v := range result {
		result[i] = strings.TrimSpace(v)
	}

	return result
}

func BuildDocument(repo repository.Repository, num int, vaiant ...string) (*Document, error) {
	// get hymn meta data

	metadata, err := repo.GetHymnMetaData(context.Background(), num, vaiant...)
	if err != nil {
		return nil, err
	}

	cats := map[string][]category.Entry{}

	doc := &Document{
		ID:            strconv.Itoa(metadata.HymnID),
		Title_ngram:   metadata.Title,
		Title:         metadata.Title,
		HymnNo:        metadata.Number,
		BE:            int(metadata.RefBE.Int16),
		NR:            int(metadata.RefNR.Int16),
		Copyright:     metadata.Copyright.String,
		MusicCredit:   metadata.Music,
		OriginalTitle: GetOriginalTitle(metadata),
		Catergory:     []category.Entry{}, // Initialize as []category.Entry
	}

	if len(vaiant) > 0 {
		doc.Variant = vaiant[0]
		doc.ID += vaiant[0]
	}
	verses := make([]VerseDoc, len(metadata.Verse)+1)
	for verse := 1; verse <= len(metadata.Verse)+1; verse++ {
		currCats := category.Lookup(num, verse)
		for _, curr := range currCats {
			if cats[curr.H1] == nil {
				cats[curr.H1] = []category.Entry{}
			}

			cats[curr.H1] = append(cats[curr.H1], curr)
		}
		verses[verse-1] = VerseDoc{No: verse, Content: BuildContent(context.Background(), repo, metadata, verse)}
	}

	h2s := map[string]bool{} // Tracks unique H2 values
	for _, cat := range cats {
		for _, entry := range cat {
			if h2s[entry.H2] {
				continue
			}
			h2s[entry.H2] = true
			doc.Catergory = append(doc.Catergory, entry)
		}
	}

	doc.Content, verses = buildContentWithOffsets(verses)
	doc.Verses = verses

	return doc, nil

}
