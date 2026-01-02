package main

import (
	"fmt"

	"github.com/rankuw/VoltKV/server"
	"github.com/rankuw/VoltKV/store"
)

func main() {
	store := store.NewStore()
	server := server.NewServer(store)

	if err := server.ListenAndServe(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
