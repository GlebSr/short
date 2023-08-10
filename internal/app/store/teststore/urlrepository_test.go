package teststore_test

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlRepository_Create(t *testing.T) {
	s := teststore.NewUrlStore()
	assert.NoError(t, s.Url().Create(model.TestUrl(t)))
}

func TestUrlRepository_FindByID(t *testing.T) {
	s := teststore.NewUrlStore()
	url1 := model.TestUrl(t)
	s.Url().Create(url1)
	url2, err := s.Url().FindByID(url1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, url2)
}

func TestUrlRepository_Delete(t *testing.T) {
	s := teststore.NewUrlStore()
	url1 := model.TestUrl(t)
	s.Url().Create(url1)
	s.Url().Delete(url1.ID)
	url2, err := s.Url().FindByID(url1.ID)
	assert.Error(t, err)
	assert.Nil(t, url2)
	err = s.Url().Delete(url1.ID)
	assert.Error(t, err)
}

func TestUrlRepository_RegisterRedirect(t *testing.T) {
	s := teststore.NewUrlStore()
	url1 := model.TestUrl(t)
	s.Url().Create(url1)
	err := s.Url().RegisterRedirect(url1)
	assert.NoError(t, err)
	assert.Equal(t, url1.RedirectsNumber, 1)
}

func TestUrlRepository_FindByUserId(t *testing.T) {
	s := teststore.NewUrlStore()
	url1 := model.TestUrl(t)
	s.Url().Create(url1)
	sl, _ := s.Url().FindByUserId(url1.UserId)
	assert.NotNil(t, sl)
}
