package message

import (
	"fmt"
	"strconv"

	"github.com/angles-n-daemons/popsql/pkg/db/sql/execution"
)

/*
backend.go contains the messages sent by the server to the client.

It specifically contains the structs that represent the messages,
the utility functions for serializing them, and the common `Dumpable`
interface they must all satisfy.
*/

type Dumpable interface {
	Type() Type
	Dump() Buffer
}

/*
AuthenticationOk (B)
Byte1('R')
Identifies the message as an authentication request.

Int32(8)
Length of message contents in bytes, including self.

Int32(0)
Specifies that the authentication was successful.
*/
type AuthenticationOk struct{}

func (a *AuthenticationOk) Type() Type {
	return M_AuthenticationOk
}

func (a *AuthenticationOk) Dump() Buffer {
	data := Buffer{}
	data.AddInt32(0)
	return data
}

/*
ReadyForQuery (B)
Byte1('Z')
Identifies the message type. ReadyForQuery is sent whenever the backend is ready for a new query cycle.

Int32(5)
Length of message contents in bytes, including self.

Byte1
Current backend transaction status indicator. Possible values are 'I' if idle (not in a transaction block); 'T' if in a transaction block; or 'E' if in a failed transaction block (queries will be rejected until block is ended).
*/

type ReadyForQuery struct{}

func (r *ReadyForQuery) Type() Type {
	return M_ReadyForQuery
}

func (r *ReadyForQuery) Dump() Buffer {
	return []byte{'I'}
}

/*
RowDescription (B)
Byte1('T')
Identifies the message as a row description.

Int32
Length of message contents in bytes, including self.

Int16
Specifies the number of fields in a row (can be zero).

Then, for each field, there is the following:

String
The field name.

Int32
If the field can be identified as a column of a specific table, the object ID of the table; otherwise zero.

Int16
If the field can be identified as a column of a specific table, the attribute number of the column; otherwise zero.

Int32
The object ID of the field's data type.

Int16
The data type size (see pg_type.typlen). Note that negative values denote variable-width types.

Int32
The type modifier (see pg_attribute.atttypmod). The meaning of the modifier is type-specific.

Int16
The format code being used for the field. Currently will be zero (text) or one (binary). In a RowDescription returned from the statement variant of Describe, the format code is not yet known and will always be zero.
*/

type RowDescription struct {
	Columns   []string
	SampleRow execution.Row
}

func (r *RowDescription) Type() Type {
	return M_RowDescription
}

func (r *RowDescription) Dump() Buffer {
	data := Buffer{}
	data.AddInt16(len(r.Columns))

	for i, col := range r.Columns {
		data.AddString(col)
		// skip table id
		data.AddInt32(0)
		// skip column offset
		data.AddInt16(0)

		var dt Oid
		var tl int
		switch r.SampleRow[i].(type) {
		case float64:
			tl = 8
			dt = T_float8
		case string:
			tl = -1
			dt = T_text
		case bool:
			tl = 1
			dt = T_bool
		}

		data.AddInt32(int(dt))
		data.AddInt16(tl)

		// ignore type modifiers and format.
		data.AddInt32(-1)
		// default to text format for the column.
		data.AddInt16(0)
	}
	return data
}

/*
DataRow (B)
Byte1('D')
Identifies the message as a data row.

Int32
Length of message contents in bytes, including self.

Int16
The number of column values that follow (possibly zero).

Next, the following pair of fields appear for each column:

Int32
The length of the column value, in bytes (this count does not include itself). Can be zero. As a special case, -1 indicates a NULL column value. No value bytes follow in the NULL case.

Byten
The value of the column, in the format indicated by the associated format code. n is the above length.
*/

type DataRow struct {
	Row execution.Row
}

func (d *DataRow) Type() Type {
	return M_DataRow
}

func (d *DataRow) Dump() Buffer {
	data := Buffer{}
	data.AddInt16(len(d.Row))

	for _, raw := range d.Row {
		var valStr string
		switch v := raw.(type) {
		case float64:
			valStr = strconv.FormatFloat(v, 'g', -1, 64)
		case string:
			valStr = v
		case bool:
			if v {
				valStr = "t"
			} else {
				valStr = "f"
			}
		default:
			panic(fmt.Errorf("unexpected type when serializing data row %T", v))
		}
		b := []byte(valStr)
		data.AddInt32(len(b))
		data.AddBytes(b)
	}

	return data
}

type CommandComplete struct {
	Command string
	Count   int
}

func (c *CommandComplete) Type() Type {
	return M_CommandComplete
}

func (c *CommandComplete) Dump() Buffer {
	data := Buffer{}
	data.AddString(c.Command + " " + strconv.Itoa(c.Count))
	return data
}

/*
ErrorResponse (B)
Byte1('E')
Identifies the message as an error.

Int32
Length of message contents in bytes, including self.

The message body consists of one or more identified fields, followed by a zero byte as a terminator. Fields can appear in any order. For each field there is the following:

Byte1
A code identifying the field type; if zero, this is the message terminator and no string follows. The presently defined field types are listed in Section 53.8. Since more field types might be added in future, frontends should silently ignore fields of unrecognized type.

String
The field value.
*/

type ErrorResponse struct {
	Error error
}

func (e *ErrorResponse) Type() Type {
	return M_ErrorResponse
}

func (e *ErrorResponse) Dump() Buffer {
	data := Buffer{}
	data.AddByte(E_Severity)
	data.AddString("ERROR")
	data.AddByte(E_Message)
	data.AddString(e.Error.Error())
	data.AddNull()
	return data
}
