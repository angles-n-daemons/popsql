package wire

import (
	"encoding/binary"
	"io"

	"github.com/angles-n-daemons/popsql/pkg/api/message"
)

func readMessage(r io.Reader) (message.Type, []byte, error) {
	tB := make([]byte, 1)
	_, err := io.ReadFull(r, tB)
	if err != nil {
		return 0, nil, err
	}

	t := message.Type(tB[0])
	msg, err := readMessageRaw(r)
	return t, msg, err
}

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

func writeMessage(w io.Writer, m message.Dumpable) error {
	t := m.Type()
	b := m.Dump()
	msgLen := len(b) + 4
	lenB := make([]byte, 4)
	binary.BigEndian.PutUint32(lenB, uint32(msgLen))

	data := append([]byte{byte(t)}, append(lenB, b...)...)
	_, err := w.Write(data)
	return err
}
