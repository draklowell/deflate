package compress

import (
	"deflate/lz77"
	"io"
)

// Compress message from the input stream using deflate and output blocks to the output stream
func Compress(
	streamIn io.Reader,
	streamOut io.Writer,
	blockSize int,
	bufferInSize, bufferOutSize int,
	staticBlockThreshold int,
) error {
	bufferIn := lz77.NewRingBuffer(min(uint16(bufferInSize), maxBufferInSize))
	bufferOut := lz77.NewRingBuffer(min(uint16(bufferOutSize), maxBufferOutSize))
	blockSize = max(blockSize, 2)

	streamMid := &encoderLiteralWriter{
		staticBlockThreshold: staticBlockThreshold,
		stream:               NewIOBitWriter(streamOut),

		buffer:     make([]uint16, 0, blockSize),
		bufferSize: blockSize - 1,
	}

	if err := lz77.Encode(streamIn, streamMid, bufferIn, bufferOut); err != nil {
		return err
	}

	if err := streamMid.flush(true); err != nil {
		return err
	}

	return streamMid.stream.flush()
}

// Decompress deflate blocks from the input stream and output message to the output stream
func Decompress(
	streamIn io.Reader,
	streamOut io.Writer,
) error {
	buffer := lz77.NewRingBuffer(maxBufferOutSize)
	streamMid := NewIOBitReader(streamIn)

	for {
		isLast, err := decodeBlock(streamMid, streamOut, buffer)
		if err != nil {
			return err
		}

		if isLast {
			break
		}
	}

	return nil
}
