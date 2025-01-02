package driver

import (
	"bufio"
	"go_echo/internal/util/helper"
	"go_echo/internal/writer/data"
	"os"

	"github.com/pkg/errors"
)

type TXTWriter struct {
	FileName string
	wCount   int
	ff       *os.File
	fw       *bufio.Writer
}

func (f *TXTWriter) Init(options data.FileOptions) error {
	var err error
	f.FileName = options.FileName
	f.ff, err = os.Create(f.FileName)
	if err != nil {
		return errors.Wrap(err, "Unable to open txt file "+f.FileName)
	}
	f.fw = bufio.NewWriterSize(f.ff, 65536) //nolint:mnd // Buffer size is 64K
	return nil
}

func (f *TXTWriter) WriteAll(rec [][]interface{}) error {
	var err error
	for _, r := range rec {
		if err = f.Write(r); err != nil {
			return err
		}
	}
	return nil
}

func (f *TXTWriter) Write(rec []interface{}) error {
	var (
		err error
		res string
	)
	for _, r := range rec {
		res = helper.AnyToString(r)
		if _, err = f.fw.WriteString(res + "\n"); err != nil {
			return errors.Wrap(err, "could not write txt file"+f.FileName+", line "+res)
		}
	}
	if f.wCount == 1000 { //nolint:mnd // 1000 lines is flush size
		err = f.fw.Flush()
		if err != nil {
			return errors.Wrap(err, "could not flush txt file"+f.FileName)
		}
		f.wCount = 0
	}
	f.wCount++
	return nil
}

func (f *TXTWriter) Close() error {
	var err error
	err = f.fw.Flush()
	if err != nil {
		return errors.Wrap(err, "could not flush txt file"+f.FileName)
	}
	err = f.ff.Close()
	if err != nil {
		return errors.Wrap(err, "could not close txt file"+f.FileName)
	}
	f.fw = nil
	f.ff = nil
	f.wCount = 0
	return nil
}
