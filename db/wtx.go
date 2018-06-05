package db

import (
	"io"
	"fmt"
	"os"
	"path/filepath"
)

type writeTransaction struct {
	partition *Partition
	reader    *readTransaction
}

func (tx *writeTransaction) Err() error {
	return tx.reader.Err()
}

//creates a new write and aquires a write lock
func newWriteTransaction(partition *Partition) *writeTransaction {
	partition.rwLock.Lock()
	return &writeTransaction{partition, &readTransaction{partition, true, nil}}
}

//just releases the write lock
func (tx *writeTransaction) Commit() error {
	tx.reader.check()
	tx.reader.alive = false
	tx.partition.rwLock.Unlock();
	return nil;
}

func (tx *writeTransaction) Rollback() error {
	tx.reader.check()
	return fmt.Errorf("implementation does not support rollback")
}

//writes into a temporary file and moves afterwards to avoid truncated files in case of process crashes - the rest is up to the filesystem
func (tx *writeTransaction) Put(key PK, src io.Reader) (int64, error) {
	tx.reader.check()
	fname := tx.partition.fanout(key)

	//write into tmp, do not overwrite
	tmp := fname + ".tmp"
	file, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY, permOwnerOnly);
	if err != nil {
		//retry by creating the parent dirs
		_ = os.MkdirAll(filepath.Dir(tmp), permOwnerOnly)
		f2, err2 := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY, permOwnerOnly);
		if err2 != nil {
			tx.reader.noteErr(err)
			return 0, err
		}
		//second retry worked
		file = f2
	}
	defer file.Close()

	n, err := io.Copy(file, src)

	if err != nil {
		tx.reader.noteErr(err)
		return n, err
	}

	//delete target file
	err = os.Remove(fname)
	if !os.IsNotExist(err) {
		tx.reader.noteErr(err)
		return n, err
	}

	//rename to target file
	err = os.Rename(tmp, fname)
	if err != nil {
		tx.reader.noteErr(err)
		return n, err
	}
	return n, nil
}

func (tx *writeTransaction) Get(key PK, dst io.Writer) (int64, error) {
	return tx.reader.Get(key, dst)
}

func (tx *writeTransaction) Has(key PK) bool {
	return tx.reader.Has(key)
}

func (tx *writeTransaction) Delete(key PK) error {
	tx.reader.check()
	fname := tx.partition.fanout(key)
	err := os.Remove(fname)
	if err != nil {
		if _, err2 := os.Stat(fname); err2 != nil {
			if os.IsNotExist(err2) {
				return nil
			} else {
				//returns the original err, probably a permission problem
				return err
			}
		}
	}
	return nil
}

func (tx *writeTransaction) GetAll() *Cursor {
	return tx.reader.GetAll()
}

func (tx *writeTransaction) NextKey() PK {
	return tx.reader.NextKey()
}
