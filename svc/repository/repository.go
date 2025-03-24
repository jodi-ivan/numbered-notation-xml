package repository

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type Repository interface {
	GetHymnMetaData(ctx context.Context, hymnNum int) (*HymnMetadata, error)
	GetMusicXML(ctx context.Context, filepath string) (musicxml.MusicXML, error)
	InsertVerse(ctx context.Context, hymn, verse, style, col, row int, content string) (int, error)
}

type repository struct {
	db *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func fillNull(val int) sql.NullInt32 {
	result := sql.NullInt32{}

	if val != 0 {
		result.Valid = true
		result.Int32 = int32(val)
	}
	return result
}

func (r *repository) InsertVerse(ctx context.Context, hymn, verse, style, col, row int, content string) (int, error) {

	styleQL := fillNull(style)
	colQL := fillNull(col)
	rowQL := fillNull(row)

	var newID int

	query := sqlx.Rebind(sqlx.QUESTION, qryInsertVerse)

	err := r.db.QueryRow(query, hymn, verse, content, styleQL, colQL, rowQL).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *repository) GetHymnMetaData(ctx context.Context, hymnNum int) (*HymnMetadata, error) {

	query := sqlx.Rebind(sqlx.QUESTION, qryHymnData)
	rows := []*HymnDB{}
	err := r.db.Select(&rows, query, hymnNum)

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, ErrHymnNotFound
	}

	result := &HymnMetadata{
		HymnData: rows[0].HymnData,
		Verse:    []HymnVerse{},
	}

	for _, verse := range rows {
		if verse.VerseID.Valid {
			result.Verse = append(result.Verse, verse.HymnVerse)
		}
	}

	return result, nil
}

func (r *repository) GetMusicXML(ctx context.Context, filepath string) (musicxml.MusicXML, error) {
	xmlFile, err := os.Open(filepath)
	if err != nil {
		return musicxml.MusicXML{}, err
	}
	defer xmlFile.Close()
	content, _ := io.ReadAll(xmlFile)

	var music musicxml.MusicXML
	err = xml.Unmarshal(content, &music)
	return music, err
}
