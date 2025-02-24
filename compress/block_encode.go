package compress

import (
	"deflate/huffman"
)

type encoderLiteralWriter struct {
	staticBlockThreshold int
	stream               *IOBitWriter

	buffer     []uint16
	bufferSize int

	letterFrequencies   [286]uint64
	distanceFrequencies [30]uint64
}

func findIndexByValue(value uint16, offsets []uint16) uint16 {
	valueIndex := len(offsets) - 1
	for index, offset := range offsets {
		if offset > value {
			valueIndex = index - 1
			break
		}
	}

	return uint16(valueIndex)
}

func (lw *encoderLiteralWriter) Write(char byte) error {
	if len(lw.buffer) >= lw.bufferSize {
		if err := lw.flush(false); err != nil {
			return err
		}
	}

	lw.letterFrequencies[char] += 1
	lw.buffer = append(lw.buffer, uint16(char))

	return nil
}

func (lw *encoderLiteralWriter) WriteBackreference(length uint16, distance uint16) error {
	if len(lw.buffer) >= lw.bufferSize {
		if err := lw.flush(false); err != nil {
			return err
		}
	}

	if length > lengthMax {
		return ErrInvalidBackreference
	}
	lengthIndex := findIndexByValue(length, lengthOffsets)

	if distance > distanceMax {
		return ErrInvalidBackreference
	}
	distanceIndex := findIndexByValue(distance, distanceOffsets)

	lw.letterFrequencies[letterEndOfBlock+lengthIndex+1] += 1
	lw.buffer = append(lw.buffer, letterEndOfBlock+length)

	lw.distanceFrequencies[distanceIndex] += 1
	lw.buffer = append(lw.buffer, bufferDistanceOffset+distance)

	return nil
}

func (lw *encoderLiteralWriter) generateHeader(isLast bool) (huffman.EncodingTree, huffman.EncodingTree, error) {
	if len(lw.buffer) <= lw.staticBlockThreshold {
		if err := encodeHeader(isLast, blockTypeStatic, lw.stream); err != nil {
			return nil, nil, err
		}

		return staticLetterEncodingTree, staticDistanceEncodingTree, nil
	}

	if err := encodeHeader(isLast, blockTypeDynamic, lw.stream); err != nil {
		return nil, nil, err
	}
	letterLengths := huffman.GenerateLengths(lw.letterFrequencies[:], 15)
	distanceLengths := huffman.GenerateLengths(lw.distanceFrequencies[:], 15)

	letterTree, err := huffman.NewEncodingTree(letterLengths)
	if err != nil {
		return nil, nil, err
	}

	distanceTree, err := huffman.NewEncodingTree(distanceLengths)
	if err != nil {
		return nil, nil, err
	}

	return letterTree, distanceTree, encodeDynamicHeader(letterLengths, distanceLengths, lw.stream)

}

func (lw *encoderLiteralWriter) flush(isLast bool) error {
	lw.letterFrequencies[letterEndOfBlock] = 1

	// Build header
	letterTree, distanceTree, err := lw.generateHeader(isLast)
	if err != nil {
		return err
	}

	// Write block
	for _, letter := range lw.buffer {
		if letter < letterEndOfBlock {
			if err := huffman.Encode(letter, lw.stream, letterTree); err != nil {
				return err
			}

			continue
		}

		if letter >= bufferDistanceOffset {
			distance := letter - bufferDistanceOffset
			distanceIndex := findIndexByValue(distance, distanceOffsets)
			if err := huffman.Encode(distanceIndex, lw.stream, distanceTree); err != nil {
				return err
			}

			lw.stream.WriteBits(uint32(distance-distanceOffsets[distanceIndex]), distanceExtraBits[distanceIndex])

			continue
		}

		length := letter - letterEndOfBlock
		lengthIndex := findIndexByValue(length, lengthOffsets)
		if err := huffman.Encode(lengthIndex+letterEndOfBlock+1, lw.stream, letterTree); err != nil {
			return err
		}

		lw.stream.WriteBits(uint32(length-lengthOffsets[lengthIndex]), lengthExtraBits[lengthIndex])
	}

	// Reset buffer slice
	lw.buffer = lw.buffer[:0]

	return huffman.Encode(letterEndOfBlock, lw.stream, letterTree)
}
