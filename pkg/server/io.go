package server

import (
	"encoding/binary"
	"io"

	"github.com/angles-n-daemons/popsql/pkg/server/message"
)

func readMessage(r io.Reader) (message.Type, message.Buffer, error) {
	tB := make([]byte, 1)
	_, err := io.ReadFull(r, tB)
	if err != nil {
		return 0, nil, err
	}

	t := message.Type(tB[0])
	msg, err := readMessageRaw(r)
	return t, msg, err
}

func readMessageRaw(r io.Reader) (message.Buffer, error) {
	lenB := make([]byte, 4)
	_, err := io.ReadFull(r, lenB)
	if err != nil {
		return nil, err
	}

	msgLen := binary.BigEndian.Uint32(lenB) - 4
	data := make(message.Buffer, msgLen)

	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeMessage(w io.Writer, m message.Dumpable) error {
	data := message.Buffer{}
	t := m.Type()
	body := m.Dump()
	msgLen := len(body) + 4

	data.AddType(t)
	data.AddInt32(msgLen)
	data.AddBytes(body)
	_, err := w.Write(data)
	return err
}
