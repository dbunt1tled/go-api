package reader

import (
	"encoding/json"
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
	Mapper  *data.Mapper
}

func NewFileReader(options data.FileOptions) (*FileReader, error) {
	reader := &FileReader{
		Options: options,
		Parsers: map[string]data.Parser{
			"txt":  &driver.TXTParser{},
			"csv":  &driver.CSVParser{},
			"xlsx": &driver.XLSXParser{},
		},
		Mapper: options.Mapper,
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

func (f *FileReader) ReadMap() (<-chan map[string]string, <-chan error) {
	resOutCh := make(chan map[string]string)
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
		isColumnsSet := false
		i := 0
		for row := range iterator {
			i++
			if f.Mapper == nil {
				resOutCh <- map[string]string{"row": strings.Join(row, "|")}
				continue
			}
			if !isColumnsSet {
				isColumnsSet = f.Mapper.SetColumns(row)
				if !isColumnsSet && i == 10 {
					resErrCh <- errors.New("columns not set")
					return
				}
				continue
			}
			f.Mapper.Values = data.SliceToSliceInterface(row)
			res := make(map[string]string)
			for key := range f.Mapper.MappedFields {
				k, e := f.Mapper.GetValue(key)
				if e != nil {
					resErrCh <- e
					return
				}
				res[key] = k
			}
			resOutCh <- res
		}
		if err := <-errChan; err != nil {
			resErrCh <- err
		}
	}()

	return resOutCh, resErrCh
}

func (f *FileReader) ReadLine() (<-chan string, <-chan error) {
	resOutCh := make(chan string)
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
		iterator, errChan := f.ReadMap()
		for row := range iterator {
			r, ok := row["row"]
			if !ok {
				b, err := json.Marshal(row)
				if err != nil {
					resErrCh <- err
					return
				}
				resOutCh <- string(b)
			} else {
				resOutCh <- r
			}
		}
		if err := <-errChan; err != nil {
			resErrCh <- err
		}
	}()
	return resOutCh, resErrCh
}
