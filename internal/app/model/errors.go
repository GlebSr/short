package model

import "errors"

var (
	ErrMakingShort = errors.New("error at making url short")
	ErrEncrypting  = errors.New("error at encrypting password")
)
