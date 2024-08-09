package main

import (
	"fmt"
	"net"
	"strings"
	"time"
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
			return
		}

		// The redis commands will always be in the form of an array. eg: SET key value
		// If it's not, then it's an invalid request
		if value.typ != "ARRAY" {
			fmt.Println("Invalid request, expected an array. Value:", value)
			time.Sleep(5 * time.Second)
			continue
		}
		if len(value.array) == 0 { // If the array is empty, then cannot parse the command.
			fmt.Println("Invalid request, expected array with length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		commandArgs := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Printf("Invalid command received. command: %v", value)
			writer.Write(Value{typ: "STRING", str: ""})
			continue
		}

		result := handler(commandArgs)
		writer.Write(result)
	}
}
