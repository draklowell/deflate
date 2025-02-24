package huffman

func generateCodes(lengths []uint8) []uint32 {
	if len(lengths) > 65536 {
		return nil
	}

	maxLength := uint8(0)
	for _, length := range lengths {
		maxLength = max(maxLength, length)
	}

	// Count occurances of each length among the codes
	lengthCounters := make([]uint16, maxLength+1)
	for _, length := range lengths {
		lengthCounters[length]++
	}

	// Code per length is basically base for code
	// to help compute it later
	lengthCounters[0] = 0
	codePerLength := make([]uint32, maxLength+1)
	code := uint32(0)
	for length := range maxLength {
		code += uint32(lengthCounters[length])
		code <<= 1
		codePerLength[length+1] = code
	}

	// Use base for code to evaluate each code
	// and increase base each iteration
	codes := make([]uint32, len(lengths))
	for symbol, length := range lengths {
		if length == 0 {
			continue
		}

		codes[symbol] = codePerLength[length]
		codePerLength[length]++
	}

	return codes
}
