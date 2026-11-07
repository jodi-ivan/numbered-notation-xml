package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/char/regexp"
	_ "github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/v2/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	bleve_unicode "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/stem"
	_ "github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/stem"

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
		verses[verse-1] = VerseDoc{No: verse, Content: BuildContent(repo, metadata, verse)}
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

func dump(m mapping.IndexMapping, analyzer, text string) {
	an := m.AnalyzerNamed(analyzer)
	var terms []string
	for _, t := range an.Analyze([]byte(text)) {
		terms = append(terms, fmt.Sprintf("%q", t.Term))
	}
	fmt.Printf("%-20s %-10s -> %v\n", analyzer, text, terms)
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

	var charFilters []string
	for _, r := range stem.ElisionRules {
		err := indexMapping.AddCustomCharFilter(r.Name, map[string]interface{}{
			"type":    regexp.Name,
			"regexp":  r.Pattern,
			"replace": r.Replace,
		})
		if err != nil {
			log.Fatal(err)

		}
		charFilters = append(charFilters, r.Name)
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

	if err := indexMapping.AddCustomAnalyzer("indonesian", map[string]interface{}{
		"type":          custom.Name,
		"char_filters":  charFilters,
		"tokenizer":     bleve_unicode.Name,
		"token_filters": []string{lowercase.Name, "indonesian_stem"},
	}); err != nil {
		panic(err)
	}

	docMapping := bleve.NewDocumentMapping()

	contentMapping := bleve.NewTextFieldMapping()
	contentMapping.Store = true
	contentMapping.IncludeTermVectors = true
	contentMapping.Analyzer = "indonesian"

	docMapping.AddFieldMappingsAt("content", contentMapping)
	versesMapping := bleve.NewDocumentMapping()

	noMapping := bleve.NewNumericFieldMapping()
	noMapping.Store = true

	verseContentMapping := bleve.NewTextFieldMapping()
	verseContentMapping.Store = true
	verseContentMapping.Analyzer = "indonesian"

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
	titleMap.Analyzer = "indonesian"

	// ngram for autocomplete
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

	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	// for _, w := range []string{"telah", "t'lah", "cahaya", "cah'ya"} {
	// 	dump(indexMapping, "indonesian", w)
	// 	dump(indexMapping, "part_word_analyzer", w)
	// }

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
		222: []string{"a", "b"},
		229: []string{"a", "b"},
		240: []string{"a", "b"},
		248: []string{"a", "b"},
		250: []string{"a", "b"},
		287: []string{"a", "b"},
		303: []string{"a", "b"},
		311: []string{"a", "b"},
		312: []string{"a", "b"},
		365: []string{"a", "b", "c"},
		369: []string{"a", "b"},
		466: []string{"a", "b"},
		473: []string{"a", "b"},
		478: []string{"a", "b", "c"},
	}
	// Add documents
	documents := []*Document{}

	for no := 1; no <= 478; no++ {
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
				documents = append(documents, doc)

			}
		} else {
			doc, err = BuildDocument(repo, no)
			if err != nil {
				log.Fatalf("Failed to build docs: %s", err.Error())
				return
			}
			documents = append(documents, doc)

		}

	}

	// Iterate and index the documents
	batch := index.NewBatch()
	for _, doc := range documents {
		// raw, _ := json.Marshal(doc.Verses)
		// log.Printf("%d. %s \n %s \n %s\n\n", doc.HymnNo, doc.Title, doc.Content, string(raw))
		err := batch.Index(doc.ID, doc)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := index.Batch(batch); err != nil {
		log.Fatal(err)
	}
	stats, _ := index.Stats().MarshalJSON()
	log.Println("success to build the index. ", string(stats))
}
