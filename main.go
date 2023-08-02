package main

import (
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

	fmt.Println("Cache server listening on :6379")

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error occurred:", err.Error())
			continue
		}
		fmt.Println("New connection:", connection)

		go cache.HandleConnection(connection)
	}

}
