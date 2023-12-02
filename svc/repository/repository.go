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
	GetHymnMetaData(ctx context.Context, hymnNum int) (*HymnData, error)
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

func (r *repository) GetHymnMetaData(ctx context.Context, hymnNum int) (*HymnData, error) {

	query := sqlx.Rebind(sqlx.QUESTION, qryHymnData)
	row := r.db.QueryRowx(query, hymnNum)
	result := &HymnData{}

	err := row.StructScan(result)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrHymnNotFound
		}

		return nil, err
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
