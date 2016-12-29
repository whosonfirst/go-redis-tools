package main

import (
	"flag"
	"fmt"
	"gopkg.in/redis.v1"
	"log"
	"os"
)

func main() {

	var redis_host = flag.String("redis-host", "localhost", "The Redis host to connect to.")
	var redis_port = flag.Int("redis-port", 6379, "The Redis port to connect to.")
	var redis_channel = flag.String("redis-channel", "", "The Redis channel to publish to.")

	flag.Parse()

	if *redis_channel == "" {
		log.Fatal("Missing channel")
	}

	redis_endpoint := fmt.Sprintf("%s:%d", *redis_host, *redis_port)

	redis_client := redis.NewTCPClient(&redis.Options{
		Addr: redis_endpoint,
	})

	defer redis_client.Close()

	_, err := redis_client.Ping().Result()

	if err != nil {
		log.Fatal("Failed to ping Redis server ", err)
	}

	pubsub_client := redis_client.PubSub()
	defer pubsub_client.Close()

	err = pubsub_client.Subscribe(*redis_channel)

	if err != nil {
		msg := fmt.Sprintf("Failed to subscribe to channel %s, because %s", *redis_channel, err)
		log.Fatal(msg)
	}

	for {

		i, _ := pubsub_client.Receive()

		if msg, _ := i.(*redis.Message); msg != nil {
			log.Println(msg.Payload)
		}
	}

	// please for to add signal handlers here...

	err = pubsub_client.Unsubscribe(*redis_channel)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
