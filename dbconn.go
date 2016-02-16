package sse

type DBConn interface {
	Init() error
	Get(table, id string) ([]byte, error)
	Put(table, id string, value []byte) error
	Delete(table, id string) error
	Close()
}
