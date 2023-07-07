package lib

type reader interface {
	Read() (record []string, err error)
}
