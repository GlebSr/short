package store

import "github.com/GlebSr/nuhaiShort/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	FindByID(int) (*model.User, error)
}

type UrlRepository interface {
	FindByID(string) (*model.Url, error)
	Create(url *model.Url) error
	Delete(string) error
	RegisterRedirect(*model.Url) error
	FindByUserId(int) ([]model.Url, error)
}
