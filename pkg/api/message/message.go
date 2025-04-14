package message

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
