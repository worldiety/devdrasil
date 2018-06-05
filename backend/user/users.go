package user

import (
	"github.com/worldiety/devdrasil/db"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

//Contains User objects as json
const TABLE_USER = "user"

//Contains binary streams of jpg or png files
const TABLE_AVATAR = "avatar"

//the admin key is hardcoded
var ADMIN = db.NewPK("admin")

const ADMIN_LOGIN = "admin"
const ADMIN_PWD = "admin"

type User struct {
	//unique entity id, e.g. "abc38293"
	Id db.PK

	//Abbreviation and/or Login, something like "tschinke", used for the login
	Login string

	//e.g. Torben
	Firstname string

	//e.g. Schinke
	Lastname string

	//and the additional password hash, bcrypt
	PasswordHash []byte

	//flag if user is active or not, without deleting him
	Active bool

	//reference to an optional avatar image
	AvatarImage *db.PK

	//list of connected email addresses e.g. tschinke@domain.com, torben.schinke@otherdomain.com, ...
	EMailAddresses []string

	//reference to an optional company
	Company *db.PK

	//the groups, which this user is a member of. This determines his actual permissions.
	Groups []db.PK
}

//Sets the password hash by calculating a bcrypt hash
func (u *User) SetPassword(pwd string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.PasswordHash = hash
}

//Compares the password hash with the given plaintext
func (u *User) PasswordEquals(plainTextPasswordToCompare string) bool {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(plainTextPasswordToCompare))
	if err != nil {
		return false
	}
	return true
}

//the users repository
type Users struct {
	db   *db.Database
	crud *db.CRUD
}

func NewUsers(d *db.Database) (*Users, error) {
	users := &Users{d, db.NewCRUD(d)}
	tx := users.db.Partition(TABLE_USER).Begin(true)
	defer tx.Commit()

	_, err := users.Get(ADMIN)
	if err != nil {
		if db.IsEntityNotFound(err) {
			//insert default configuration
			adminUser := &User{Id: ADMIN, Login: ADMIN_LOGIN, Active: true}
			adminUser.SetPassword(ADMIN_PWD)
			err = users.Add(adminUser)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return users, nil
}

func (r *Users) List() ([]*User, error) {
	res := make([]*User, 0)
	err := r.crud.List(TABLE_USER, "", &res)
	return res, err
}

func (r *Users) Get(id db.PK) (*User, error) {
	user := &User{Id: id}
	err := r.crud.Read(TABLE_USER, user)
	return user, err
}

func (r *Users) Delete(id db.PK) error {
	return r.crud.Delete(TABLE_USER, id)
}

func (r *Users) Add(user *User) error {
	tx := r.db.Partition(TABLE_USER).Begin(true)
	defer tx.Commit()
	res := make([]*User, 0)
	err := r.crud.List(TABLE_USER, "", &res)
	if res != nil {
		return err
	}

	//rewrite login to be case insensitive
	user.Login = strings.ToLower(user.Login)

	//find login
	for _, usr := range res {
		if usr.Login == user.Login {
			return &db.NotUnique{user.Login}
		}
	}
	//create user
	return r.crud.CreateTX(tx, user)
}

func (r *Users) Update(user *User) error {
	tx := r.db.Partition(TABLE_USER).Begin(true)
	defer tx.Commit()

	res := make([]*User, 0)
	err := r.crud.List(TABLE_USER, "", &res)
	if res != nil {
		return err
	}

	//rewrite login to be case insensitive
	user.Login = strings.ToLower(user.Login)

	//find login
	for _, usr := range res {
		if usr.Login == user.Login {
			return &db.NotUnique{user.Login}
		}
	}
	//update user
	return r.crud.Update(TABLE_USER, user)
}

func (r *Users) FindByLogin(login string) (*User, error) {
	//rewrite login to be case insensitive
	login = strings.ToLower(login)

	list, err := r.List()
	if err != nil {
		return nil, err
	}
	//find login
	var foundUsr *User
	for _, usr := range list {
		if usr.Login == login {
			foundUsr = usr
			break
		}
	}
	if foundUsr == nil {
		return nil, &db.EntityNotFound{login}
	}

	return foundUsr, nil
}
