package document

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
)

func BuildContent(ctx context.Context, repo repository.Repository, metadata *repository.HymnMetadata, verse int) string {
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
			measure.Build(ctx)
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
	whole := [][]entity.LyricWordVerse{}

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
