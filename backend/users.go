package backend

import (
	"net/http"
	"github.com/worldiety/devdrasil/backend/user"
	"github.com/worldiety/devdrasil/backend/session"
	"github.com/worldiety/devdrasil/db"
	"strings"
)

type userListDTO struct {
	Users []*userDTO
}

type userDTO struct {
	//unique entity id, e.g. "abc38293"
	Id *db.PK

	//Abbreviation and/or Login, something like "tschinke", used for the login
	Login *string

	//e.g. Torben
	Firstname *string

	//e.g. Schinke
	Lastname *string

	//the new password
	Password *string

	//flag if user is active or not, without deleting him
	Active *bool

	//reference to an optional avatar image
	AvatarImage *db.PK

	//list of connected email addresses e.g. tschinke@domain.com, torben.schinke@otherdomain.com, ...
	EMailAddresses *[]string

	//reference to an optional company
	Company *db.PK

	//the groups, which this user is a member of. This determines his actual permissions.
	Groups *[]db.PK
}

func newUserDTO(user *user.User) *userDTO {
	return &userDTO{Id: &user.Id, Login: &user.Login, Firstname: &user.Firstname, Lastname: &user.Lastname, Active: &user.Active, AvatarImage: user.AvatarImage, EMailAddresses: &user.EMailAddresses, Groups: &user.Groups}
}

type EndpointUsers struct {
	mux         *http.ServeMux
	sessions    *session.Sessions
	users       *user.Users
	permissions *user.Permissions
}

func NewEndpointUsers(mux *http.ServeMux, sessions *session.Sessions, users *user.Users, permissions *user.Permissions) *EndpointUsers {
	endpoint := &EndpointUsers{mux: mux, users: users, sessions: sessions, permissions: permissions}
	mux.HandleFunc("/users/", endpoint.userVerbs)
	mux.HandleFunc("/users/permissions/", endpoint.permissionsVerbs)
	mux.HandleFunc("/users", endpoint.usersVerbs)
	return endpoint
}

