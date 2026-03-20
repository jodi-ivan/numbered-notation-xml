package usage

import (
	"fmt"
	"strconv"
	"strings"
)

// Entry represents a single hymn number (with optional verse) mapped to its categories.
type Entry struct {
	H1    string `json:"category"`     // top-level section (all caps)
	H2    string `json:"sub_category"` // sub-category
	Num   int    `json:"-"`
	Verse int    `json:"verse,omitempty"` // 0 means "any verse" / not specified
}

// rawData holds the index as [h1, h2, comma-separated token list].
// Tokens may be:
//   - "287"       → single number
//   - "1-22"      → inclusive range
//   - "287:3"     → number with specific verse
var rawData = [][3]string{
	{"MENGHADAP ALLAH", "Puji-pujian dan Pembukaan Ibadah", "1-22,47,57,60,62,64,191,242-244,246,287-288,290-292,294,295,391,452,454,456,459,464"},
	{"MENGHADAP ALLAH", "Pengakuan dan Pengampunan Dosa", "23-41,169,174,179,185,213,286,293,297,300,305,309,351-359,361,362,380,381,386-388,392,394,398,404,411,434,435,453,459,467"},
	{"MENGHADAP ALLAH", "Kyrie dan Gloria", "42-48,5,7,13,242-244,287:3,303,347"},
	{"PELAYANAN FIRMAN", "Pembacaan Alkitab", "49-59,73,145,150,228-231,233,236-238,240,423,464,473,474"},
	{"PELAYANAN FIRMAN", "Penciptaan dan Pemeliharaan", "60-69,1,3,4,10,19,243,244,287,289,291,298,334,337,385,389,449,461,467"},
	{"PELAYANAN FIRMAN", "Perjanjian Lama", "70-75,4,9,24,46,52,53,69,146,283,285,288,291,292,330,333,334,377,412"},
	{"PELAYANAN FIRMAN", "Penantian Mesias dan Masa Adven", "76-91,73,139,162,260,273,284"},
	{"PELAYANAN FIRMAN", "Kelahiran Yesus dan Masa Natal", "92-127,19,77,87,90,91,136,137,297"},
	{"PELAYANAN FIRMAN", "Akhir Masa Natal dan Epifania", "128-143,19,92,107,110,113-116,118,119,121-123,127,248,281,284,286,293,297,386"},
	{"PELAYANAN FIRMAN", "Kisah Pelayanan Yesus", "144-154,30,74,138,140,283,285,294,298,351-360,370,377,385,398,407,415,418,422,428,431,434,440,451,459,464,468"},
	{"PELAYANAN FIRMAN", "Masa Prapaskah", "155-163,24,28,32,33,46,87,91,152,174,179,286,309,314,372,375,376,381,412,430,443,463"},
	{"PELAYANAN FIRMAN", "Sengsara Yesus dan Jumat Agung", "164-186,32-34,152,156-158,160,286,311-313,368,382,394,404,435,460"},
	{"PELAYANAN FIRMAN", "Kebangkitan Yesus dan Masa Paskah", "187-217,1,5,7,19,65,72,139,152,181,222,226,242,243,246-248,250,281,283,285,286,291,292,295,323,370,373,383,386,394,397,398,404,407,415,435,443"},
	{"PELAYANAN FIRMAN", "Hari Kenaikan", "218-227,5,7,19,41,181,194,202,242,244,247,248,281,284,293,308,383,345,435"},
	{"PELAYANAN FIRMAN", "Roh Kudus dan Hari Pentakosta", "228-241,7,8,16,45,55,56,58,74,242-244,246,257,403"},
	{"PELAYANAN FIRMAN", "Allah Tritunggal dan Hari Trinitatis", "242-246,8,13,16,45,47,48,287:3,303,347,348"},
	{"PELAYANAN FIRMAN", "Gereja dan Kerajaan Allah", "247-261,1,4-7,46,72,74,88,194,213,220,222,224,226,242,243,269,282,330,339,340,345,372,391,434"},
	{"PELAYANAN FIRMAN", "Kehidupan Sorgawi", "262-271,2,5,7,72,219,221,222,224,252,282,283,285,330,355,377,398"},
	{"PELAYANAN FIRMAN", "Akhir Zaman dan Penggenapan", "272-279,5,7,72,74,220,224,225,247,248,260-270,281,282,293,323,340"},
	{"RESPONS TERHADAP PELAYANAN FIRMAN", "Pernyataan Keyakinan Iman", "280-285,5,7,38,46,72,77,194,222,247,248,250,252,305,306,308,309,314,356,367,370,374,376,377,380,383,386-388,392,394,396,405,415"},
	{"RESPONS TERHADAP PELAYANAN FIRMAN", "Pengucapan Syukur dan Persembahan", "286-303,1-12,19,60,62,65,72,77,194,220,249,258,259,291,309,341,361-365,367,393,405,424,433,437,444,450"},
	{"PELAYANAN KHUSUS", "Baptisan Kudus dan Peneguhan Sidi", "304-309,6,19,32,33,36,38,40,41,55,56,58,72,141,146,152,154,228,229,233,236-240,247-261,269,280-285,287,293,298,300,314,330,339,340,344,351,352,356,360,361-376,377,380,381,382,388,392,396,398,402,421,423,430,434,437,446,453,457,466"},
	{"PELAYANAN KHUSUS", "Perjamuan Kudus", "310-315,2,6,9,17,18,32-41,55,56,58,72,74,77,88,128,139,146,153,156,157,160,163,169,174,194,207,210,215,220,222,224,226-240,247-261,262,263,272,273,275,276,279,282,283,285,287-289,291,293,298,305,306,309,323,358,359,361-376,377,380-382,386,388,392,396,398,404,406-412,434,443,464,466,470"},
	{"PELAYANAN KHUSUS", "Pernikahan", "316-318,1,3,9,10,14,16,18,55,56,60,62,65,233,237,239,283,285,287,288,291,295,298,314,330,350,367,370,377,399,407,415,419,444,447,450,461,466"},
	{"PELAYANAN KHUSUS", "Peristiwa Istimewa Gerejawi", "319-320,1-18,22,49-56,60,72,194,213,220,222,229,233,237,239,240,242,243,247-260,272,284,287,291,292,303,314,330,338-343,367,372,376,391,405,422-437,446"},
	{"PELAYANAN KHUSUS", "Pemakaman", "32,37,46,72,128,194,202,207-210,222,261-270,274-279,282,283,285,291,300,329-331,370,372,376,377,380,388,394,412,417,438,445,453,457"},
	{"WAKTU DAN MUSIM", "Pagi dan Siang", "321-323,1-7,19,60-66,194,213,227,237,239,243,245,248,291,298,309,337,344,356,365,384,385,389,390,393,405,407,414,420,421,423,424,437,444,446-448,450,452-470"},
	{"WAKTU DAN MUSIM", "Petang dan Malam", "324-329,8,23,25,29,41,51,60,68,86,148,227,245,282,286,290,291,309,345,383,384,388,389,393,398,405,406,410,411,417,420,422,440-442,445,451-470"},
	{"WAKTU DAN MUSIM", "Pergantian Tahun", "330-335,1-16,23,47,60,72,121,136,138,220,242-246,248,250,255,260,283,285,286-291,298,305,321,328,329,337,343,345,365,370,377,379,383,393,405,406-421,440-442,445,453,457,461,466"},
	{"WAKTU DAN MUSIM", "Musim dan Panen", "333-335,1,3,4,10,19,60-66,74,243,244,287,289-291,295,298,299,302,303,322,337,444,449,461,469,470"},
	{"WAKTU DAN MUSIM", "Bangsa dan Negara", "336-337,1,3,9,15,60,67,260,287,289,291,292,295,298,299,322,330,334,399,444,461,470"},
	{"PENUTUPAN IBADAH", "Pengutusan", "338-344,49,50,128,163,231,243,244,247-260,281,287:3,303,329,345,346,365,372,374-376,402,405-408,414,421-438,440,448,450,457,458,462,466,470"},
	{"PENUTUPAN IBADAH", "Berkat", "345-350,5:7,7:4,478"},
	{"PENUTUPAN IBADAH", "Haleluya, Amin dan Lain-lain", "417-477,1,46,193,205,243,349"},
}

