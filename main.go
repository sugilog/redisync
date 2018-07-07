package main

import (
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
)

// get messages
// [set a 1 EX 5]
// [DEL a]
// [del a]

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

			if message == "PING" {
				// nothing to do
			} else if strings.Contains(message, "redis-ver") {
				// nothing to do
			} else {
				fmt.Println(message)
			}
		case []interface{}:
			messages := make([]string, len(v))

			for i, item := range v {
				message := item.([]byte)
				messages[i] = string(message[:])
			}

			switch messages[0] {
			case "PING":
			case "SELECT":
			default:
				fmt.Println(messages)
			}
		default:
			panic(fmt.Sprintf("unsupported message, %v", v))
		}
	}
}
