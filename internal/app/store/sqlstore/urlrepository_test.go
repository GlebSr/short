package sqlstore_test

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("urls")
	s := sqlstore.NewUrlStore(db)
	assert.NoError(t, s.Url().Create(model.TestUrl(t)))
}

func TestUrlRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("urls")

	s := sqlstore.NewUrlStore(db)
	u1 := model.TestUrl(t)
	s.Url().Create(u1)

	u2, err := s.Url().FindByID(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUrlRepository_Delete(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("urls")

	s := sqlstore.NewUrlStore(db)
	u1 := model.TestUrl(t)
	s.Url().Create(u1)
	err := s.Url().Delete(u1.ID)
	assert.NoError(t, err)
	u2, err := s.Url().FindByID(u1.ID)
	assert.Error(t, err)
	assert.Nil(t, u2)
}

func TestUrlRepository_RegisterRedirect(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("urls")
	s := sqlstore.NewUrlStore(db)
	u1 := model.TestUrl(t)
	s.Url().Create(u1)
	err := s.Url().RegisterRedirect(u1)
	assert.NoError(t, err)
	assert.Equal(t, 1, u1.RedirectsNumber)
	u2 := &model.Url{}
	err = s.Url().RegisterRedirect(u2)
	assert.Error(t, err)
	assert.Equal(t, 0, u2.RedirectsNumber)
}
