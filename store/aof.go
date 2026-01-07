package store

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rankuw/VoltKV/resp"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
	buf  *bufio.Writer
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
		buf:  bufio.NewWriter(f),
	}

	go func(aof *Aof) {
		for {
			aof.mu.Lock()
			aof.buf.Flush()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}(aof)

	return aof, nil
}

func (aof *Aof) Close() {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Close()
}

func (aof *Aof) Write(v resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.buf.Write(v.Marshal())

	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) LoadData(kv *Store) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Seek(0, 0)
	reader := resp.NewResp(aof.file)
	for {
		value, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		if command == "SET" {
			key := args[0].Bulk
			val := args[1].Bulk
			var ttl time.Duration

			if len(args) > 2 {
				for i := 2; i < len(args); i++ {
					arg := strings.ToUpper(args[i].Bulk)
					if arg == "EX" {
						if i+1 < len(args) {
							seconds, err := strconv.Atoi(args[i+1].Bulk)
							if err == nil {
								ttl = time.Duration(seconds) * time.Second
							}
						}
						i++
					}
				}
			}
			kv.Set(key, val, ttl)

		}
	}
	return nil
}
