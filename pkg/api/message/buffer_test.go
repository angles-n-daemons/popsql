package message

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"
)

type EmbeddedBuffer struct {
	i int
	b []byte
}

func (m *EmbeddedBuffer) AddUint16(val int) {
	arr := make([]byte, 2)
	binary.BigEndian.PutUint16(arr, uint16(val))
	m.b = append(m.b, arr...)
}

func (m *EmbeddedBuffer) ReadUint16() uint16 {
	intB := binary.BigEndian.Uint16(m.b[m.i : m.i+2])
	m.i += 2
	return intB
}

func TestMessageBuffer(t *testing.T) {
	m := Buffer{}
	m.AddUint16(1)
	m.AddUint32(2)
	if len(m) != 6 {
		t.Errorf("Expected length 6, got %d", len(m))
	}
	if binary.BigEndian.Uint16(m[:2]) != 1 {
		t.Errorf("Expected first 2 bytes to be 1, got %d", binary.BigEndian.Uint16(m[:2]))
	}
	if binary.BigEndian.Uint32(m[2:]) != 2 {
		t.Errorf("Expected last 4 bytes to be 2, got %d", binary.BigEndian.Uint32(m[2:]))
	}
}

func BenchmarkInlineBuffer(b *testing.B) {
	for e := range 7 {
		b.Run(fmt.Sprintf("PutInt16 e=%d", e), func(b *testing.B) {
			buf := Buffer{}
			for range b.N {
				for range int(1 * math.Pow(10, float64(e))) {
					buf.AddUint16(1)
				}
				for range int(1 * math.Pow(10, float64(e))) {
					buf.ReadUint16()
				}
			}
		})
	}
}

func BenchmarkEmbeddedBuffer(b *testing.B) {
	for e := range 7 {
		b.Run(fmt.Sprintf("PutInt16 e=%d", e), func(b *testing.B) {
			buf := EmbeddedBuffer{}
			for range b.N {
				for range int(1 * math.Pow(10, float64(e))) {
					buf.AddUint16(1)
				}
				for range int(1 * math.Pow(10, float64(e))) {
					buf.ReadUint16()
				}
			}
		})
	}
}

// Curious, performance is a bit better using the embedded buffer.
// The inline buffer, while a bit slower, is easier for me to use.
// BenchmarkInlineBuffer/PutInt16_e=0-11           195745138                5.935 ns/op
// BenchmarkInlineBuffer/PutInt16_e=1-11           25965454                45.76 ns/op
// BenchmarkInlineBuffer/PutInt16_e=2-11            4980732               241.5 ns/op
// BenchmarkInlineBuffer/PutInt16_e=3-11             462692              2523 ns/op
// BenchmarkInlineBuffer/PutInt16_e=4-11              46969             26716 ns/op
// BenchmarkInlineBuffer/PutInt16_e=5-11               4712            245777 ns/op
// BenchmarkInlineBuffer/PutInt16_e=6-11                468           2558346 ns/op
// BenchmarkEmbeddedBuffer/PutInt16_e=0-11         276754176                4.017 ns/op
// BenchmarkEmbeddedBuffer/PutInt16_e=1-11         61591104                21.41 ns/op
// BenchmarkEmbeddedBuffer/PutInt16_e=2-11          5288402               232.5 ns/op
// BenchmarkEmbeddedBuffer/PutInt16_e=3-11           503234              2314 ns/op
// BenchmarkEmbeddedBuffer/PutInt16_e=4-11            53121             22909 ns/op
// BenchmarkEmbeddedBuffer/PutInt16_e=5-11             4650            230988 ns/op
// BenchmarkEmbeddedBuffer/PutInt16_e=6-11              514           2328108 ns/op
