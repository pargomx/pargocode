package sqlitedb

import (
	"context"
	"database/sql"
)

type Transaccion struct {
	tx  *sql.Tx
	log bool
}

// ================================================================ //
// ================================================================ //

func (s *SqliteDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.log {
		logSQL(query, args...)
	}
	return s.db.QueryContext(context.Background(), query, args...)
}

func (s *SqliteDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if s.log {
		logSQL(query, args...)
	}
	return s.db.QueryRow(query, args...)
}

func (s *SqliteDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if s.log {
		logSQL(query, args...)
	}
	return s.db.Exec(query, args...)
}

// ================================================================ //
// ================================================================ //

func (s *SqliteDB) Begin() (*Transaccion, error) {
	if s.log {
		logSQL("BEGIN TRANSACTION")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaccion{
		tx:  tx,
		log: s.log,
	}, nil
}

func (s *Transaccion) Commit() error {
	if s.log {
		logSQL("COMMIT")
	}
	return s.tx.Commit()
}

func (s *Transaccion) Rollback() error {
	if s.log {
		logSQL("ROLLBACK")
	}
	return s.tx.Rollback()
}

func (s *Transaccion) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.log {
		logSQL(query, args...)
	}
	return s.tx.QueryContext(context.Background(), query, args...)
}

func (s *Transaccion) QueryRow(query string, args ...interface{}) *sql.Row {
	if s.log {
		logSQL(query, args...)
	}
	return s.tx.QueryRow(query, args...)
}

func (s *Transaccion) Exec(query string, args ...interface{}) (sql.Result, error) {
	if s.log {
		logSQL(query, args...)
	}
	return s.tx.Exec(query, args...)
}
