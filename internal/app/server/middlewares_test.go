package server

import (
	"fmt"
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store/teststore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_AuthenticateUser(t *testing.T) {
	store := teststore.NewUserStore()
	u := model.TestUser(t)
	store.User().Create(u)

	testCases := []struct {
		name         string
		cookieValue  map[any]any
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[any]any{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "no authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}
	s := newServer(store, teststore.NewUrlStore(), sessions.NewCookieStore([]byte("secret")))
	sc := securecookie.New([]byte("secret"), nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.authenticateUser(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
