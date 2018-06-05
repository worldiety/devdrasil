package user

import "github.com/worldiety/devdrasil/db"

//permissions for treating users
const TABLE_USER_PERMISSION = "user_permission"

var LIST_USERS = db.NewPK("LIST_USERS")
var CREATE_USER = db.NewPK("CREATE_USER")
var DELETE_USER = db.NewPK("DELETE_USER")
var UPDATE_USER = db.NewPK("UPDATE_USER")
var GET_USER = db.NewPK("GET_USER")

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
	ensureEntities := []db.PK{LIST_USERS, CREATE_USER, DELETE_USER, UPDATE_USER, GET_USER}
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
