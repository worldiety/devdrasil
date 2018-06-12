package group

import (
	"github.com/worldiety/devdrasil/db"
	"strings"
)

const TABLE_GROUP = "group"

//A group, attached to users. Use users
type Group struct {
	//Entity id, e.g. "xyz1234"
	Id db.PK

	//Name of the group, e.g. 'My Employees'
	Name string
}

type Groups struct {
	db   *db.Database
	crud *db.CRUD
}

func NewGroups(d *db.Database) *Groups {
	return &Groups{d, db.NewCRUD(d)}
}

func (r *Groups) List() ([]*Group, error) {
	tx := r.db.Partition(TABLE_GROUP).Begin(false)
	defer tx.Commit()
	return r.list(tx)
}

func (r *Groups) list(tx db.Transaction) ([]*Group, error) {
	res := make([]*Group, 0)
	err := r.crud.ListTX(tx, "", &res)
	return res, err
}

func (r *Groups) Add(group *Group) error {
	tx := r.db.Partition(TABLE_GROUP).Begin(true)
	defer tx.Commit()

	groups, err := r.List()
	if err != nil {
		return err
	}
	//ensure unique name
	myLowerCaseName := strings.ToLower(group.Name)
	for _, grp := range groups {
		if strings.ToLower(grp.Name) == myLowerCaseName {
			return &db.NotUnique{group.Name}
		}
	}

	group.Id = tx.NextKey()
	json := db.NewJSONDecorator(tx)
	return json.Put(group)
}

func (r *Groups) Update(group *Group) error {
	tx := r.db.Partition(TABLE_GROUP).Begin(true)
	defer tx.Commit()

	groups, err := r.List()
	if err != nil {
		return err
	}

	//ensure unique name
	myLowerCaseName := strings.ToLower(group.Name)

	//find other
	for _, grp := range groups {
		if strings.ToLower(grp.Name) == myLowerCaseName && group.Id != grp.Id {
			return &db.NotUnique{group.Name}
		}
	}
	//update user
	return r.crud.Update(TABLE_GROUP, group)
}

func (r *Groups) Delete(id db.PK) error {
	return r.crud.Delete(TABLE_GROUP, id)
}

func (r *Groups) Get(id db.PK) (*Group, error) {
	group := &Group{Id: id}
	err := r.crud.Read(TABLE_GROUP, group)
	return group, err
}
