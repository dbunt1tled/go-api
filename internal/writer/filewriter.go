package writer

import (
	"fmt"
	"go_echo/internal/writer/data"
	"go_echo/internal/writer/driver"
	"strings"
)

type FileWriter struct {
	Options data.FileOptions
	Writers map[string]data.FileWriter
	Writer  *data.FileWriter
}

func NewFileWriter(options data.FileOptions) (*FileWriter, error) {
	writer := &FileWriter{
		Options: options,
		Writers: map[string]data.FileWriter{
			"xlsx": &driver.XLSXWriter{},
			"csv":  &driver.CSVWriter{},
			"txt":  &driver.TXTWriter{},
		},
	}
	if err := writer.setWriter(); err != nil {
		return nil, err
	}

	return writer, nil
}

func (f *FileWriter) setWriter() error {
	ext := strings.ToLower(f.Options.FileName[strings.LastIndex(f.Options.FileName, ".")+1:])
	writer, exists := f.Writers[ext]
	if !exists {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
	err := writer.Init(f.Options)
	if err != nil {
		return err
	}
	f.Writer = &writer
	return nil
}

func (f *FileWriter) Write(rec []any) error {
	return (*f.Writer).Write(rec)
}

func (f *FileWriter) WriteAll(rec [][]any) error {
	return (*f.Writer).WriteAll(rec)
}

func (f *FileWriter) Close() error {
	err := (*f.Writer).Close()
	f.Writer = nil
	return err
}
