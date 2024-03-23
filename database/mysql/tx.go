package mysql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Tx struct
type Tx struct {
	*sqlx.Tx
	logMode bool
	context context.Context
}

// Query queries the database and returns an *sql.Rows.
func (db *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {

	Log(db.logMode, query, args...)
	if db.context != nil {
		return db.Tx.QueryContext(db.context, query, args...)
	}

	return db.Tx.Query(query, args...)
}

// QueryRow queries the database and returns an *sqlx.Row.
func (db *Tx) QueryRow(query string, args ...interface{}) *sql.Row {

	Log(db.logMode, query, args...)
	if db.context != nil {
		return db.Tx.QueryRowContext(db.context, query, args...)
	}

	return db.Tx.QueryRow(query, args...)
}

// Queryx queries the databas	e and returns an *sqlx.Rows.
func (db *Tx) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {

	Log(db.logMode, query, args...)
	if db.context != nil {
		return db.Tx.QueryxContext(db.context, query, args...)
	}

	return db.Tx.Queryx(query, args...)
}

// QueryRowx queries the database and returns an *sqlx.Row.
func (db *Tx) QueryRowx(query string, args ...interface{}) *sqlx.Row {

	Log(db.logMode, query, args...)
	if db.context != nil {
		return db.Tx.QueryRowxContext(db.context, query, args...)
	}

	return db.Tx.QueryRowx(query, args...)
}

// Exec using master db
func (db *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {

	Log(db.logMode, query, args...)
	if db.context != nil {
		return db.Tx.ExecContext(db.context, query, args...)
	}

	return db.Tx.Exec(query, args...)
}

// Select using slave db.
func (db *Tx) Select(dest interface{}, query string, args ...interface{}) error {

	Log(db.logMode, query, args...)
	if db.context != nil {
		return db.Tx.SelectContext(db.context, dest, query, args...)
	}

	return db.Tx.Select(dest, query, args...)
}

// Get using slave db.
func (db *Tx) Get(dest interface{}, query string, args ...interface{}) error {

	Log(db.logMode, query, args...)
	if db.context != nil {
		return db.Tx.GetContext(db.context, dest, query, args...)
	}

	return db.Tx.Get(dest, query, args...)
}
