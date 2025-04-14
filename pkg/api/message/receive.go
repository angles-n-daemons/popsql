package message

import "encoding/binary"

const SSLRequestCode = 80877103 // Magic code for SSL request

type Startup struct {
	Version uint32
	Data    map[string]string
}

func (s *Startup) Load(b []byte) error {
	var i int
	s.Version, i = NextUint32(b[:4], i)
	s.Data = make(map[string]string)

	for i < len(b)-1 && b[i] != 0 {
		var key, value string
		key, i = NextString(b, i)
		value, i = NextString(b, i)
		s.Data[key] = value
	}

	return nil
}

func NextString(b []byte, i int) (string, int) {
	end := i
	for b[end] != 0 {
		end++
	}

	return string(b[i:end]), end + 1
}

func NextUint32(b []byte, i int) (uint32, int) {
	return binary.BigEndian.Uint32(b[i : i+4]), i + 4
}
