package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/rankuw/VoltKV/resp"
	"github.com/rankuw/VoltKV/store"
)

type Server struct {
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{
		store: store,
	}
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

	command := data.Array[0].Bulk
	arugments := data.Array[1:]

	switch command {
	case "PING":
		if err := writer.Write(resp.Value{Type: resp.STRING, Str: "PONG"}); err != nil {
			fmt.Println(err)
		}
		return

	case "GET":
		if len(arugments) != 1 {
			if err := writer.Write(resp.Value{Type: resp.ERROR, Str: "ERR wrong number of arugments for GET"}); err != nil {
				fmt.Println(err)
			}
		}

		key := arugments[0].Bulk
		val, ok := s.store.Get(key)

		if !ok {
			if err := writer.Write(resp.Value{Type: resp.BULK, IsNull: true}); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := writer.Write(resp.Value{Type: resp.BULK, Bulk: val}); err != nil {
				fmt.Println(err)
			}
		}

		return

	case "SET":
		if len(arugments) < 2 {
			if err := writer.Write(resp.Value{Type: resp.ERROR, Str: "ERR wrong number of arugments for SET"}); err != nil {
				fmt.Println(err)
			}
		}

		key := arugments[0].Bulk
		val := arugments[1].Bulk

		var ttl time.Duration
		if len(arugments) >= 4 && arugments[2].Bulk == "EX" {
			seconds, _ := strconv.Atoi(arugments[3].Bulk)
			ttl = time.Duration(seconds) * time.Second
		}

		s.store.Set(key, val, ttl)
		if err := writer.Write(resp.Value{Type: resp.STRING, Str: "OK"}); err != nil {
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
