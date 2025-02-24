package lz77

type LiteralWriter interface {
	Write(char byte) (err error)
	WriteBackreference(length uint16, distance uint16) (err error)
}
