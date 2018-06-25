package backend

import (
	"net/http"
	"github.com/worldiety/devdrasil/backend/session"
	"github.com/worldiety/devdrasil/backend/user"
	"github.com/worldiety/devdrasil/backend/store"
	"github.com/worldiety/devdrasil/backend/plugin"
	"strings"
	"os"
	"time"
)

type PluginInfo struct {
	plugin.PluginVersionInfo
}

type EndpointMarket struct {
	mux           *http.ServeMux
	sessions      *session.Sessions
	users         *user.Users
	permissions   *user.Permissions
	store         *store.Store
	pluginManager *plugin.PluginManager
}

func NewEndpointStore(mux *http.ServeMux, sessions *session.Sessions, users *user.Users, permissions *user.Permissions, pluginManager *plugin.PluginManager) *EndpointMarket {
	endpoint := &EndpointMarket{mux: mux, users: users, sessions: sessions, permissions: permissions, store: &store.Store{}, pluginManager: pluginManager}
	mux.HandleFunc("/market/index", endpoint.getIndex)
	mux.HandleFunc("/plugins/", endpoint.pluginsVerbs)
	return endpoint
}

func (e *EndpointMarket) pluginsVerbs(writer http.ResponseWriter, request *http.Request) {
	pluginId := strings.TrimPrefix(request.URL.Path, "/plugins/")
	switch request.Method {
	case "GET":
		e.getPluginInfo(writer, request, pluginId)
	case "POST":
		e.installPlugin(writer, request, pluginId)
	case "PUT":
		e.updatePlugin(writer, request, pluginId)
	case "DELETE":
		e.deletePlugin(writer, request, pluginId)
	default:
		http.Error(writer, request.Method, http.StatusMethodNotAllowed);
		return
	}
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

// Requires permissions UPDATE_PLUGIN
//  @Path PUT /plugins/{id}
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/PluginInfo
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointMarket) updatePlugin(writer http.ResponseWriter, request *http.Request, pluginId string) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.UPDATE_PLUGIN)
	if usr == nil {
		return
	}

	//some sleep for nice visualization
	time.Sleep(2 * time.Second)

	index, err := e.store.GetIndex()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	plg := index.GetPlugin(pluginId)
	if plg == nil {
		http.Error(writer, pluginId, http.StatusNotFound)
		return
	}
	if plg.Source.Type != "git" {
		http.Error(writer, "sources of type "+plg.Source.Type+" are not supported", http.StatusInternalServerError)
		return
	}
	err = e.pluginManager.Update(pluginId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	version, err := e.pluginManager.GetVersion(pluginId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSONBody(writer, &PluginInfo{*version})

}

// Requires permissions INSTALL_PLUGIN
//  @Path POST /plugins/{id}
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/PluginInfo
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointMarket) installPlugin(writer http.ResponseWriter, request *http.Request, pluginId string) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.INSTALL_PLUGIN)
	if usr == nil {
		return
	}

	//some sleep for nice visualization
	time.Sleep(2 * time.Second)

	index, err := e.store.GetIndex()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	plg := index.GetPlugin(pluginId)
	if plg == nil {
		http.Error(writer, pluginId, http.StatusNotFound)
		return
	}
	if plg.Source.Type != "git" {
		http.Error(writer, "sources of type "+plg.Source.Type+" are not supported", http.StatusInternalServerError)
		return
	}
	err = e.pluginManager.Install(pluginId, plg.Source.Url)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	version, err := e.pluginManager.GetVersion(pluginId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSONBody(writer, &PluginInfo{*version})
}

// Requires only an authenticated user. By design every user can query installed plugins, so that later UI components can fit themself properly
//  @Path GET /plugins/{id}
//  @Header sid string
//	@Return 200 github.com/worldiety/devdrasil/backend/PluginInfo
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointMarket) getPluginInfo(writer http.ResponseWriter, request *http.Request, pluginId string) {
	_, usr := GetSessionAndUser(e.sessions, e.users, writer, request)
	if usr == nil {
		return
	}

	version, err := e.pluginManager.GetVersion(pluginId)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	WriteJSONBody(writer, &PluginInfo{*version})
}

// Requires permissions REMOVE_PLUGIN
//  @Path DELETE /plugins/{id}
//  @Header sid string
//	@Return 200
//  @Return 403 (if session id is invalid | if session user is inactive | if session user is absent)
//  @Return 500 (for any other error)
func (e *EndpointMarket) deletePlugin(writer http.ResponseWriter, request *http.Request, pluginId string) {
	_, usr := validate(e.sessions, e.users, e.permissions, writer, request, user.UPDATE_PLUGIN)
	if usr == nil {
		return
	}

	//some sleep for nice visualization
	time.Sleep(2 * time.Second)

	_, err := e.pluginManager.GetVersion(pluginId)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = e.pluginManager.Remove(pluginId, false)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
