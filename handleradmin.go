package main

import "net/http"

type adminHandler struct {
	server *Devdrasil
}

func installAdminHandler(server *Devdrasil) {
	handler := &adminHandler{}
	server.mux.HandleFunc("/admin", handler.handle)
	server.mux.HandleFunc("/admin/", handler.handle)
}

func (h *adminHandler) handle(writer http.ResponseWriter, request *http.Request) {

	writer.Write([]byte("hello admin area"))
}
