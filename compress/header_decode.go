package compress

import "deflate/huffman"

func readCode(lengths []uint8, code uint16, stream *IOBitReader) ([]uint8, error) {
	// Length code
	if code < 16 {
		return append(lengths, uint8(code)), nil
	}

	var (
		value  uint8
		repeat uint32
		err    error
	)

	// Repeat codes
	switch code {
	case 16:
		if len(lengths) == 0 {
			return lengths, ErrInvalidHeader
		}

		value = lengths[len(lengths)-1]

		repeat, err = stream.ReadBits(2)
		if err != nil {
			return lengths, err
		}

		repeat += 3
	case 17:
		repeat, err = stream.ReadBits(3)
		if err != nil {
			return lengths, err
		}

		repeat += 3
	case 18:
		repeat, err = stream.ReadBits(7)
		if err != nil {
			return lengths, err
		}

		repeat += 11
	default:
		return lengths, ErrInvalidHeader
	}

	for ; repeat > 0; repeat-- {
		lengths = append(lengths, value)
	}

	return lengths, nil
}

func decodeDynamicHeader(stream *IOBitReader) (huffman.DecodingTree, huffman.DecodingTree, error) {
	// Read tree sizes
	letterNumber, err := stream.ReadBits(5)
	if err != nil {
		return nil, nil, err
	}
	letterNumber += 257

	distanceNumber, err := stream.ReadBits(5)
	if err != nil {
		return nil, nil, err
	}
	distanceNumber += 1

	codeNumber, err := stream.ReadBits(4)
	if err != nil {
		return nil, nil, err
	}
	codeNumber += 4

	// Read code tree
	codeLengths := make([]uint8, len(codeAlphabetDecode))
	for index := range codeNumber {
		length, err := stream.ReadBits(3)
		if err != nil {
			return nil, nil, err
		}

		codeLengths[codeAlphabetDecode[index]] = uint8(length)
	}
	codeTree, err := huffman.NewDecodingTree(codeLengths)
	if err != nil {
		return nil, nil, err
	}

	// Read lengths for letter and distance trees
	lengths := make([]uint8, 0, letterNumber+distanceNumber)
	for len(lengths) < int(letterNumber+distanceNumber) {
		code, err := huffman.Decode(stream, codeTree)
		if err != nil {
			return nil, nil, err
		}

		lengths, err = readCode(lengths, code, stream)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(lengths) > int(letterNumber+distanceNumber) {
		return nil, nil, ErrInvalidHeader
	}

	// Create letter and distance trees
	letterTree, err := huffman.NewDecodingTree(lengths[:letterNumber])
	if err != nil {
		return nil, nil, err
	}

	distanceTree, err := huffman.NewDecodingTree(lengths[letterNumber:])
	if err != nil {
		return nil, nil, err
	}

	return letterTree, distanceTree, nil
}

func decodeHeader(stream *IOBitReader) (bool, int, error) {
	isLast, err := stream.ReadBits(1)
	if err != nil {
		return false, 0, err
	}

	blockType, err := stream.ReadBits(2)
	if err != nil {
		return false, 0, err
	}

	return isLast == 1, int(blockType), nil
}
