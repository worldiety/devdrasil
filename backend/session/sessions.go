package session

import (
	"github.com/worldiety/devdrasil/db"
)

const TABLE_SESSION = "session"

type Session struct {
	//the actual unique session id
	Id db.PK

	//user id
	User db.PK

	//for gc'ing the session, the number of seconds elapsed since January 1, 1970 UTC.
	CreatedAt int64

	//for gc'ing the session, the number of seconds elapsed since January 1, 1970 UTC.
	LastUsedAt int64

	//the last remote address
	LastRemoteAddr string

	//the user agent string
	LastUserAgent string
}

type Sessions struct {
	db   *db.Database
	crud *db.CRUD
}

func NewSessions(d *db.Database) *Sessions {
	r := &Sessions{d, db.NewCRUD(d)}
	return r
}

func (s *Sessions) Get(pk db.PK) (*Session, error) {
	session := &Session{Id: pk}
	err := s.crud.Read(TABLE_SESSION, session)
	return session, err
}

func (s *Sessions) Delete(pk db.PK) error {
	return s.crud.Delete(TABLE_SESSION, pk)
}

func (s *Sessions) Update(session *Session) error {
	return s.crud.Update(TABLE_SESSION, session)
}

func (s *Sessions) Create(session *Session) error {
	return s.crud.Create(TABLE_SESSION, session)
}

func (s *Sessions) List() ([]*Session, error) {
	res := make([]*Session, 0)
	err := s.crud.List(TABLE_SESSION, "", res)
	return res, err
}
