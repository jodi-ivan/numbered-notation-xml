package repository

import (
	"context"
	"encoding/xml"
	"io"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type Repository interface {
	GetHymnMetaData(ctx context.Context, hymnNum int) (*HymnMetadata, error)
	GetMusicXML(ctx context.Context, filepath string) (musicxml.MusicXML, error)
}

type repository struct {
	db *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
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
