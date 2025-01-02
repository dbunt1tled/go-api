package data

type FileWriter interface {
	WriteAll(rec [][]interface{}) error
	Write(rec []interface{}) error
	Init(options FileOptions) error
	Close() error
}

type FileOptions struct {
	FileName  string
	SheetNum  *int
	Delimiter *rune
}
