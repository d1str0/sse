package sse

type DBConn interface {
	Init() error
	Get(table string, id []byte) ([]byte, error)
	Put(table string, id, value []byte) error
	Delete(table string, id []byte) error
	Close()
}
