package driver

import (
	"fmt"
	"go_echo/internal/reader/data"

	"github.com/pkg/errors"
	"github.com/thedatashed/xlsxreader"
)

type XLSXParser struct {
	FileName string
	SheetNum int
}

func (f *XLSXParser) Init(options data.FileOptions) error {
	f.FileName = options.FileName
	f.SheetNum = 0
	if options.SheetNum != nil {
		f.SheetNum = *options.SheetNum
	}
	return nil
}

func (f *XLSXParser) Read() (<-chan []string, <-chan error) {
	outCh := make(chan []string)
	errCh := make(chan error, 1)
	go func() {
		defer close(outCh)
		defer close(errCh)
		var (
			i          int
			rec        []string = make([]string, 0)
			fileReader *xlsxreader.XlsxFileCloser
			err        error
		)

		if fileReader, err = xlsxreader.OpenFile(f.FileName); err != nil {
			errCh <- errors.Wrap(err, "could not open xlsx file")
			return
		}
		defer fileReader.Close()
		for row := range fileReader.ReadRows(fileReader.Sheets[f.SheetNum]) {
			if row.Error != nil {
				errCh <- fmt.Errorf("error reading xlsx file: %w", err)
				return
			}
			rec = make([]string, cap(row.Cells))
			for _, cell := range row.Cells {
				i = cell.ColumnIndex()
				rec[i] = cell.Value
			}
			outCh <- rec
		}
	}()

	return outCh, errCh
}
