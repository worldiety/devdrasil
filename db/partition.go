package db

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Partition struct {
	parent *Database
	name   string
	mutex  sync.Mutex
}

func (p *Partition) Put(obj interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	id, e := GetId(obj)
	if e != nil {
		panic(e)
	}
	fname := p.parent.fanout(p.name, id)
	return Write(fname, obj)
}

func (p *Partition) Get(key string, obj interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	fname := p.parent.fanout(p.name, key)
	return Read(fname, obj)
}

func (p *Partition) Delete(key string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	fname := p.parent.fanout(p.name, key)
	return os.Remove(fname)
}

func (p *Partition) GetAll() *Cursor {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	partionDir := filepath.Join(p.parent.dir, p.name)
	fanoutsFolders, e := ioutil.ReadDir(partionDir)
	if e != nil {
		return &Cursor{}
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
	return &Cursor{parent: p, files: entities, idx: -1}
}

/*
Supported format of query is
ORDER BY <field>
ORDER BY <field> ASC
ORDER BY <field> DESC
*/
func (p *Partition) Query(query string) *Cursor {
	q := parse(query)
	if q.orderByField == ""  {
		return p.GetAll()
	}

	tmp := make([]genericJson, 0)
	tmpCursor := p.GetAll()

	for tmpCursor.Next() {
		json := make(map[string]interface{})
		e := tmpCursor.Scan(&json)
		if e != nil {
			log.Println(e)
			continue
		}
		tmp = append(tmp, genericJson{tmpCursor.files[tmpCursor.idx], json})
	}

	asc := q.orderDir == "ASC"
	sort.Sort(&byCustomField{tmp, q.orderByField, asc})

	entities := make([]string, len(tmp))
	for i, g := range tmp {
		entities[i] = g.fname
	}
	return &Cursor{parent: p, files: entities, idx: -1}
}
