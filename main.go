package main

import (
	"fmt"

	"github.com/rankuw/VoltKV/server"
	"github.com/rankuw/VoltKV/store"
)

func main() {
	kv := store.NewStore()

	aof, err := store.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()
	server := server.NewServer(kv, aof)

	if err := server.ListenAndServe(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
