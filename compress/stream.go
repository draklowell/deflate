package compress

import (
	"io"
)

type IOBitReader struct {
	stream    io.Reader
	data      []byte
	dataIndex int
}

func NewIOBitReader(stream io.Reader) *IOBitReader {
	return &IOBitReader{
		stream:    stream,
		data:      make([]byte, 1),
		dataIndex: 8,
	}
}

func (br *IOBitReader) ReadBit() (uint8, error) {
	if br.dataIndex == 8 {
		br.dataIndex = 0
		n, err := br.stream.Read(br.data)
		if n == 0 {
			return 0, io.EOF
		}
		if err != nil {
			return 0, err
		}
	}

	bit := (br.data[0] >> br.dataIndex) & 1
	br.dataIndex++

	return bit, nil
}

func (br *IOBitReader) ReadBits(n uint8) (uint32, error) {
	var result uint32
	for i := uint8(0); i < n; i++ {
		bit, err := br.ReadBit()
		if err != nil {
			return 0, err
		}

		result |= uint32(bit) << i
	}

	return result, nil
}

type IOBitWriter struct {
	stream    io.Writer
	data      []byte
	dataIndex int
}

func NewIOBitWriter(stream io.Writer) *IOBitWriter {
	return &IOBitWriter{
		stream:    stream,
		data:      make([]byte, 1),
		dataIndex: 0,
	}
}

func (bw *IOBitWriter) flush() error {
	if bw.dataIndex == 0 {
		return nil
	}

	_, err := bw.stream.Write(bw.data)
	bw.data[0] = 0
	bw.dataIndex = 0
	return err
}

func (bw *IOBitWriter) WriteBit(bit uint8) error {
	if bw.dataIndex == 8 {
		if err := bw.flush(); err != nil {
			return err
		}
	}

	bw.data[0] |= byte((bit & 1) << bw.dataIndex)
	bw.dataIndex++
	return nil
}

func (bw *IOBitWriter) WriteBits(bits uint32, n uint8) error {
	for i := uint8(0); i < n; i++ {
		bit := uint8(bits >> i)
		if err := bw.WriteBit(bit); err != nil {
			return err
		}
	}

	return nil
}
