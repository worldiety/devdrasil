package main

import (
	"flag"
	"github.com/worldiety/devdrasil/db"
	"github.com/worldiety/devdrasil/backend"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"github.com/worldiety/devdrasil/backend/session"
	"github.com/worldiety/devdrasil/backend/user"
	"github.com/worldiety/devdrasil/backend/group"
	"github.com/worldiety/devdrasil/backend/company"
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

	mutex sync.Mutex

	restSessions  *backend.EndpointSessions
	restUsers     *backend.EndpointUsers
	restGroups    *backend.EndpointGroups
	restCompanies *backend.EndpointCompanies
	restMarket    *backend.EndpointMarket
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

	installFrontendHandler(devdrasil)

	users, err := user.NewUsers(devdrasil.db)
	if err != nil {
		panic(err)
	}

	permissions, err := user.NewPermissions(devdrasil.db)
	if err != nil {
		panic(err)
	}
	sessions := session.NewSessions(devdrasil.db)

	groups := group.NewGroups(devdrasil.db)

	companies := company.NewCompanies(devdrasil.db)

	devdrasil.restUsers = backend.NewEndpointUsers(devdrasil.mux, sessions, users, permissions)
	devdrasil.restSessions = backend.NewEndpointSessions(devdrasil.mux, sessions, users)
	devdrasil.restGroups = backend.NewEndpointGroups(devdrasil.mux, sessions, users, permissions, groups)
	devdrasil.restCompanies = backend.NewEndpointCompanies(devdrasil.mux, sessions, users, permissions, companies)
	devdrasil.restMarket = backend.NewEndpointStore(devdrasil.mux, sessions, users, permissions)

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

//Start actually launches the server
func (s *Devdrasil) Start() {

	log.Printf("current working directory is %s\n", s.cwd)
	log.Printf("workspace is %s\n", s.workspace)
	log.Printf("plugins located at %s\n", s.plugins)
	log.Printf("starting devdrasil at %s:%d...\n", s.host, s.port)
	log.Fatal(http.ListenAndServe(s.host+":"+strconv.Itoa(s.port), nil))
}
