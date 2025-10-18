package driver

import (
	"github.com/dbunt1tled/go-api/internal/writer/data"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

type XLSXWriter struct {
	FileName string
	row      int
	fl       *excelize.File
	sw       *excelize.StreamWriter
}

func (f *XLSXWriter) Init(options data.FileOptions) error {
	var err error
	f.FileName = options.FileName
	f.row = 1
	f.fl = excelize.NewFile()
	f.sw, err = f.fl.NewStreamWriter("Sheet1")
	if err != nil {
		return errors.Wrap(err, "failed to create xlsx stream "+f.FileName)
	}
	return nil
}

func (f *XLSXWriter) WriteAll(rec [][]interface{}) error {
	var err error
	for _, r := range rec {
		if err = f.Write(r); err != nil {
			return err
		}
	}
	return nil
}

func (f *XLSXWriter) Write(rec []interface{}) error {
	var err error
	cell, _ := excelize.CoordinatesToCellName(1, f.row)
	if err = f.sw.SetRow(cell, rec); err != nil {
		return errors.Wrap(err, "failed to write row file "+f.FileName)
	}
	f.row++
	return nil
}

func (f *XLSXWriter) Close() error {
	var err error
	if err = f.sw.Flush(); err != nil {
		return errors.Wrap(err, "failed to flush xlsx stream "+f.FileName)
	}
	if err = f.fl.SaveAs(f.FileName); err != nil {
		return errors.Wrap(err, "failed to save file: "+f.FileName)
	}
	f.sw = nil
	f.fl = nil
	f.row = 1
	return nil
}