var masterIndex []Entry

func init() {
	if masterIndex == nil {
		masterIndex = buildIndex()
	}
}

// buildIndex parses rawData into a slice of Entry.
func buildIndex() []Entry {
	var entries []Entry
	for _, row := range rawData {
		h1, h2, tokens := row[0], row[1], row[2]
		for _, tok := range strings.Split(tokens, ",") {
			tok = strings.TrimSpace(tok)
			if tok == "" {
				continue
			}
			if strings.Contains(tok, ":") {
				// e.g. "287:3"
				parts := strings.SplitN(tok, ":", 2)
				num, err1 := strconv.Atoi(parts[0])
				verse, err2 := strconv.Atoi(parts[1])
				if err1 == nil && err2 == nil {
					entries = append(entries, Entry{h1, h2, num, verse})
				}
			} else if strings.Contains(tok, "-") {
				// e.g. "1-22"
				parts := strings.SplitN(tok, "-", 2)
				a, err1 := strconv.Atoi(parts[0])
				b, err2 := strconv.Atoi(parts[1])
				if err1 == nil && err2 == nil {
					for n := a; n <= b; n++ {
						entries = append(entries, Entry{h1, h2, n, 0})
					}
				}
			} else {
				num, err := strconv.Atoi(tok)
				if err == nil {
					entries = append(entries, Entry{h1, h2, num, 0})
				}
			}
		}
	}
	return entries
}

