package reader

import (
	"errors"
	"fmt"
	"go_echo/internal/reader/data"
	"go_echo/internal/reader/driver"
	"strings"
)

type FileReader struct {
	Options data.FileOptions
	Parsers map[string]data.Parser
	Parser  *data.Parser
	Mapper  *Mapper
}

func NewFileReader(options data.FileOptions) (*FileReader, error) {
	reader := &FileReader{
		Options: options,
		Parsers: map[string]data.Parser{
			"csv":  &driver.CSVParser{},
			"xlsx": &driver.XLSXParser{},
		},
	}
	if err := reader.setParser(); err != nil {
		return nil, err
	}

	return reader, nil
}

func (f *FileReader) setParser() error {
	ext := strings.ToLower(f.Options.FileName[strings.LastIndex(f.Options.FileName, ".")+1:])
	parser, exists := f.Parsers[ext]
	if !exists {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
	err := parser.Init(f.Options)
	if err != nil {
		return err
	}
	f.Parser = &parser
	return nil
}

func (f *FileReader) Read() (<-chan []string, <-chan error) {
	resOutCh := make(chan []string)
	resErrCh := make(chan error, 1)
	if f.Parser == nil {
		defer close(resOutCh)
		defer close(resErrCh)
		resErrCh <- errors.New("parser is not set")
		return resOutCh, resErrCh
	}

	go func() {
		defer close(resOutCh)
		defer close(resErrCh)
		iterator, errChan := (*f.Parser).Read()
		for row := range iterator {
			if err := ValidateRow(row); err != nil {
				// Логируем ошибку (если нужно)
				fmt.Println("Validation failed:", err)
				continue // Пропускаем невалидные строки
			}
			resOutCh <- row
		}

		// Передаём ошибки чтения из оригинального итератора
		if err := <-errChan; err != nil {
			resErrCh <- err
		}
	}()

	return resOutCh, resErrCh
}

func ValidateRow(row []string) error {
	if len(row) < 3 {
		return errors.New("row has less than 3 columns")
	}
	for _, field := range row {
		if len(strings.TrimSpace(field)) == 0 {
			return errors.New("row contains empty field")
		}
	}
	return nil
}
