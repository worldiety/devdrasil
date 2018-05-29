package session

import (
	"github.com/worldiety/devdrasil/db"
	"github.com/worldiety/devdrasil/plugin"
	"github.com/worldiety/devdrasil/tools/rand"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"sync"
	"time"
)

const TABLE_USER = "user"
const TABLE_SESSION = "session"

const ADMIN_LOGIN = "admin"
const ADMIN_PWD = "admin"

const PLUGIN_ID = "com.worldiety.session"

type Plugin struct {
	mutex sync.Mutex
	db    *db.Database
}

func NewSessionAPI(db *db.Database) *Plugin {
	return &Plugin{db: db}
}

func NewPlugin(db *db.Database) *plugin.Description {
	plg := NewSessionAPI(db)
	plg.ensureAdminUser()

	desc := plugin.NewDescription(PLUGIN_ID)
	desc.Doc("The Authentication & Session plugin, providing the basic API.")
	//desc.AddEndpoint("user/authenticate", plg.AuthenticateUser).Doc("Authenticates a user").DocParam(&AuthUserRequest{}).DocReturn(&AuthUserResponse{})
	desc.AddEndpoint("user/add", plg.AddUser).Doc("adds a user. Requires a user with ROLE_ADD_USER").DocParam(&AddUserRequest{}).DocReturn(&AuthUserResponse{})
	//desc.AddEndpoint("user/list", plg.ListUser).Doc("lists all user. Requires a user with ROLE_LIST_USER").DocParam(&plugin.Void{}).DocReturn(&ListUserResponse{})
	return desc
}

func (s *Plugin) ensureAdminUser() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.db.Partition(TABLE_USER).Has(ADMIN_LOGIN) {
		user := &DBUser{}
		hash, err := bcrypt.GenerateFromPassword([]byte(ADMIN_PWD), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		user.Id = ADMIN_LOGIN
		user.PasswordHash = hash
		user.Active = true
		user.Properties = plugin.Properties{}
		user.Properties[ROLE_LIST_USER] = nil
		user.Properties[ROLE_ADD_USER] = nil

		err = s.db.Partition(TABLE_USER).Put(user)
		if err != nil {
			panic(err)
		}
		log.Printf("created default user with password '%':'%s'", ADMIN_LOGIN, ADMIN_PWD)
	}
}

func (s *Plugin) AuthenticateUser(login string, password string, client string) (plugin.SessionId, error) {

	//for now, we ignore the client
	request := &AuthUserRequest{}

	user := &DBUser{}
	err := s.db.Partition(TABLE_USER).Get(request.Login, user)
	if err != nil {
		log.Printf("failed to authenticate '%s': %s\n", request.Login, err)
		return "", ErrPermissionDenied("authentication failure")
	}

	if !user.Active {
		return "", ErrPermissionDenied("User is inactive")
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password))
	if err != nil {
		log.Printf("failed to compare '%s': %s\n", request.Login, err)
		return "", ErrPermissionDenied("authentication failure")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	//generate a guaranteed unique session id
	sessionId := rand.SecureRandomHex(16)
	for ; s.db.Partition(TABLE_SESSION).Has(sessionId); sessionId = rand.SecureRandomHex(16) {

	}

	//insert session into the db
	t := time.Now().Unix()
	session := &DBSession{Id: sessionId, Uid: user.Id, CreatedAt: t, LastUsedAt: t}
	err = s.db.Partition(TABLE_SESSION).Put(session)
	if err != nil {
		return "", err
	}

	return plugin.SessionId(session.Id), nil
}

func (s *Plugin) AddUser(ctx plugin.Context) (interface{}, error) {
	sessionUser := ctx.User()
	if !sessionUser.Properties.Has(ROLE_ADD_USER) {
		return nil, ErrPermissionDenied(ROLE_ADD_USER)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	request := &AddUserRequest{}
	ctx.Parameter(request)
	request.Login = strings.TrimSpace(strings.ToLower(request.Login))
	request.Password = strings.TrimSpace(request.Password)

	if request.Login == "" {
		return nil, ErrInvalidParameter("Login is empty")
	}

	if len(request.Password) < 4 {
		return nil, ErrInvalidParameter("Password to short")
	}

	if s.db.Partition(TABLE_USER).Has(request.Login) {
		return nil, ErrInvalidParameter("user exists")
	}

	newUser := &DBUser{}
	hash, err := bcrypt.GenerateFromPassword([]byte(ADMIN_PWD), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	newUser.Id = request.Login
	newUser.PasswordHash = hash
	newUser.Active = true
	err = s.db.Partition(TABLE_USER).Put(newUser)
	if err != nil {
		return nil, err
	}

	return &plugin.Success{}, nil
}

func (s *Plugin) ListUser() ([]*DBUser, error) {

	res := make([]*DBUser, 0)

	cursor := s.db.Partition(TABLE_USER).Query("ORDER BY Id")
	defer cursor.Close()
	for cursor.Next() {
		user := &DBUser{}
		e := cursor.Scan(user)
		if e != nil {
			log.Printf("cannot read user: %s\n", e)
			continue
		}
		res = append(res, user)
	}

	return res, nil
}

func (s *Plugin) GetSession(id string) (*DBSession, error) {
	session := &DBSession{}
	e := s.db.Partition(TABLE_SESSION).Get(id, session)
	if e != nil {
		return nil, e
	}
	return session, nil
}

func (s *Plugin) GetUser(id string) (*DBUser, error) {
	user := &DBUser{}
	e := s.db.Partition(TABLE_USER).Get(id, user)
	if e != nil {
		return nil, e
	}
	return user, nil
}

//internal user
type DBUser struct {
	//the login id
	Id string

	//and the additional password hash
	PasswordHash []byte

	//roles and such
	Properties plugin.Properties

	//the allowed plugin instance ids
	Plugins plugin.Plugins

	//flag if user is active or not, without deleting him
	Active bool
}

//internal session
type DBSession struct {
	//the actual unique session id
	Id string

	//user id
	Uid string

	//for gc'ing the session, the number of seconds elapsed since January 1, 1970 UTC.
	CreatedAt int64

	//for gc'ing the session, the number of seconds elapsed since January 1, 1970 UTC.
	LastUsedAt int64
}
