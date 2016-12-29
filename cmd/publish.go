package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-redis-tools/pubsub"
	"gopkg.in/redis.v1"
	"log"
	"os"
	"strings"
)

func main() {

	var redis_host = flag.String("redis-host", "localhost", "Redis host")
	var redis_port = flag.Int("redis-port", 6379, "Redis port")
	var redis_channel = flag.String("redis-channel", "", "Redis channel to publish to")
	var pubsubd = flag.Bool("pubsubd", false, "...")

	flag.Parse()

	if *redis_channel == "" {
		log.Fatal("Missing channel")
	}

	redis_endpoint := fmt.Sprintf("%s:%d", *redis_host, *redis_port)

	if *pubsubd {

		server, err := pubsub.NewServer(*redis_host, *redis_port)

		if err != nil {
			log.Fatal(err)
		}

		go func() {

			err := server.ListenAndServe()

			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	redis_client := redis.NewTCPClient(&redis.Options{
		Addr: redis_endpoint,
	})

	defer redis_client.Close()

	_, err := redis_client.Ping().Result()

	if err != nil {
		log.Fatal("Failed to ping Redis server ", err)
	}

	args := flag.Args()

	if len(args) == 1 && args[0] == "-" {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			msg := scanner.Text()
			err := redis_client.Publish(*redis_channel, msg)

			if err != nil {
				log.Println(err)
			}
		}

	} else {

		msg := strings.Join(args, " ")
		redis_client.Publish(*redis_channel, msg)
	}

	os.Exit(0)
}
