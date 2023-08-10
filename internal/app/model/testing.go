package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestUrl(t *testing.T) *Url {
	return &Url{
		LongUrl:         "https://ya.ru",
		RedirectsNumber: 0,
		UserId:          1,
	}
}
