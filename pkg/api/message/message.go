package message

import (
	"encoding/binary"

	"github.com/angles-n-daemons/popsql/pkg/sql/execution"
)

type Type byte

const (
	M_Query            = 'Q'
	M_No               = 'N'
	M_AuthenticationOk = 'R'
	M_ReadyForQuery    = 'Z'
	M_Terminate        = 'X'
	M_RowDescription   = 'T'
)

type Oid uint32

const (
	T_bool   Oid = 16
	T_text   Oid = 25
	T_float8 Oid = 701
)

// FRONTEND MESSAGES

type Parseable[P any] interface {
	Load(b Buffer) (P, error)
}

func Parse[P Parseable[P]](b Buffer) (P, error) {
	var p P
	return p.Load(b)
}

type Startup struct {
	Version uint32
	Data    map[string]string
}

func (s Startup) Load(b Buffer) (Startup, error) {
	var o Startup
	o.Version = b.ReadUint32()
	o.Data = b.ReadObject()
	return o, nil
}

// BACKEND MESSAGES

type Dumpable interface {
	Type() Type
	Dump() []byte
}

type AuthenticationOk struct{}

func (a *AuthenticationOk) Type() Type {
	return M_AuthenticationOk
}

func (a *AuthenticationOk) Dump() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, 0)
	return b
}

type ReadyForQuery struct{}

func (r *ReadyForQuery) Type() Type {
	return M_ReadyForQuery
}

func (r *ReadyForQuery) Dump() []byte {
	return []byte{'I'}
}

type RowDescription struct {
	Columns   []string
	SampleRow execution.Row
}

func (r *RowDescription) Type() Type {
	return M_RowDescription
}

func (r *RowDescription) Dump() []byte {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, uint16(len(r.Columns)))

	for i, col := range r.Columns {
		result = append(result, []byte(col)...)
		// add null terminator
		result = append(result, 0)

		// skip table id
		result = append(result, 0, 0, 0, 0)
		// skip column offset
		result = append(result, 0, 0)

		var dt Oid
		switch r.SampleRow[i].(type) {
		case float64:
			dt = T_float8
		case string:
			dt = T_text
		case bool:
			dt = T_bool
		}

		dtBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(dtBytes, uint32(dt))
		result = append(result, dtBytes...)
	}
	return nil
}

type ColumnDescription struct {
}
