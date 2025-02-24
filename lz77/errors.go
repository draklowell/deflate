package lz77

import "errors"

var (
	ErrInvalidBackreference = errors.New("backreference with invalid distance occured")
)
