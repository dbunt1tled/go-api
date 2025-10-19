package driver

import (
	"encoding/csv"
	"os"

	"github.com/dbunt1tled/go-api/internal/util/helper"
	"github.com/dbunt1tled/go-api/internal/writer/data"

	"github.com/pkg/errors"
)

type CSVWriter struct {
	FileName  string
	Delimiter rune
	wCount    int
	ff        *os.File
	fw        *csv.Writer
}

func (f *CSVWriter) Init(options data.FileOptions) error {
	var err error
	f.FileName = options.FileName
	f.Delimiter = ','
	if options.Delimiter != nil {
		f.Delimiter = *options.Delimiter
	}
	f.ff, err = os.Create(f.FileName)
	if err != nil {
		return errors.Wrap(err, "Unable to open csv file "+f.FileName)
	}
	f.fw = csv.NewWriter(f.ff)

	return nil
}

func (f *CSVWriter) WriteAll(rec [][]interface{}) error {
	var err error
	for _, r := range rec {
		if err = f.Write(r); err != nil {
			return err
		}
	}
	f.fw.Flush()
	return nil
}

func (f *CSVWriter) Write(rec []interface{}) error {
	var err error
	var res []string
	for _, r := range rec {
		res = append(res, helper.AnyToString(r))
	}
	if err = f.fw.Write(res); err != nil {
		return errors.Wrap(err, "unable to write csv file "+f.FileName)
	}
	if f.wCount == 1000 { //nolint:mnd // 1000 lines is flush size
		f.fw.Flush()
		f.wCount = 0
	}
	f.wCount++
	return nil
}

func (f *CSVWriter) Close() error {
	f.fw.Flush()
	err := f.ff.Close()
	if err != nil {
		return errors.Wrap(err, "could not close csv file"+f.FileName)
	}
	f.fw = nil
	f.ff = nil
	return nil
}
