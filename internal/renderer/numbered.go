package renderer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func googlefont(f string) []byte {
	empty := []byte{}
	r, err := http.Get(gwfURI + url.QueryEscape(f))
	log.Println("error call", err)
	if err != nil {
		return empty
	}
	defer r.Body.Close()
	b, rerr := ioutil.ReadAll(r.Body)
	log.Println(rerr, r.Status, string(b))
	if rerr != nil || r.StatusCode != http.StatusOK {
		return empty
	}

	return b
}

func RenderNumbered(w http.ResponseWriter, r *http.Request, music musicxml.MusicXML) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(200)
	s := canvas.NewCanvas(svg.New(w))
	s.Start(LAYOUT_WIDTH, 1000)
	ctx := r.Context()
	s.Def()
	fmt.Fprintf(s.Writer(), fontfmt, string(googlefont("Caladea|Old Standard TT|Noto Music")))
	s.DefEnd()

	relativeY := 100
	// render title
	workTitle := strings.ToUpper(music.Work.Title)
	if workTitle == "" {
		workTitle = strings.ToUpper(music.Credit.Words)
	}
	titleWidth := lyric.CalculateLyricWidth(workTitle)
	titleX := (LAYOUT_WIDTH / 2) - (titleWidth * 0.5)
	s.Text(int(titleX), relativeY, workTitle)

	// render key signature
	relativeY += 25

	keySignature := keysig.NewKeySignature(music.Part.Measures[0].Attribute.Key)
	timeSignature := timesig.NewTimeSignatures(ctx, music.Part.Measures)

	humanizedKeySignature := keySignature.String()

	s.Text(LAYOUT_INDENT_LENGTH, relativeY, keySignature.String())

	// render time signature
	// TODO: check the time signature on github issue
	// TODO: time signature changing happens on the top and not on the measure
	/*
		time signatures
		4/4
		3/4
		6/4
		1/4
		6/8 (shown as 3 x 2)
		2/4
	*/
	beat := music.Part.Measures[0].Attribute.Time
	s.Text(LAYOUT_INDENT_LENGTH+(len(humanizedKeySignature)*LOWERCASE_LENGTH), relativeY, fmt.Sprintf("%d ketuk", beat.Beats))
	relativeY += 50

	// RenderMeasures(r.Context(), s, LAYOUT_INDENT_LENGTH, relativeY, music.Part)
	staff := SplitLines(ctx, music.Part)
	x := LAYOUT_INDENT_LENGTH
	for _, st := range staff {
		multiline, marginBottom, marginLeft := RenderStaff(ctx, s, x, relativeY, keySignature, timeSignature, st)
		relativeY = relativeY + 70 + marginBottom
		if multiline {
			x = marginLeft
		} else {
			x = LAYOUT_INDENT_LENGTH
		}
	}
	s.End()
}

// SplitLines split the measure in the lines manner
func SplitLines(ctx context.Context, part musicxml.Part) [][]musicxml.Measure {
	result := [][]musicxml.Measure{}
	currentLine := []musicxml.Measure{}
	for _, measure := range part.Measures {

		if measure.Print != nil && measure.Print.NewSystem == musicxml.PrintNewSystemTypeYes {
			finishLine := make([]musicxml.Measure, len(currentLine))
			copy(finishLine, currentLine)

			result = append(result, finishLine)

			currentLine = []musicxml.Measure{}
		}
		currentLine = append(currentLine, measure)

	}

	return append(result, currentLine)
}
