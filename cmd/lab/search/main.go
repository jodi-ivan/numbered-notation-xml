package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	_ "github.com/blevesearch/bleve/v2/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/category"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/storage"
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

func BuildContent(repo repository.Repository, metadata *repository.HymnMetadata, verse int) string {
	if verse == 1 {
		path := fmt.Sprintf("files/scores/musicxml/kj-%03d.musicxml", metadata.Number)
		if metadata.Variant.String != "" {
			path = fmt.Sprintf("files/scores/musicxml/kj-%03d%s.musicxml", metadata.Number, metadata.Variant.String)
		}
		xmls, err := repo.GetMusicXML(context.Background(), path)
		if err != nil {
			log.Println("failed to get the xml data, err:", err.Error())
			return ""
		}

		lMapper := map[int]string{}
		prevTotalLyric := -1
		for _, measure := range xmls.Part.Measures {
			measure.Build()
			for _, note := range measure.Notes {
				if len(note.Lyric) == 0 {
					continue
				}

				if prevTotalLyric != len(note.Lyric) && prevTotalLyric != -1 {
					lMapper[1] += lMapper[2]
					lMapper[2] = ""
				}

				for _, l := range note.Lyric {
					syl := ""

					for _, s := range l.Text {
						syl += s.Value
					}
					// li := lyric.NewLyric()
					if unicode.IsDigit(rune(syl[0])) {
						syl = syl[2:]
					}
					lMapper[l.Number] += syl

					if l.Syllabic == musicxml.LyricSyllabicTypeEnd || l.Syllabic == musicxml.LyricSyllabicTypeSingle {
						lMapper[l.Number] += " "
					}

				}

				prevTotalLyric = len(note.Lyric)
			}
		}

		for part := 2; part <= 4; part++ {
			if lMapper[part] != "" {
				lMapper[1] += lMapper[part] + " "
			}
		}

		return lMapper[1]

	}

	result := ""

	if _, ok := metadata.Verse[verse]; !ok {
		return ""
	}
	whole := [][]lyric.LyricWordVerse{}

	err := json.Unmarshal([]byte(metadata.Verse[verse].Content.String), &whole)
	if err != nil {
		log.Println("[RenderVerse] failed to unmarshal, err ", err)
	}

	for _, line := range whole {
		for _, word := range line {
			result += word.Word + " "
		}
	}

	return result
}

func buildContentWithOffsets(verses []VerseDoc) (string, []VerseDoc) {
	var sb strings.Builder
	for i := range verses {
		verses[i].Start = sb.Len()
		sb.WriteString(verses[i].Content)
		verses[i].End = sb.Len()
	}
	return sb.String(), verses
}

func BuildDocument(repo repository.Repository, num int, vaiant ...string) (*Document, error) {
	// get hymn meta data

	metadata, err := repo.GetHymnMetaData(context.Background(), num, vaiant...)
	if err != nil {
		return nil, err
	}

	cats := map[string][]string{}

	doc := &Document{
		ID:            strconv.Itoa(metadata.HymnID),
		Title:         metadata.Title,
		Title_ngram:   metadata.Title,
		HymnNo:        metadata.Number,
		BE:            int(metadata.RefBE.Int16),
		NR:            int(metadata.RefNR.Int16),
		Copyright:     metadata.Copyright.String,
		MusicCredit:   metadata.Music,
		OriginalTitle: GetOriginalTitle(metadata),
		Catergory:     []category.Entry{},
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
				cats[curr.H1] = []string{}
			}

			cats[curr.H1] = append(cats[curr.H1], curr.H2)
		}
		verses[verse-1] = VerseDoc{No: verse, Content: BuildContent(repo, metadata, verse)}
	}

	h2s := map[string]bool{}
	for h1, cat := range cats {
		for _, h2 := range cat {
			if h2s[h2] {
				continue
			}
			h2s[h2] = true
			doc.Catergory = append(doc.Catergory, category.Entry{H1: h1, H2: h2})
		}
	}

	doc.Content, verses = buildContentWithOffsets(verses)
	doc.Verses = verses

	return doc, nil

}