func (endpoint *EndpointUsers) permissionsVerbs(writer http.ResponseWriter, request *http.Request) {
	userId, err := db.ParsePK(strings.TrimPrefix(request.URL.Path, "/users/permissions/"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	switch request.Method {
	case "GET":
		endpoint.queryPermissions(writer, request, userId)
	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
		return
	}
}

func (endpoint *EndpointUsers) userVerbs(writer http.ResponseWriter, request *http.Request) {
	userId, err := db.ParsePK(strings.TrimPrefix(request.URL.Path, "/users/"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	switch request.Method {
	case "GET":
		endpoint.getUser(writer, request, userId)
	case "PUT":
		endpoint.updateUser(writer, request, userId)
	case "DELETE":
		endpoint.deleteUser(writer, request, userId)

	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
		return
	}
}

// A session user can always request it's own user object, but others require the correct permission (which is GET_USER)
//  @Path GET /users/{id} (id is hex encoded user PK)
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/userDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointUsers) getUser(writer http.ResponseWriter, request *http.Request, userId db.PK) {
	_, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return
	}

	//a user can always request himself
	if usr.Id == userId {
		//return the user dto
		WriteJSONBody(writer, newUserDTO(usr))
		return
	} else {
		//check if the permission is available
		allowed, err := e.permissions.IsAllowed(user.GET_USER, usr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}

		otherUser, err := e.users.Get(userId)
		if err != nil {
			if db.IsEntityNotFound(err) {
				http.Error(writer, err.Error(), http.StatusNotFound)
				return
			} else {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		//return the other user dto
		WriteJSONBody(writer, newUserDTO(otherUser))
		return
	}
}

// A user can always update his data on his own. Inactive users cannot make changes, but they can make themself inactive. Also users can do so, when having UPDATE permission.
//  @Path PUT /users/{id}
//  @Header sid string
//	@Body github.com/worldiety/devdrasil/backend/userDTO
//	@Return 200 github.com/worldiety/devdrasil/backend/userDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointUsers) updateUser(writer http.ResponseWriter, request *http.Request, userId db.PK) {
	_, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return
	}

	dto := &userDTO{}
	err := ReadJSONBody(writer, request, dto)
	if err != nil {
		return
	}

	var userToUpdate *user.User
	//a user can always update himself
	if usr.Id == userId {
		userToUpdate = usr

	} else {
		//check if the permission is available
		allowed, err := e.permissions.IsAllowed(user.UPDATE_USER, usr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}

		otherUser, err := e.users.Get(userId)
		if err != nil {
			if db.IsEntityNotFound(err) {
				http.Error(writer, err.Error(), http.StatusNotFound)
			} else {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		userToUpdate = otherUser
	}

	//actually transfer affected fields
	e.updateUserFields(userToUpdate, dto)

	//rewrite
	err = e.users.Update(userToUpdate)
	if err != nil {
		if db.IsNotUnique(err) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	//return the newly data
	WriteJSONBody(writer, newUserDTO(userToUpdate))

}

func (e *EndpointUsers) updateUserFields(usr *user.User, dto *userDTO) {
	dto.Id = &usr.Id
	if dto.Active != nil {
		usr.Active = *dto.Active
	}
	if dto.Groups != nil {
		usr.Groups = *dto.Groups
	}

	if dto.EMailAddresses != nil {
		usr.EMailAddresses = *dto.EMailAddresses
	}

	if dto.AvatarImage != nil {
		usr.AvatarImage = dto.AvatarImage
	}

	if dto.Lastname != nil {
		usr.Lastname = *dto.Lastname
	}

	if dto.Firstname != nil {
		usr.Firstname = *dto.Firstname
	}

	if dto.Login != nil {
		usr.Login = *dto.Login
	}

	if dto.Company != nil {
		usr.Company = dto.Company
	}

	if dto.Password != nil {
		usr.SetPassword(*dto.Password)
	}
}

func (endpoint *EndpointUsers) usersVerbs(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		endpoint.listUsers(writer, request);
	case "POST":
		endpoint.addUser(writer, request);
	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
		return
	}
}

//if returns nil, just return, because the request has been cancelled
func (e *EndpointUsers) validate(writer http.ResponseWriter, request *http.Request, kind db.PK) (*session.Session, *user.User) {
	ses, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return nil, nil
	}

	allowed, err := e.permissions.IsAllowed(kind, usr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return nil, nil
	}

	if !allowed {
		http.Error(writer, err.Error(), http.StatusForbidden)
		return nil, nil
	}
	return ses, usr
}

// A user can list all other users, if he has the permission
//  @Path GET /users
//  @Header sid string
//	@Body []github.com/worldiety/devdrasil/backend/userListDTO
//	@Return 200
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointUsers) listUsers(writer http.ResponseWriter, request *http.Request) {
	_, usr := e.validate(writer, request, user.LIST_USERS)
	if usr == nil {
		return
	}

	users, err := e.users.List()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &userListDTO{}
	for _, u := range users {
		res.Users = append(res.Users, newUserDTO(u))
	}
	WriteJSONBody(writer, res)
}

// A user can add another user, if he has the permission
//  @Path POST /users
//  @Header sid string
//	@Body github.com/worldiety/devdrasil/backend/userDTO
//	@Return 200 github.com/worldiety/devdrasil/backend/userDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointUsers) addUser(writer http.ResponseWriter, request *http.Request) {
	_, usr := e.validate(writer, request, user.CREATE_USER)
	if usr == nil {
		return
	}

	dto := &userDTO{}
	err := ReadJSONBody(writer, request, dto)
	if err != nil {
		return
	}

	newUser := &user.User{}
	e.updateUserFields(newUser, dto)

	err = e.users.Add(newUser)
	if err != nil {
		if db.IsNotUnique(err) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	WriteJSONBody(writer, newUserDTO(newUser))
}

// A user can delete another user, if he has the permission
//  @Path DELETE /users/{id}
//  @Header sid string
//	@Return 200
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointUsers) deleteUser(writer http.ResponseWriter, request *http.Request, userId db.PK) {
	_, usr := e.validate(writer, request, user.DELETE_USER)
	if usr != nil {
		return
	}

	err := e.users.Delete(userId)
	if err != nil {
		if db.IsEntityNotFound(err) {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	WriteOK(writer)
}

type userPermissionDTO struct {
	ListUsers  bool
	CreateUser bool
	DeleteUser bool
	UpdateUser bool
	GetUser    bool
}

// A user can request the permissions, always from his own account, or if he has the list permission
//  @Path GET /users/permissions/{id}
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/userPermissionDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointUsers) queryPermissions(writer http.ResponseWriter, request *http.Request, userId db.PK) {
	_, usr := e.validate(writer, request, user.LIST_USERS)
	if usr == nil {
		return
	}

	listUser, err := e.permissions.IsAllowed(user.LIST_USERS, usr)
	if AnyErrorAsInternalError(err, writer) {
		return
	}

	createUser, err := e.permissions.IsAllowed(user.CREATE_USER, usr)
	if AnyErrorAsInternalError(err, writer) {
		return
	}

	deleteUser, err := e.permissions.IsAllowed(user.DELETE_USER, usr)
	if AnyErrorAsInternalError(err, writer) {
		return
	}

	updateUser, err := e.permissions.IsAllowed(user.UPDATE_USER, usr)
	if AnyErrorAsInternalError(err, writer) {
		return
	}

	getUser, err := e.permissions.IsAllowed(user.GET_USER, usr)
	if AnyErrorAsInternalError(err, writer) {
		return
	}

	dto := &userPermissionDTO{}
	dto.ListUsers = listUser
	dto.CreateUser = createUser
	dto.DeleteUser = deleteUser
	dto.UpdateUser = updateUser
	dto.GetUser = getUser

	WriteJSONBody(writer, dto)
}
