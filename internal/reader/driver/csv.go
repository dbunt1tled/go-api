package driver

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"go_echo/internal/reader/data"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type CSVParser struct {
	FileName  string
	Delimiter rune
}

func (f *CSVParser) Init(options data.FileOptions) error {
	var err error
	f.FileName = options.FileName
	if options.Delimiter == nil {
		f.Delimiter, err = f.detectDelimiter(f.FileName)
		if err != nil {
			return err
		}
	} else {
		f.Delimiter = *options.Delimiter
	}
	return nil
}

func (f *CSVParser) Read() (<-chan []string, <-chan error) {
	outCh := make(chan []string)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)
		var (
			rec        []string
			fileReader *os.File
			err        error
		)
		if fileReader, err = os.Open(f.FileName); err != nil {
			errCh <- errors.Wrap(err, "could not open csv file")
			return
		}
		defer fileReader.Close()

		r := csv.NewReader(fileReader)
		r.Comma = f.Delimiter

		for {
			rec, err = r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- fmt.Errorf("error reading csv file: %w", err)
			}
			outCh <- rec
		}
	}()

	return outCh, errCh
}

func (f *CSVParser) detectDelimiter(filePath string) (rune, error) {
	var detectedDelimiter rune
	delimiters := []rune{',', ';', '\t', '|'}
	maxCount := 0
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return 0, errors.New("file is empty or unreadable")
	}
	line := scanner.Text()
	for _, delim := range delimiters {
		count := strings.Count(line, string(delim))
		if count > maxCount {
			maxCount = count
			detectedDelimiter = delim
		}
	}
	if detectedDelimiter == 0 {
		return 0, errors.New("could not detect delimiter")
	}
	return detectedDelimiter, nil
}
