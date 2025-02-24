package compress

import (
	"deflate/huffman"
	"deflate/lz77"
	"io"
)

func decodeValue(stream *IOBitReader, index uint16, offsets []uint16, extraBits []uint8) (uint16, error) {
	if int(index) > len(offsets) {
		return 0, ErrInvalidLetter
	}

	offset, err := stream.ReadBits(extraBits[index])
	if err != nil {
		return 0, err
	}

	return offsets[index] + uint16(offset), nil
}

func decodeBlockBody(
	streamIn *IOBitReader, streamOut io.Writer,
	letterTree, distanceTree huffman.DecodingTree,
	buffer *lz77.RingBuffer,
) error {
	for {
		letter, err := huffman.Decode(streamIn, letterTree)
		if err != nil {
			return err
		}

		if letter == letterEndOfBlock {
			return nil
		}

		if letter < letterEndOfBlock {
			if err := lz77.Decode(byte(letter), streamOut, buffer); err != nil {
				return err
			}

			continue
		}

		length, err := decodeValue(streamIn, letter-letterEndOfBlock-1, lengthOffsets, lengthExtraBits)
		if err != nil {
			return err
		}

		distance, err := huffman.Decode(streamIn, distanceTree)
		if err != nil {
			return err
		}

		distance, err = decodeValue(streamIn, distance, distanceOffsets, distanceExtraBits)
		if err != nil {
			return err
		}

		if err := lz77.DecodeBackreference(length, distance, streamOut, buffer); err != nil {
			return err
		}
	}
}

func decodeBlock(streamIn *IOBitReader, streamOut io.Writer, buffer *lz77.RingBuffer) (bool, error) {
	isLast, blockType, err := decodeHeader(streamIn)
	if err != nil {
		return false, err
	}

	switch blockType {
	case blockTypeUncompressed:
		return false, ErrUncompressedNotSupported
	case blockTypeStatic:
		return isLast, decodeBlockBody(
			streamIn, streamOut, staticLetterDecodingTree,
			staticDistanceDecodingTree, buffer,
		)
	case blockTypeDynamic:
		letterTree, distanceTree, err := decodeDynamicHeader(streamIn)
		if err != nil {
			return false, err
		}

		return isLast, decodeBlockBody(
			streamIn, streamOut, letterTree, distanceTree, buffer,
		)
	}

	return false, ErrInvalidHeader
}
