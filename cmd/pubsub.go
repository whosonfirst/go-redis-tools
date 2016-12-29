package main

import (
	"flag"
	"github.com/whosonfirst/go-redis-tools/pubsub"
	"log"
	"os"
)

func main() {

	var redis_host = flag.String("redis-host", "localhost", "Redis host")
	var redis_port = flag.Int("redis-port", 6379, "Redis port")

	flag.Parse()

	server, err := pubsub.NewServer(*redis_host, *redis_port)

	if err != nil {
		log.Fatal(err)
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
