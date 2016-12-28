package main

import (
	"bufio"
	"flag"
	"fmt"
	server "github.com/docker/go-redis-server"
	redis "gopkg.in/redis.v1"
	"log"
	"os"
	"strings"
)

type MyHandler struct {
	server.DefaultHandler
}

func main() {

	var redis_server = flag.Bool("redis-server", false, "...")
	var redis_host = flag.String("redis-host", "localhost", "Redis host")
	var redis_port = flag.Int("redis-port", 6379, "Redis port")
	var redis_channel = flag.String("redis-channel", "", "Redis channel to publish to")

	flag.Parse()

	if *redis_channel == "" {
		log.Fatal("Missing channel")
	}

	redis_endpoint := fmt.Sprintf("%s:%d", *redis_host, *redis_port)

	if *redis_server {

		daemon, err := server.NewServer(server.DefaultConfig())

		if err != nil {
			log.Fatal("Failed to create daemon", err)
		}

		go daemon.ListenAndServe()
	}

	redis_client := redis.NewTCPClient(&redis.Options{
		Addr: redis_endpoint,
	})

	defer redis_client.Close()

	args := flag.Args()

	if len(args) == 1 && args[0] == "-" {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			redis_client.Publish(*redis_channel, scanner.Text())
		}

	} else {

		msg := strings.Join(args, " ")
		redis_client.Publish(*redis_channel, msg)
	}

	os.Exit(0)
}
