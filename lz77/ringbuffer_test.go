package lz77

import (
	"testing"
)

func NewTestRingBuffer(size uint16, chars []byte) *RingBuffer {
	buffer := NewRingBuffer(size)

	for _, char := range chars {
		buffer.Push(char)
	}

	return buffer
}

func TestRingBufferCapacityLength(t *testing.T) {
	rb := NewRingBuffer(5)

	if capacity := rb.Capacity(); capacity != 5 {
		t.Errorf("Expected capacity 5, got %d", capacity)
	}

	if length := rb.Len(); length != 0 {
		t.Errorf("Expected length 0, got %d", length)
	}

	rb.Push('A')

	if length := rb.Len(); length != 1 {
		t.Errorf("Expected length 1, got %d", length)
	}
}

func TestRingBufferPushPop(t *testing.T) {
	rb := NewTestRingBuffer(3, []byte("ABC"))

	if rb.Len() != 3 {
		t.Errorf("Expected length 3, got %d", rb.Len())
	}

	char, ok := rb.Pop()
	if !ok || char != 'A' {
		t.Errorf("Expected 'A', got '%c'", char)
	}

	char, ok = rb.Pop()
	if !ok || char != 'B' {
		t.Errorf("Expected 'B', got '%c'", char)
	}

	char, ok = rb.Pop()
	if !ok || char != 'C' {
		t.Errorf("Expected 'C', got '%c'", char)
	}

	_, ok = rb.Pop()
	if ok {
		t.Errorf("Expected buffer to be empty, but Pop returned a value")
	}
}

func TestRingBufferPushPopMany(t *testing.T) {
	rb := NewTestRingBuffer(10, []byte("ABC"))
	data := []byte("DEF")

	rb.PushMany(data)
	if rb.Len() != 6 {
		t.Errorf("Expected buffer length 6, got %d", rb.Len())
	}

	data = make([]byte, 4)
	rb.PopMany(data)
	if string(data) != "ABCD" {
		t.Errorf("Expected popped data 'ABCD', got '%s'", string(data))
	}

	if rb.Len() != 2 {
		t.Errorf("Expected buffer length 3 after pop, got %d", rb.Len())
	}

	rb.PushMany([]byte("EFGH"))
	if rb.Len() != 6 {
		t.Errorf("Expected buffer length 7 after push, got %d", rb.Len())
	}
}

func TestRingBufferOverflow(t *testing.T) {
	rb := NewTestRingBuffer(3, []byte("XYZA"))

	char, _ := rb.Pop()
	if char != 'Y' {
		t.Errorf("Expected 'Y' after overflow, got '%c'", char)
	}

	char, _ = rb.Pop()
	if char != 'Z' {
		t.Errorf("Expected 'Z' after overflow, got '%c'", char)
	}

	char, _ = rb.Pop()
	if char != 'A' {
		t.Errorf("Expected 'A' after overflow, got '%c'", char)
	}

	_, ok := rb.Pop()
	if ok {
		t.Errorf("Expected buffer to be empty, but Pop returned a value")
	}
}

func TestRingBufferRead(t *testing.T) {
	rb := NewTestRingBuffer(6, []byte("MNOP"))

	char, _ := rb.Read(0)
	if char != 'M' {
		t.Errorf("Expected 'M' at index 0, got '%c'", char)
	}

	char, _ = rb.Read(2)
	if char != 'O' {
		t.Errorf("Expected 'O' at index 2, got '%c'", char)
	}

	_, ok := rb.Read(5)
	if ok {
		t.Errorf("Expected to fail, but Read returned a value")
	}

	char, _ = rb.Read(-1)
	if char != 'P' {
		t.Errorf("Expected 'D' at backward index 1, got '%c'", char)
	}

	char, _ = rb.Read(-3)
	if char != 'N' {
		t.Errorf("Expected 'B' at backward index 3, got '%c'", char)
	}

	_, ok = rb.Read(-5)
	if ok {
		t.Errorf("Expected to fail, but ReadBackwards returned a value")
	}
}

func TestRingBufferUnderflow(t *testing.T) {
	rb := NewRingBuffer(2)
	_, ok := rb.Pop()
	if ok {
		t.Errorf("Expected underflow condition, but Pop returned a value")
	}
}
