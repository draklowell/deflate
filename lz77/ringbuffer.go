package lz77

// Ring buffer implementation to be used in LZ77 algorithm.
type RingBuffer struct {
	buffer   []byte
	writeIdx int
	readIdx  int
}

func (rb *RingBuffer) Capacity() int {
	return len(rb.buffer) - 1
}

func (rb *RingBuffer) Len() int {
	// Evaluate length of data in buffer
	// length = writeIdx - readIdx (mod size)
	return (rb.writeIdx - rb.readIdx + len(rb.buffer)) % len(rb.buffer)
}

func (rb *RingBuffer) Push(char byte) {
	// Generate index of next empty cell
	// that is next to write pointer
	// index = writeIdx + 1 (mod size)
	index := (rb.writeIdx + 1) % len(rb.buffer)

	// If next cell is at the read pointer then
	// discard that cell and move read pointer right
	if index == rb.readIdx {
		rb.readIdx = (rb.readIdx + 1) % len(rb.buffer)
	}

	// Write char and move pointer
	rb.buffer[rb.writeIdx] = char
	rb.writeIdx = index
}

func (rb *RingBuffer) PushMany(chars []byte) {
	for _, char := range chars {
		rb.Push(char)
	}
}

func (rb *RingBuffer) Pop() (byte, bool) {
	// Check whether buffer is empty
	if rb.readIdx == rb.writeIdx {
		return 0, false
	}

	char := rb.buffer[rb.readIdx]
	// Move read pointer right
	// readIdx = readIdx + 1 (mod size)
	rb.readIdx = (rb.readIdx + 1) % len(rb.buffer)
	return char, true
}

func (rb *RingBuffer) PopMany(chars []byte) {
	for i := 0; i < len(chars); i++ {
		chars[i], _ = rb.Pop()
	}
}

func (rb *RingBuffer) Read(index int) (byte, bool) {
	// Boundaries check
	if index >= rb.Len() || index < -rb.Len() {
		return 0, false
	}

	if index >= 0 {
		index += rb.readIdx // From the beginning
	} else {
		index += rb.writeIdx // From the end
	}

	// Apply index = index (mod size)
	index = (index + len(rb.buffer)) % len(rb.buffer)

	return rb.buffer[index], true
}

func NewRingBuffer(size uint16) *RingBuffer {
	return &RingBuffer{
		buffer:   make([]byte, size+1),
		readIdx:  0,
		writeIdx: 0,
	}
}
