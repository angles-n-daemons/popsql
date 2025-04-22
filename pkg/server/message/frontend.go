package message

/*
frontend.go contains the messages sent by the client to the server.

It specifically contains the structs that represent the messages,
the utility functions for parsing them on the server side, and
the common `Parseable` interface they must all satisfy.
*/

type Parseable[P any] interface {
	Load(b Buffer) P
}

func Parse[P Parseable[P]](b Buffer) P {
	var p P
	return p.Load(b)
}

/*
StartupMessage (F)
Int32
Length of message contents in bytes, including self.

Int32(196608)
The protocol version number. The most significant 16 bits are the major version number (3 for the protocol described here). The least significant 16 bits are the minor version number (0 for the protocol described here).

The protocol version number is followed by one or more pairs of parameter name and value strings. A zero byte is required as a terminator after the last name/value pair. Parameters can appear in any order. user is required, others are optional. Each parameter is specified as:

String
The parameter name. Currently recognized names are:

user
The database user name to connect as. Required; there is no default.

database
The database to connect to. Defaults to the user name.

options
Command-line arguments for the backend. (This is deprecated in favor of setting individual run-time parameters.) Spaces within this string are considered to separate arguments, unless escaped with a backslash (\); write \\ to represent a literal backslash.

replication
Used to connect in streaming replication mode, where a small set of replication commands can be issued instead of SQL statements. Value can be true, false, or database, and the default is false. See Section 53.4 for details.

In addition to the above, other parameters may be listed. Parameter names beginning with _pq_. are reserved for use as protocol extensions, while others are treated as run-time parameters to be set at backend start time. Such settings will be applied during backend start (after parsing the command-line arguments if any) and will act as session defaults.

String
The parameter value.
*/
type Startup struct {
	Version int
	Data    map[string]string
}

func (s Startup) Load(b Buffer) Startup {
	var o Startup
	o.Version = b.ReadInt32()
	o.Data = b.ReadObject()
	return o
}

/*
Byte1('Q')
Identifies the message as a simple query.

Int32
Length of message contents in bytes, including self.

String
The query string itself.
*/

type Query struct {
	Query string
}

func (q Query) Load(b Buffer) Query {
	var o Query
	o.Query = b.ReadString()
	return o
}
