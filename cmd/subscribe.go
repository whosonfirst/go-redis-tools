package main

import (
       "flag"
       "fmt"	
	"gopkg.in/redis.v1"
	"log"
	"os"		
)

func main() {

	var redis_host = flag.String("redis-host", "localhost", "Redis host")
	var redis_port = flag.Int("redis-port", 6379, "Redis port")
	var redis_channel = flag.String("redis-channel", "", "Redis channel")

	if *redis_channel == "" {
		log.Fatal("Missing channel")
	}
	
	redis_endpoint := fmt.Sprintf("%s:%d", *redis_host, *redis_port)

	redis_client := redis.NewTCPClient(&redis.Options{
		Addr: redis_endpoint,
	})

	defer redis_client.Close()

	pubsub_client := redis_client.PubSub()
	defer pubsub_client.Close()

	err := pubsub_client.Subscribe(*redis_channel)

	if err != nil {
		log.Fatal("Failed to subscribe to channel")
	}
	
	for {

		i, _ := pubsub_client.Receive()

		if msg, _ := i.(*redis.Message); msg != nil {
			log.Println(msg.Payload)
		}
	}

	os.Exit(0)
}
