package usecase

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/renderer"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
)

type Usecase interface {
	RenderHymn(ctx context.Context, canv canvas.Canvas, hymnNum int, variant ...string) error
}

type interactor struct {
	config   config.Config
	repo     repository.Repository
	renderer renderer.Renderer
}

func New(config config.Config, repo repository.Repository, renderer renderer.Renderer) Usecase {
	return &interactor{
		config:   config,
		repo:     repo,
		renderer: renderer,
	}
}

// Helper
func getSyllablesBefore(start int, measures []musicxml.Measure) int {
	count := 0
	for i := 0; i < start; i++ {
		for _, a := range measures[i].Appendix {
			if n, err := a.ParseAsNote(); err == nil {
				count += len(n.Lyric)
			}
		}
	}
	return count
}

type RepeatInfo struct {
	Type           musicxml.RepeatInfoType `json:"type"`
	SyllableCount  int                     `json:"syllableCount"`
	OffsetSyllable int                     `json:"offsetSyllable"` // NEW
}

func ProcessRepeats(music *musicxml.MusicXML) {
	bi := barline.NewBarline()
	measures := music.Part.Measures

	// First pass: calculate cumulative syllables from the beginning of the song
	cumulativeSyllables := make([]int, len(measures)+1) // cumulativeSyllables[i] = syllables before measure i

	for i := 0; i < len(measures); i++ {
		count := 0
		for _, a := range measures[i].Appendix {
			if n, err := a.ParseAsNote(); err == nil {
				count += len(n.Lyric)
			}
		}
		cumulativeSyllables[i+1] = cumulativeSyllables[i] + count
	}

	// Now detect repeat sections and assign info
	i := 0
	for i < len(measures) {
		measure := &measures[i]

		lb, _ := bi.GetRendererLeftBarline(*measure, 0, nil)
		_, rb := bi.GetRendererRightBarline(*measure, 0)

		isOpening := lb != nil && lb.Barline != nil &&
			lb.Barline.Repeat != nil &&
			lb.Barline.Repeat.Direction == musicxml.BarLineRepeatDirectionForward

		isClosing := rb != nil && rb.Barline != nil &&
			rb.Barline.Repeat != nil &&
			rb.Barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward

		// Default: no repeat
		if !isOpening && !isClosing {
			// Still set basic info
			measure.RepeatInfo = &musicxml.RepeatInfo{
				Type:                 musicxml.RepeatInfoTypeMiddle, // or Middle if you prefer
				OffsetSyllable:       0,
				SectionSyllableCount: 0,
				StartIndex:           cumulativeSyllables[i],
			}
			i++
			continue
		}

		// Found a new repeat section starting here
		if isOpening || i == 0 {
			startMeasure := i
			startIndex := cumulativeSyllables[i] // absolute start in flattenSyllable

			// Find the closing repeat
			repeatSyllableCount := 0
			j := i
			for j < len(measures) {
				m := &measures[j]

				// Count syllables in this measure
				measSyl := 0
				for _, a := range m.Appendix {
					if n, err := a.ParseAsNote(); err == nil {
						measSyl += len(n.Lyric)
					}
				}
				repeatSyllableCount += measSyl

				_, rbj := bi.GetRendererRightBarline(*m, 0)
				isEnd := rbj != nil && rbj.Barline != nil &&
					rbj.Barline.Repeat != nil &&
					rbj.Barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward

				if isEnd {
					// Found the end of this repeat section
					sectionLength := repeatSyllableCount

					// Now assign to all measures in this section
					for k := startMeasure; k <= j; k++ {
						mk := &measures[k]
						offsetInSection := cumulativeSyllables[k] - cumulativeSyllables[startMeasure]

						mk.RepeatInfo = &musicxml.RepeatInfo{
							Type:                 getRepeatType(k, startMeasure, j),
							OffsetSyllable:       offsetInSection,
							SectionSyllableCount: sectionLength,
							StartIndex:           startIndex,
						}
					}
					i = j + 1
					break
				}
				j++
			}

			if j == len(measures) {
				// No closing found - treat as open-ended (rare)
				i++
			}
			continue
		}

		i++
	}
}

func getRepeatType(idx, start, end int) musicxml.RepeatInfoType {
	if idx == start {
		return musicxml.RepeatInfoTypeOpening
	}
	if idx == end {
		return musicxml.RepeatInfoTypeClosing
	}
	return musicxml.RepeatInfoTypeMiddle
}

