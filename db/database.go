package db

import (
	"path/filepath"
	"encoding/hex"
	"fmt"
	"encoding/base64"
)

const permOwnerOnly = 0700

var NIL PK
/*
A very simple file based "database". Actually it is just a wrapper to handle local files (more) correctly. Note that
a local file system is already a very good hierarchical database but without much useful guarantees. This wrapper
provides some comfort functions, like serialized (pseudo)transactions, which guarantee non-racy read/writes, however
it does not provide ACID features.

It uses a single level of fanout to distribute the files evenly and to treat the local fs gracefully with up-to 1 million entries
per partition, which will result in around 4000 entries per directory, which is something reasonable. The actual
key is encoded as hex within the fanout. Empty keys are not supported.
 */
type Database struct {
	dir string
}

//open the database, performs no I/O. Do not share the same directory across multiple instances.
func Open(dir string) *Database {
	return &Database{dir}
}

//The primary key definition is a fixed length byte array
type PK [16]byte

func (p PK) IsNIL() bool {
	return p == NIL
}

//returns the base64 encoding of PK
func (p PK) String() string {
	return base64.StdEncoding.EncodeToString(p[:])
}

//decodes a base64 encoding, as returned by PK.String()
func ParsePK(base64str string) (PK, error) {
	tmp, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return NIL, err
	}

	var pk PK
	copy(pk[:], tmp)
	return pk, nil
}

//throws if bytes is longer than 16, rest is zero. This is only useful for hardcoded keys
func NewPK(str string) PK {
	if len(str) > 16 {
		panic("PK truncated")
	}
	var pk PK
	copy(pk[:], str)
	return pk
}

//throws if bytes is longer than 16, rest is zero. This is only useful for hardcoded keys
func NewPKFromArray(str []byte) PK {
	if len(str) > 16 {
		panic("PK truncated")
	}
	var pk PK
	copy(pk[:], str)
	return pk
}

//get a partition, performs no I/O.
func (d *Database) Partition(name string) *Partition {
	return &Partition{parent: d, name: name}
}

//creates a fanout of any byte sequence used as a PK. Beyond the fanout dir, the key is encoded as hex to be case insensitive to avoid ugly collisions
func (d *Database) fanout(partName string, key PK) string {
	strKey := hex.EncodeToString(key[:])
	fname := filepath.Join(d.dir, partName, strKey[0:2], strKey[2:])
	return fname
}

type EntityNotFound struct {
	What interface{}
}

func (e *EntityNotFound) Error() string {
	return "EntityNotFound: " + fmt.Sprintf("%v", e.What)
}

func IsEntityNotFound(err error) bool {
	_, ok := err.(*EntityNotFound)
	return ok
}

func IsNotUnique(err error) bool {
	_, ok := err.(*NotUnique)
	return ok
}

type NotUnique struct {
	What interface{}
}

func (e *NotUnique) Error() string {
	return "NotUnique: " + fmt.Sprintf("%v", e.What)
}

func AssertNotNIL(pk PK) {
	if pk == NIL {
		panic("pk may not be NIL [0,0,...]")
	}
}
