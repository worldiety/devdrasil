package db

import (
	"io"
	"fmt"
	"path/filepath"
	"os"
	"encoding/hex"
)

type Transaction interface {
	//Commits the transaction
	Commit() error

	//rolls the transaction back
	Rollback() error

	//Transfers the given bytes from reader into the partition
	Put(key PK, src io.Reader) (int64, error)

	//Reads given bytes from key into the given writer
	Get(key PK, dst io.Writer) (int64, error)

	//checks if an entry addressed with the given key exists. Failures are interpreted as false
	Has(key PK) bool

	//deletes the entry addressed by the given key. Deleting a non-existing entry is not considered as a failure.
	Delete(key PK) error

	//returns a cursor to read one entry after the other
	GetAll() *Cursor

	//returns the first error occured, which may be caused by any I/O failure. Reading a file which does not exist, is not considered an Error, so that is not tracked.
	Err() error

	//generates a new secure key, which is guaranteed not to collide with any existing key. It is 16 bytes long.
	NextKey() PK
}

type Cursor struct {
	tx    *readTransaction
	files []string
	idx   int
}

func (c *Cursor) check() {
	c.tx.check()
}
func (c *Cursor) checkPos() error {
	if c.idx < 0 || c.idx >= len(c.files) {
		return fmt.Errorf("cursor is out of bound")
	}
	return nil
}

//returns the key of the current cursor position.
func (c *Cursor) Key() (PK, error) {
	c.check()
	err := c.checkPos()
	if err != nil {
		return NIL, err
	}
	fanoutHex := filepath.Base(filepath.Dir(c.files[c.idx]))
	key, err := hex.DecodeString(fanoutHex + filepath.Base(c.files[c.idx]))
	if err != nil {
		return NIL, err
	}
	return NewPKFromArray(key), nil

}

//reads the entry at the current cursor position and returns the amount of transferred bytes
func (c *Cursor) Get(dst io.Reader) (int64, error) {
	c.check()
	err := c.checkPos()
	if err != nil {
		return 0, err
	}
	file, err := os.OpenFile(c.files[c.idx], os.O_RDONLY, permOwnerOnly)
	if err != nil {
		if os.IsNotExist(err) {
			key, kErr := c.Key()
			if kErr != nil {
				panic(kErr)
			}
			//here we track also NotExist as error, because our transaction isolation should always protect us
			c.noteErr(err)
			return 0, &EntityNotFound{key}
		}
		c.noteErr(err)
		return 0, err
	}
	defer file.Close()

	n, err := io.Copy(file, dst)
	c.noteErr(err)
	return n, err
}

//returns the amount of bytes at the current cursor position
func (c *Cursor) Length() (int64, error) {
	c.check()
	err := c.checkPos()
	if err != nil {
		return 0, err
	}
	stat, err := os.Stat(c.files[c.idx])
	c.noteErr(err)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func (c *Cursor) noteErr(err error) error {
	c.check()
	return c.tx.noteErr(err)
}

func (c *Cursor) Next() bool {
	c.check()
	c.idx++
	return c.idx < len(c.files)
}

//the transaction will also close the cursor
func (c *Cursor) Close() {
	c.check()
	//no-op
}

//returns the first error
func (c *Cursor) Err() error {
	c.check()
	return c.tx.firstErr
}

//returns the amount of entries
func (c *Cursor) Size() int {
	c.check()
	return len(c.files)
}
