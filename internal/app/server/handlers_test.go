package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store/teststore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	s := newServer(teststore.NewUserStore(), teststore.NewUrlStore(), sessions.NewCookieStore([]byte("secret")))
	testCases := []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "valid@valid.valid",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email":    "valid@valid",
				"password": "pas",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "empty",
			expectedCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest("POST", "/v1/registration", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleSessionsCreate(t *testing.T) {
	u := model.TestUser(t)
	store := teststore.NewUserStore()
	store.User().Create(u)
	testCases := []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid@valid.ru",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Password,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "empty",
			expectedCode: http.StatusUnauthorized,
		},
	}
	s := newServer(store, teststore.NewUrlStore(), sessions.NewCookieStore([]byte("secret")))
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest("POST", "/v1/login", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUrlRedirection(t *testing.T) {
	store := teststore.NewUrlStore()
	url := model.TestUrl(t)
	store.Url().Create(url)
	s := newServer(teststore.NewUserStore(), store, sessions.NewCookieStore([]byte("secret")))
	{
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", url.ID), nil)
		s.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, url.LongUrl, rec.Header().Get("Location"))
	}
	{
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/invalid"), nil)
		s.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.NotEqual(t, url.LongUrl, rec.Header().Get("Location"))
	}

}

func TestServer_HandleNewUrl(t *testing.T) {
	u := model.TestUser(t)
	s, cookieStr := TestServerWithCookie(t, u)
	url := model.TestUrl(t)
	testCases := []struct {
		name         string
		payload      any
		cookie       string
		withCookie   bool
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"long_url": url.LongUrl,
			},
			cookie:       cookieStr,
			withCookie:   true,
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid json",
			payload: map[string]string{
				"url": url.LongUrl,
			},
			cookie:       cookieStr,
			withCookie:   true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid cookie",
			payload: map[string]string{
				"long_url": url.LongUrl,
			},
			cookie:       "cookieStr",
			withCookie:   true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "without cookie",
			payload: map[string]string{
				"long_url": url.LongUrl,
			},
			cookie:       cookieStr,
			withCookie:   false,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid url",
			payload: map[string]string{
				"long_url": "bad.url",
			},
			cookie:       cookieStr,
			withCookie:   true,
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			req, _ := http.NewRequestWithContext(ctx, "POST", "/v1/account/new-url", b)
			if tc.withCookie {
				req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, tc.cookie))
			}
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleNewCoolUrl(t *testing.T) {
	url := model.TestUrl(t)
	url.ID = "test"
	store := teststore.NewUserStore()
	u := model.TestUser(t)
	store.User().Create(u)
	s := newServer(store, teststore.NewUrlStore(), sessions.NewCookieStore([]byte("secret")))
	sc := securecookie.New([]byte("secret"), nil)
	cookieStr, _ := sc.Encode(sessionName, map[any]any{
		"user_id": u.ID,
	})
	rec := httptest.NewRecorder()
	b := &bytes.Buffer{}
	rawData := &struct {
		Key      string `json:"key"`
		LongUrl  string `json:"long_url"`
		ShortUrl string `json:"short_url"`
	}{
		Key:      "key",
		LongUrl:  url.LongUrl,
		ShortUrl: url.ID,
	}
	json.NewEncoder(b).Encode(rawData)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	req, _ := http.NewRequestWithContext(ctx, "POST", "/v1/account/new-special-url", b)
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	rawData.Key = ""
	rawData.ShortUrl = "test2"
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestServer_HandleDeleteUrl(t *testing.T) {
	url := model.TestUrl(t)
	url.ID = "test"
	userStore := teststore.NewUserStore()
	urlStore := teststore.NewUrlStore()
	u := model.TestUser(t)
	userStore.User().Create(u)
	urlStore.Url().Create(url)
	s := newServer(userStore, urlStore, sessions.NewCookieStore([]byte("secret")))
	sc := securecookie.New([]byte("secret"), nil)
	cookieStr, _ := sc.Encode(sessionName, map[any]any{
		"user_id": u.ID,
	})
	rec := httptest.NewRecorder()
	b := &bytes.Buffer{}
	rawData := &struct {
		Url string `json:"short_url"`
	}{
		Url: url.ID,
	}
	json.NewEncoder(b).Encode(rawData)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	req, _ := http.NewRequestWithContext(ctx, "DELETE", "/v1/account/delete-url", b)
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestServer_HandleOneUrlInformation(t *testing.T) {
	url := model.TestUrl(t)
	url.ID = "test"
	userStore := teststore.NewUserStore()
	urlStore := teststore.NewUrlStore()
	u := model.TestUser(t)
	userStore.User().Create(u)
	urlStore.Url().Create(url)
	s := newServer(userStore, urlStore, sessions.NewCookieStore([]byte("secret")))
	sc := securecookie.New([]byte("secret"), nil)
	cookieStr, _ := sc.Encode(sessionName, map[any]any{
		"user_id": u.ID,
	})
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/account/urls-information/%s", url.ID), nil)
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	url2 := &model.Url{}
	json.NewDecoder(rec.Body).Decode(url2)
	assert.Equal(t, url2.LongUrl, url.LongUrl)
	assert.Equal(t, url2.ID, url.ID)
	assert.Equal(t, url2.RedirectsNumber, url.RedirectsNumber)
}

func TestServer_HandleUrlsInformation(t *testing.T) {

}
