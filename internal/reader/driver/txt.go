package driver

import (
	"bufio"
	"go_echo/internal/reader/data"
	"os"

	"github.com/pkg/errors"
)

type TXTParser struct {
	FileName string
}

func (f *TXTParser) Init(options data.FileOptions) error {
	f.FileName = options.FileName
	return nil
}

func (f *TXTParser) Read() (<-chan []string, <-chan error) {
	outCh := make(chan []string)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)
		var (
			fileReader *os.File
			err        error
		)
		if fileReader, err = os.Open(f.FileName); err != nil {
			errCh <- errors.Wrap(err, "could not open csv file")
			return
		}
		defer fileReader.Close()
		scanner := bufio.NewScanner(fileReader)
		for scanner.Scan() {
			outCh <- []string{scanner.Text()}
		}
	}()

	return outCh, errCh
}
