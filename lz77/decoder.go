package lz77

import "io"

// Decode (store in the buffer) single char and write it to the stream.
func Decode(char byte, stream io.Writer, buffer *RingBuffer) error {
	buffer.Push(char)
	_, err := stream.Write([]byte{char})
	return err
}

// Decode backreference pointer and write chars to the stream.
func DecodeBackreference(length, distance uint16, stream io.Writer, buffer *RingBuffer) error {
	for ; length > 0; length-- {
		char, ok := buffer.Read(-int(distance))
		if !ok {
			return ErrInvalidBackreference
		}

		buffer.Push(char)
		if _, err := stream.Write([]byte{char}); err != nil {
			return err
		}
	}

	return nil
}
