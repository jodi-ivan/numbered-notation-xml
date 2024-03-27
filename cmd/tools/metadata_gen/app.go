package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
	"github.com/jodi-ivan/numbered-notation-xml/utils/storage"
)

type KJMetadata struct {
	XMLName       xml.Name  `xml:"kj-metadata"`
	Number        int       `xml:"number,attr"`
	NumberVariant int       `xml:"number-variant,attr,omitempty"`
	KidsStarred   bool      `xml:"kids-starred,attr,omitempty"`
	Title         Title     `xml:"title"`
	Credit        Credit    `xml:"credit"`
	CrossRef      CrossRef  `xml:"cross-reference,omitempty"`
	Category      *Category `xml:"category,omitempty"`
	Verses        Verses    `xml:"verses"`
	Internal      *Internal `xml:"internal,omitempty"`
}

type Title struct {
	Value    string  `xml:"value"`
	Footnote *string `xml:"footnote"`
}

type Credit struct {
	Lyric string `xml:"lyric"`
	Music string `xml:"music"`
}

type CrossRef struct {
	NR []CrossRefItem `xml:"nr"`
	BE []CrossRefItem `xml:"be"`
}

type CrossRefItem struct {
	Number int `xml:"number,attr"`
}

type Category struct {
	ParentID int    `xml:"parent_id,attr,omitempty"`
	Title    string `xml:"title"`
}

type Verses struct {
	Verse []Verse `xml:"verse"`
}

type Verse struct {
	Row   int    `xml:"row,attr,omitempty"`
	No    int    `xml:"no,attr,omitempty"`
	X     int    `xml:"x,attr,omitempty"`
	Y     int    `xml:"y,attr,omitempty"`
	Lines []Line `xml:"line"`
}

type Line struct {
	Word []Word `xml:"word"`
}

type Word struct {
	Text      string      `xml:"text,attr"`
	Breakdown []Breakdown `xml:"breakdown"`
}

type Breakdown struct {
	Type      string      `xml:"syllabic-type,attr,omitempty"`
	Combine   bool        `xml:"combine,attr,omitempty"`
	Underline bool        `xml:"underline,attr,omitempty"`
	Text      string      `xml:"text"`
	Breakdown []Breakdown `xml:"breakdown"`
}

type Internal struct {
	Breaklines *Breaklines `xml:"breaklines"`
	// TODO: Add fields for hidden notes data
}

type Breaklines struct {
	Breakline []Breakline `xml:"breakline"`
}

type Breakline struct {
	Measure   int    `xml:"measure"`
	Note      int    `xml:"note"`
	WordLyric string `xml:"word-lyric"`
}

func main() {

	// get the hymn information
	// verse
	// breakdown

	db, err := storage.NewStorage(context.Background(), "C:\\Users\\jodiv\\go\\src\\github.com\\jodi-ivan\\numbered-notation-xml\\files\\database\\kidung-jemaat.db")
	if err != nil {
		log.Fatalf("Failed to connect to storage: %s", err.Error())
		return
	}

	defer func() {
		db.Close()
	}()

	repo := repository.New(context.Background(), db, &config.Config{
		Metadata: config.MetadataConfig{
			UseXMLFile: true,
		},
	})

	hymnMeta, err := repo.GetHymnMetaData(context.Background(), 2)

	if err != nil {
		log.Fatalf("Failed to get the metadata: %s", err.Error())
		return
	}

	res := KJMetadata{
		Number: hymnMeta.Number,
		Title: Title{
			Value: hymnMeta.Title,
		},
		Credit: Credit{
			Music: hymnMeta.Music,
			Lyric: hymnMeta.Lyric,
		},
		Verses: Verses{
			Verse: make([]Verse, 0),
		},
	}

	verses := map[int][][]lyric.LyricWordVerse{}

	for _, v := range hymnMeta.Verse {
		curr := [][]lyric.LyricWordVerse{}
		err = json.Unmarshal([]byte(v.Content.String), &curr)
		if err != nil {
			log.Println("[RenderVerse] failed to unmarshal, err ", err)
		}
		verses[int(v.VerseNum.Int32)] = curr
	}

	for _, v := range hymnMeta.Verse {
		curr := [][]lyric.LyricWordVerse{}
		err = json.Unmarshal([]byte(v.Content.String), &curr)
		if err != nil {
			log.Println("[RenderVerse] failed to unmarshal, err ", err)
		}
		resVerse := Verse{
			No:    int(v.VerseNum.Int32),
			Lines: []Line{},
			Row:   int(v.StyleRow.Int32),
			// TODO: sync with db
			X: 1,
			Y: int(v.VerseNum.Int32) - 1,
		}

		for _, line := range curr {
			words := []Word{}

			for _, w := range line {
				currWord := Word{
					Text: w.Word,
				}

				for _, bd := range w.Breakdown {
					currBd := Breakdown{
						Type:    string(bd.Type),
						Combine: bd.Combine,
						Text:    bd.Text,
					}

					if currBd.Combine {
						currBd.Breakdown = []Breakdown{}
						for _, lv2db := range bd.Breakdown {
							currBd.Breakdown = append(currBd.Breakdown, Breakdown{
								Text:      lv2db.Text,
								Underline: lv2db.Underline,
							})
						}
					}
					currWord.Breakdown = append(currWord.Breakdown, currBd)
				}

				words = append(words, currWord)
			}
			resVerse.Lines = append(resVerse.Lines, Line{
				Word: words,
			})
		}

		res.Verses.Verse = append(res.Verses.Verse, resVerse)

	}

	xmlRaw, _ := xml.MarshalIndent(res, "", "    ")

	fmt.Println(string(xmlRaw))

	return
}
