package usecase

import (
	"context"
	"fmt"

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

func collectRepeat(measures []musicxml.Measure) [][2]int {

	result := [][2]int{}
	// check if there is any repeat at all
	for i := 0; i < len(measures); i++ {
		for _, b := range measures[i].Barline {
			if b.Repeat != nil {
				switch b.Repeat.Direction {
				case musicxml.BarLineRepeatDirectionForward:
					result = append(result, [2]int{measures[i].Number})
				case musicxml.BarLineRepeatDirectionBackward:
					// closing
					if len(result) == 0 {
						result = append(result, [2]int{1}) // beginning of the measure
					}
					lastKnown := result[len(result)-1]
					lastKnown[1] = measures[i].Number
					result[len(result)-1] = lastKnown
				}
			}
		}
	}
	return result
}

func ProcessRepeats(music *musicxml.MusicXML) {

	repeats := collectRepeat(music.Part.Measures)

	if len(repeats) == 0 {
		return
	}
	measureMap := map[int]*musicxml.Measure{}
	syllCountMeasure := map[int][2]int{}

	for i, measure := range music.Part.Measures {
		measureMap[measure.Number] = &music.Part.Measures[i]
		count := 0
		for _, a := range measure.Appendix {
			if n, err := a.ParseAsNote(); err == nil && len(n.Lyric) > 0 {
				count++
			}
		}
		lastMeasureCount := 1
		if len(syllCountMeasure) > 0 {
			lastMeasureCount = syllCountMeasure[measure.Number-1][1] + 1
		}
		syllCountMeasure[measure.Number] = [2]int{lastMeasureCount, lastMeasureCount + count - 1}
	}

	for _, repeat := range repeats {
		for start := repeat[0]; start <= repeat[1]; start++ {
			repeatType := musicxml.RepeatInfoTypeMiddle
			switch start {
			case repeat[0]:
				repeatType = musicxml.RepeatInfoTypeOpening
			case repeat[1]:
				repeatType = musicxml.RepeatInfoTypeClosing
			}
			measureMap[start].RepeatInfo = &musicxml.RepeatInfo{
				Type:         repeatType,
				SyllCntStart: syllCountMeasure[start][0],
				SyllCntEnd:   syllCountMeasure[start][1],
				OffsetStart:  syllCountMeasure[repeat[1]][1],
			}

		}
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
