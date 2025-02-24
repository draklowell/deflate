package huffman

type encodingNode struct {
	code   uint32
	length uint8
}

type EncodingTree []encodingNode

// Generate encoding tree using Shannon-Fano algorithm.
func NewEncodingTree(lengths []uint8) (EncodingTree, error) {
	codes := generateCodes(lengths)
	if codes == nil {
		return nil, ErrFailedToBuildTree
	}

	tree := make(EncodingTree, len(lengths))
	for symbol, code := range codes {
		tree[symbol] = encodingNode{
			code:   code,
			length: lengths[symbol],
		}
	}

	return tree, nil
}

// Encode symbol using previously generated encoding tree.
func Encode(symbol uint16, stream BitWriter, tree EncodingTree) error {
	if symbol >= uint16(len(tree)) {
		return ErrInvalidSymbol
	}

	node := tree[symbol]
	for i := node.length; i > 0; i-- {
		stream.WriteBit(uint8(node.code >> (i - 1)))
	}
	return nil
}
