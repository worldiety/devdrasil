package main
/*
import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/worldiety/devdrasil/backend"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type pluginProxy struct {
	server *Devdrasil
}

//registers for all requests which are constructed like /rpc/* calls where /rpc/<plugin instance id>/<endpoint path>
func installPluginProxy(server *Devdrasil) {
	handler := &pluginProxy{}
	server.mux.HandleFunc("/rpc/*", handler.handle)
}

func (h *pluginProxy) handle(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	segments := strings.Split(path, "/")
	instanceId := backend.InstanceId(segments[1])
	methodId := path[5+len(instanceId):]
	sessionId := request.Header.Get("sid")
	params := request.Header.Get("params")

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

	if !user.Plugins.HasInstance(instanceId) {
		http.Error(writer, string(instanceId), http.StatusForbidden)
		return
	}

	managedInstance := h.server.GetManagedPlugin(instanceId)
	if managedInstance == nil {
		http.Error(writer, string(instanceId), http.StatusNotFound)
		return
	}

	req, e := http.NewRequest("GET", "http://"+managedInstance.Host+":"+strconv.Itoa(managedInstance.Port)+"/"+managedInstance.ContextPath, nil)
	if e != nil {
		http.Error(writer, e.Error(), http.StatusInternalServerError)
		return
	}

	userJson := toJson(user)

	mac := hmac.New(sha256.New, []byte(managedInstance.Secret))
	mac.Write([]byte(sessionId))
	mac.Write([]byte(params))
	mac.Write([]byte(methodId))
	mac.Write([]byte(userJson))

	req.Header.Add("sid", sessionId)
	req.Header.Add("params", params)
	req.Header.Add("mid", methodId)
	req.Header.Add("user", userJson)
	req.Header.Add("hmac", base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	client := &http.Client{}
	res, e := client.Do(req)

	if e != nil {
		http.Error(writer, e.Error(), http.StatusServiceUnavailable)
		return
	}

	defer res.Body.Close()
	copyHeader(writer.Header(), res.Header)
	writer.WriteHeader(res.StatusCode)
	_, e = io.Copy(writer, res.Body)
	if e != nil {
		http.Error(writer, e.Error(), http.StatusServiceUnavailable)
		return
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func toJson(obj interface{}) string {
	b, e := json.Marshal(obj)
	if e != nil {
		panic(e)
	}
	return string(b)
}
*/