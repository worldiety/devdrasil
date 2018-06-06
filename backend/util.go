package backend

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/worldiety/devdrasil/backend/session"
	"github.com/worldiety/devdrasil/db"
	"github.com/worldiety/devdrasil/backend/user"
	"io/ioutil"
)

func ErrPermissionDenied(which string) error {
	return fmt.Errorf("permissioned denied: " + which)
}

func ErrInvalidParameter(which string) error {
	return fmt.Errorf("invalid parameter: " + which)
}

func WriteOK(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
}

//respond with a json body
func WriteJSONBody(writer http.ResponseWriter, obj interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Write(b)
}

//parse a json body, spitting a server error, so you simply must return
func ReadJSONBody(writer http.ResponseWriter, request *http.Request, obj interface{}) error {
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return err
	}
	err = json.Unmarshal(b, obj)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

//returns either a session and a user or both nil. In the latter, you may not write anything more to the outputstream
func GetSessionAndUser(sessions *session.Sessions, users *user.Users, writer http.ResponseWriter, request *http.Request) (*session.Session, *user.User) {
	sessionId, err := db.ParsePK(request.Header.Get("sid"))
	if err != nil {
		http.Error(writer, "invalid session id format", http.StatusForbidden)
		return nil, nil
	}

	//check session
	session, err := sessions.Get(sessionId)
	if err != nil {
		http.Error(writer, "invalid session id", http.StatusForbidden)
		return nil, nil
	}

	//re-check if user actually still exists
	user, err := users.Get(session.User)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusForbidden)
		return nil, nil
	}

	//check inactive
	if !user.Active {
		http.Error(writer, "user is inactive", http.StatusForbidden)
		return nil, nil
	}

	return session, user
}

func AnyErrorAsInternalError(err error, writer http.ResponseWriter) bool {
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
