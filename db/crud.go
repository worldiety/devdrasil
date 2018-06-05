package db

import (
	"reflect"
	"fmt"
)

//Create Read Update Delete helper class to avoid boilerplate code using composition
type CRUD struct {
	db *Database
}

func NewCRUD(db *Database) *CRUD {
	return &CRUD{db}
}

//Reads the given object from the partition
func (c *CRUD) Read(partition string, obj interface{}) error {
	tx := c.db.Partition(partition).Begin(false)
	defer tx.Commit()

	return c.ReadTX(tx, obj)
}

func (c *CRUD) ReadTX(tx Transaction, obj interface{}) error {
	json := NewJSONDecorator(tx)

	err := json.Get(obj)
	if err != nil {
		return err
	}
	return nil
}

//Generates a new unique id and writes it into the partition
func (c *CRUD) Create(partition string, obj interface{}) error {
	tx := c.db.Partition(partition).Begin(true)
	defer tx.Commit()

	return c.CreateTX(tx, obj)
}

func (c *CRUD) CreateTX(tx Transaction, obj interface{}) error {
	err := SetId(obj, tx.NextKey())
	if err != nil {
		return err
	}
	json := NewJSONDecorator(tx)
	return json.Put(obj)
}

//Updates the entity in the partition.
func (c *CRUD) Update(partition string, obj interface{}) error {
	tx := c.db.Partition(partition).Begin(true)
	defer tx.Commit()

	return c.UpdateTX(tx, obj)
}

func (c *CRUD) UpdateTX(tx Transaction, obj interface{}) error {
	json := NewJSONDecorator(tx)
	return json.Put(obj)
}

//Delete the given key. Ignores not existing entries
func (c *CRUD) Delete(partition string, key PK) error {
	tx := c.db.Partition(partition).Begin(true)
	defer tx.Commit()

	return tx.Delete(key)
}

//Has convenience method
func (c *CRUD) Has(partition string, key PK) bool {
	tx := c.db.Partition(partition).Begin(false)
	defer tx.Commit()

	return tx.Has(key)
}

/*
Loads all entities from the partition into the target slice, e.g.

var arr := make([]*MyEntity, 0)

crud.List("myTable",&arr)
 */
func (c *CRUD) List(partition string, query string, v interface{}) error {
	tx := c.db.Partition(partition).Begin(false)
	defer tx.Commit()

	return c.ListTX(tx, query, v)
}

func (c *CRUD) ListTX(tx Transaction, query string, v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	} else {
		return fmt.Errorf("input param '%v' is not a slice", v)
	}

	sl := reflect.ValueOf(v)

	if t.Kind() == reflect.Ptr {
		sl = sl.Elem()
	}

	st := sl.Type()

	sliceType := st.Elem()
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}

	json := NewJSONDecorator(tx)
	cursor := json.Query(query)
	for cursor.Next() {
		newItem := reflect.New(sliceType)
		err := cursor.Read(newItem.Interface())
		if err != nil {
			return err
		}
		sl.Set(reflect.Append(sl, newItem))
	}
	return tx.Err()
}
