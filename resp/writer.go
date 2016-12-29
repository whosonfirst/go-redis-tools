package resp

// Adapted from:
// https://www.redisgreen.net/blog/beginners-guide-to-redis-protocol/
// https://www.redisgreen.net/blog/reading-and-writing-redis-protocol/

// Maybe use this:
// https://godoc.org/github.com/fzzy/radix/redis/resp

// Either way this code will be refactored soon (20161229/thisisaaronland)

import (
	"bufio"
	"fmt"
	"io"
	_ "os"
	"strconv"
)

var (
	stringPrefixSlice     = []byte{'+'}
	numberPrefixSlice     = []byte{':'}
	arrayPrefixSlice      = []byte{'*'}
	bulkStringPrefixSlice = []byte{'$'}
	lineEndingSlice       = []byte{'\r', '\n'}
	errorPrefixSlice      = []byte{'-', 'E', 'R', 'R'}
)

type RESPWriter struct {
	*bufio.Writer
}

func NewRESPWriter(writer io.Writer) *RESPWriter {

	writers := []io.Writer{
		writer,
		// this is useful for debugging but otherwise unnecessary (20161229/thisisaaronland)
		// os.Stdout,
	}

	multi := io.MultiWriter(writers...)

	return &RESPWriter{
		Writer: bufio.NewWriter(multi),
	}
}

func (w *RESPWriter) WriteStringMessage(str string) error {

	w.Write(stringPrefixSlice)
	w.WriteString(str)
	w.Write(lineEndingSlice)

	return w.Flush()
}

func (w *RESPWriter) WriteNullMessage() error {

	w.Write(bulkStringPrefixSlice)
	w.WriteString("-1")
	w.Write(lineEndingSlice)

	return w.Flush()
}

func (w *RESPWriter) WriteSubscribeMessage(channels []string) error {

	for i, ch := range channels {

		w.Write(arrayPrefixSlice)
		w.WriteString("3")
		w.Write(lineEndingSlice)

		w.WriteString("$9")
		w.Write(lineEndingSlice)

		w.WriteString("subscribe")
		w.Write(lineEndingSlice)

		w.Write(bulkStringPrefixSlice)
		w.WriteString(strconv.Itoa(len(ch)))
		w.Write(lineEndingSlice)

		w.WriteString(ch)
		w.Write(lineEndingSlice)

		w.Write(numberPrefixSlice)
		w.WriteString(strconv.Itoa(i + 1))
		w.Write(lineEndingSlice)
	}

	return w.Flush()
}

func (w *RESPWriter) WriteUnsubscribeMessage(channels []string) error {

	i := len(channels) - 1

	for _, ch := range channels {

		w.Write(arrayPrefixSlice)
		w.WriteString("3")
		w.Write(lineEndingSlice)

		w.WriteString("$11")
		w.Write(lineEndingSlice)

		w.WriteString("unsubscribe")
		w.Write(lineEndingSlice)

		w.Write(bulkStringPrefixSlice)
		w.WriteString(strconv.Itoa(len(ch)))
		w.Write(lineEndingSlice)

		w.WriteString(ch)
		w.Write(lineEndingSlice)

		w.Write(numberPrefixSlice)
		w.WriteString(strconv.Itoa(i - 1))
		w.Write(lineEndingSlice)
	}

	return w.Flush()
}

func (w *RESPWriter) WritePublishMessage(channel string, msg string) error {

	w.Write(arrayPrefixSlice)
	w.WriteString("3")
	w.Write(lineEndingSlice)

	w.Write(bulkStringPrefixSlice)
	w.WriteString("7")
	w.Write(lineEndingSlice)

	w.WriteString("message")
	w.Write(lineEndingSlice)

	w.Write(bulkStringPrefixSlice)
	w.WriteString(strconv.Itoa(len(channel)))
	w.Write(lineEndingSlice)

	w.WriteString(channel)
	w.Write(lineEndingSlice)

	w.Write(bulkStringPrefixSlice)
	w.WriteString(strconv.Itoa(len(msg)))
	w.Write(lineEndingSlice)

	w.WriteString(msg)
	w.Write(lineEndingSlice)

	return w.Flush()
}

func (w *RESPWriter) WriteErrorMessage(err error) error {

	w.Write(errorPrefixSlice)
	w.WriteString(fmt.Sprintf("%s", err))
	w.Write(lineEndingSlice)

	return w.Flush()
}
