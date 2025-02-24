package huffman

type decodingNode struct {
	letter uint16
	isLeaf bool
	left   *decodingNode // Higher frequency, gets 0
	right  *decodingNode // Lower frequency, gets 1
}

type DecodingTree *decodingNode

func insertCode(tree DecodingTree, code uint32, length uint8, letter uint16) bool {
	// Offset is a cursor that points to the bit in the
	// code and moves from left to right
	for offset := uint32(1 << (length - 1)); offset > 0; offset >>= 1 {
		if tree.isLeaf {
			return false
		}
		// Go into respective subtree and create node
		// if there is no subtree
		if code&offset > 0 {
			if tree.right == nil {
				tree.right = &decodingNode{}
			}

			tree = tree.right
			continue
		}

		if tree.left == nil {
			tree.left = &decodingNode{}
		}

		tree = tree.left
	}

	tree.letter = letter
	tree.isLeaf = true
	return true
}

// Generate decoding tree using Shannon-Fano algorithm.
func NewDecodingTree(lengths []uint8) (DecodingTree, error) {
	codes := generateCodes(lengths)
	if codes == nil {
		return nil, ErrFailedToBuildTree
	}

	tree := &decodingNode{}
	for letter, code := range codes {
		if lengths[letter] == 0 {
			continue
		}

		if !insertCode(tree, code, lengths[letter], uint16(letter)) {
			return nil, ErrFailedToBuildTree
		}
	}

	return tree, nil
}

// Read bits from stream and try to decode symbol using previously generated decoding tree.
func Decode(stream BitReader, tree DecodingTree) (uint16, error) {
	// Traverse the code tree
	for tree != nil && !tree.isLeaf {
		bit, err := stream.ReadBit()
		if err != nil {
			return 0, err
		}

		if bit == 1 {
			tree = tree.right
		} else {
			tree = tree.left
		}
	}

	if tree == nil {
		return 0, ErrInvalidCode
	}

	return tree.letter, nil
}
