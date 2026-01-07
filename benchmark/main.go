package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/rankuw/VoltKV/resp"
)

var (
	concurrency = flag.Int("c", 100, "concurrency level")
	total       = flag.Int("f", 40000, "total number of requests")
)

func main() {
	flag.Parse()
	fmt.Println("Concurrency:", *concurrency)
	fmt.Println("Total Requests:", *total)

	requestPerClient := *total / *concurrency

	var wg sync.WaitGroup
	wg.Add(*concurrency)

	latency := make(chan time.Duration, *total)
	start := time.Now()
	for i := 0; i < *concurrency; i++ {
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", ":8081")

			if err != nil {
				fmt.Println("Error in creating connection ", err)
				return
			}
			defer conn.Close()
			reader := resp.NewResp(conn)

			data := resp.Value{
				Type: resp.ARRAY,
				Array: []resp.Value{
					{Type: resp.BULK, Bulk: "SET"},
					{Type: resp.BULK, Bulk: fmt.Sprintf("key-%d", rand.Int())},
					{Type: resp.BULK, Bulk: "value"},
				},
			}

			dataBytes := data.Marshal()

			for j := 0; j < requestPerClient; j++ {
				fmt.Println("Sending new request ")
				start := time.Now()
				if _, err := conn.Write(dataBytes); err != nil {
					fmt.Println("Error in write ", err)
					return
				}

				_, err = reader.Read()

				if err != nil {
					fmt.Println("Read error:", err)
					return
				} else {
					fmt.Println("fine")
				}

				latency <- time.Since(start)
			}
		}()
	}

	wg.Wait()
	close(latency)

	duration := time.Since(start)

	fmt.Printf("Completed in %v s \n", duration)
	fmt.Printf("Throughput: %.2f requests/sec\n", float64(*total)/duration.Seconds())

	var totalLat time.Duration
	var count int
	for lat := range latency {
		totalLat += lat
		count++
	}
	fmt.Printf("Average Latency: %v\n", totalLat/time.Duration(count))

}
