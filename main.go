package main

import (
	"fmt"
	"net"
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
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Printf("Error in reading resp value. Err: %v\n", err)
		}
		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))
	}
}
