package backend

import (
	"net/http"
	"github.com/worldiety/devdrasil/backend/session"
	"github.com/worldiety/devdrasil/backend/user"
	"github.com/worldiety/devdrasil/backend/store"
)

type EndpointMarket struct {
	mux         *http.ServeMux
	sessions    *session.Sessions
	users       *user.Users
	permissions *user.Permissions
	store       *store.Store
}

func NewEndpointStore(mux *http.ServeMux, sessions *session.Sessions, users *user.Users, permissions *user.Permissions) *EndpointMarket {
	endpoint := &EndpointMarket{mux: mux, users: users, sessions: sessions, permissions: permissions, store: &store.Store{}}
	mux.HandleFunc("/market/index", endpoint.getIndex)
	mux.HandleFunc("/market/install/", endpoint.installPlugin)
	return endpoint
}

// Requires permissions LIST_MARKET
//  @Path GET /market/plugins
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/store/Index
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointMarket) getIndex(writer http.ResponseWriter, request *http.Request) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.LIST_MARKET)
	if usr == nil {
		return
	}

	index, err := e.store.GetIndex()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSONBody(writer, index)
}

// Requires permissions INSTALL_PLUGIN
//  @Path POST /market/install/{id}
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/?
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointMarket) installPlugin(writer http.ResponseWriter, request *http.Request) {

}
