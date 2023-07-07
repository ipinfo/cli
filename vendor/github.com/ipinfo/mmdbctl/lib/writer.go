package lib

type writer interface {
	Write(record []string) error
	Flush()
	Error() error
}
