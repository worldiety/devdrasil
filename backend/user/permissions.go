package user

import "github.com/worldiety/devdrasil/db"

//permissions for treating users
const TABLE_USER_PERMISSION = "user_permission"

var LIST_USERS = db.NewPK("LIST_USERS")
var CREATE_USER = db.NewPK("CREATE_USER")
var DELETE_USER = db.NewPK("DELETE_USER")
var UPDATE_USER = db.NewPK("UPDATE_USER")
var GET_USER = db.NewPK("GET_USER")

var INSTALL_PLUGIN = db.NewPK("INSTALL_PLUGIN")
var REMOVE_PLUGIN = db.NewPK("REMOVE_PLUGIN")
var LIST_MARKET = db.NewPK("LIST_MARKET")

var LIST_GROUPS = db.NewPK("LIST_GROUPS")
var CREATE_GROUP = db.NewPK("CREATE_GROUP")
var DELETE_GROUP = db.NewPK("DELETE_GROUP")
var UPDATE_GROUP = db.NewPK("UPDATE_GROUP")
var GET_GROUP = db.NewPK("GET_GROUP")

var LIST_COMPANIES = db.NewPK("LIST_COMPANIES")
var CREATE_COMPANY = db.NewPK("CREATE_COMPANY")
var DELETE_COMPANY = db.NewPK("DELETE_COMPANY")
var UPDATE_COMPANY = db.NewPK("UPDATE_COMPANY")
var GET_COMPANY = db.NewPK("GET_COMPANY")

type Permission struct {
	//unique entity id, e.g. "0xaccc32"
	Id db.PK

	//allowed groups
	AllowedGroups []db.PK

	//allowed users
	AllowedUsers []db.PK
}

//checks if either(!) the user or(!) the group is allowed. Do not simply pass an unkown id to both, because groups and users do not share the same ID space, collisions may false grant access!
func (p *Permission) IsAllowed(user *db.PK, group *db.PK) bool {
	if user != nil {
		id := *user
		for _, k := range p.AllowedUsers {
			if k == id {
				return true
			}
		}
	}

	if group != nil {
		id := *group
		for _, k := range p.AllowedGroups {
			if k == id {
				return true
			}
		}
	}

	return false
}

//the permissions repository
type Permissions struct {
	db   *db.Database
	crud *db.CRUD
}

func NewPermissions(d *db.Database) (*Permissions, error) {
	perms := &Permissions{d, db.NewCRUD(d)}

	tx := perms.db.Partition(TABLE_USER_PERMISSION).Begin(true)
	defer tx.Commit()
	json := db.NewJSONDecorator(tx)

	//ensure that at least for each permission, an empty entity is available
	ensureEntities := []db.PK{LIST_USERS, CREATE_USER, DELETE_USER, UPDATE_USER, GET_USER, INSTALL_PLUGIN, LIST_MARKET, LIST_GROUPS, CREATE_GROUP, DELETE_GROUP, UPDATE_GROUP, GET_GROUP, LIST_COMPANIES, CREATE_COMPANY, DELETE_COMPANY, UPDATE_COMPANY, GET_COMPANY}
	for _, id := range ensureEntities {
		if tx.Has(id) {
			continue
		}
		perm := &Permission{Id: id}
		//always include the admin in every permission
		perm.AllowedUsers = append(perm.AllowedUsers, ADMIN)
		err := json.Put(perm)
		if err != nil {
			return nil, err
		}
	}
	return perms, nil

}

func (r *Permissions) Get(kind db.PK) (*Permission, error) {
	perm := &Permission{Id: kind}
	err := r.crud.Read(TABLE_USER_PERMISSION, perm)
	return perm, err
}

func (r *Permissions) Update(perm *Permission) error {
	return r.crud.Update(TABLE_USER_PERMISSION, perm)
}

func (r *Permissions) IsAllowed(kind db.PK, user *User) (bool, error) {
	perm, err := r.Get(kind)
	if err != nil {
		return false, err
	}

	for _, k := range perm.AllowedUsers {
		if k == user.Id {
			return true, nil
		}
	}

	for _, allowedGroup := range perm.AllowedGroups {
		for _, userGroup := range user.Groups {
			if allowedGroup == userGroup {
				return true, nil
			}
		}
	}
	return false, nil
}
