package wire

import (
	"encoding/binary"
	"io"
)

func readMessageRaw(r io.Reader) ([]byte, error) {
	lenB := make([]byte, 4)
	_, err := io.ReadFull(r, lenB)
	if err != nil {
		return nil, err
	}

	msgLen := binary.BigEndian.Uint32(lenB) - 4
	msg := make([]byte, msgLen)

	_, err = io.ReadFull(r, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
