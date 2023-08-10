package sqlstore

import (
	"database/sql"
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store"
)

type UrlRepository struct {
	store *UrlStore
}

func (r *UrlRepository) Create(url *model.Url) error {
	if len(url.ID) == 0 {
		err := url.MakeShort()
		if err != nil {
			return err
		}
	}
	return r.store.db.QueryRow(
		"INSERT INTO urls (id, long_url, user_id, redirects_number) VALUES ($1, $2, $3, $4) RETURNING id",
		url.ID, url.LongUrl, url.UserId, url.RedirectsNumber,
	).Scan(&url.ID)
}

func (r *UrlRepository) FindByID(id string) (*model.Url, error) {
	u := &model.Url{}
	if err := r.store.db.QueryRow(
		"SELECT id, long_url, user_id, redirects_number FROM urls WHERE id = $1",
		id,
	).Scan(
		&u.ID, &u.LongUrl, &u.UserId, &u.RedirectsNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UrlRepository) Delete(id string) error {
	_, err := r.store.db.Query("DELETE FROM urls WHERE ID = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UrlRepository) RegisterRedirect(url *model.Url) error {
	if err := r.store.db.QueryRow(
		"UPDATE urls SET redirects_number = $1 WHERE id = $2 RETURNING redirects_number",
		url.RedirectsNumber+1, url.ID,
	).Scan(
		&url.RedirectsNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (r *UrlRepository) FindByUserId(id int) ([]model.Url, error) {
	urls := make([]model.Url, 0)
	rows, err := r.store.db.Query(
		"SELECT * FROM urls WHERE user_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for rows.Next() {
		url := model.Url{}
		rows.Scan(&url.ID, &url.LongUrl, &url.RedirectsNumber, &url.UserId)
		urls = append(urls, url)
	}
	return urls, nil
}
