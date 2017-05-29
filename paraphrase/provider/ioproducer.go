package provider

import (
	"bufio"
	"fmt"
	"io"
)

func NewFileListProducer(namespace string, reader io.Reader) (DocumentProducer, error) {
	output := make(DocumentProducer, 10)

	go func() {
		defer close(output)

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			path := scanner.Text()
			d := Document{path, namespace, readFileCallback(path)}
			output <- d
		}
	}()

	return output, nil
}

// NewDummyProducer creates a producer that consumes elements and writes
// their paths to the given writer. It does not pass any elements through.
func NewDummyProducer(producer DocumentProducer, writer io.Writer) DocumentProducer {
	output := make(DocumentProducer, 10)

	go func() {
		defer close(output)

		for doc := range producer {
			fmt.Fprintf(writer, "%s\t%s", doc.Namespace(), doc.Path())
			fmt.Fprintln(writer)
		}
	}()

	return output
}
