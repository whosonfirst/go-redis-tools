package resp

import (
	"bufio"
	_ "bytes"
	"fmt"
	"io"
	"strconv"
)

var (
	arrayPrefixSlice      = []byte{'*'}
	bulkStringPrefixSlice = []byte{'$'}
	lineEndingSlice       = []byte{'\r', '\n'}
	errorSlice            = []byte{'-', 'E', 'R', 'R'}
)

type RESPWriter struct {
	*bufio.Writer
}

func NewRESPWriter(writer io.Writer) *RESPWriter {
	return &RESPWriter{
		Writer: bufio.NewWriter(writer),
	}
}

func (w *RESPWriter) WriteResponse(args []string) error {

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

	w.Write(arrayPrefixSlice)
	w.WriteString(fmt.Sprintf("%s", err))
	w.Write(lineEndingSlice)

	return w.Flush()
}
