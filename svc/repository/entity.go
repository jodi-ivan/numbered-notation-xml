package repository

import (
	"database/sql"
	"errors"
)

var ErrHymnNotFound = errors.New("hymn not found")

type HymnDB struct {
	HymnData
	HymnVerse
}

type HymnMetadata struct {
	HymnData
	Verse []HymnVerse
}

type HymnData struct {
	HymnID         int            `db:"hymn_id"`
	Number         int            `db:"hymn_number"`
	Variant        sql.NullInt16  `db:"hymn_variant"`
	Title          string         `db:"title"`
	Footnotes      sql.NullString `db:"footnotes"`
	TitleFootnotes sql.NullString `db:"footnotes_title"`
	Lyric          string         `db:"lyric"`
	Music          string         `db:"music"`
	RefNR          sql.NullInt16  `db:"nr_number"`
	RefBE          sql.NullInt16  `db:"be_number"`
	Copyright      sql.NullString `db:"copyright"`
	IsForKids      sql.NullInt16  `db:"kids_starred"`
}

type HymnVerse struct {
	VerseID  sql.NullInt32  `db:"verse_id"`
	Number   sql.NullInt32  `db:"hymn_num"`
	VerseNum sql.NullInt32  `db:"verse_num"`
	StyleRow sql.NullInt32  `db:"style_row"`
	Content  sql.NullString `db:"content"`
}
