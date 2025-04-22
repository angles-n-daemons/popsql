package server

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/angles-n-daemons/popsql/pkg/db"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/execution"
	"github.com/angles-n-daemons/popsql/pkg/server/message"
)

const SSLRequestCode = 80877103 // Magic code for SSL request

func Run() {
	srv := NewServer()
	err := srv.ListenAndServe(":5432")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

type Server struct {
	db *db.Engine
}

func NewServer() *Server {
	return &Server{
		db: db.GetEngine(),
	}
}

func (srv *Server) ListenAndServe(address string) error {
	fmt.Println("listening on", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}

		if err != nil {
			return err
		}

		go func() {
			err = srv.serve(conn)
			if err != nil && err != io.EOF {
				fmt.Printf("unexpected error serving connection %s", err)
			}
		}()
	}
}

func (srv *Server) serve(conn net.Conn) error {
	defer conn.Close()
	err := connInit(conn)
	if err != nil {
		return err
	}

	return srv.loop(conn)
}

func connInit(conn net.Conn) error {
	data, err := readMessageRaw(conn)
	if err != nil {
		return err
	}

	// Special case, check for SSL escalation request.
	// If found, send a no response, and read the next message
	// as startup.
	num := data.PeekUint32()
	if num == SSLRequestCode {
		conn.Write([]byte{message.M_No})
		data, err = readMessageRaw(conn)
		if err != nil {
			return err
		}
	}

	// Read the startup message, ignore it for now.
	_ = message.Parse[message.Startup](data)

	err = writeMessage(conn, &message.AuthenticationOk{})
	if err != nil {
		return nil
	}

	err = writeMessage(conn, &message.ReadyForQuery{})
	if err != nil {
		return nil
	}
	return nil
}

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
			messages := resultToMessages(result, err)
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

func resultToMessages(result *execution.Result, err error) []message.Dumpable {
	msgs := []message.Dumpable{}

	if err == nil {
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
	} else {
		msgs = append(msgs, &message.ErrorResponse{Error: err})
	}

	return msgs
}
