package main

import (
	"bufio"
	"flag"
	"fmt"
	redis "gopkg.in/redis.v1"
	"log"
	"os"
	"strings"
)

func main() {

	var redis_host = flag.String("redis-host", "localhost", "Redis host")
	var redis_port = flag.Int("redis-port", 6379, "Redis port")
	var redis_channel = flag.String("redis-channel", "", "Redis channel to publish to")

	flag.Parse()

	if *redis_channel == "" {
		log.Fatal("Missing channel")
	}

	redis_endpoint := fmt.Sprintf("%s:%d", *redis_host, *redis_port)

	/*
		import "github.com/docker/go-redis-server"

		type LocalHandler struct {
		     server.DefaultHandler
		}

		if *redis_server {

			handler := &LocalHandler{}

			cfg := server.DefaultConfig()
			cfg.Host(*redis_host)
			cfg.Port(*redis_port)
			cfg.Handler(handler)

			daemon, err := server.NewServer(cfg)

			if err != nil {
				log.Fatal("Failed to create daemon ", err)
			}

			go func() {

				err := daemon.ListenAndServe()

				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	*/

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

			log.Println("SEND", msg)
			i := redis_client.Publish(*redis_channel, msg)
			log.Println("RECEIVE", i)
		}

	} else {

		msg := strings.Join(args, " ")
		redis_client.Publish(*redis_channel, msg)
	}

	os.Exit(0)
}
