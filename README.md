# go-redis-tools

A Go port of the Python redis-tools package.

## Tools

### pubsubd

```
./bin/pubsubd -h
Usage of ./bin/pubsubd:
  -host string
    	The hostname to listen on. (default "localhost")
  -port int
    	The port number to listen on. (default 6379)
```

This will launch a daemon to support most (but not all) of the [Redis Publish/Subscribe protocol](https://redis.io/topics/pubsub). It has not been tested for load or scale but it works. The following commands are supported: `PING, SUBSCRIBE, UNSUBSCRIBE, PUBLISH`

### publish

```
./bin/publish  -h
Usage of ./bin/publish:
  -pubsubd
    	Invoke a local pubsubd server that publish and subscribe clients will connect to. This may be useful when you don't have a local copy of Redis around.
  -redis-channel string
    	The Redis channel to publish to.
  -redis-host string
    	The Redis host to connect to. (default "localhost")
  -redis-port int
    	The Redis port to connect to. (default 6379)
```

Publish a message to PubSub channel. If the message is `-` then the client will read and publish all subsequent input from STDIN.

### subscribe

```
./bin/subscribe -h
Usage of ./bin/subscribe:
  -redis-channel string
    	The Redis channel to publish to.
  -redis-host string
    	The Redis host to connect to. (default "localhost")
  -redis-port int
    	The Redis port to connect to. (default 6379)
```

Subscribe to a PubSub channel and print the result to `STDOUT` using the Go [log package](https://golang.org/pkg/log/) (other outputs to follow).

## See also

* https://github.com/whosonfirst/redis-tools
* https://redis.io/topics/pubsub
