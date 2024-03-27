package repository

import (
	"database/sql"
	"encoding/xml"
	"errors"
)

var (
	ErrHymnNotFound         = errors.New("hymn not found")
	ErrHymnMetadataNotFound = errors.New("hymn metadata not found")
)

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
	Variant        sql.NullInt16  `db:"hymn_variant"` // FIXME: does not the variant supposed to be a / b /c? why is it int?
	Title          string         `db:"title"`
	Footnotes      sql.NullString `db:"footnotes"` // TODO: add support for the XML metadata
	TitleFootnotes sql.NullString `db:"footnotes_title"`
	Lyric          string         `db:"lyric"` // for the credit
	Music          string         `db:"music"` // also for the credit
	RefNR          sql.NullInt16  `db:"nr_number"`
	RefBE          sql.NullInt16  `db:"be_number"`
	Copyright      sql.NullString `db:"copyright"` // TODO: add support for the XML metadata
	IsForKids      sql.NullInt16  `db:"kids_starred"`
}

type HymnVerse struct {
	VerseID  sql.NullInt32  `db:"verse_id"`
	Number   sql.NullInt32  `db:"hymn_num"`
	VerseNum sql.NullInt32  `db:"verse_num"`
	StyleRow sql.NullInt32  `db:"style_row"`
	Content  sql.NullString `db:"content"`
}

// XML metadata files
type KJMetadata struct {
	XMLName        xml.Name      `xml:"kj-metadata"`
	Number         int           `xml:"number,attr"`
	NumberVariant  int           `xml:"number-variant,attr,omitempty"`
	KidsStarred    bool          `xml:"kids-starred,attr,omitempty"`
	Title          TitleType     `xml:"title"`
	Credit         CreditType    `xml:"credit"`
	CrossReference *CrossRefType `xml:"cross-reference,omitempty"`
	Category       *CategoryType `xml:"category,omitempty"`
	Verses         VersesType    `xml:"verses"`
	Internal       *InternalType `xml:"internal,omitempty"`
}

type TitleType struct {
	Value    string `xml:"value"`
	Footnote string `xml:"footnote,omitempty"`
}

type CreditType struct {
	Lyric string `xml:"lyric"`
	Music string `xml:"music"`
}

type CrossRefType struct {
	Nr []CrossRefItemType `xml:"nr"`
	Be []CrossRefItemType `xml:"be"`
}

type CrossRefItemType struct {
	Number int `xml:"number,attr"`
}

type CategoryType struct {
	Title    string `xml:"title"`
	ParentID int    `xml:"parent_id,attr,omitempty"`
}

type VersesType struct {
	Verse []VerseType `xml:"verse"`
}

type VerseType struct {
	Line []LineType `xml:"line"`
	Row  int        `xml:"row,attr,omitempty"`
	No   int        `xml:"no,attr,omitempty"`
	X    int        `xml:"x,attr,omitempty"`
	Y    int        `xml:"y,attr,omitempty"`
}

type LineType struct {
	Word []WordType `xml:"word"`
}

type WordType struct {
	Text      string          `xml:"text,attr"`
	Breakdown []BreakdownType `xml:"breakdown"`
}

type BreakdownType struct {
	Text         string          `xml:"text"`
	SubBreakdown []BreakdownType `xml:"breakdown"`
	SyllabicType string          `xml:"syllabic-type,attr,omitempty"`
	Combine      bool            `xml:"combine,attr,omitempty"`
	Underline    bool            `xml:"underline,attr,omitempty"`
}

type InternalType struct {
	Breaklines []BreaklineType `xml:"breaklines>breakline"`
}

type BreaklineType struct {
	Measure   int    `xml:"measure"`
	Note      int    `xml:"note"`
	WordLyric string `xml:"word-lyric"`
}
