package main

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/Av3nger123/mini-redis/core"
)

func main() {

	cache, err := core.NewCache("cache.txt", 5*time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	defer cache.StopClearingRecords()

	fmt.Println("Cache server listening on :6379")

	for {
		conn, err := listener.Accept()
		connection := core.NewConnection(conn)
		if err != nil {
			fmt.Println("Error occurred:", err.Error())
			continue
		}
		fmt.Println("New connection:", connection)

		go subscribeEvents(connection)
		go cache.HandleConnection(connection)
	}

}

func subscribeEvents(conn core.Connection) {
	writer := bufio.NewWriter(conn.Connection)
	for {
		message := <-conn.Channel
		writer.WriteString(message + "\n")
		fmt.Println("Publishing the message", message)
		writer.Flush()
	}
}
