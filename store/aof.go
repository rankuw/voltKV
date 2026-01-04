package store

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rankuw/VoltKV/resp"
)

type Aof struct {
	file *os.File
	mu   sync.Mutex
	rd   *bufio.Reader
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	go func(aof *Aof) {
		aof.mu.Lock()
		defer aof.mu.Unlock()
		aof.file.Sync()
		time.Sleep(time.Second)
	}(aof)

	return aof, nil
}

func (aof *Aof) Close() {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Close()
}

func (aof *Aof) Write(v resp.Value) error {
	fmt.Println("I'm in write", v, aof)
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(v.Marshal())

	if err != nil {
		return err
	}

	return nil
}
