package main

import (
	"bufio"
	"fmt"
	"net"
)

type Cache struct {
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (cache *Cache) handleConnection(connection net.Conn) {
	defer connection.Close()

	reader := bufio.NewReader(connection)
	writer := bufio.NewWriter(connection)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Client closed the connection.")
				break
			}
			fmt.Println("Error reading command:", err)
		}
		response := cache.handleCommand(command)
		bytesWritten, err := writer.WriteString(response + "\n")
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		fmt.Println(bytesWritten)

		writer.Flush()

	}
}

func (cache *Cache) handleCommand(command string) string {

	// TODO: Handle commands like a cache
	return command
}

func main() {

	cache := NewCache()
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	fmt.Println("Cache server listening on :6379")

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error occurred:", err.Error())
			continue
		}
		fmt.Println("New connection:", connection)

		go cache.handleConnection(connection)
	}

}
