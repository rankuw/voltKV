package server

import (
	"fmt"
	"io"
	"net"

	"github.com/rankuw/VoltKV/resp"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ListenAndServe(address string) error {
	ln, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	fmt.Println("Server started on ", address)
	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	respReader := resp.NewResp(conn)
	writer := resp.NewWriter(conn)

	data, err := respReader.Read()

	if err != nil {
		if err != io.EOF {
			fmt.Println(err)
		}
		return
	}

	if data.Type != resp.ARRAY {
		fmt.Println("Invalid request, exptected an array")
		if err := writer.Write(resp.Value{Type: resp.ERROR, Str: "ERR request must be an ARRAY"}); err != nil {
			fmt.Println(err)
		}
		return
	}

	if len(data.Array) == 0 {
		fmt.Println("Invalid request, no arguments given")
		if err := writer.Write(resp.Value{Type: resp.ERROR, Str: "ERR request must contains something"}); err != nil {
			fmt.Println(err)
		}
		return
	}

	command := data.Array[0].Str
	// arugments := data.Array[1:]

	// fmt.Println(command, arugments, "HELLO")

	fmt.Println("This is the command -> ", command)

	switch command {
	case "PING":
		if err := writer.Write(resp.Value{Type: resp.STRING, Str: "PONG"}); err != nil {
			fmt.Println(err)
		}
		return

	case "GET":
		fmt.Println("here hu m.")
		if err := writer.Write(resp.Value{Type: resp.STRING, Str: "GETRESPONSE"}); err != nil {
			fmt.Println(err)
		}
		return

	case "SET":
		if err := writer.Write(resp.Value{Type: resp.STRING, Str: "SETRESPONSE"}); err != nil {
			fmt.Println(err)
		}
		return
	}

	err = writer.Write(data)

	if err != nil {
		fmt.Println(err)
		return
	}

}
