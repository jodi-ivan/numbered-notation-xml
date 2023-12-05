package barline

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

var (
	// only support noto-music font
	barlineWidth = map[musicxml.BarLineStyle]float64{
		musicxml.BarLineStyleRegular:    4.16,
		musicxml.BarLineStyleLightHeavy: 7.7,
		musicxml.BarLineStyleLightLight: 6.28,
		musicxml.BarLineStyleHeavyHeavy: 8,
		musicxml.BarLineStyleHeavyLight: 7.7,
	}

	unicode = map[musicxml.BarLineStyle]string{
		musicxml.BarLineStyleRegular:    `&#x01D100;`,
		musicxml.BarLineStyleLightHeavy: `&#x01D102;`,
		musicxml.BarLineStyleLightLight: `&#x01D101;`,
		musicxml.BarLineStyleHeavyHeavy: `&#x01D101;`,
		musicxml.BarLineStyleHeavyLight: `&#x01D103;`,
	}
)

type BarlineInfo struct {
	XIncrement int
}
