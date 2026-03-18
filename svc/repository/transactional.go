package repository

import "database/sql"

type Transactional interface {
	Commit() error
	Rollback() error
}

type sqlTx struct {
	tx *sql.Tx
}

func (stx *sqlTx) Commit() error {
	return stx.tx.Commit()
}
func (stx *sqlTx) Rollback() error {
	return stx.tx.Rollback()
}

// FIXME: move everything that use this to its own repository
func GetSqlTx(tx Transactional) *sql.Tx {
	return (tx.(*sqlTx)).tx
}
