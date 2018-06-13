package company

import (
	"github.com/worldiety/devdrasil/db"
	"strings"
)

const TABLE_COMPANY = "company"

//A company, attached to users. Use users to connect it.
type Company struct {
	//Entity id, e.g. "xyz1234"
	Id db.PK

	//Name of the company, e.g. 'My company'
	Name string

	//the primary color
	ThemePrimaryColor string
}

type Companies struct {
	db   *db.Database
	crud *db.CRUD
}

func NewCompanies(d *db.Database) *Companies {
	return &Companies{d, db.NewCRUD(d)}
}

func (r *Companies) List() ([]*Company, error) {
	tx := r.db.Partition(TABLE_COMPANY).Begin(false)
	defer tx.Commit()
	return r.list(tx)
}

func (r *Companies) list(tx db.Transaction) ([]*Company, error) {
	res := make([]*Company, 0)
	err := r.crud.ListTX(tx, "", &res)
	return res, err
}

func (r *Companies) Add(group *Company) error {
	tx := r.db.Partition(TABLE_COMPANY).Begin(true)
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

func (r *Companies) Update(group *Company) error {
	tx := r.db.Partition(TABLE_COMPANY).Begin(true)
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
	return r.crud.Update(TABLE_COMPANY, group)
}

func (r *Companies) Delete(id db.PK) error {
	return r.crud.Delete(TABLE_COMPANY, id)
}

func (r *Companies) Get(id db.PK) (*Company, error) {
	group := &Company{Id: id}
	err := r.crud.Read(TABLE_COMPANY, group)
	return group, err
}
