package sse

const (
	DOCUMENTS = "documents"
	COUNTS    = "counts"
	INDEX     = "index"
)

type DBConn interface {
	Init() error
	Get(table string, id []byte) ([]byte, error)
	Put(table string, id, value []byte) error
	Delete(table string, id []byte) error
	Close()
}
