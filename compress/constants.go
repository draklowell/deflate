package compress

import (
	"errors"
)

const (
	letterEndOfBlock = 256

	// See https://datatracker.ietf.org/doc/html/rfc1951#page-9 (3.2.3)
	blockTypeUncompressed = 0b00
	blockTypeStatic       = 0b01
	blockTypeDynamic      = 0b10

	codeOffsetRepeat  = 16
	codeOffsetZRepeat = codeOffsetRepeat + 4

	bufferDistanceOffset = 32767

	maxBufferInSize  = 258
	maxBufferOutSize = 32768
)

var (
	// See https://datatracker.ietf.org/doc/html/rfc1951#page-13 (3.2.7)
	codeAlphabetDecode = []uint8{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}
	// See https://datatracker.ietf.org/doc/html/rfc1951#page-13 (3.2.7)
	codeAlphabetEncode = []uint8{3, 17, 15, 13, 11, 9, 7, 5, 4, 6, 8, 10, 12, 14, 16, 18, 0, 1, 2}
)

var (
	ErrInvalidHeader            = errors.New("invalid block header")
	ErrInvalidLetter            = errors.New("invalid letter occured")
	ErrInvalidBackreference     = errors.New("invalid backreference")
	ErrUncompressedNotSupported = errors.New("uncompressed blocks are not supported")
)
