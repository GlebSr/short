package model_test

import (
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrl_MakeShort(t *testing.T) {
	url := model.TestUrl(t)
	assert.NoError(t, url.MakeShort())
	assert.NotEmpty(t, url.ID)
}

func TestUrl_Validation(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model.Url
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.Url {
				return model.TestUrl(t)
			},
			isValid: true,
		},
		{
			name: "empty",
			u: func() *model.Url {
				u := model.Url{}
				return &u
			},
			isValid: false,
		},
		{
			name: "empty url",
			u: func() *model.Url {
				u := model.TestUrl(t)
				u.LongUrl = ""
				return u
			},
			isValid: false,
		},
		{
			name: "0 UserId",
			u: func() *model.Url {
				u := model.TestUrl(t)
				u.UserId = 0
				return u
			},
			isValid: false,
		},
		{
			name: "invalid url",
			u: func() *model.Url {
				u := model.TestUrl(t)
				u.LongUrl = "google.ru"
				return u
			},
			isValid: false,
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
