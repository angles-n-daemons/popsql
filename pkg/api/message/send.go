package message

import (
	"encoding/binary"

	"github.com/angles-n-daemons/popsql/pkg/sql/execution"
)

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
