package sqlstore

import (
	"database/sql"
	"github.com/GlebSr/nuhaiShort/internal/app/store"
	_ "github.com/lib/pq"
)

type UserStore struct {
	db         *sql.DB
	repository *UserRepository
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) User() store.UserRepository {
	if s.repository != nil {
		return s.repository
	}
	s.repository = &UserRepository{
		store: s,
	}
	return s.repository
}
