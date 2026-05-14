package lyric

import "github.com/jodi-ivan/numbered-notation-xml/internal/constant"

const (
	MAX_VERSE_IN_MUSIC          = 4
	MAX_LINE_PER_VERSE_IN_MUSIC = 2

	LINE_BETWEEN_LYRIC     = 20
	DISTANCE_NOTE_TO_LYRIC = 25

	HYPHEN_LEFT_INDENT = 30 + 10 + constant.LAYOUT_INDENT_LENGTH // clef + padding
)
