package pubsub

// https://redis.io/topics/protocol
// https://redis.io/topics/pubsub
// https://www.redisgreen.net/blog/reading-and-writing-redis-protocol/

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
)

var (
	arrayPrefixSlice      = []byte{'*'}
	bulkStringPrefixSlice = []byte{'$'}
	lineEndingSlice       = []byte{'\r', '\n'}
)

type RESPWriter struct {
	*bufio.Writer
}

func NewRESPWriter(writer io.Writer) *RESPWriter {
	return &RESPWriter{
		Writer: bufio.NewWriter(writer),
	}
}

type RESPReader struct {
	*bufio.Reader
}

func NewReader(reader io.Reader) *RESPReader {
	return &RESPReader{
		Reader: bufio.NewReaderSize(reader, 32*1024),
	}
}

type Conns map[string]net.Conn

type Clients map[string]bool

type Server struct {
	host string
	port int
	subs map[string]Client
}

func NewServer(host string, port string) (*Server, error) {

	conns := make(map[string]net.Conn)
	subs := make(map[string]Clients)

	s := Server{
		host: host,
		port: port,
		subs: subs,
	}

	return &s, nil
}

func (s *Server) ListenAndServe() error {

	address := fmt.Sprintf("%s:%d", s.host, s.port)
	daemon, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	defer daemon.Close()

	for {

		conn, err := daemon.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go s.Receive(conn)
	}

	return nil
}

func (s *Server) Receive(conn net.Conn) error {

}

func (s *Server) Send() error {

}

func (s *Server) Subscribe(conn net.Conn, channels []string) error {

	remote := conn.RemoteAddr()

	_, ok := s.conns[remote]

	if !ok {
		s.conns[remote] = conn
	}

	for _, ch := range channels {

		clients, ok := s.subs[ch]

		if !ok {
			clients = make(map[string]bool)
			s.subs[ch] = clients
		}

		s.subs[ch][remote] = true
	}

	return nil
}

func (s *Server) Unsubscribe(conn net.Conn, channels []string) error {

	remote := conn.RemoteAddr()

	_, ok := s.conns[remote]

	if !ok {
		return errors.New("Can not find connection thingy for %s", remote)
	}

	for _, ch := range channels {

		clients, ok := s.subs[ch]

		if !ok {
			continue
		}

		_, ok = s.subs[ch][remote]

		if !ok {
			continue
		}

		delete(s.subs[ch][remote])
	}

	count := 0

	for _, ch := range channels {

		for addr, _ := range ch {

			if addr == remote {
				count += 1
			}
		}
	}

	if count == 0 {

		conn.Close()
		delete(s.conns, remote)
	}

	return nil
}

func (s *Server) Publish(channel string, message []byte) error {

	sub, ok := s.subs[channel]

	if !ok {
		return nils
	}

	for remote, _ := range sub {

		conn, ok := s.conns[remote]

		if !ok {
			continue
		}

		go conn.Write(message)
	}

	return nil
}
