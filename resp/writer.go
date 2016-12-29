package resp

// https://godoc.org/github.com/fzzy/radix/redis/resp

import (
	"bufio"
	"fmt"
	"io"
	_ "log"
	"os"
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
		os.Stdout,
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

func (w *RESPWriter) WriteArray(args []string) error {

	w.Write(arrayPrefixSlice)
	w.WriteString(strconv.Itoa(len(args)))
	w.Write(lineEndingSlice)

	for _, arg := range args {
		w.Write(bulkStringPrefixSlice)
		w.WriteString(strconv.Itoa(len(arg)))
		w.Write(lineEndingSlice)
		w.WriteString(arg)
		w.Write(lineEndingSlice)
	}

	return w.Flush()
}

func (w *RESPWriter) WriteError(err error) error {

	w.Write(errorPrefixSlice)
	w.WriteString(fmt.Sprintf("%s", err))
	w.Write(lineEndingSlice)

	return w.Flush()
}
