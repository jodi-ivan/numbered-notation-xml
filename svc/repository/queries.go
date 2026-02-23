package repository

const (
	qryHymnData = `
	SELECT 
    a.ID as hymn_id,
    a.hymn_number,
    a.hymn_variant,
    a.title,
    a.footnotes,
    a.footnotes_title,
    a.lyric,
    a.music,
    a.nr_number,
    a.be_number,
    a.copyright,
    a.kids_starred,
    b.ID as verse_id,
    b.verse_num,
    b.style_row,
    b.column_pos,
    b.row_pos,
    b.content,
    c.id as footnote_id,
    c.line_pos,
    c.footnote_marker,
    c.marker_style,
    c.footnote
FROM jdy_hymn a 
LEFT JOIN jdy_hymn_verces b 
    ON a.hymn_number = b.hymn_num 
    AND (a.hymn_variant = b.hymn_variant OR (a.hymn_variant IS NULL AND b.hymn_variant IS NULL))
LEFT JOIN verse_footnotes c
    ON b.hymn_num = c.hymne_num 
    AND (b.hymn_variant = c.hymne_variant OR (b.hymn_variant IS NULL AND c.hymne_variant IS NULL))
    AND b.verse_num = c.verse_num
WHERE a.hymn_number = ?
	`

	qryHymnHasVariant = `
	SELECT 
		a.ID as hymn_id,
		a.hymn_number,
		a.hymn_variant
	FROM 
		jdy_hymn a 
	WHERE 
		hymn_number = ? AND hymn_variant IS NOT NULL
	
	`
)

const (
	qryInsertVerse = `
		INSERT INTO jdy_hymn_verces 
		(
			hymn_num,
			verse_num,
			content,
			style_row,
			column_pos,
			row_pos
		)
		VALUES
		(
			?,
			?,
			?,
			?,
			?,
			?
		)
		RETURNING id
	`
)
