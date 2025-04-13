package message

import "encoding/binary"

type Startup struct {
	Version uint32
	Data    map[string]string
}

func (s *Startup) Load(b []byte) error {
	var i int
	s.Version, i = nextUint32(b[:4], i)
	s.Data = make(map[string]string)

	for i < len(b) {
		var key, value string
		key, i = nextString(b, i)
		value, i = nextString(b, i)
		s.Data[key] = value
	}

	return nil
}

func nextString(b []byte, i int) (string, int) {
	end := i
	for b[end] != 0 {
		end++
	}

	return string(b[i:end]), end + 1
}

func nextUint32(b []byte, i int) (uint32, int) {
	return binary.BigEndian.Uint32(b[i : i+4]), i + 4
}
