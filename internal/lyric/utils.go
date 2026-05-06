package lyric

import (
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
)

// Only for Caladea font
var width = map[string]float64{
	"A": 9.59, "a": 7.52, "1": 9.28,
	"B": 9.27, "b": 8.32, "2": 7.55,
	"C": 8.1, "c": 6.74, "3": 7.43,
	"D": 10, "d": 8.32, "4": 8.57,
	"E": 8.65, "e": 7.06, "5": 7.61,
	"F": 8.15, "f": 5.87, "6": 7.53,
	"G": 8.63, "g": 7.35, "7": 7.53,
	"H": 11.15, "h": 8.86, "8": 8,
	"I": 5.49, "i": 4.44, "9": 7.65,
	"J": 4.99, "j": 4.76, "0": 8.57,
	"K": 10.08, "k": 8.43, ",": 3.28,
	"L": 8.02, "l": 4.34, "'": 3.28,
	"M": 14.21, "m": 13.01, ".": 3.3,
	"N": 11.09, "n": 8.94, "!": 4.58,
	"O": 9.59, "o": 7.69, ";": 4.23,
	"P": 8.53, "p": 8.32, " ": 4,
	"Q": 9.59, "q": 8.02, "-": 5.27,
	"R": 9.81, "r": 6.34, "—": 16,
	"S": 7.25, "s": 6.28, "*": 6.83,
	"T": 8.92, "t": 5.21, "\"": 8,
	"U": 11, "u": 8.74, "+": 8.86,
	"V": 9.57, "v": 8.08, ":": 4.23,
	"W": 14.23, "w": 12.08,
	"X": 9.95, "x": 7.78,
	"Y": 8.92, "y": 8.18,
	"Z": 8.11, "z": 6.85,
}

func (li *lyricInteractor) CalculateOverallWidth(ls []entity.Lyric) float64 {
	result := 0.0

	for _, v := range ls {
		result = math.Max(result, li.CalculateLyricWidth(entity.LyricVal(v.Text).String()))
	}
	return result
}

func (li *lyricInteractor) CalculateLyricWidth(txt string) float64 {

	res := 0.0

	for _, l := range txt {
		res += width[string(l)]
	}

	return res
}
