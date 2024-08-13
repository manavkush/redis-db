package main

import (
	"log"
	"log/slog"
	"net"
	"strings"
)

type Peer struct {
	conn net.Conn
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{conn: conn}
}

// readLoop starts a loop to read value continously from the peer
// After reading the value, the value is parsed and validated to be of proper structure.

func (p *Peer) readLoop(aof *Aof) {
	resp := NewResp(p.conn)
	for {
		value, err := resp.Read()
		if err != nil {
			slog.Error("Error while reading from peer.", "remoteAddr", p.conn.RemoteAddr())
			break
		}

		log.Println(value)

		if value.typ != TYPE_ARRAY {
			log.Println("ERROR: Invalid value type found. Expected array found ", value.typ)
			continue
		}
		if len(value.array) < 1 {
			slog.Error("Invalid array found. Expected non-zero length for array")
			continue
		}

		writer := NewWriter(p.conn)

		command := COMMAND(strings.ToUpper(value.array[0].bulk))
		commandArgs := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			slog.Error("Unsupported command.", "command", command)
			writer.Write(Value{typ: TYPE_STRING, str: ""})
			continue
		}

		result := handler(commandArgs)
		err = writer.Write(result)
		if err != nil {
			slog.Error("Error in writing response to the client")
		}

		if command == COMMAND_HSET || command == COMMAND_SET {
			aof.Write(value)
		}
	}
}
