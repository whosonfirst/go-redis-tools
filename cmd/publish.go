package main

import (
       "bufio"
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

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		redis_client.Publish(*redis_channel, scanner.Text())
	}

	os.Exit(0)
}
