package db

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

type Database struct {
	dir string
}

func Open(dir string) *Database {
	return &Database{dir}
}

type PK string

func (d *Database) Partition(name string) *Partition {
	return &Partition{parent: d, name: name}
}

//returns the id or panics
func GetId(obj interface{}) (string, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		return "", fmt.Errorf("must be a pointer: " + val.String())
	}
	field := val.Elem().FieldByName("Id")

	if field.Kind() != reflect.Invalid {
		if t, ok := field.Interface().(string); ok {
			return t, nil
		}
	}
	return "", fmt.Errorf("field 'Id' of type string is required for " + reflect.ValueOf(obj).Elem().String())
}

func SetId(obj interface{}, value string) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("must be a pointer: " + val.String())
	}

	field := val.Elem().FieldByName("Id")

	if field.Kind() != reflect.Invalid {
		if _, ok := field.Interface().(string); ok {
			field.SetString(value)
			return nil
		}
	}
	return fmt.Errorf("field 'Id' of type string is required for " + reflect.ValueOf(obj).Elem().String())
}

func Write(fname string, obj interface{}) error {
	b, e := json.Marshal(obj)
	if e != nil {
		return e
	}

	e = ioutil.WriteFile(fname, b, os.ModePerm)
	if e != nil {
		//fanout dir
		os.MkdirAll(filepath.Dir(fname), os.ModePerm)
		re := ioutil.WriteFile(fname, b, os.ModePerm)
		if re != nil {
			return e
		}
	}
	return nil
}

func Read(fname string, obj interface{}) error {
	b, e := ioutil.ReadFile(fname)
	if e != nil {
		return e
	}
	e = json.Unmarshal(b, obj)
	if e != nil {
		return e
	}
	return nil
}

func HashId(id string) PK {
	tmp := sha512.Sum512_224([]byte(id))
	return PK(hex.EncodeToString(tmp[:]))
}

//creates a fanout of any text used as a PK
func (d *Database) fanout(partName string, text string) string {
	id := string(HashId(text))
	fname := filepath.Join(d.dir, partName, id[0:2], id[2:])
	return fname
}
