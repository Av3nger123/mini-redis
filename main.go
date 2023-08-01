package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(connection net.Conn) {
	defer connection.Close()

	reader := bufio.NewReader(connection)
	writer := bufio.NewWriter(connection)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
		}

		bytesWritten, err := writer.WriteString(command)
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		fmt.Println(bytesWritten)

		writer.Flush()

	}
}

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error occurred:", err.Error())
			continue
		}

		go handleConnection(connection)
	}

}
