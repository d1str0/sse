# go-sse
Searchable symmetric encryption.


Overview
=

go-sse is a *mostly* database agnostic searchable symmetric encryption (SSE)
solution. The implementation is based off of the following paper:
http://www.internetsociety.org/doc/dynamic-searchable-encryption-very-large-databases-data-structures-and-implementation
by David Cash, Joseph Jaeger, Stanislaw Jarecki, Charanjit Jutla, Hugo Krawczyk,
Marcel-Catalin Ros, and Michael Steiner.


Implementation
=

Data Structures
-

Three persistent data structures are required for our implementation of SSE.

Document Table:
First, we will have a dictionary of documents mapped by a document ID. This
document ID will be a randomly generated string that must be unique.

Count Table:
Second, a second dictionary will be used to track the number of matching
documents for a given keyword. The key for this dictionary will be an encrypted
version of the keyword and the value will be the count (TODO: possible encrypt
the count).

Block Table (Index):
Third, a third dictionary will be used to track "blocks". The key for these
blocks will be created deterministically given several different variables.
These will be encrypted blobs containing an array to matching document IDs.


Constructing the Document Table
-

The document table is simply a key value store relating document IDs to document
data. The ID of the document is derived from the document by taking an
HMAC-SHA256 of the document data after encryption (Encrypt-then-MAC). The
encrypted document is then stored in the table with the coinciding key.

Using the standard library's 'crypto/hmac' package, we will generate the ID of
documents like so:

    // Where key and document are byte slices (arrays).
    mac := hmac.New(sha256.New, key)
    mac.Write(document)
    docId := mac.Sum(nil)

Documents will be encrypted with AES-256. Key's will be derived using PBKDF2
with a work factor or 20,000 + random number, and a unique salt. The salt and
work factor will be stored client side.


Constructing the Count Table
-

This table will hold the hashed keyword and a count of documents that match.

Keywords will be HMACd using the keyword and the key as the document IDs were
HMACd above.

    table[hash] = count


Constructing the Index
-

Also known as the Block Table, this will be another key value store with
seemingly random keys and encrypted blobs for values.

Each block will be holding an array of size B containing document IDs matching a
given keyword w.

From K, we derive a K1 and a K2.

    K1 = HMAC(K, 1||w)
    K2 = HMAC(K, 2||w)

    c = Count[w]
    i = floor(c/B) // 5 matches / 10 ids per block = .5 -> Block 0
    li = HMAC(K1, c) // Where c is the count in the array
    d = Enc(A, K2) // Where A is the array.
    Store[li] = d

Note: the entire array is stored as an encrypted blob. There is no ability to
tell that it is an array at this point.
