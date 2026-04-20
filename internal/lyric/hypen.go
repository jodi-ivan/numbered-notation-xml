package lyric

import (
	"context"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (li *lyricInteractor) CalculateHypen(ctx context.Context, prevLyric, currentLyric *LyricPosition) (location []entity.Coordinate) {

	if prevLyric.Lyrics.Syllabic == musicxml.LyricSyllabicTypeEnd || prevLyric.Lyrics.Syllabic == musicxml.LyricSyllabicTypeSingle {
		return nil
	}

	hypenWidth := li.CalculateLyricWidth("-")

	lyricText := entity.LyricVal(prevLyric.Lyrics.Text).String()

	startPosition := prevLyric.Coordinate.X + li.CalculateLyricWidth(lyricText)
	endPostion := currentLyric.Coordinate.X
	distance := endPostion - startPosition
	if distance < 4 {

		// force add hyphen at the end if the lyric near the end margin
		if endPostion == float64(constant.LAYOUT_WIDTH-constant.LAYOUT_INDENT_LENGTH) {
			return []entity.Coordinate{
				entity.NewCoordinate(startPosition+2, currentLyric.Coordinate.Y),
			}
		}
		return nil
	}

	// every 1/6 of layout has 2 hypen
	container := (constant.LAYOUT_WIDTH - (2 * constant.LAYOUT_INDENT_LENGTH)) / 6
	if distance < float64(container) {
		offset := (distance / 2) - hypenWidth
		if offset < 0 {
			offset = 0
		}
		return []entity.Coordinate{
			entity.Coordinate{
				X: startPosition + offset,
				Y: currentLyric.Coordinate.Y,
			},
		}
	} else {
		result := []entity.Coordinate{}
		totalContainer := math.Ceil((distance - (2 * hypenWidth)) / float64(container))
		totalHypen := (totalContainer * 2)
		if lyricText == "" {
			result = append(result, entity.NewCoordinate(constant.LAYOUT_INDENT_LENGTH, currentLyric.Coordinate.Y))
			totalHypen += 1
		}
		startPosition += (distance / totalHypen)
		for i := float64(0); i < totalHypen-1; i++ {
			result = append(result,
				entity.NewCoordinate(startPosition+(i*(distance/totalHypen)), currentLyric.Coordinate.Y))
		}

		return result
	}
}

// RenderHypen writes the hypen
// @measure :is the notes for the whole staff (flatten across measures)
func (li *lyricInteractor) RenderHypen(ctx context.Context, canv canvas.Canvas, measure []*entity.NoteRenderer) {
	pos := map[int][2]*LyricPosition{}

	// for tracking the pair of begin to end
	hs := NewHypenStack()
	baseYPos := float64(0)
	var lastLyric []entity.Lyric
	hypenLocation := []entity.Coordinate{}
	hasLyricBefore := false
	for _, n := range measure {

		if len(n.Lyric) == 0 {
			continue
		}

		yPos := float64(n.PositionY) + 25

		spacing := func(y float64, i int) float64 {
			if len(n.Lyric) > MAX_VERSE_IN_MUSIC {
				y = y + (math.Trunc(float64(i)/MAX_LINE_PER_VERSE_IN_MUSIC) * LINE_BETWEEN_LYRIC)
			}

			return y
		}
		lastLyric = n.Lyric
		for i, l := range n.Lyric {
			hs.Process(ctx, l.Syllabic)
			if len(pos[i]) == 0 {
				pos[i] = [2]*LyricPosition{}
			}
			pair := pos[i]
			switch l.Syllabic {
			case musicxml.LyricSyllabicTypeBegin:
				// prev
				pair[0] = &LyricPosition{
					Coordinate: entity.NewCoordinate(float64(n.PositionX), spacing(yPos, i)+float64(i*LINE_BETWEEN_LYRIC)),
					Lyrics:     l,
				}
			case musicxml.LyricSyllabicTypeEnd:
				// curr
				pair[1] = &LyricPosition{
					Coordinate: entity.NewCoordinate(float64(n.PositionX), spacing(yPos, i)+float64(i*LINE_BETWEEN_LYRIC)),
					Lyrics:     l,
				}
				// do the calculation here
				start := pair[0]
				if start == nil {
					start = &LyricPosition{
						Coordinate: entity.NewCoordinate(float64(measure[0].PositionX), spacing(yPos, i)+float64(i*LINE_BETWEEN_LYRIC)),
					}
				}
				hypenLocation = append(li.CalculateHypen(ctx, start, pair[1]), hypenLocation...)
				pair = [2]*LyricPosition{}

			case musicxml.LyricSyllabicTypeMiddle:
				if pair[0] == nil {
					pair[0] = &LyricPosition{
						Coordinate: entity.NewCoordinate(float64(n.PositionX), spacing(yPos, i)+float64(i*LINE_BETWEEN_LYRIC)),
						Lyrics:     l,
					}

					// some bug happening that for somereason it has middle sylable
					// without other middle or start ever present
					// need to check the lyric existence manually
					if hasLyricBefore {
						pos[i] = pair
						continue
					} else {
						empty := &LyricPosition{
							Coordinate: entity.NewCoordinate(float64(constant.LAYOUT_INDENT_LENGTH), spacing(yPos, i)+float64(i*LINE_BETWEEN_LYRIC)),
						}

						hypenLocation = append(li.CalculateHypen(ctx, empty, pair[0]), hypenLocation...)

					}
				}
				if pair[1] == nil {
					pair[1] = &LyricPosition{
						Coordinate: entity.NewCoordinate(float64(n.PositionX), spacing(yPos, i)+float64(i*LINE_BETWEEN_LYRIC)),
						Lyrics:     l,
					}
					hypenLocation = append(li.CalculateHypen(ctx, pair[0], pair[1]), hypenLocation...)
					pair = [2]*LyricPosition{
						pair[1],
						nil,
					}
				}
			}
			pos[i] = pair
			baseYPos = yPos
		}
		hasLyricBefore = hasLyricBefore || len(n.Lyric) > 0

	}

	if len(pos) > 0 { // add unpaired syllable before move on to next staff
		for i, p := range pos {
			if p[0] != nil && p[1] == nil && i < len(lastLyric) { // append to end of file
				lastXHypen := float64(constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH)

				pEnd := LyricPosition{
					Coordinate: entity.NewCoordinate(lastXHypen, baseYPos+float64(i*LINE_BETWEEN_LYRIC)),
					Lyrics:     lastLyric[i],
				}
				hypenLocation = append(li.CalculateHypen(ctx, p[0], &pEnd), hypenLocation...)

			}
		}
	}
	canv.Group("hyphens")
	for _, hl := range hypenLocation {
		canv.TextUnescaped(hl.X, hl.Y, "-") // use the Unescaped because of the floating number
	}
	canv.Gend()
}
