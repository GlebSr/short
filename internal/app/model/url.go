package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"strings"
)

type Url struct {
	ID              string `json:"short_url"`
	LongUrl         string `json:"long_url"`
	UserId          int    `json:"-"`
	RedirectsNumber int    `json:"redirects_number"`
}

var alphabet = []byte("vuXzSehWigb476NpVP2CxH3ETGcnAmMLDUa8jqJ9Zto1drKYF5QwkfyBRs")

func (u *Url) MakeShort() error {
	id := int(uuid.New().ID())
	short := make([]int, 0, 6)
	for id != 0 {
		short = append(short, id%len(alphabet))
		id /= len(alphabet)
	}
	var url strings.Builder
	for pos, _ := range short {
		err := url.WriteByte(alphabet[short[len(short)-1-pos]])
		if err != nil {
			return ErrMakingShort
		}
	}
	u.ID = url.String()
	return nil
}

func (u *Url) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.LongUrl, validation.Required, is.URL, is.RequestURL),
		validation.Field(&u.UserId, validation.Required, validation.Min(1)),
	)
}
