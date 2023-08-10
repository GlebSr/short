package teststore

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store"
)

type UserRepository struct {
	store *UserStore
	users map[int]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if len(u.EncryptedPassword) == 0 {
		err := u.EncryptPassword()
		if err != nil {
			return err
		}
	}
	u.ID = len(r.users) + 1
	r.users[u.ID] = u
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	u, ok := r.users[id]
	if ok {
		return u, nil
	}
	return nil, store.ErrRecordNotFound
}
