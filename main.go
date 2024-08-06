package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var port = ":6379"

func main() {
	fmt.Printf("Starting listening on port: %v\n", port)
	list, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error in listening at port %v. Err: %v", port, err)
		return
	}

	conn, err := list.Accept()
	if err != nil {
		fmt.Printf("Error in accepting connection at port %v. Err: %v", port, err)
		return
	}
	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error in reading bytes from connection. Err", err)
			os.Exit(1)
		}

		conn.Write([]byte("+OK\r\n"))
	}
}
