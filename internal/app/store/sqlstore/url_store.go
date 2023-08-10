package sqlstore

import (
	"database/sql"
	"github.com/GlebSr/nuhaiShort/internal/app/store"
)

type UrlStore struct {
	db         *sql.DB
	repository *UrlRepository
}

func NewUrlStore(db *sql.DB) *UrlStore {
	return &UrlStore{
		db: db,
	}
}

func (s *UrlStore) Url() store.UrlRepository {
	if s.repository != nil {
		return s.repository
	}
	s.repository = &UrlRepository{
		store: s,
	}
	return s.repository
}
