package db

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
	"io"
	"fmt"
	"crypto/rand"
)

type readTransaction struct {
	partition *Partition
	alive     bool
	firstErr  error
}

func (tx *readTransaction) Err() error {
	return tx.firstErr
}

func (tx *readTransaction) noteErr(err error) error {
	tx.check()
	if err != nil && tx.firstErr == nil {
		tx.firstErr = err
	}
	return err
}

//creates a new read and aquires a read lock
func newReadTransaction(partition *Partition) *readTransaction {
	partition.rwLock.RLock()
	return &readTransaction{partition, true, nil}
}

func (tx *readTransaction) check() {
	if !tx.alive {
		panic("transaction invalid")
	}
}

//just unlocks the read lock
func (tx *readTransaction) Commit() error {
	tx.check()
	tx.alive = false
	tx.partition.rwLock.RUnlock();
	return nil;
}

//read transactions have nothing to rollback, just delegates to commit
func (tx *readTransaction) Rollback() error {
	return tx.Commit()
}

func (tx *readTransaction) Put(key PK, reader io.Reader) (int64, error) {
	AssertNotNIL(key)
	tx.check()
	return 0, fmt.Errorf("unsupported operation: readonly");
}

func (tx *readTransaction) Delete(key PK) error {
	tx.check()
	return fmt.Errorf("unsupported operation: readonly");
}

func (tx *readTransaction) Get(key PK, dst io.Writer) (int64, error) {
	tx.check()
	fname := tx.partition.fanout(key)
	file, err := os.OpenFile(fname, os.O_RDONLY, permOwnerOnly)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, &EntityNotFound{key}
		}
		tx.noteErr(err)
		return 0, err
	}
	defer file.Close()

	n, err := io.Copy(dst, file)
	tx.noteErr(err)
	return n, err

}

func (tx *readTransaction) Has(key PK) bool {
	tx.check()
	fname := tx.partition.parent.fanout(tx.partition.name, key)
	if _, err := os.Stat(fname); err != nil {
		return false
	}
	return true
}

func (tx *readTransaction) GetAll() *Cursor {
	tx.check()
	partionDir := filepath.Join(tx.partition.parent.dir, tx.partition.name)
	fanoutsFolders, e := ioutil.ReadDir(partionDir)
	if e != nil {
		return &Cursor{tx: tx, idx: -1}
	}
	entities := make([]string, 0)
	for _, folder := range fanoutsFolders {
		if folder.IsDir() && !strings.HasPrefix(folder.Name(), ".") {
			folderPath := filepath.Join(partionDir, folder.Name())
			files, e := ioutil.ReadDir(folderPath)
			if e == nil {
				for _, file := range files {
					if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
						entities = append(entities, filepath.Join(folderPath, file.Name()))
					}
				}
			}
		}
	}
	return &Cursor{tx: tx, files: entities, idx: -1}
}

func (tx *readTransaction) NextKey() PK {
	tx.check()
	var key PK
	n, e := rand.Read(key[:])
	if e != nil || n != len(key) {
		panic(e)
	}
	//generate recursively, very unlike to ever enter
	if (tx.Has(key)) {
		return tx.NextKey()
	}

	return key
}
