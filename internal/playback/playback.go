package playback

import (
	"time"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
)

func CalculateDuration(steps *[]Step, tempo int, totalBeat map[int][]float64, rect map[int][][2]entity.Coordinate) time.Duration {
	if tempo == 0 {
		tempo = 85
	}
	totalOverallBeat := 0.0
	for i := 0; i < len(*steps); i++ {
		dStep := *steps
		currStep := dStep[i]

		totalOverallBeat += totalBeat[currStep.MeasureNumber][0]
		if len(totalBeat[currStep.MeasureNumber]) > 1 {
			for i := 1; i < len(totalBeat[currStep.MeasureNumber]); i++ {
				totalOverallBeat += totalBeat[currStep.MeasureNumber][i]
			}
		}
	}

	overallDuration := (totalOverallBeat / float64(tempo)) * 60.0
	durationEachBeat := (overallDuration / totalOverallBeat) * float64(time.Second)
	totalSteps := len(*steps)
	for i := 0; i < totalSteps; i++ {
		dStep := *steps
		currStep := dStep[i]
		if len(rect[currStep.MeasureNumber]) == 1 {

			currStep.TotalBeat = totalBeat[currStep.MeasureNumber][0]
			currStep.Duration = time.Duration(durationEachBeat * currStep.TotalBeat)
			currStep.DurationStr = currStep.Duration.String()
			currStep.Rect = rect[currStep.MeasureNumber][0]

			dStep[i] = currStep
			*steps = dStep
			continue
		} else {

			firstStep := currStep
			firstStep.TotalBeat = totalBeat[currStep.MeasureNumber][0]
			firstStep.Duration = time.Duration(durationEachBeat * firstStep.TotalBeat)
			firstStep.DurationStr = firstStep.Duration.String()
			firstStep.Rect = rect[currStep.MeasureNumber][0]

			newSteps := []Step{}

			for i := 1; i < len(totalBeat[currStep.MeasureNumber]); i++ {
				dur := time.Duration(durationEachBeat * totalBeat[currStep.MeasureNumber][i])
				newSteps = append(newSteps, Step{
					TotalBeat:     totalBeat[currStep.MeasureNumber][i],
					Duration:      dur,
					DurationStr:   dur.String(),
					Rect:          rect[currStep.MeasureNumber][i],
					MeasureNumber: currStep.MeasureNumber,
					LyricPart:     firstStep.LyricPart,
				})
			}

			dStep[i] = firstStep
			dStep = append(dStep[:i+1], append(newSteps, dStep[i+1:]...)...)
			i += len(newSteps)
			totalSteps += len(newSteps)
			*steps = dStep
		}

	}

	return time.Duration(overallDuration * float64(time.Second))
}
