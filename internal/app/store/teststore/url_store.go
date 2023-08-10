package teststore

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store"
)

type UrlStore struct {
	repository *UrlRepository
}

func NewUrlStore() *UrlStore {
	return &UrlStore{}
}

func (s *UrlStore) Url() store.UrlRepository {
	if s.repository != nil {
		return s.repository
	}
	s.repository = &UrlRepository{
		store: s,
		urls:  make(map[string]*model.Url),
	}
	return s.repository
}
