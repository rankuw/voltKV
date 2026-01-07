package server

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/rankuw/VoltKV/resp"
	"github.com/rankuw/VoltKV/store"
)

type Server struct {
	store *store.Store
	aof   *store.Aof
}

func NewServer(store *store.Store, aof *store.Aof) *Server {
	return &Server{
		store: store,
		aof:   aof,
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

	for {
		data, err := respReader.Read()

		if err != nil {
			return
		}

		if data.Type != resp.ARRAY {
			return
		}

		if len(data.Array) == 0 {
			return
		}

		command := data.Array[0].Bulk
		arugments := data.Array[1:]

		switch command {
		case "PING":
			writer.Write(resp.Value{Type: resp.STRING, Str: "PONG"})

		case "GET":
			if len(arugments) != 1 {
				writer.Write(resp.Value{Type: resp.ERROR, Str: "ERR wrong number of arugments for GET"})
			}

			key := arugments[0].Bulk
			val, ok := s.store.Get(key)

			if !ok {
				writer.Write(resp.Value{Type: resp.BULK, IsNull: true})
			} else {
				writer.Write(resp.Value{Type: resp.BULK, Bulk: val})
			}

		case "SET":
			if len(arugments) < 2 {
				writer.Write(resp.Value{Type: resp.ERROR, Str: "ERR wrong number of arugments for SET"})
			}

			key := arugments[0].Bulk
			val := arugments[1].Bulk

			var ttl time.Duration
			if len(arugments) >= 4 && arugments[2].Bulk == "EX" {
				seconds, _ := strconv.Atoi(arugments[3].Bulk)
				ttl = time.Duration(seconds) * time.Second
			}

			if err := s.aof.Write(data); err != nil {
				writer.Write(resp.Value{Type: resp.ERROR, Str: "ERR persistance error"})
				continue
			}
			s.store.Set(key, val, ttl)
			writer.Write(resp.Value{Type: resp.STRING, Str: "OK"})

		default:
			writer.Write(resp.Value{Type: resp.ERROR, Str: "ERR unknown command " + command})
		}
	}

}
