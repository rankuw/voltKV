package server

import (
	"fmt"
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
		fmt.Println(err)
		return
	}

	fmt.Println(data)
	err = writer.Write(data)

	if err != nil {
		fmt.Println(err)
		return
	}

}
