package wire

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/angles-n-daemons/popsql/pkg/api/message"
	"github.com/angles-n-daemons/popsql/pkg/db"
)

const SSLRequestCode = 80877103 // Magic code for SSL request

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

	for {
		t, data, err := readMessage(conn)
		if err != nil {
			fmt.Println("unexpected read error", err)
		}
		switch t {
		case message.M_Query:
			fmt.Println("query", string(data))
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

	_, err = message.Parse[message.Startup](data)
	if err != nil {
		return err
	}

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
