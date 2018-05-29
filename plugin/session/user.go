package session

import "github.com/worldiety/devdrasil/plugin"

type AuthUserRequest struct {
	//The login
	Login string
	//the plain text password for authentication
	Password string
	//the client token
	Token string
}

type AuthUserResponse struct {
	//the session id
	Session string
}

type AddUserRequest struct {
	Login    string
	Password string
}

type ListUserResponse struct {
	Users []*plugin.User
}
