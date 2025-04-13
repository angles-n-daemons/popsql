package wire

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/angles-n-daemons/popsql/pkg/api/message"
)

type Server struct{}

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
	data, err := readMessageRaw(conn)
	if err != nil {
		return err
	}
	msg := &message.Startup{}
	err = msg.Load(data)
	if err != nil {
		return err
	}
	fmt.Println("msg", msg)
	fmt.Println("err", err)
	return nil
}
