package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rdnply/url-shortener/internal/counter"
)

var _ counter.Storage = &CounterStorage{}

type CounterStorage struct {
	statementStorage

	incrementStmt *sqlx.Stmt
}

func NewCounterStorage(db *DB) (*CounterStorage, error) {
	s := &CounterStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: incCounterQuery, Dst: &s.incrementStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const initCounterQuery = `INSERT INTO counter(value)
						  SELECT 0 
						  WHERE NOT EXISTS (SELECT * FROM counter)`

func (s *CounterStorage) Init() error {
	if _, err := s.db.Session.Exec(initCounterQuery); err != nil {
		return errors.Wrap(err, "can't init counter")
	}

	return nil
}

const incCounterQuery = `UPDATE counter SET value=value+1 RETURNING value`

func (s *CounterStorage) Increment() (uint, error) {
	var value uint
	if err := s.incrementStmt.Get(&value); err != nil {
		return 0, errors.Wrap(err, "can't exec query")
	}

	return value, nil
}
