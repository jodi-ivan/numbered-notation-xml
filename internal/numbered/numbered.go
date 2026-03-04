package numbered

import (
	"context"
	"math"
	"slices"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type NoteLength struct {
	IsDotted bool
	Type     musicxml.NoteLength
}

type Numbered interface {
	GetLengthNote(ctx context.Context, ts timesig.TimeSignature, measure int, noteLength float64) []NoteLength
	RenderOctave(ctx context.Context, canv canvas.Canvas, octave int, pos entity.Coordinate)
	SplitNote(ctx context.Context, noteLength float64, ts timesig.Time, flag, next musicxml.NoteLength) []NoteLength
	RenderStrikethrough(ctx context.Context, canv canvas.Canvas, strikethrough bool, pos entity.Coordinate)
}

type numberedInteractor struct{}

func (ni *numberedInteractor) SplitNote(ctx context.Context, noteLength float64, ts timesig.Time, flag, next musicxml.NoteLength) []NoteLength {
	base := map[musicxml.NoteLength]float64{
		musicxml.NoteLengthWhole:   4,
		musicxml.NoteLengthHalf:    2,
		musicxml.NoteLengthQuarter: 1,
		musicxml.NoteLengthEighth:  0.5,
		musicxml.NoteLength16th:    0.25,
	}
	flags := map[float64]musicxml.NoteLength{
		4:    musicxml.NoteLengthWhole,
		2:    musicxml.NoteLengthHalf,
		1:    musicxml.NoteLengthQuarter,
		0.5:  musicxml.NoteLengthEighth,
		0.25: musicxml.NoteLength16th,
	}

	unit := base[flag]
	if base[flag] > float64(ts.BeatType)/4 {
		unit = base[flags[float64(ts.BeatType)/4]]
	}

	results := []NoteLength{}

	fullNotes := math.Floor(noteLength / unit)
	remaining := noteLength - (fullNotes * unit)

	// currently tailored for kj-047 & kj-093
	shouldMerge := unit*2 == base[next]

	if fullNotes >= 2 && remaining == 0 && shouldMerge {
		// Add the full notes to the slice
		// kj-45 has 3 half beat on the quater type time sig
		// since this is for readabilty,
		// the the previous note before the ties start is half
		// and after the last ties is quater.
		// hence the remaining follows the quater.
		results = []NoteLength{
			NoteLength{
				Type:     flag,
				IsDotted: true,
			},
		}
		return append(results, slices.Repeat([]NoteLength{
			NoteLength{
				Type:     flags[unit*2],
				IsDotted: true,
			},
		}, int(fullNotes/2))...)
	}

	// kj-007, kj-075 and kj-093
	// show the quater notes as double 8th notes on ties.
	// since it is follow the readabilty:
	// the notes after the last ties is 8th notes (next) for uniformity
	for i := 0; i < int(fullNotes); i++ {
		results = append(results, NoteLength{
			Type:     flag,
			IsDotted: true,
		})
	}

	// kj-064, display the remaining as the dot with proper flag.
	// 2.5 notes in quater time sig would be:
	// [note quater + dot quater + dot eitght]
	dotValue := unit * 0.5

	if _, ok := flags[dotValue]; ok && remaining >= dotValue {
		results = append(results, NoteLength{
			Type:     flags[dotValue],
			IsDotted: true,
		})

	}

	return results
}

func (ni *numberedInteractor) GetLengthNote(ctx context.Context, ts timesig.TimeSignature, measure int, noteLength float64) []NoteLength {
	currentTimeSig := ts.GetTimesignatureOnMeasure(ctx, measure)

	// 1. Setup Base and Subdivision Types
	var baseType, halfType, quarterType musicxml.NoteLength
	switch currentTimeSig.BeatType {
	case 8:
		baseType, halfType, quarterType = musicxml.NoteLengthEighth, musicxml.NoteLength16th, musicxml.NoteLength32nd
	case 4, 2:
		baseType, halfType, quarterType = musicxml.NoteLengthQuarter, musicxml.NoteLengthEighth, musicxml.NoteLength16th
	default:
		return nil
	}

	// 2. Handle Standalone Special Cases (Less than 1 beat)
	if noteLength < 1.0 {
		switch noteLength {
		case 0.75:
			return []NoteLength{{Type: halfType}, {IsDotted: true, Type: quarterType}}
		case 0.5:
			return []NoteLength{{Type: halfType}}
		case 0.25:
			return []NoteLength{{Type: quarterType}}
		default:
			return []NoteLength{{Type: baseType}} // Fallback
		}
	}

	// 3. Handle Full Beats (First note is head, others are dots)
	result := []NoteLength{{Type: baseType}}
	fullBeats := int(math.Trunc(noteLength))
	for i := 1; i < fullBeats; i++ {
		result = append(result, NoteLength{IsDotted: true, Type: baseType})
	}

	// 4. Handle Decimal Extensions (The "Tail")
	if tail := noteLength - float64(fullBeats); tail > 0 {
		var extensionType musicxml.NoteLength

		// Based on your original logic:
		// In BT 8, a decimal always appended a 16th.
		// In BT 4, a decimal always appended an 8th.
		if currentTimeSig.BeatType == 8 {
			extensionType = musicxml.NoteLength16th
		} else {
			extensionType = musicxml.NoteLengthEighth
		}

		result = append(result, NoteLength{IsDotted: true, Type: extensionType})
	}

	return result
}

func New() Numbered {
	return &numberedInteractor{}
}
