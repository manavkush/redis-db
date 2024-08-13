package main

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		slog.Error("Unable to open migration file.", "err", err)
		return nil, err
	}
	aof := &Aof{
		file: file,
		rd:   bufio.NewReader(file),
	}

	// Start a go routine to sync the file buffer changes to disk every second.
	go func(*Aof) {
		aof.mu.Lock()
		aof.file.Sync()
		aof.mu.Unlock()

		time.Sleep(time.Second)
	}(aof)

	return aof, nil
}

func (aof *Aof) Read(fn func(Value)) error {
	resp := NewResp(aof.file)
	for {
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			slog.Error("Unable to read the Aof file", "err", err)
			return err
		}

		fn(value)
	}
}

func (aof *Aof) Write(val Value) error {
	bytes := val.Marshal()

	_, err := aof.file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	err := aof.file.Close()
	if err != nil {
		return err
	}

	return nil
}
