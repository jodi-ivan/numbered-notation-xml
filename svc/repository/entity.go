package repository

import (
	"database/sql"
	"errors"
)

var ErrHymnNotFound = errors.New("hymn not found")

type HymnData struct {
	ID             int            `db:"id"`
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
