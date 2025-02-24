package compress

import "deflate/huffman"

func writeCode(codes []uint8, code uint8, repeat uint8, codeFrequencies []uint64) []uint8 {
	// Repeat code 18
	if repeat >= 11 && code == 0 {
		codeFrequencies[18] += 1
		return append(codes, codeOffsetZRepeat+repeat-3)
	}

	// Repeat code 17
	if repeat >= 3 && code == 0 {
		codeFrequencies[17] += 1
		return append(codes, codeOffsetZRepeat+repeat-3)
	}

	// Write heading length
	codeFrequencies[code] += 1
	codes = append(codes, code)
	repeat--

	// Repeat code 16
	if repeat >= 3 {
		codeFrequencies[16] += 1
		return append(codes, codeOffsetRepeat+repeat-3)
	}

	// If no repeat code available, write remaining lengths
	for ; repeat > 0; repeat-- {
		codeFrequencies[code] += 1
		codes = append(codes, code)
	}

	return codes
}

func compressLengths(lengths []uint8) ([]uint8, []uint64) {
	codes := make([]uint8, 0, len(lengths))
	codeFrequencies := make([]uint64, len(codeAlphabetEncode))

	for index := 0; index < len(lengths); {
		code := lengths[index]
		repeat := 0
		// Search how many times does this length repeat
		for ; index+repeat < len(lengths); repeat++ {
			if lengths[index+repeat] != code {
				break
			}

			if repeat == 6 && code != 0 {
				break
			}

			if repeat == 137 && code == 0 {
				break
			}
		}

		// Write code (repeat+1) times
		codes = writeCode(codes, code, uint8(repeat), codeFrequencies[:])
		index += repeat
	}

	return codes, codeFrequencies
}

func findLastNonZero(array []uint8) int {
	current := 0
	for index, value := range array {
		if value != 0 {
			current = index
		}
	}

	return current
}

func encodeDynamicHeader(letterLengths, distanceLengths []uint8, stream *IOBitWriter) error {
	letterLengthsLast := max(findLastNonZero(letterLengths)+1, 257)
	distanceLengthLast := max(findLastNonZero(distanceLengths)+1, 1)

	// Compress trees
	lengths := make([]uint8, letterLengthsLast+distanceLengthLast)

	copy(lengths, letterLengths[:letterLengthsLast])
	copy(lengths[letterLengthsLast:], distanceLengths[:distanceLengthLast])

	codes, codeFrequencies := compressLengths(lengths)

	// Generate code tree
	codeLengths := huffman.GenerateLengths(codeFrequencies, 7)
	codeTree, err := huffman.NewEncodingTree(codeLengths)
	if err != nil {
		return err
	}

	// Reorder code lengths
	codesLengthsOrdered := make([]uint8, len(codeAlphabetEncode))
	for code, length := range codeLengths {
		codesLengthsOrdered[codeAlphabetEncode[code]] = length
	}
	codesLengthsLast := max(findLastNonZero(codesLengthsOrdered)+1, 4)

	// Write tree sizes
	if err := stream.WriteBits(uint32(letterLengthsLast-257), 5); err != nil {
		return err
	}
	if err := stream.WriteBits(uint32(distanceLengthLast-1), 5); err != nil {
		return err
	}
	if err := stream.WriteBits(uint32(codesLengthsLast-4), 4); err != nil {
		return err
	}

	// Write code tree
	for _, length := range codesLengthsOrdered[:codesLengthsLast] {
		if err := stream.WriteBits(uint32(length), 3); err != nil {
			return err
		}
	}

	// Write lengths for letter and distance trees
	for _, code := range codes {
		extraData := uint32(0)
		extraBits := uint8(0)
		if code < codeOffsetRepeat {
		} else if code < codeOffsetZRepeat {
			extraData = uint32(code - codeOffsetRepeat)
			extraBits = 2
			code = 16
		} else {
			extraData = uint32(code - codeOffsetZRepeat)

			if extraData < 8 {
				code = 17
				extraBits = 3
			} else {
				extraData -= 8
				code = 18
				extraBits = 7
			}
		}

		if err := huffman.Encode(uint16(code), stream, codeTree); err != nil {
			return err
		}

		if err := stream.WriteBits(extraData, extraBits); err != nil {
			return err
		}
	}

	return nil
}

func encodeHeader(isLast bool, blockType int, stream *IOBitWriter) error {
	var isLastBit uint32 = 0
	if isLast {
		isLastBit = 1
	}

	if err := stream.WriteBits(isLastBit, 1); err != nil {
		return err
	}

	return stream.WriteBits(uint32(blockType), 2)
}
