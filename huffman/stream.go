package huffman

type BitReader interface {
	ReadBit() (bit uint8, err error)
}

type BitWriter interface {
	WriteBit(bit uint8) (err error)
}