func (i *interactor) RenderHymn(ctx context.Context, canv canvas.Canvas, hymnNum int, variant ...string) error {
	filepath := fmt.Sprintf("%s%s-%03d.musicxml", i.config.MusicXML.Path, i.config.MusicXML.FilePrefix, hymnNum)
	if len(variant) > 0 {
		filepath = fmt.Sprintf("%s%s-%03d%s.musicxml", i.config.MusicXML.Path, i.config.MusicXML.FilePrefix, hymnNum, variant[0])
	}
	music, err := i.repo.GetMusicXML(ctx, filepath)
	if err != nil {
		flow := canv.Delegator().OnError(err)
		if flow != canvas.DelegatorErrorFlowControlIgnore {
			return err
		}
	}

	ProcessRepeats(&music)

	// bi := barline.NewBarline()
	// repeatInfo := map[int]musicxml.RepeatInfo{}
	// hasRepeat := false
	// lastMeasureRepeat := 0
	// currentStartIndex := 0
	// // --- Pass 1: detect repeats and syllable counts (backward) ---
	// for i := len(music.Part.Measures) - 1; i >= 0; i-- {
	// 	measure := music.Part.Measures[i]
	// 	_, rb := bi.GetRendererRightBarline(measure, 0)

	// 	if rb.Barline != nil && rb.Barline.Repeat != nil && rb.Barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
	// 		hasRepeat = true
	// 		repeatInfo[measure.Number] = musicxml.RepeatInfo{
	// 			Type: musicxml.RepeatInfoTypeClosing,
	// 		}
	// 		lastMeasureRepeat = measure.Number
	// 	}
	// 	measureSyllableCount := 0
	// 	if hasRepeat {
	// 		info := repeatInfo[lastMeasureRepeat]
	// 		for _, a := range measure.Appendix {
	// 			n, err := a.ParseAsNote()
	// 			if err != nil {
	// 				continue
	// 			}
	// 			if len(n.Lyric) > 0 {
	// 				info.SyllableCount++
	// 				measureSyllableCount++
	// 			}

	// 			log.Println("current start index", currentStartIndex)
	// 		}
	// 		info.StartIndex = currentStartIndex
	// 		repeatInfo[lastMeasureRepeat] = info
	// 	}

	// 	lb, _ := bi.GetRendererLeftBarline(measure, 0, nil)
	// 	openingRepeat := lb != nil && lb.Barline != nil && lb.Barline.Repeat != nil && lb.Barline.Repeat.Direction == musicxml.BarLineRepeatDirectionForward

	// 	if openingRepeat || measure.Number == 1 {
	// 		info := repeatInfo[lastMeasureRepeat]

	// 		// Handle the rare case where opening == closing (single measure repeat)
	// 		if lb != nil && lastMeasureRepeat == lb.MeasureNumber {
	// 			measure.RepeatInfo = &musicxml.RepeatInfo{
	// 				Type:          musicxml.RepeatInfoTypeBoth,
	// 				SyllableCount: info.SyllableCount,
	// 			}
	// 		} else {
	// 			measure.RepeatInfo = &musicxml.RepeatInfo{
	// 				Type:          musicxml.RepeatInfoTypeOpening,
	// 				SyllableCount: info.SyllableCount,
	// 			}
	// 		}
	// 		music.Part.Measures[i] = measure

	// 		// Mark middle and closing measures
	// 		for idx := i + 1; idx < len(music.Part.Measures); idx++ {
	// 			m := music.Part.Measures[idx]
	// 			repeatType := musicxml.RepeatInfoTypeMiddle
	// 			if m.Number == lastMeasureRepeat {
	// 				repeatType = musicxml.RepeatInfoTypeClosing
	// 			}
	// 			m.RepeatInfo = &musicxml.RepeatInfo{
	// 				Type:          repeatType,
	// 				SyllableCount: info.SyllableCount,
	// 			}
	// 			music.Part.Measures[idx] = m
	// 			if m.Number == lastMeasureRepeat {
	// 				break
	// 			}
	// 		}

	// 		repeatInfo = map[int]musicxml.RepeatInfo{}
	// 		lastMeasureRepeat = 0
	// 		hasRepeat = false // ← important reset for sequential repeats
	// 	}

	// 	currentStartIndex += measureSyllableCount
	// }

	// // --- Pass 2: assign OffsetSyllable (forward) ---
	// cumulativeSyllables := 0
	// // track which repeat block we last finished, to avoid double-counting
	// lastRepeatClosingMeasure := -1

	// for i := 0; i < len(music.Part.Measures); i++ {
	// 	m := music.Part.Measures[i]

	// 	if m.RepeatInfo != nil {
	// 		m.RepeatInfo.OffsetSyllable = cumulativeSyllables
	// 		music.Part.Measures[i] = m

	// 		// Accumulate syllables measure by measure
	// 		measureSyllables := 0
	// 		for _, a := range m.Appendix {
	// 			n, err := a.ParseAsNote()
	// 			if err != nil {
	// 				continue
	// 			}
	// 			if len(n.Lyric) > 0 {
	// 				measureSyllables++
	// 			}
	// 		}
	// 		cumulativeSyllables += measureSyllables

	// 		// After the closing barline of a repeat block, the next block
	// 		// continues from where this one ended (no double-count)
	// 		_ = lastRepeatClosingMeasure // used for future nested repeat support
	// 		if m.RepeatInfo.Type == musicxml.RepeatInfoTypeClosing || m.RepeatInfo.Type == musicxml.RepeatInfoTypeBoth {
	// 			lastRepeatClosingMeasure = m.Number
	// 		}

	// 	} else {
	// 		// Non-repeat measure: still accumulate
	// 		for _, a := range m.Appendix {
	// 			n, err := a.ParseAsNote()
	// 			if err != nil {
	// 				continue
	// 			}
	// 			if len(n.Lyric) > 0 {
	// 				cumulativeSyllables++
	// 			}
	// 		}
	// 	}
	// }
	metaData, err := i.repo.GetHymnMetaData(ctx, hymnNum, variant...)
	if err != nil {
		flow := canv.Delegator().OnError(err)
		if flow != canvas.DelegatorErrorFlowControlIgnore {
			return err
		}
	}

	canv.Delegator().OnBeforeStartWrite()

	i.renderer.Render(ctx, music, canv, metaData)

	return nil
}
