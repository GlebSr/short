package teststore

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store"
)

type UserStore struct {
	repository *UserRepository
}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (s *UserStore) User() store.UserRepository {
	if s.repository != nil {
		return s.repository
	}
	s.repository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}
	return s.repository
}
