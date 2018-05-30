package main

import (
	"encoding/json"
	"github.com/worldiety/devdrasil/plugin"
	"github.com/worldiety/devdrasil/plugin/session"
	"net/http"
)

type handlerSession struct {
	server *Devdrasil
}

/*
	Register endpoints for session handling
*/
func installHandlerSession(server *Devdrasil) {
	handler := &handlerSession{server}
	server.mux.HandleFunc("/session/auth", handler.login)
	server.mux.HandleFunc("/users", handler.listUsers)
}

//request header must "login" and "password" and "client"
func (h *handlerSession) login(writer http.ResponseWriter, request *http.Request) {
	login := request.Header.Get("login")
	pwd := request.Header.Get("password")
	client := request.Header.Get("client")

	sessionId, e := h.server.sessionPlugin.AuthenticateUser(login, pwd, client)
	if e != nil {
		http.Error(writer, e.Error(), http.StatusForbidden)
		return
	}

	type wrapper struct {
		SessionId string
	}

	buf, e := json.Marshal(&wrapper{string(sessionId)})
	if e != nil {
		http.Error(writer, e.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(buf)
}

//request header must be "sid"
func (h *handlerSession) listUsers(writer http.ResponseWriter, request *http.Request) {
	sessionId := request.Header.Get("sid")
	activeSession, e := h.server.sessionPlugin.GetSession(sessionId)
	if e != nil {
		http.Error(writer, e.Error(), http.StatusForbidden)
		return
	}

	user, e := h.server.sessionPlugin.GetUser(activeSession.Uid)
	if e != nil {
		http.Error(writer, e.Error(), http.StatusForbidden)
		return
	}

	if !user.Properties.Has(session.ROLE_LIST_USER) {
		http.Error(writer, e.Error(), http.StatusForbidden)
		return
	}

	users, e := h.server.sessionPlugin.ListUser()
	if e != nil {
		http.Error(writer, e.Error(), http.StatusInternalServerError)
		return
	}

	type userDTO struct {
		//the login id
		Id string

		//roles and such
		Properties plugin.Properties

		//the allowed plugin instance ids
		Plugins plugin.Plugins

		//flag if user is active or not, without deleting him
		Active bool
	}

	type wrapper struct {
		Users []*userDTO
	}

	res := &wrapper{}
	for _, user := range users {
		res.Users = append(res.Users, &userDTO{user.Id, user.Properties, user.Plugins, user.Active})
	}
	buf, e := json.Marshal(res)
	if e != nil {
		http.Error(writer, e.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(buf)
}
