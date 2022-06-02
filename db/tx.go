// region: packages

package db

import (
	"database/sql"

	"github.com/SandorMiskey/TEx-kit/log"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// endregion: packages
// region: types

type Tx struct {
	db      *Db
	history History
	session *sql.Tx
}

// endregion: types
// region: begin

func Begin(db *Db) (*Tx, error) {
	s := &Statement{SQL: "BEGIN"}

	session, e := db.Conn().Begin()
	if e != nil {
		s.Err = e
		db.appendHistory(s)
		log.Out(db.Config().Logger, *db.Config().Loglevel, e)
		return nil, e
	}

	tx := Tx{db: db, session: session}
	tx.history = make(History, 0, *db.Config().History)
	tx.appendHistory(s)
	return &tx, nil
}

func (db *Db) Begin() (*Tx, error) {
	return Begin(db)
}

// endregion: begin
// region: commit

func Commit(tx *Tx) error {
	s := &Statement{SQL: "COMMIT"}
	e := tx.Session().Commit()
	if e != nil {
		s.Err = e
		tx.appendHistory(s)
		log.Out(tx.Db().Config().Logger, *tx.Db().Config().Loglevel, e)
		return e
	}

	tx.appendHistory(s)
	return nil
}

func (tx *Tx) Commit() error {
	return Commit(tx)
}

// endregion: commit
// region: exec

func (tx *Tx) Exec(s *Statement) error {
	return Exec(tx, s)
}

// endregion: exec
// region: getters

func (tx *Tx) Config() *Config {
	return tx.Db().Config()
}

func (tx *Tx) Db() *Db {
	return tx.db
}

func (tx *Tx) exec() interface{} {
	return tx.session
}

func (tx *Tx) History() History {
	return tx.history
}

func (tx *Tx) Session() *sql.Tx {
	return tx.session
}

// endregion: getters
// region: history

func (tx *Tx) appendHistory(s *Statement) {
	appendHistory(tx, s)
	appendHistory(tx.Db(), s)
}

func (tx *Tx) setHistory(h *History) {
	tx.history = *h
}

// endregion: history
// region: rollback

func Rollback(tx *Tx) error {
	s := &Statement{SQL: "ROLLBACK"}
	e := tx.Session().Rollback()
	if e != nil {
		s.Err = e
		tx.appendHistory(s)
		log.Out(tx.Db().Config().Logger, *tx.Db().Config().Loglevel, e)
		return e
	}

	tx.appendHistory(s)
	return nil
}

func (tx *Tx) Rollback() error {
	return Rollback(tx)
}

// endregion: rollback
