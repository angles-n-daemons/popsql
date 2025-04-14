package message

import "encoding/binary"

type Buffer []byte

func (m *Buffer) AddType(t Type) {
	*m = append(*m, []byte{byte(t)}...)
}

func (m *Buffer) AddUint16(val int) {
	arr := make([]byte, 2)
	binary.BigEndian.PutUint16(arr, uint16(val))
	*m = append(*m, arr...)
}

func (m *Buffer) AddUint32(val int) {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(val))
	*m = append(*m, arr...)
}

func (m *Buffer) AddString(s string) {
	*m = append(*m, []byte(s)...)
	*m = append(*m, 0) // null-terminate the string
}

func (m *Buffer) AddObject(obj map[string]string) {
	for k, v := range obj {
		m.AddString(k)
		m.AddString(v)
	}
	*m = append(*m, 0) // null-terminate the object
}

func (m *Buffer) AddBytes(arrs ...[]byte) {
	for _, b := range arrs {
		*m = append(*m, b...)
	}
}

func (m *Buffer) ReadUint16() uint16 {
	intB := binary.BigEndian.Uint16((*m)[:2])
	*m = (*m)[2:]
	return intB
}

func (m *Buffer) ReadUint32() uint32 {
	intB := binary.BigEndian.Uint32((*m)[:4])
	*m = (*m)[4:]
	return intB
}

func (m *Buffer) ReadString() string {
	end := 0
	for _, c := range *m {
		if c == 0 {
			break
		}
		end++
	}
	str := string((*m)[:end])
	*m = (*m)[len(str)+1:]
	return str
}

func (m *Buffer) ReadObject() map[string]string {
	obj := make(map[string]string)
	for {
		key := m.ReadString()
		if key == "" {
			break
		}
		value := m.ReadString()
		obj[key] = value
	}
	return obj
}

func (m *Buffer) PeekUint32() uint32 {
	return binary.BigEndian.Uint32((*m)[:4])
}
