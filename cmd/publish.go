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

func (h *MyHandler) Publish(key string, value []byte) (int, error) {

	log.Printf("key: %s value: %s\n", key, value)

	return h.DefaultHandler.Publish(key, value)
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

		handler := &MyHandler{}

		cfg := server.DefaultConfig()
		cfg.Host(*redis_host)
		cfg.Port(*redis_port)
		cfg.Handler(handler)

		daemon, err := server.NewServer(cfg)

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
