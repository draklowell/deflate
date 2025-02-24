package huffman

import "errors"

var (
	ErrFailedToBuildTree = errors.New("failed to build tree")
	ErrInvalidCode       = errors.New("invalid code")
	ErrInvalidSymbol     = errors.New("invalid symbol")
)
