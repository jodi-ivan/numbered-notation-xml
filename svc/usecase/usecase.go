package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
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
	bli := barline.NewBarline()

	measureMap := map[int]*musicxml.Measure{}
	syllCountMeasure := map[int][2]int{}

	// lastSyllBefore := 0
	// lastMeasure := 0
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

		bl, _ := bli.GetRendererLeftBarline(measure, constant.LAYOUT_INDENT_LENGTH, nil)
		if bl == nil {
			_, bl = bli.GetRendererRightBarline(measure, constant.LAYOUT_INDENT_LENGTH)
		}

		syllCountMeasure[measure.Number] = [2]int{lastMeasureCount, lastMeasureCount + count - 1}

	}
	for _, repeat := range repeats {
		for start := repeat[0]; start <= repeat[1]; start++ {
			var barlineEnding *musicxml.BarlineEnding

			bl, _ := bli.GetRendererLeftBarline(*measureMap[start], constant.LAYOUT_INDENT_LENGTH, nil)
			if bl == nil {
				_, bl = bli.GetRendererRightBarline(*measureMap[start], constant.LAYOUT_INDENT_LENGTH)
			}

			if bl != nil {
				barlineEnding = bl.Barline.Ending
			}
			repeatType := musicxml.RepeatInfoTypeMiddle
			switch start {
			case repeat[0]:
				repeatType = musicxml.RepeatInfoTypeOpening
			case repeat[1]:
				repeatType = musicxml.RepeatInfoTypeClosing
			}
			idx := start
			if barlineEnding != nil {
				measureOffset, _ := strconv.Atoi(barlineEnding.Number)
				idx -= measureOffset
			}
			measureMap[start].RepeatInfo = &musicxml.RepeatInfo{
				Type:          repeatType,
				SyllCntStart:  syllCountMeasure[idx][0],
				SyllCntEnd:    syllCountMeasure[idx][1],
				OffsetStart:   syllCountMeasure[repeat[1]][1],
				MeasureNumber: start,
				BarlineEnding: barlineEnding,
			}

		}

		if measureMap[repeat[1]].RepeatInfo.BarlineEnding != nil {
			nextMeasure := repeat[1] + 1

			measureMap[nextMeasure].RepeatInfo = &musicxml.RepeatInfo{
				Type:          musicxml.RepeatInfoTypeClosing,
				SyllCntStart:  syllCountMeasure[nextMeasure][0],
				SyllCntEnd:    syllCountMeasure[nextMeasure][1],
				OffsetStart:   syllCountMeasure[nextMeasure][1],
				MeasureNumber: nextMeasure,
				BarlineEnding: &musicxml.BarlineEnding{
					Number: "2",
				},
			}
		}
	}

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
