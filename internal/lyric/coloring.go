package lyric

var coloringOpacity = map[int]string{
	0: `style="opacity:0.6"`,
}

func getColoringStyle(verse, totalLyric int) string {
	if totalLyric == 1 {
		return ""
	}

	return coloringOpacity[verse]

}
