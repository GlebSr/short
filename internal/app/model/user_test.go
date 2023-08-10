package model_test

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "empty user",
			u: func() *model.User {
				u := model.User{}
				return &u
			},
			isValid: false,
		},
		{
			name: "empty email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = "test"
				return u
			},
			isValid: false,
		},
		{
			name: "long password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
				return u
			},
			isValid: false,
		},
		{
			name: "Only encrypted password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				u.EncryptedPassword = "password"
				return u
			},
			isValid: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_EncryptPassword(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.EncryptPassword())
	assert.NotEmpty(t, u.EncryptedPassword)
}