func main() {

	db, err := storage.NewStorage(context.Background(), "files/database/kidung-jemaat.db")
	if err != nil {
		log.Fatalf("Failed to connect to storage: %s", err.Error())
		return
	}

	repo := repository.New(context.Background(), db)

	indexPath := "files/index/kj.bleve"
	// Create a new index
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	// 1. Create a custom token filter for n-grams (3 to 6 characters is usually sweet spot)
	err = indexMapping.AddCustomTokenFilter("my_ngram",
		map[string]interface{}{
			"type": "ngram",
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

	docMapping := bleve.NewDocumentMapping()

	contentMapping := bleve.NewTextFieldMapping()
	contentMapping.Store = true
	contentMapping.IncludeTermVectors = true
	docMapping.AddFieldMappingsAt("content", contentMapping)

	versesMapping := bleve.NewDocumentMapping()

	noMapping := bleve.NewNumericFieldMapping()
	noMapping.Store = true

	verseContentMapping := bleve.NewTextFieldMapping()
	verseContentMapping.Store = true

	startMapping := bleve.NewNumericFieldMapping()
	startMapping.Store = true

	endMapping := bleve.NewNumericFieldMapping()
	endMapping.Store = true

	variantMapping := bleve.NewTextFieldMapping()
	variantMapping.Store = true

	hymnMap := bleve.NewNumericFieldMapping()
	hymnMap.Store = true

	titleMap := bleve.NewTextFieldMapping()
	titleMap.Store = true

	titleNgram := bleve.NewTextFieldMapping()
	titleNgram.Store = true
	titleNgram.Analyzer = "part_word_analyzer"

	ogMap := bleve.NewTextFieldMapping()
	ogMap.Store = true

	versesMapping.AddFieldMappingsAt("no", noMapping)
	versesMapping.AddFieldMappingsAt("content", verseContentMapping)
	versesMapping.AddFieldMappingsAt("start", startMapping)
	versesMapping.AddFieldMappingsAt("end", endMapping)
	docMapping.AddSubDocumentMapping("verses", versesMapping)
	docMapping.AddFieldMappingsAt("variant", variantMapping)
	docMapping.AddFieldMappingsAt("hymn_no", hymnMap)
	docMapping.AddFieldMappingsAt("title", titleMap)
	docMapping.AddFieldMappingsAt("original_title", ogMap)
	docMapping.AddFieldMappingsAt("title_ngram", titleNgram)

	indexMapping.DefaultMapping = docMapping

	var variants = map[int][]string{
		24:  []string{"a", "b"},
		30:  []string{"a", "b"},
		31:  []string{"a", "b"},
		37:  []string{"a", "b"},
		50:  []string{"a", "b"},
		95:  []string{"a", "b"},
		144: []string{"a", "b"},
		146: []string{"a", "b"},
		168: []string{"a", "b", "c"},
		174: []string{"a", "b"},
	}
	// Add documents
	documents := []*Document{}

	for no := 1; no <= 214; no++ {
		var doc *Document
		var err error
		if variant, ok := variants[no]; ok {
			for _, vs := range variant {
				doc, err = BuildDocument(repo, no, vs)
				if err != nil {
					if err != nil {
						log.Fatalf("Failed to build docs: %s", err.Error())
						return
					}
				}

				raw, _ := json.MarshalIndent(doc, "", "    ")
				log.Println(string(raw))
			}
		} else {
			doc, err = BuildDocument(repo, no)
			if err != nil {
				log.Fatalf("Failed to build docs: %s", err.Error())
				return
			}
		}

		documents = append(documents, doc)
	}

	// Iterate and index the documents
	batch := index.NewBatch()
	for _, doc := range documents {
		batch.Index(doc.ID, doc)
	}
	if err := index.Batch(batch); err != nil {
		log.Fatal(err)
	}

}
