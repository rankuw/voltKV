package main

import (
	"fmt"

	"github.com/rankuw/VoltKV/server"
)

func main() {
	server := server.NewServer()

	if err := server.ListenAndServe(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
