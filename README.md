A Learning go project

<!---
 will update this though 
-->

Features:
 - Basic functionality set and get
 - TTL 
 - persistance (saves to the disk for every x seconds)
 - pub/sub 

Running the code
```sh
go build main.go
./main
```

Here is the client code to connect to the cache

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to the cache server:", err)
		return
	}
	defer conn.Close()

	go receiveMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()

		_, err := fmt.Fprintf(conn, command+"\n")
		if err != nil {
			fmt.Println("Error sending command to the cache server:", err)
			continue
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving message:", err)
			break
		}
		handleMessage(response)
	}
}

func handleMessage(message string) {
	fmt.Println("Received message:", message)
}

```


Usage:

```
>set key value 100
>Received message: value

>get key
>Received message: value

>subscribe topic
>Received message: Subscribed to channel topic

>send topic message

>Received message: message (to the subscribed client)
```