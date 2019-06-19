package pubsub

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Subscriber struct {
	redis_client *redis.Client
}

func NewSubscriber(host string, port int) (*Subscriber, error) {

	redis_endpoint := fmt.Sprintf("%s:%d", host, port)

	redis_client := redis.NewClient(&redis.Options{
		Addr: redis_endpoint,
	})

	s := Subscriber{
		redis_client:  redis_client,
	}

	return &s, nil
}

func (s *Subscriber) Subscribe(channel string, message_ch chan string) error {

	pubsub_client := s.redis_client.PSubscribe(channel)
	defer pubsub_client.Close()

	for {

		i, _ := pubsub_client.Receive()

		if msg, _ := i.(*redis.Message); msg != nil {
			message_ch <- msg.Payload
		}
	}

	return nil
}

func (s *Subscriber) Close() error {

	var err error

	err = s.redis_client.Close()

	if err != nil {
		return err
	}

	return nil
}
