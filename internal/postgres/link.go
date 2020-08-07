package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rdnply/url-shortener/internal/link"
)

var _ link.Storage = &LinkStorage{}

type LinkStorage struct {
	statementStorage

	addStmt           *sqlx.Stmt
	getByShortIDStmt  *sqlx.Stmt
	incrementCntrStmt *sqlx.Stmt
}

func NewLinkStorage(db *DB) (*LinkStorage, error) {
	s := &LinkStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: addLinkQuery, Dst: &s.addStmt},
		{Query: getLinkByShortIDQuery, Dst: &s.getByShortIDStmt},
		{Query: incrementLinkCntrQuery, Dst: &s.incrementCntrStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const linkFields = `url, short_link_id, short_link_id_int, count_clicks`
const addLinkQuery = `INSERT INTO links(` + linkFields + `)
					  VALUES($1, $2, $3, $4) RETURNING link_id`

func (s *LinkStorage) AddLink(l *link.Link) (uint, error) {
	var id uint
	if err := s.addStmt.QueryRowx(l.URL, l.ShortID, l.ShortIDInt, l.Clicks).Scan(&id); err != nil {
		return 0, errors.Wrap(err, "can't exec query")
	}

	return id, nil
}

const getLinkByShortIDQuery = `SELECT link_id, ` + linkFields + ` FROM links WHERE short_link_id=$1`

func (s *LinkStorage) GetLinkByShortID(shortID string) (*link.Link, error) {
	var link link.Link
	if err := s.getByShortIDStmt.Get(&link, shortID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "can't get link by short id")
	}

	return &link, nil
}

const incrementLinkCntrQuery = `UPDATE links SET count_clicks=count_clicks+1  
								WHERE short_link_id=$1 RETURNING count_clicks`

func (s *LinkStorage) IncrementLinkCounter(l *link.Link) (uint, error) {
	var count uint
	if err := s.incrementCntrStmt.Get(&count, l.ShortID); err != nil {
		return 0, errors.Wrap(err, "can't increment clicks counter")
	}

	return count, nil
}
