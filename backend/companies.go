package backend

import (
	"net/http"
	"github.com/worldiety/devdrasil/backend/session"
	"github.com/worldiety/devdrasil/backend/user"
	"github.com/worldiety/devdrasil/db"
	"strings"
	"github.com/worldiety/devdrasil/backend/company"
)

type EndpointCompanies struct {
	mux         *http.ServeMux
	users       *user.Users
	sessions    *session.Sessions
	companies   *company.Companies
	permissions *user.Permissions
}

type companyListDTO struct {
	List []*companyDTO
}

func newCompanyDTO(users *user.Users, company *company.Company) *companyDTO {
	list, _ := users.List()
	tmp := make([]db.PK, 0)
	for _, usr := range list {
		if usr.HasCompany(company.Id) {
			tmp = append(tmp, usr.Id)
		}
	}
	return &companyDTO{Id: company.Id, Name: company.Name, Users: tmp, ThemePrimaryColor: company.ThemePrimaryColor}
}

type companyDTO struct {
	//unique entity id, e.g. "abc38293"
	Id db.PK

	//the name of the group
	Name string

	//the primary color
	ThemePrimaryColor string

	//all users within this group
	Users []db.PK
}

func NewEndpointCompanies(mux *http.ServeMux, sessions *session.Sessions, users *user.Users, permissions *user.Permissions, companies *company.Companies) *EndpointCompanies {
	endpoint := &EndpointCompanies{mux: mux, sessions: sessions, permissions: permissions, users: users, companies: companies}
	mux.HandleFunc("/companies/", endpoint.companyVerbs)
	mux.HandleFunc("/companies", endpoint.companiesVerbs)
	return endpoint
}

func (e *EndpointCompanies) companyVerbs(writer http.ResponseWriter, request *http.Request) {
	groupId, err := db.ParsePK(strings.TrimPrefix(request.URL.Path, "/companies/"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	switch request.Method {
	case "GET":
		e.getCompany(writer, request, groupId)
	case "PUT":
		e.updateCompany(writer, request, groupId)
	case "DELETE":
		e.deleteCompany(writer, request, groupId)

	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
		return
	}
}

// A user can list all companies, if he has the permission list user permissions
//  @Path GET /companies
//  @Header sid string
//	@Body []github.com/worldiety/devdrasil/backend/companyListDTO
//	@Return 200
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointCompanies) listCompanies(writer http.ResponseWriter, request *http.Request) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.LIST_USERS)
	if usr == nil {
		return
	}

	companies, err := e.companies.List()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &companyListDTO{}
	res.List = make([]*companyDTO, 0)
	for _, g := range companies {
		res.List = append(res.List, newCompanyDTO(e.users, g))
	}
	WriteJSONBody(writer, res)
}

func (e *EndpointCompanies) companiesVerbs(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		e.listCompanies(writer, request);
	case "POST":
		e.addCompany(writer, request);
	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
		return
	}
}

// A user can add another company, if he has the permission for adding user
//  @Path POST /companies
//  @Header sid string
//	@Body github.com/worldiety/devdrasil/backend/companyDTO
//	@Return 200 github.com/worldiety/devdrasil/backend/companyDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointCompanies) addCompany(writer http.ResponseWriter, request *http.Request) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.CREATE_USER)
	if usr == nil {
		return
	}

	dto := &companyDTO{}
	err := ReadJSONBody(writer, request, dto)
	if err != nil {
		return
	}

	newCompany := &company.Company{}

	updateModelFromDTO(dto, newCompany)

	err = e.companies.Add(newCompany)
	if err != nil {
		if db.IsNotUnique(err) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = e.updateAllUsers(newCompany.Id, dto.Users)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONBody(writer, newCompanyDTO(e.users, newCompany))
}

//loops through all users and removes the group reference from all users which are not in the given list and adds the group to all users given.
func (e *EndpointCompanies) updateAllUsers(companyId db.PK, users []db.PK) error {
	allUsers, err := e.users.List()
	if err != nil {
		return err
	}

	for _, user := range allUsers {
		userShouldBeInCompany := false
		for _, usr := range users {
			if user.Id == usr {
				userShouldBeInCompany = true
				break
			}
		}
		if userShouldBeInCompany && !user.HasCompany(companyId) {
			user.Company = &companyId
			err = e.users.Update(user)
			if err != nil {
				return err
			}
		} else
		if !userShouldBeInCompany && user.HasCompany(companyId) {
			user.Company = nil
			err = e.users.Update(user)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// A user can delete another company, if he has the permission (delete user)
//  @Path DELETE /companies/{id}
//  @Header sid string
//	@Return 200
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent | if user has not the permission)
//  @Return 500 (for any other error)
func (e *EndpointCompanies) deleteCompany(writer http.ResponseWriter, request *http.Request, companyId db.PK) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.DELETE_USER)
	if usr == nil {
		return
	}

	e.updateAllUsers(companyId, nil)
	err := e.companies.Delete(companyId)
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

// A user needs the GET_USER permission
//  @Path GET /companies/{id} (id is hex encoded group PK)
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/companyDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointCompanies) getCompany(writer http.ResponseWriter, request *http.Request, companyId db.PK) {
	_, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return
	}

	allowed := false
	//check if it is the company of the current user
	if usr.Company != nil && *usr.Company == companyId {
		allowed = true
	}

	//check if the permission is available
	if !allowed {
		tmp, err := e.permissions.IsAllowed(user.GET_USER, usr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		allowed = tmp
	}

	if !allowed {
		http.Error(writer, "", http.StatusForbidden)
		return
	}

	company, err := e.companies.Get(companyId)
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
	WriteJSONBody(writer, newCompanyDTO(e.users, company))
	return
}

func updateModelFromDTO(src *companyDTO, dst *company.Company) {
	dst.Name = src.Name
	dst.ThemePrimaryColor = src.ThemePrimaryColor
}

// A user needs the UPDATE permission.
//  @Path PUT /companies/{id}
//  @Header sid string
//	@Body github.com/worldiety/devdrasil/backend/companyDTO
//	@Return 200 github.com/worldiety/devdrasil/backend/companyDTO
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointCompanies) updateCompany(writer http.ResponseWriter, request *http.Request, groupId db.PK) {
	_, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return
	}

	dto := &companyDTO{}
	err := ReadJSONBody(writer, request, dto)
	if err != nil {
		return
	}

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

	otherCompany, err := e.companies.Get(groupId)
	if err != nil {
		if db.IsEntityNotFound(err) {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	updateModelFromDTO(dto, otherCompany)

	//rewrite
	err = e.companies.Update(otherCompany)
	if err != nil {
		if db.IsNotUnique(err) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	err = e.updateAllUsers(otherCompany.Id, dto.Users)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	//return the newly data
	WriteJSONBody(writer, newCompanyDTO(e.users, otherCompany))

}
