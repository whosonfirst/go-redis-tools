package pubsub

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Publisher struct {
	redis_client *redis.Client
}

func NewPublisher(host string, port int) (*Publisher, error) {

	redis_endpoint := fmt.Sprintf("%s:%d", host, port)

	redis_client := redis.NewClient(&redis.Options{
		Addr: redis_endpoint,
	})

	p := Publisher{
		redis_client: redis_client,
	}

	return &p, nil
}

func (p *Publisher) Publish(channel string, message string) error {

	rsp := p.redis_client.Publish(channel, message)
	return rsp.Err()
}

func (p *Publisher) Close() error {
	return p.redis_client.Close()
}
