package playback

func GenerateSteps(hasFine, lastMeasure int, repeat [][3]int) []Step {

	result := []Step{}

	lastPart := lastMeasure
	if len(repeat) > 0 {
		lastPart = repeat[0][0]
	}

	for i := 1; i <= lastPart && lastPart != 1; i++ {
		result = append(result, Step{MeasureNumber: i, LyricPart: 1})
	}

	for _, r := range repeat {
		start, end, jump := r[0], r[1], r[2]
		for part := 1; part <= 2; part++ {
			for i := start; i <= end; i++ {
				measureNumber := i
				if part == 2 && i == end && jump > 0 {
					measureNumber += jump
				}
				result = append(result, Step{MeasureNumber: measureNumber, LyricPart: part})

			}
		}

		// TODO: fill gap between repeats.
	}

	if len(repeat) > 0 {
		for i := repeat[len(repeat)-1][1] + 1; i <= lastMeasure; i++ {
			result = append(result, Step{MeasureNumber: i, LyricPart: 1})
		}
	}

	if hasFine > 0 {
		jump := len(repeat) > 0 && (repeat[0][1] == hasFine-repeat[0][2])

		for i := 1; i <= hasFine; i++ {
			if jump && i == hasFine-1 {
				continue
			}
			result = append(result, Step{MeasureNumber: i, LyricPart: 1})
		}
	}
	return result
}
