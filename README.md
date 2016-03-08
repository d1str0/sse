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

Crypto 
-

    AES-CBC + HMAC

    MasterKey = PBKDF2(password, salt, iter, size)
    AES-Key = HMAC(MasterKey, "AES-Key" | 0x01)
    MAC-Key = HMAC(MasterKey, "MAC-Key" | 0x01)


Data Structures
-

Three persistent data structures are required for our implementation of SSE.

Document Dictionary:
We will have a dictionary of documents mapped by document ID. The document ID is
a procedurally generated string that will be unique.

Count Dictionary:
This dictionary will be used to track the number of matching documents for a
given keyword. The key for this dictionary will be an encrypted version of the
keyword and the value will be the count (TODO: possible encrypt the count).

Index Dictionary:
This disctionary is slightly more complex than the previous two as the value
held by the key will be an encrypted array. Each of these encrypted arrays are
also referred to as blocks and will be an encrypted blob to be stored in the
dictionary. The key for this dictionary will be an HMAC of several different
variables.


Constructing the Document Dictionary
-

The document dictionary is simply a key value store relating document IDs to
document data. The ID of the document is derived from the document by taking an
HMAC-SHA256 of the document data after encryption (Encrypt-then-MAC). The
encrypted document is then stored in the table with the coinciding key.

Using the standard library's 'crypto/hmac' package, we will generate the ID of
documents like so:

    // Where key and document are byte slices (arrays).
    mac := hmac.New(sha256.New, key)
    mac.Write(document)
    docId := mac.Sum(nil)

Documents will be encrypted with AES-256. Key's will be derived using PBKDF2
with a work factor of 20,000 + random number, and a unique salt. The salt and
work factor will be stored client side.


Constructing the Count Dictionary
-

This dictionary will hold the hashed keyword and a count of documents that match.

The indice of this table is computed as the HMAC of the keyword and the key.

    mac := hmac.New(sha256.New, key)
    mac.Write(keyword)
    hash := mac.Sum(nil)
    table[hash] = count


Constructing the Index
-

Also known as the Block Dictionary, this will be another key value store with
seemingly random keys and encrypted blobs for values. These encrypted blobs are
also referred to as blocks.

Each block is an array of size B containing document IDs matching a given
keyword w.

From K, we derive a K1 and a K2.

    K1 = HMAC(K, 1||w)
    K2 = HMAC(K, 2||w)

These two keys are used later to derive the block ID and to encrypt the block
respectively.

    c = Count[w]
    i = floor(c/B) // 5 matches / 10 ids per block = .5 -> Block 0
    l = HMAC(K1, c) // Where c is the count in the array
    d = Enc(A, K2) // Where A is the array.
    Store[l] = d

Note: the entire array is stored as an encrypted blob. There is no ability to
tell that it is an array at this point.

Installation
=

Requirements
-

This package requires github.com/d1str0/pkcs7 for PKCS#7 Padding.