// lookup finds all categories for a given hymn number and optional verse.
// verse == 0 means "no filter" — returns all matches regardless of verse.
func Lookup(num, verse int) []Entry {
	seen := map[string]bool{}
	var results []Entry
	for _, e := range masterIndex {
		// number must match
		if e.Num != num {
			continue
		}
		// if entry has a specific verse AND caller specified a verse, they must match
		if e.Verse != 0 && verse != 0 && e.Verse != verse {
			continue
		}
		key := fmt.Sprintf("%s|%s|%d", e.H1, e.H2, e.Verse)
		if !seen[key] {
			seen[key] = true
			results = append(results, e)
		}
	}
	return results
}

// func printResults(results []Entry, num, verse int) {
// 	label := strconv.Itoa(num)
// 	if verse != 0 {
// 		label += fmt.Sprintf(":%d", verse)
// 	}
// 	if len(results) == 0 {
// 		fmt.Printf("No categories found for hymn %s.\n", label)
// 		return
// 	}
// 	fmt.Printf("Hymn %s appears in %d category(ies):\n\n", label, len(results))
// 	prevH1 := ""
// 	for _, e := range results {
// 		if e.H1 != prevH1 {
// 			fmt.Printf("  [ %s ]\n", e.H1)
// 			prevH1 = e.H1
// 		}
// 		verseNote := ""
// 		if e.Verse != 0 {
// 			verseNote = fmt.Sprintf(" (verse %d)", e.Verse)
// 		}
// 		fmt.Printf("      • %s%s\n", e.H2, verseNote)
// 	}
// 	fmt.Println()
// }

// func parseArg(s string) (num, verse int, err error) {
// 	if strings.Contains(s, ":") {
// 		parts := strings.SplitN(s, ":", 2)
// 		num, err = strconv.Atoi(parts[0])
// 		if err != nil {
// 			return
// 		}
// 		verse, err = strconv.Atoi(parts[1])
// 		return
// 	}
// 	num, err = strconv.Atoi(s)
// 	return
// }

// func interactive(index []Entry) {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Println("Hymn Index Lookup — type a number (e.g. 287 or 287:3), or 'q' to quit.")
// 	for {
// 		fmt.Print("> ")
// 		if !scanner.Scan() {
// 			break
// 		}
// 		input := strings.TrimSpace(scanner.Text())
// 		if input == "q" || input == "quit" || input == "exit" {
// 			break
// 		}
// 		if input == "" {
// 			continue
// 		}
// 		num, verse, err := parseArg(input)
// 		if err != nil {
// 			fmt.Println("Invalid input. Enter a number like 287 or 287:3.")
// 			continue
// 		}
// 		printResults(lookup(index, num, verse), num, verse)
// 	}
// }

// func main() {
// 	index := buildIndex()

// 	if len(os.Args) > 1 {
// 		// Non-interactive: hymn [287] or [287:3]
// 		for _, arg := range os.Args[1:] {
// 			num, verse, err := parseArg(arg)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "Invalid argument %q: %v\n", arg, err)
// 				continue
// 			}
// 			printResults(lookup(index, num, verse), num, verse)
// 		}
// 		return
// 	}

// 	interactive(index)
// }
