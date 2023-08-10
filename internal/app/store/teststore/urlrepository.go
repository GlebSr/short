package teststore

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store"
)

type UrlRepository struct {
	store *UrlStore
	urls  map[string]*model.Url
}

func (r *UrlRepository) Create(u *model.Url) error {
	if len(u.ID) == 0 {
		err := u.MakeShort()
		if err != nil {
			return err
		}
	}
	_, ok := r.urls[u.ID]
	if ok {
		return store.ErrShortUrlЕaken
	}
	r.urls[u.ID] = u
	return nil
}

func (r *UrlRepository) FindByID(id string) (*model.Url, error) {
	u, ok := r.urls[id]
	if ok {
		return u, nil
	}
	return nil, store.ErrRecordNotFound
}

func (r *UrlRepository) Delete(id string) error {
	_, ok := r.urls[id]
	if !ok {
		return store.ErrRecordNotFound
	}
	delete(r.urls, id)
	return nil
}

func (r *UrlRepository) RegisterRedirect(url *model.Url) error {
	u, ok := r.urls[url.ID]
	if !ok {
		return store.ErrRecordNotFound
	}
	u.RedirectsNumber += 1
	r.urls[url.ID] = u
	return nil
}

func (r *UrlRepository) FindByUserId(id int) ([]model.Url, error) {
	ans := make([]model.Url, 0)
	for _, data := range r.urls {
		if data.UserId == id {
			ans = append(ans, *data)
		}
	}
	return ans, nil
}
