package sse

import (
	"encoding/json"
	"math"
	"strconv"
)

type Client struct {
	DB  DBConn
	key []byte
}

const (
	BlobSize = 10 // Size of array holding doc IDs
)

var One = []byte{0x01} // Needed for different keys.
var Two = []byte{0x02}

func NewClient(db DBConn) *Client {
	db.Init()
	return &Client{DB: db}
}

// Get a document back from the store with the given ID.
func (c *Client) Get(id string) ([]byte, error) {
	// Get the encrypted doc from the database.
	edoc, err := c.DB.Get(DOCUMENTS, []byte(id))

	// Decrypt the doc with the client key.
	doc, err := Decrypt(edoc, c.key)
	return doc, err
}

// Put a document into the store with the given ID.
func (c *Client) Put(id string, doc []byte) error {
	// Let's encrypt the document with the client key.
	edoc, err := Encrypt(doc, c.key)
	if err != nil {
		return err
	}

	// Now that it's encrypted, store it in the DB.
	err = c.DB.Put(DOCUMENTS, []byte(id), edoc)
	return err
}

// Get the current count value for a given keyword.
func (c *Client) Count(keyword string) (int, error) {
	// HMAC the keyword.
	ekey := HMAC([]byte(keyword), c.key)

	count, err := c.DB.Get(COUNTS, ekey)
	if err != nil {
		return 0, err
	}
	var i int
	if len(count) == 0 {
		i = 0
		err := c.SetCount(keyword, i)
		if err != nil {
			return 0, err
		}
	} else {
		i, err = strconv.Atoi(string(count))
		if err != nil {
			return 0, err
		}

	}

	return i, nil
}

func (c *Client) SetCount(keyword string, count int) error {
	// HMAC the keyword.
	ekey := HMAC([]byte(keyword), c.key)

	i := strconv.Itoa(count)
	err := c.DB.Put(COUNTS, ekey, []byte(i))
	return err
}

// Find all document IDs associated with the given keyword.
func (c *Client) Search(keyword string) (ids []string, err error) {
	kw1 := append([]byte(keyword), One[:]...) // We append a constant to each keyword before MAC
	kw2 := append([]byte(keyword), Two[:]...)
	k1 := HMAC(kw1, c.key) // Generate two separate keys.
	k2 := HMAC(kw2, c.key)

	// Get count
	count, err := c.Count(keyword)
	if err != nil {
		return
	}

	max := int(math.Floor(float64(count) / float64(BlobSize)))
	for i := 0; i <= max; i++ {

		// Generate the id of this block using k1.
		h := HMAC(append([]byte("COUNT"), byte(i)), k1)

		// Get the encrypted blob.
		ejson, err2 := c.DB.Get(INDEX, h)
		if err2 != nil {
			err = err2
			return
		}

		// Decrypt the blob with k2.
		pjson, err2 := Decrypt(ejson, k2)
		if err2 != nil {
			err = err2
			return
		}

		var block []string
		err = json.Unmarshal(pjson, &block)
		if err != nil {
			return
		}

		ids = append(ids, block...)
	}

	return
}

/*
// TODO: Use an argument list here or not? (ie. keywords ...string)
func (c *Client) AddDocsToKeyword(keyword string, docs []string) error {

}

// TODO: Use an argument list here or not? (ie. keywords ...string)
func (c *Client) AddKeywordsToDoc(doc string, keywords []string) error {

}
*/

func (c *Client) AddDocToKeyword(keyword, doc string) error {
	kw1 := append([]byte(keyword), One[:]...) // We append a constant to each keyword before MAC
	kw2 := append([]byte(keyword), Two[:]...)
	k1 := HMAC(kw1, c.key) // Generate two separate keys.
	k2 := HMAC(kw2, c.key)

	// Get count
	count, err := c.Count(keyword)
	if err != nil {
		return err
	}

	max := int(math.Floor(float64(count) / float64(BlobSize)))

	// Generate the id of this block using k1.
	h := HMAC(append([]byte("COUNT"), byte(max)), k1)

	// Get the encrypted blob.
	ejson, err := c.DB.Get(INDEX, h)
	if err != nil {
		return err
	}

	var block []string
	var pjson []byte
	if len(ejson) > 0 {
		// Decrypt the blob with k2.
		pjson, err = Decrypt(ejson, k2)
		if err != nil {
			return err
		}
		err = json.Unmarshal(pjson, &block)
		if err != nil {
			return err
		}

	}

	// If we are overflowing this block
	if len(block) > 10 {
		block = make([]string, 10)
		max := max + 1
		h = HMAC(append([]byte("COUNT"), byte(max)), k1)
	}
	block = append(block, doc)

	// Get the json for the array
	pjson, err = json.Marshal(block)
	if err != nil {
		return err
	}

	ejson, err = Encrypt(pjson, k2)
	if err != nil {
		return err
	}

	err = c.DB.Put(INDEX, h, ejson)
	if err != nil {
		return err
	}

	c.SetCount(keyword, count+1)
	return err

}

// Set the key for the client.
func (c *Client) SetKey(passphrase, salt string, iter int) {
	c.key = Key([]byte(passphrase), []byte(salt), iter)
}
