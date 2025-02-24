package lz77

import "io"

func findMatch(bufferIn, bufferOut *RingBuffer) (uint16, uint16) {
	var maxLength, maxDistance int

	for index := range bufferOut.Len() {
		// If maxLength is already larger than all possible
		// next lengths, then we don't need to iterate
		if maxLength >= bufferOut.Len()-index {
			break
		}

		offsetIn := 0
		offsetOut := 0
		for {
			charIn, ok := bufferIn.Read(offsetIn)
			if !ok {
				break
			}

			charOut, ok := bufferOut.Read(index + offsetOut)
			// If we overflow bufferOut then reset offset
			if !ok {
				offsetOut = 0
				continue
			}

			if charIn != charOut {
				break
			}

			offsetIn++
			offsetOut++
		}

		if offsetIn >= maxLength && offsetIn >= 3 {
			// Since index is from the beginning of the
			// bufferOut, we need to convert it to distance
			// from the end
			maxLength, maxDistance = offsetIn, bufferOut.Len()-index
		}
	}

	return uint16(maxLength), uint16(maxDistance)
}

func Encode(streamIn io.Reader, streamOut LiteralWriter, bufferIn, bufferOut *RingBuffer) error {
	// streamIn -> bufferIn
	chars := make([]uint8, bufferIn.Capacity())
	charsRead, err := streamIn.Read(chars)
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	bufferIn.PushMany(chars[:charsRead])

	for bufferIn.Len() > 0 {
		length, distance := findMatch(bufferIn, bufferOut)

		// bufferIn -> streamOut
		if length == 0 {
			char, _ := bufferIn.Read(0)
			streamOut.Write(char)
			length = 1
		} else {
			streamOut.WriteBackreference(length, distance)
		}

		// bufferIn -> bufferOut
		chars = make([]byte, length)
		bufferIn.PopMany(chars)
		bufferOut.PushMany(chars)

		// streamIn -> bufferIn
		charsRead, err := streamIn.Read(chars)
		if err != nil && err != io.EOF {
			return err
		}
		bufferIn.PushMany(chars[:charsRead])
	}

	return nil
}
