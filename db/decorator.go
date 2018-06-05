package db

import (
	"log"
	"reflect"
	"fmt"
	"encoding/json"
	"bytes"
	"sort"
)

type JSONCursor struct {
	cursor *Cursor
}

type JSONDecorator struct {
	//may be nil
	wTx *writeTransaction
	//always available
	rTx *readTransaction
}

func NewJSONDecorator(tx Transaction) *JSONDecorator {

	d := &JSONDecorator{}
	switch t := tx.(type) {
	case *readTransaction:
		d.rTx = t
		d.wTx = nil
	case *writeTransaction:
		d.wTx = t
		d.rTx = t.reader
	default:
		panic(t)
	}
	return d
}

/*
Supported format of query is

''

ORDER BY <field>

ORDER BY <field> ASC

ORDER BY <field> DESC

*/
func (p *JSONDecorator) Query(query string) *JSONCursor {
	if len(query) == 0 {
		return &JSONCursor{p.rTx.GetAll()}
	}
	q := parse(query)
	if q.orderByField == "" {
		return &JSONCursor{p.rTx.GetAll()}
	}

	tmp := make([]genericJson, 0)
	tmpCursor := p.rTx.GetAll()

	//naive approach is to simply load everything into memory
	buf := &bytes.Buffer{}
	for tmpCursor.Next() {
		jsonMap := make(map[string]interface{})
		buf.Reset()
		_, e := tmpCursor.Get(buf)
		if e != nil {
			log.Println(e)
			continue
		}
		e = json.Unmarshal(buf.Bytes(), jsonMap)
		if e != nil {
			log.Println(e)
			continue
		}
		tmp = append(tmp, genericJson{tmpCursor.files[tmpCursor.idx], jsonMap})
	}

	asc := q.orderDir == "ASC"
	sort.Sort(&byCustomField{tmp, q.orderByField, asc})

	entities := make([]string, len(tmp))
	for i, g := range tmp {
		entities[i] = g.fname
	}
	return &JSONCursor{&Cursor{tx: p.rTx, files: entities, idx: -1}}
}

func (p *JSONDecorator) Put(obj interface{}) error {
	id, err := getId(obj)
	if err != nil {
		p.rTx.noteErr(err)
		return err
	}

	if p.wTx == nil {
		_, err = p.rTx.Put(id, nil)
		if err != nil {
			return err
		}
	}
	b, e := json.Marshal(obj)
	if e != nil {
		return e
	}

	_, e = p.wTx.Put(id, bytes.NewReader(b))
	if e != nil {
		return e
	}
	return nil
}

func (p *JSONDecorator) Get(obj interface{}) error {
	id, err := getId(obj)
	if err != nil {
		p.rTx.noteErr(err)
		return err
	}

	buf := &bytes.Buffer{}

	_, err = p.rTx.Get(id, buf)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf.Bytes(), obj)
	p.rTx.noteErr(err)
	return err
}

//returns the id or errors
func getId(obj interface{}) (PK, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		return NIL, fmt.Errorf("must be a pointer: " + val.String())
	}
	field := val.Elem().FieldByName("Id")

	if field.Kind() != reflect.Invalid {
		if t, ok := field.Interface().(PK); ok {
			return t, nil
		}
	}

	return NIL, fmt.Errorf("field 'Id' of type PK is required for " + reflect.ValueOf(obj).Elem().String())
}

func SetId(obj interface{}, value PK) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("must be a pointer: " + val.String())
	}
	field := val.Elem().FieldByName("Id")
	if field.Kind() != reflect.Invalid {
		if _, ok := field.Interface().(PK); ok {
			v := reflect.ValueOf(value)
			field.Set(v)
			return nil
		}
	}
	return fmt.Errorf("field 'Id' of type db.PK is required for " + reflect.ValueOf(obj).Elem().String())
}

func (c *JSONCursor) Next() bool {
	return c.cursor.Next()
}

func (c *JSONCursor) Read(obj interface{}) error {
	buf := &bytes.Buffer{}
	_, err := c.cursor.Get(buf)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf.Bytes(), obj)
	c.cursor.noteErr(err)
	return err
}
