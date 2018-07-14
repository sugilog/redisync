package main

import (
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func main() {
	conn, err := redis.Dial("tcp", ":6379")

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// https://godoc.org/github.com/gomodule/redigo/redis#hdr-Publish_and_Subscribe
	conn.Send("SYNC")
	conn.Flush()

	for {
		reply, err := conn.Receive()

		if err != nil {
			panic(err)
		}

		switch v := reply.(type) {
		case []byte:
			message := string(v[:])

			switch message {
			case "PING":
			case "SELECT":
			case "0":
			default:
				if strings.Contains(message, "redis-ver") {
				} else {
					panic(fmt.Sprintf("unsupported message, %s", message))
				}
			}
		case []interface{}:
			messages := make([]string, len(v))

			for i, item := range v {
				message := item.([]byte)
				messages[i] = string(message[:])
			}

			command := strings.ToUpper(messages[0])
			args := messages[1:]

			switch messages[0] {
			case "PING":
			case "SELECT":
			default:
				// get messages
				// [set a 1 EX 5]
				// [DEL a]
				// [del a]
				fmt.Printf("command: %s, args: %v\n", command, args)
			}
		default:
			panic(fmt.Sprintf("unsupported message, %v", v))
		}
	}
}
