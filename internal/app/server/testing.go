package server

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/GlebSr/nuhaiShort/internal/app/store/teststore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"testing"
)

func TestServerWithCookie(t *testing.T, user *model.User) (*server, string) {
	store := teststore.NewUserStore()
	store.User().Create(user)
	s := newServer(store, teststore.NewUrlStore(), sessions.NewCookieStore([]byte("secret")))
	sc := securecookie.New([]byte("secret"), nil)
	cookieStr, _ := sc.Encode(sessionName, map[any]any{
		"user_id": user.ID,
	})
	return s, cookieStr
}
