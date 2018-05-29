package db

import "fmt"

type Cursor struct {
	parent   *Partition
	files    []string
	idx      int
	firstErr error
}

//it is a total valid situation that (due to missing isolation) entities become inavailable while iterating
func (c *Cursor) Scan(obj interface{}) error {
	c.parent.mutex.Lock()
	defer c.parent.mutex.Unlock()
	if c.idx < 0 || c.idx >= len(c.files) {
		return fmt.Errorf("cursor is out of bound")
	}
	e := Read(c.files[c.idx], obj)
	if e != nil && c.firstErr == nil {
		c.firstErr = e
	}
	return e
}

func (c *Cursor) Next() bool {
	c.idx++
	return c.idx < len(c.files)
}

func (c *Cursor) Close() {

}

func (c *Cursor) Err() error {
	return c.firstErr
}

func (c *Cursor) Size() int {
	return len(c.files)
}
