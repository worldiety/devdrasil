package main

import (
	"flag"
	"github.com/worldiety/devdrasil/db"
	"github.com/worldiety/devdrasil/plugin"
	"github.com/worldiety/devdrasil/plugin/session"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

type Devdrasil struct {
	//Host is the ip or host name where the server listens
	host string

	//Port is the actual port where the server listens
	port int

	//the folder where all data and the (un-isolated) plugins are organized
	plugins string

	//the folder where our home is
	workspace string

	//the devdrasil configuration database
	db *db.Database

	//the http server
	mux *http.ServeMux

	//current working dir
	cwd string

	//plugins are organized through their managers. Each plugin get's it's own manager.
	pluginManagers []plugin.Manager

	//each plugin is accessed by their instance id, but may be located at a different host
	reverseInstanceLookup map[plugin.InstanceId]*plugin.ManagedInstance

	mutex sync.Mutex

	//the built-in session management
	sessionPlugin *session.Plugin
}

func NewDevdrasil() *Devdrasil {
	if runtime.GOOS == "windows" {
		log.Fatal("devdrasil currently does not support windows.")
	}

	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("the 'HOME' environment variable is missing")
	}

	cwd, e := os.Getwd()
	if e != nil {
		log.Fatalf("failed to get the current working dir: %s\n", e)
	}

	devdrasil := &Devdrasil{}
	devdrasil.workspace = filepath.Join(home, ".devdrasil")
	ensureDir(devdrasil.workspace)

	devdrasil.plugins = filepath.Join(devdrasil.workspace, "plugins")
	ensureDir(devdrasil.plugins)

	devdrasil.cwd = cwd

	devdrasil.db = db.Open(ensureDir(filepath.Join(devdrasil.workspace, "db")))
	devdrasil.host = *flag.String("host", "0.0.0.0", "A host name or ip address to which devdrasil is bound")
	devdrasil.port = *flag.Int("port", 8080, "The port on which devdrasil listens")

	devdrasil.mux = http.DefaultServeMux
	devdrasil.sessionPlugin = session.NewSessionAPI(devdrasil.db)

	installFrontendHandler(devdrasil)
	installPluginProxy(devdrasil)
	return devdrasil
}

func ensureDir(dir string) string {
	//only the owner can read/write/execute
	os.MkdirAll(dir, 0700)
	stat, e := os.Stat(dir)
	if e != nil {
		log.Fatalf("failed to ensure directory '%s': %s\n", dir, e)
	}
	if !stat.IsDir() {
		log.Fatalf("not a directory '%s'\n", dir)
	}
	return dir
}

func (s *Devdrasil) GetManagedPlugin(instanceId plugin.InstanceId) *plugin.ManagedInstance {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.reverseInstanceLookup[instanceId]
}

//Start actually launches the server
func (s *Devdrasil) Start() {

	log.Printf("current working directory is %s\n", s.cwd)
	log.Printf("workspace is %s\n", s.workspace)
	log.Printf("plugins located at %s\n", s.plugins)
	log.Printf("starting devdrasil at %s:%d...\n", s.host, s.port)
	log.Fatal(http.ListenAndServe(s.host+":"+strconv.Itoa(s.port), nil))
}
