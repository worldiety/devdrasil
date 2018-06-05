package db

import (
	"sync"
)

type Partition struct {
	parent *Database
	name   string
	rwLock sync.RWMutex
}

func (p *Partition) fanout(key PK) string {
	return p.parent.fanout(p.name, key)
}

//begins a transaction, always ensure to commit or rollback the transaction
func (p *Partition) Begin(writeable bool) Transaction {
	if writeable {
		return newWriteTransaction(p)
	} else {
		return newReadTransaction(p)
	}
}

