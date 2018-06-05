package backend

import (
	"net/http"
	"github.com/worldiety/devdrasil/backend/session"
	"github.com/worldiety/devdrasil/backend/user"
	"github.com/worldiety/devdrasil/db"
	"time"
	"strings"
)

var allowedClients = []string{"web-client-1.0"}

type sessionDTO struct {
	//base64 encoded session id
	Id db.PK

	//base64 encoded user id
	User db.PK
}

type EndpointSessions struct {
	mux      *http.ServeMux
	sessions *session.Sessions
	users    *user.Users
}

func NewEndpointSessions(mux *http.ServeMux, sessions *session.Sessions, users *user.Users) *EndpointSessions {
	endpoint := &EndpointSessions{mux: mux, users: users, sessions: sessions}
	mux.HandleFunc("/sessions", endpoint.sessionsVerbs)
	mux.HandleFunc("/sessions/", endpoint.sessionVerbs)
	return endpoint
}

func (endpoint *EndpointSessions) sessionsVerbs(writer http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		endpoint.listSessions(writer, request)
	case "POST":
		endpoint.auth(writer, request)

	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
	}
}

func (endpoint *EndpointSessions) sessionVerbs(writer http.ResponseWriter, request *http.Request) {
	sessionId, err := db.ParsePK(strings.TrimPrefix(request.URL.Path, "/session/"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	switch request.Method {
	case "DELETE":
		endpoint.deleteSession(writer, request, sessionId);

	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
	}
}

// Everybody can delete a session for logout.
//  @Path DELETE /sessions/{id}
//	@Return 200
//  @Return 404 (if session id is invalid )
//  @Return 500 (for any other error)
func (e *EndpointSessions) deleteSession(writer http.ResponseWriter, request *http.Request, sessionId db.PK) {
	err := e.sessions.Delete(sessionId)
	if err != nil {
		if db.IsEntityNotFound(err) {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

// The admin user can list all available sessions.
//  @Path GET /sessions
//  @Header sid string
//	@Return 200 []*github.com/worldiety/devdrasil/backend/session/Session
//  @Return 403 (if session id is invalid | if session user is not admin)
//  @Return 500 (for any other error)
func (e *EndpointSessions) listSessions(writer http.ResponseWriter, request *http.Request) {
	sess, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return
	}

	//only admin users are allowed to view all sessions
	if sess.User != user.ADMIN {
		http.Error(writer, "user must be admin", http.StatusForbidden)
		return
	}

	sessions, err := e.sessions.List()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONBody(writer, sessions)
}

// Everybody can try to create a session by posting to the session resource. For security reason each request is delayed at least by 1 second.
//  @Path POST /sessions
//  @Header login string (The login)
//  @Header password string (The password)
//  @Header User-Agent string (The user agent)
//	@Header client string (A client token, to validate, or issue compatiblity or what else. Allowed values: 'web-client-1.0')
//	@Return 200 github.com/worldiety/devdrasil/backend/sessionDTO
//  @Return 403 (if any auth data is invalid or rejected, or user is inactive etc.)
//  @Return 500 (for any other error)
func (e *EndpointSessions) auth(writer http.ResponseWriter, request *http.Request) {
	login := request.Header.Get("login")
	pwd := request.Header.Get("password")
	agent := request.Header.Get("User-Agent")
	client := request.Header.Get("client")

	//try to do something against brute force attacks, see also https://www.owasp.org/index.php/Blocking_Brute_Force_Attacks
	time.Sleep(1000 * time.Second)

	//another funny idea is to return a fake session id, after many wrong login attempts

	if len(login) < 3 {
		http.Error(writer, "login too short", http.StatusForbidden)
		return
	}

	if len(agent) == 0 {
		http.Error(writer, "user agent missing", http.StatusForbidden)
		return
	}

	if len(pwd) < 4 {
		http.Error(writer, "password too short", http.StatusForbidden)
		return
	}

	if len(client) == 0 {
		http.Error(writer, "client missing", http.StatusForbidden)
		return
	}

	allowed := false
	for _, allowedClient := range allowedClients {
		if allowedClient == client {
			allowed = true
			break
		}
	}

	if !allowed {
		http.Error(writer, "client is invalid", http.StatusForbidden)
		return
	}

	usr, err := e.users.FindByLogin(login)

	if err != nil {
		if db.IsEntityNotFound(err) {
			http.Error(writer, "user unkown", http.StatusForbidden)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if !usr.Active {
		http.Error(writer, "user is inactive", http.StatusForbidden)
		return
	}

	//login is fine now, create a session
	currentTime := time.Now().Unix()
	ses := &session.Session{User: usr.Id, LastUsedAt: currentTime, CreatedAt: currentTime, LastRemoteAddr: request.RemoteAddr, LastUserAgent: agent}
	err = e.sessions.Create(ses)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONBody(writer, &sessionDTO{Id: ses.Id, User: usr.Id})
}
