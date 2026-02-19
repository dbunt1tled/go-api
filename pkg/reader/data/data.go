package data

type FileParser interface {
	Read() (<-chan []string, <-chan error)
	Init(options FileOptions) error
}

type FileOptions struct {
	FileName  string
	SheetNum  *int
	Delimiter *rune
	Mapper    *Mapper
}
