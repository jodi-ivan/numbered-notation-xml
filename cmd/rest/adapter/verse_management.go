package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	lab_verse "github.com/jodi-ivan/numbered-notation-xml/cmd/lab/verse"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/internal/verse"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/julienschmidt/httprouter"
)

type VerseManagement struct {
	VerseRepo repository.Repository
}

type Input struct {
	Style   int                      `json:"style"`
	Col     int                      `json:"col"`
	Row     int                      `json:"row"`
	Content [][]verse.LyricWordVerse `json:"content"`
}

func (vm *VerseManagement) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	hymn, err := strconv.Atoi(ps.ByName("hymn"))
	if err != nil {
		log.Printf("invalid hymn number: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL hymn"))
		return
	}

	verse, err := strconv.Atoi(ps.ByName("verse"))
	if err != nil {
		log.Printf("invalid verse number: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL verse"))
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	input := &Input{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error input: :%s", err.Error())))
		return
	}

	stringify, err := json.Marshal(input.Content)
	if err != nil {
		log.Printf("cannot stringfy content: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot stringfy content"))
		return
	}
	id, err := vm.VerseRepo.InsertVerse(ctx, nil, hymn, verse, input.Style, input.Col, input.Row, string(stringify))
	if err != nil {
		log.Printf("Failed to insert the verse: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert verse"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("sucess insert with ID: %d", id)))

}

type VerseManagementV2 struct {
	VerseRepo repository.Repository
	Db        *sqlx.DB
}

