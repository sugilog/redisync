package main

import (
	"github.com/sugilog/redisync"
)

func main() {
	client := redisync.NewClient(":6379")
	client.Start()
}
