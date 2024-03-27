package repository

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
)

type Repository interface {
	GetHymnMetaData(ctx context.Context, hymnNum int) (*HymnMetadata, error)
	GetMusicXML(ctx context.Context, filepath string) (musicxml.MusicXML, error)
}

type repository struct {
	conf *config.Config
	db   *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB, conf *config.Config) Repository {
	return &repository{
		conf: conf,
		db:   db,
	}
}

func (r *repository) getHymnMetaDataFromFile(ctx context.Context, hymnNum int) (*HymnMetadata, error) {

	path := fmt.Sprintf("%s%s-%03d.hymexml", r.conf.Metadata.Path, r.conf.Metadata.FilePrefix, hymnNum)
	xmlFile, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrHymnMetadataNotFound
		}
		return nil, fmt.Errorf("[getHymnMetaDataFromFile] Failed to open the file. Err: %s", err.Error())
	}
	defer xmlFile.Close()
	content, _ := io.ReadAll(xmlFile)

	metadaXML := KJMetadata{}
	err = xml.Unmarshal(content, &metadaXML)
	if err != nil {
		return nil, fmt.Errorf("[getHymnMetaDataFromFile] Failed to unmarshal the file. Err: %s", err.Error())
	}

	result := &HymnMetadata{
		HymnData: HymnData{
			HymnID: hymnNum,
			Number: hymnNum,
			Variant: sql.NullInt16{
				Int16: int16(metadaXML.NumberVariant),
				Valid: metadaXML.NumberVariant != 0,
			},
			Title: metadaXML.Title.Value,
			TitleFootnotes: sql.NullString{
				String: metadaXML.Title.Footnote,
				Valid:  metadaXML.Title.Footnote != "",
			},
			Lyric: metadaXML.Credit.Lyric,
			Music: metadaXML.Credit.Music,
		},
	}

	if metadaXML.KidsStarred {
		result.IsForKids = sql.NullInt16{
			Int16: 1,
			Valid: true,
		}

	}

	if metadaXML.CrossReference != nil {
		// TODO: should we make multi reference?
		if len(metadaXML.CrossReference.Be) > 0 {
			result.RefBE = sql.NullInt16{
				Int16: int16(metadaXML.CrossReference.Be[0].Number),
				Valid: true,
			}
		}

		if len(metadaXML.CrossReference.Nr) > 0 {
			result.RefNR = sql.NullInt16{
				Int16: int16(metadaXML.CrossReference.Nr[0].Number),
				Valid: true,
			}
		}
	}

	return result, nil
}

func (r *repository) GetHymnMetaData(ctx context.Context, hymnNum int) (*HymnMetadata, error) {

	if r.conf.Metadata.UseXMLFile {
		return r.getHymnMetaDataFromFile(ctx, hymnNum)
	}
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

func (r *repository) transformVerse(ctx context.Context, verse []VerseType) ([]HymnVerse, error) {

	return nil, nil
}
