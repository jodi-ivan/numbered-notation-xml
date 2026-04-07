package footnote

type VerseNoteStyle int

const (
	VerseNoteStyleDirectAppendText VerseNoteStyle = 0
	VerseNoteStyleAlignRight       VerseNoteStyle = 1
	VerseNoteStyleHeadless         VerseNoteStyle = 2
	VerseNoteStyleHeadonly         VerseNoteStyle = 3
	VerseNoteStyleForTitle         VerseNoteStyle = 4
)
