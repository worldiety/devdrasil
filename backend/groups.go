package backend

import (
	"net/http"
	"github.com/worldiety/devdrasil/backend/session"
	"github.com/worldiety/devdrasil/backend/user"
	"github.com/worldiety/devdrasil/backend/group"
	"github.com/worldiety/devdrasil/db"
	"strings"
)

type EndpointGroups struct {
	mux         *http.ServeMux
	users       *user.Users
	sessions    *session.Sessions
	groups      *group.Groups
	permissions *user.Permissions
}

type groupListDTO struct {
	List []*groupDTO
}

func newGroupDTO(users *user.Users, group *group.Group) *groupDTO {
	list, _ := users.List()
	tmp := make([]db.PK, 0)
	for _, usr := range list {
		if usr.HasGroup(group.Id) {
			tmp = append(tmp, usr.Id)
		}
	}
	return &groupDTO{Id: group.Id, Name: group.Name, Users: tmp}
}

type groupDTO struct {
	//unique entity id, e.g. "abc38293"
	Id db.PK

	//the name of the group
	Name string

	//all users within this group
	Users []db.PK
}

func NewEndpointGroups(mux *http.ServeMux, sessions *session.Sessions, users *user.Users, permissions *user.Permissions, groups *group.Groups) *EndpointGroups {
	endpoint := &EndpointGroups{mux: mux, sessions: sessions, permissions: permissions, users: users, groups: groups}
	mux.HandleFunc("/groups/", endpoint.groupVerbs)
	mux.HandleFunc("/groups", endpoint.groupsVerbs)
	return endpoint
}

func (e *EndpointGroups) groupVerbs(writer http.ResponseWriter, request *http.Request) {
	groupId, err := db.ParsePK(strings.TrimPrefix(request.URL.Path, "/groups/"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	switch request.Method {
	case "GET":
		e.getGroup(writer, request, groupId)
	case "PUT":
		e.updateGroup(writer, request, groupId)
	case "DELETE":
		e.deleteGroup(writer, request, groupId)

	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
		return
	}
}

// A user can list all groups, if he has the permission LIST_GROUPS
//  @Path GET /groups
//  @Header sid string
//	@Body []github.com/worldiety/devdrasil/backend/groupListDTO
//	@Return 200
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointGroups) listGroups(writer http.ResponseWriter, request *http.Request) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.LIST_GROUPS)
	if usr == nil {
		return
	}

	groups, err := e.groups.List()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &groupListDTO{}
	res.List = make([]*groupDTO, 0)
	for _, g := range groups {
		res.List = append(res.List, newGroupDTO(e.users, g))
	}
	WriteJSONBody(writer, res)
}

func (e *EndpointGroups) groupsVerbs(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		e.listGroups(writer, request);
	case "POST":
		e.addGroup(writer, request);
	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
		return
	}
}

// A user can add another group, if he has the permission for ADDING_GROUPS
//  @Path POST /groups
//  @Header sid string
//	@Body github.com/worldiety/devdrasil/backend/groupDTO
//	@Return 200 github.com/worldiety/devdrasil/backend/groupDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointGroups) addGroup(writer http.ResponseWriter, request *http.Request) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.CREATE_GROUP)
	if usr == nil {
		return
	}

	dto := &groupDTO{}
	err := ReadJSONBody(writer, request, dto)
	if err != nil {
		return
	}

	newGroup := &group.Group{}
	newGroup.Name = dto.Name

	err = e.groups.Add(newGroup)
	if err != nil {
		if db.IsNotUnique(err) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = e.updateAllUsers(newGroup.Id, dto.Users)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONBody(writer, newGroupDTO(e.users, newGroup))
}

//loops through all users and removes the group reference from all users which are not in the given list and adds the group to all users given.
func (e *EndpointGroups) updateAllUsers(groupId db.PK, users []db.PK) error {
	allUsers, err := e.users.List()
	if err != nil {
		return err
	}

	for _, user := range allUsers {
		userShouldBeInGroup := false
		for _, usr := range users {
			if user.Id == usr {
				userShouldBeInGroup = true
				break
			}
		}
		if userShouldBeInGroup && !user.HasGroup(groupId) {
			user.Groups = append(user.Groups, groupId)
			err = e.users.Update(user)
			if err != nil {
				return err
			}
		} else
		if !userShouldBeInGroup && user.HasGroup(groupId) {
			user.RemoveGroup(groupId)
			err = e.users.Update(user)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// A user can delete another group, if he has the permission DELETE_GROUP
//  @Path DELETE /groups/{id}
//  @Header sid string
//	@Return 200
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointGroups) deleteGroup(writer http.ResponseWriter, request *http.Request, groupId db.PK) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.DELETE_GROUP)
	if usr == nil {
		return
	}

	e.updateAllUsers(groupId, nil)
	err := e.groups.Delete(groupId)
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

// A user needs the GET_GROUP permission
//  @Path GET /groups/{id} (id is hex encoded group PK)
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/groupDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointGroups) getGroup(writer http.ResponseWriter, request *http.Request, groupId db.PK) {
	_, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return
	}

	//check if the permission is available
	allowed, err := e.permissions.IsAllowed(user.GET_GROUP, usr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(writer, err.Error(), http.StatusForbidden)
		return
	}

	group, err := e.groups.Get(groupId)
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
	WriteJSONBody(writer, newGroupDTO(e.users, group))
	return
}

// A user needs the UPDATE_GROUP permission.
//  @Path PUT /groups/{id}
//  @Header sid string
//	@Body github.com/worldiety/devdrasil/backend/groupDTO
//	@Return 200 github.com/worldiety/devdrasil/backend/groupDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointGroups) updateGroup(writer http.ResponseWriter, request *http.Request, groupId db.PK) {
	_, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return
	}

	dto := &groupDTO{}
	err := ReadJSONBody(writer, request, dto)
	if err != nil {
		return
	}

	//check if the permission is available
	allowed, err := e.permissions.IsAllowed(user.UPDATE_GROUP, usr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(writer, err.Error(), http.StatusForbidden)
		return
	}

	otherGroup, err := e.groups.Get(groupId)
	if err != nil {
		if db.IsEntityNotFound(err) {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	otherGroup.Name = dto.Name

	//rewrite
	err = e.groups.Update(otherGroup)
	if err != nil {
		if db.IsNotUnique(err) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	err = e.updateAllUsers(otherGroup.Id, dto.Users)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	//return the newly data
	WriteJSONBody(writer, newGroupDTO(e.users, otherGroup))

}
