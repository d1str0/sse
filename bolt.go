package main

import (
	"errors"
	"github.com/boltdb/bolt"
)

const (
	DOCUMENTS = "documents"
	COUNTS    = "counts"
	INDEX     = "index"
)

type BoltDB struct {
	Conn *bolt.DB
}

func BoltDBOpen() (*BoltDB, error) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	conn := &BoltDB{Conn: db}
	return conn, nil
}

func (db *BoltDB) Init() error {
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DOCUMENTS))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(COUNTS))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(INDEX))
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (db *BoltDB) Get(bucket, id string) ([]byte, error) {
	var value []byte
	// Use View() to enforce read-only access to BoltDB.
	err := db.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("BoltDB: bucket not found")
		}
		temp := b.Get([]byte(id))
		value = make([]byte, len(temp))
		copy(value, temp)
		return nil
	})

	return value, err
}

func (db *BoltDB) Put(bucket, id string, value []byte) error {
	// Use Update() to enforce read-write access to BoltDB.
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		// TODO: Put logic
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("BoltDB: bucket not found")
		}
		err := b.Put([]byte(id), value)
		return err
	})

	return err
}

func (db *BoltDB) Delete(bucket, id string) error {
	// Use Update() to enforce read-write access to BoltDB.
	err := db.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("BoltDB: bucket not found")
		}
		err := b.Delete([]byte(id))
		return err
	})

	return err
}

func (db *BoltDB) Close() {
	db.Conn.Close()
}