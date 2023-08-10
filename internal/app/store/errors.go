package store

import "errors"

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrCreatingProblem = errors.New("problem at creating url")
	ErrShortUrlЕaken   = errors.New("this short url already taken")
)