func (vmv2 *VerseManagementV2) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	hymnNo, _, err := utils.ParseHymnWithVariant(ps.ByName("number"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_ = hymnNo
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	input := struct {
		Style     int               `json:"style"`
		Breakdown []string          `json:"breakdown"`
		Generated map[string]string `json:"generated"`
	}{}

	err = json.Unmarshal(b, &input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error input: :%s", err.Error())))
		return
	}

	verses := map[int][]string{}

	currentVerse := 0
	for _, line := range input.Breakdown {
		if parsed, err := strconv.Atoi(string(line[0])); err == nil {
			offset := 0
			if unicode.IsDigit(rune(line[1])) {
				parsed, err = strconv.Atoi(string(line[0:2]))
				if err != nil {
					log.Println("failed to parse ", parsed)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(fmt.Sprintf("error input: :%s", err.Error())))
					return
				}
				offset = 1
			}
			currentVerse = parsed
			verses[currentVerse] = []string{}
			line = line[3+offset:] // 2.
		}

		verses[currentVerse] = append(verses[currentVerse], line)
	}

	collection := map[int][][]verse.LyricWordVerse{}

	for no, lines := range verses {
		param := make([][]verse.LyricWordVerse, len(lines))
		for l, line := range lines {
			words := strings.Split(line, " ")
			for _, word := range words {
				doubledash := strings.ReplaceAll(word, "--", "|-")
				breakdown := strings.Split(doubledash, "-")
				cleanup := strings.Join(breakdown, "")
				if strings.Contains(cleanup, "_") {
					cleanup = strings.ReplaceAll(cleanup, "_", "")
				}
				bRes := verse.LyricWordVerse{
					Word:      strings.ReplaceAll(cleanup, "|", "-"),
					Breakdown: []verse.LyricPartVerse{},
				}

				for i, b := range breakdown {
					syll := musicxml.LyricSyllabicTypeMiddle
					if len(breakdown) == 1 {
						syll = musicxml.LyricSyllabicTypeSingle
					} else if i == 0 {
						syll = musicxml.LyricSyllabicTypeBegin
					} else if i == len(breakdown)-1 {
						syll = musicxml.LyricSyllabicTypeEnd
					}

					if strings.Contains(b, "_") {
						b = strings.ReplaceAll(b, "_", "")
						word = strings.ReplaceAll(word, "_", "")
						start, end := -1, -1
						for ir, r := range b {
							if lab_verse.IsVowel(string(r)) || (r == 'h' && start != -1) {
								if start == -1 {
									start = ir
								} else {
									end = ir
								}
							}
						}

						partBreakdown := []verse.LyricStylePart{}
						if start == 0 {
							partBreakdown = []verse.LyricStylePart{
								{
									Text:      b[start : end+1],
									Underline: true,
								},
								{
									Text:      b[end+1:],
									Underline: false,
								},
							}
						} else if end == len(b)-1 {
							partBreakdown = []verse.LyricStylePart{
								{
									Text:      b[0:start],
									Underline: false,
								},
								{
									Text:      b[start:],
									Underline: true,
								},
							}
						} else {
							partBreakdown = []verse.LyricStylePart{
								{
									Text:      b[0:start],
									Underline: false,
								},
								{
									Text:      b[start : end+1],
									Underline: true,
								},
								{
									Text:      b[end+1:],
									Underline: false,
								},
							}

						}
						bRes.Breakdown = append(bRes.Breakdown, verse.LyricPartVerse{
							Text:      b,
							Type:      syll,
							Combine:   true,
							Breakdown: partBreakdown,
						})

					} else {
						bRes.Breakdown = append(bRes.Breakdown, verse.LyricPartVerse{
							Text: b,
							Type: syll,
						})
					}
				}

				if param[l] == nil {
					param[l] = []verse.LyricWordVerse{}
				}
				param[l] = append(param[l], bRes)
			}

		}

		collection[no] = param
	}

	param := map[int]struct {
		No      int
		Style   int
		Col     int
		Row     int
		Content string
	}{}

	verseIds := []int{}
	// maxverse := len(collection) + 1

	count := len(collection)
	half := count / 2
	for no, v := range collection {
		style := input.Style
		col := 1
		row := no - 1

		if input.Style == 6 {
			// Calculate a 0-based index for logic (since 'no' starts at 2)
			row++
			idx := no - 2
			// 1. Check if this is the "Odd One Out" (the last item in an odd list)
			if count%2 != 0 && idx == count-1 {
				style = 12
				col = 1
				row = half + 1 // Placed at the bottom across both columns
			} else {
				// 2. Vertical Column Logic
				if idx < half {
					// First Column (Vertical flow)
					col = 1
					row = idx + 1
				} else {
					// Second Column (Starts at Verse 4 if count is 5)
					col = 2
					row = (idx - half) + 1
				}
			}

		}

		stringify, err := json.Marshal(v)
		if err != nil {
			log.Printf("cannot stringfy content: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("cannot stringfy content"))
			return
		}

		param[no] = struct {
			No      int
			Style   int
			Col     int
			Row     int
			Content string
		}{
			No:      no,
			Style:   style,
			Row:     row,
			Col:     col,
			Content: string(stringify),
		}

	}

	tx, err := vmv2.VerseRepo.StartTransaction(r.Context())
	if err != nil {
		log.Printf("failed to start transaction: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to start transaction"))
		return
	}

	defer tx.Rollback()
	for no, p := range param {
		id, err := vmv2.VerseRepo.InsertVerse(r.Context(), tx, hymnNo, no, p.Style, p.Col, p.Row, p.Content)
		if err != nil {
			log.Printf("Failed to insert the verse: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to insert verse"))
			return
		}
		verseIds = append(verseIds, id)
	}

	breakdownIds := []int{}
	insertQuery := `INSERT INTO syllable_breakdown (whole, breakdown) VALUES(?, ?) RETURNING id`
	insertQueryWithElision := `INSERT INTO syllable_breakdown (whole, breakdown, elision_index) VALUES(?, ?, ?) RETURNING id`
	for key, sig := range input.Generated {
		var newID int
		args := []interface{}{}
		query := insertQuery
		if strings.Contains(key, "_") {
			eidx := -1
			key = strings.ReplaceAll(key, "_", "")
			b := strings.Split(sig, "-")
			for idx, syll := range b {
				if strings.HasPrefix(syll, "_") {
					eidx = idx
				}
			}
			sig = strings.ReplaceAll(sig, "_", "")
			query = insertQueryWithElision
			args = []interface{}{strings.ToLower(key), sig, eidx}
		} else {
			args = []interface{}{strings.ToLower(key), sig}

		}
		err = (repository.GetSqlTx(tx)).QueryRow(query, args...).Scan(&newID)
		if err != nil {
			log.Println("Failed to insert", key, sig)
		}
		breakdownIds = append(breakdownIds, newID)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to insert the commit: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to commit"))
		return
	}

	response := map[string]interface{}{
		"verses_ids":    verseIds,
		"breakdown_ids": breakdownIds,
	}

	raw, _ := json.MarshalIndent(response, "", "    ")
	w.WriteHeader(http.StatusOK)
	w.Write(raw)

}
