package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/blevesearch/bleve/v2"
	"github.com/jodi-ivan/numbered-notation-xml/cmd/lab/search/usage"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/storage"
)

var originalTitleExp *regexp.Regexp

func init() {
	if originalTitleExp == nil {
		originalTitleExp = regexp.MustCompile("<i>.*</i>")
	}
}

type VerseDoc struct {
	No      int    `json:"no"`
	Content string `json:"content"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
}
type Document struct {
	ID            string        `json:"id"`
	Title         string        `json:"title"`
	Content       string        `json:"content"`
	Catergory     []usage.Entry `json:"categories,omitempty"`
	HymnNo        int           `json:"hymn_no"`
	Variant       string        `json:"variant,omitempty"`
	Verses        []VerseDoc    `json:"verses"`
	OriginalTitle []string      `json:"original_title"`
	BE            int           `json:"be_num,omitempty"`
	NR            int           `json:"nr_num,omitempty"`
	ForKids       bool          `json:"for_kids"`
	Copyright     string        `json:"copyright,omitempty"`
	MusicCredit   string        `json:"music"`
}

var categories = map[string]map[string][2]int{
	"Menghadap Allah": map[string][2]int{
		"Puji-pujian dan Pembukaan Ibadah": [2]int{1, 22},
		"Pengakuan dan Pengampunan Dosa":   [2]int{23, 41},
		"Kyrie dan Gloria":                 [2]int{42, 48},
	},
	"Pelayanan Firman": map[string][2]int{
		"Pembacaan Alkitab":                    [2]int{49, 59},
		"Penciptaan dan Pemeliharaan":          [2]int{60, 69},
		"Pejanjian Lama":                       [2]int{70, 75},
		"Penantian Mesias dan Masa Adven":      [2]int{76, 91},
		"Kelahiran Yesus dan Masa Natal":       [2]int{92, 127},
		"Akhir Masa Natal dan Epifania":        [2]int{128, 143},
		"Kisah Pelayanan Yesus":                [2]int{144, 154},
		"Masa Prapaskah":                       [2]int{155, 163},
		"Sengsara Yesus dan Jumat Agung":       [2]int{164, 186},
		"Kebangkitan Yesus dan Masa Paskah":    [2]int{187, 217},
		"Hari Kenaikan":                        [2]int{218, 227},
		"Roh Kudus dan Hari Pentakosta":        [2]int{228, 241},
		"Allah Tritunggal dan Hari Trinitatis": [2]int{242, 246},
		"Gereja dan Kerajaan Allah":            [2]int{247, 261},
		"Kehidupan Sorgawi":                    [2]int{262, 271},
		"Akhir Zaman dan Penggenapan":          [2]int{272, 279},
	},
	"Respons Terhadap Pelayanan Firman": map[string][2]int{
		"Pernyataan Keyakinan Iman":         [2]int{280, 285},
		"Pengucapan Syukur dan Persembahan": [2]int{286, 303},
	},
	"Pelayanan Khusus": map[string][2]int{
		"Baptisan Kudus dan Peneguhan Sidi": [2]int{304, 309},
		"Perjamuan Kudus":                   [2]int{310, 315},
		"Pernikahan":                        [2]int{316, 318},
		"Peristiwa Isimewa Gerejawi":        [2]int{319, 320},
	},
	"Waktu dan Musim": map[string][2]int{
		"Pagi dan Siang":    [2]int{321, 323},
		"Petang dan Malam":  [2]int{323, 329},
		"Pergantian Tahun":  [2]int{330, 332},
		"Musim dan Panen":   [2]int{333, 335},
		"Bangsa dan Negara": [2]int{336, 337},
	},
	"Penutupan Ibadah": map[string][2]int{
		"Pungutusan": [2]int{338, 344},
		"Berkat":     [2]int{345, 350},
	},
	"Hidup Beriman Sehari-hari": map[string][2]int{
		"Panggilan Juruselamat":           [2]int{351, 360},
		"Penyerahan Diri":                 [2]int{361, 376},
		"Kebesaran Rahmat Tuhan":          [2]int{377, 390},
		"Sukacita dalam Tuhan":            [2]int{391, 399},
		"Hidup Bersama Tuhan":             [2]int{400, 405},
		"Tuntunan Tuhan":                  [2]int{406, 421},
		"Tanggung Jawab Pengikut Kristus": [2]int{422, 437},
		"Kemenangan dalam Perjuangan":     [2]int{438, 446},
		"Keluarga dan Persekutuan":        [2]int{447, 451},
		"Doa dan Setiap Waktu":            [2]int{452, 471},
	},
	"Haleluya, Amin dan Lain-lain": map[string][2]int{
		"": [2]int{472, 478},
	},
}

func GetCategory(num int) (string, string) {

	for cat, subs := range categories {
		for subcat, ranges := range subs {
			if num >= ranges[0] && num <= ranges[1] {
				return cat, subcat
			}
		}
	}

	return "", ""

}

func GetOriginalTitle(metadata *repository.HymnMetadata) []string {
	if strings.Contains(metadata.Lyric, "<i>") {
		submatch := originalTitleExp.FindStringSubmatch(metadata.Lyric)
		if len(submatch) > 0 && len(submatch[0]) > 6 {
			stripped := submatch[0][3 : len(submatch[0])-4]
			strings.TrimSuffix(stripped, ",")
			return strings.Split(stripped, "/")
		}
	}

	return []string{metadata.Title}
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
		// if i < len(verses)-1 {
		// 	sb.WriteString(" ")
		// }
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
		HymnNo:        metadata.Number,
		BE:            int(metadata.RefBE.Int16),
		NR:            int(metadata.RefNR.Int16),
		Copyright:     metadata.Copyright.String,
		MusicCredit:   metadata.Music,
		OriginalTitle: GetOriginalTitle(metadata),
		Catergory:     []usage.Entry{},
	}

	if len(vaiant) > 0 {
		doc.Variant = vaiant[0]
		doc.ID += vaiant[0]
	}
	verses := make([]VerseDoc, len(metadata.Verse)+1)
	for verse := 1; verse <= len(metadata.Verse)+1; verse++ {
		currCats := usage.Lookup(num, verse)
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
			doc.Catergory = append(doc.Catergory, usage.Entry{H1: h1, H2: h2})
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

	titleMap := bleve.NewNumericFieldMapping()
	titleMap.Store = true

	versesMapping.AddFieldMappingsAt("no", noMapping)
	versesMapping.AddFieldMappingsAt("content", verseContentMapping)
	versesMapping.AddFieldMappingsAt("start", startMapping)
	versesMapping.AddFieldMappingsAt("end", endMapping)
	docMapping.AddSubDocumentMapping("verses", versesMapping)
	docMapping.AddFieldMappingsAt("variant", variantMapping)
	docMapping.AddFieldMappingsAt("hymn_no", hymnMap)
	docMapping.AddFieldMappingsAt("title", titleMap)

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

	for no := 1; no <= 25; no++ {
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
		log.Println("indexing", doc.ID)
	}
	if err := index.Batch(batch); err != nil {
		log.Fatal(err)
	}

}
