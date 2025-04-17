package wire

import (
	"fmt"
	"net"

	"github.com/angles-n-daemons/popsql/pkg/api/message"
	"github.com/angles-n-daemons/popsql/pkg/sql/execution"
)

func (srv *Server) loop(conn net.Conn) error {
	for {
		t, data, err := readMessage(conn)
		if err != nil {
			fmt.Println("unexpected read error", err)
		}
		switch t {
		case message.M_Query:
			q := message.Parse[message.Query](data)
			result, err := srv.db.Query(q.Query, nil)
			if err != nil {
				return err
			}
			messages := resultToMessages(result)
			for _, msg := range messages {
				err = writeMessage(conn, msg)
				if err != nil {
					return err
				}
			}
		case message.M_Terminate:
			return nil
		default:
			fmt.Println("unknown message type", t)
		}
		err = writeMessage(conn, &message.ReadyForQuery{})
		if err != nil {
			return nil
		}
	}
}

func resultToMessages(result *execution.Result) []message.Dumpable {
	msgs := []message.Dumpable{}

	// No queries should return 0 rows. Assume this and proceed.
	// Column descriptions
	msgs = append(msgs, &message.RowDescription{
		Columns:   result.Columns,
		SampleRow: result.Rows[0],
	})
	// Rows
	for _, row := range result.Rows {
		msgs = append(msgs, &message.DataRow{Row: row})
	}
	// Command complete message.
	msgs = append(msgs, &message.CommandComplete{
		Command: "SELECT",
	})

	return msgs
}
