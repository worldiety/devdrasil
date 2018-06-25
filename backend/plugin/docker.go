package plugin

import (
	"path/filepath"
	"regexp"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/worldiety/devdrasil/tools"
	"github.com/worldiety/devdrasil/log"
	"sync"
)

//by default only the owner is able to read, write and exec
const defaultFilePermission = 0700

const dockerLabelPlugin = "devdrasil-plugin"
const pluginApp = "app"
const pluginData = "data"
const pluginBranch = "master"

type PluginManager struct {
	/*
		The dir where all plugins are stored, e.g.
			~/.devdrasil/plugins
				  ./de.worldiety.devdrasil.buildserver/
				 	  ./docker/
				 	  ./data/

	 */
	rootDir string
	mutex   sync.Mutex
}

func NewPluginManager(dir string) *PluginManager {
	return &PluginManager{rootDir: dir}
}

func (r *PluginManager) getGit(dir string) *tools.Git {
	git := tools.NewGit(tools.NewEnv())
	git.Env.Dir = dir
	return git
}

func (r *PluginManager) getDocker(dir string) *tools.Docker {
	docker := tools.NewDocker(tools.NewEnv())
	docker.Env.Dir = dir
	return docker
}

/*
Expects that the git is configured correctly and is able to just pull the repository without password or else.
It is a known problem for repositories which requires private/public ssh keys or user/password authentication. You
need to solve that for the devdrasil-user by hand on the devdrasil machine. We could use a go-only git implementation
but we would require to support all edge cases then - unclear if it is worth the effort.
 */
func (r *PluginManager) Install(pluginId string, gitUrl string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := validatePluginId(pluginId)
	if err != nil {
		return err
	}
	pluginDir := filepath.Join(r.rootDir, pluginId)

	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if len(files) > 0 {
			log.Default.Warn(log.New("expected an empty or non-existing directory").Put("dir", pluginDir))
			return fmt.Errorf("cannot install plugin '%s', directory is not empty", pluginId)
		}
	}

	os.MkdirAll(pluginDir, defaultFilePermission)
	stat, err := os.Stat(pluginDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		log.Default.Warn(log.New("expected a directory").Put("dir", pluginDir))
		return fmt.Errorf("cannot create plugin directory")
	}

	os.MkdirAll(pluginDir, defaultFilePermission)

	appDir := filepath.Join(pluginDir, pluginApp)
	os.Mkdir(appDir, defaultFilePermission)

	dataDir := filepath.Join(pluginDir, pluginData)
	os.Mkdir(dataDir, defaultFilePermission)

	log.Default.Info(log.New("cloning into").Put("dir", appDir))
	git := r.getGit(appDir)
	err = git.Clone(gitUrl)
	if err != nil {
		return err
	}
	err = git.Checkout(pluginBranch)
	if err != nil {
		return err
	}

	//all files are here, start dockering
	revision, err := git.GetHead()
	if err != nil {
		return err
	}
	docker := r.getDocker(appDir)
	err = docker.Build(pluginId, revision, true)
	if err != nil {
		return err
	}

	options := &tools.StartOptions{}
	options.Repository = pluginId
	options.Tag = revision
	options.Labels = map[string]string{dockerLabelPlugin: pluginId}
	options.ContainerPort = 80
	options.HostPort = 4000 //TODO
	options.RemoveOnExit = false //does not work with restart always
	options.Mounts = []*tools.Mount{{HostDir: dataDir, ContainerDir: "/" + pluginData, ReadOnly: false}}
	options.Restart = "always"
	options.HostIP = "127.0.0.1"

	cid, err := docker.Start(options)
	if err != nil {
		return err
	}

	log.Default.Info(log.New("container running").Put("containerId", cid))

	return nil
}

//checks if a plugin folder with data is available. Does not check if the plugin is functional. Returns the git hash
func (r *PluginManager) GetVersion(pluginId string) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := validatePluginId(pluginId)
	if err != nil {
		return "", err
	}
	pluginAppDir := filepath.Join(r.rootDir, pluginId, pluginApp)

	files, err := ioutil.ReadDir(pluginAppDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		} else {
			return "", err
		}
	} else {
		if len(files) > 0 {
			git := r.getGit(pluginAppDir)
			head, err := git.GetHead()
			if err != nil {
				log.Default.Error(log.New("cannot get version").Put("dir", pluginAppDir).SetError(err))
				return "n/a", nil
			}
			return head, nil
		} else {
			return "", os.ErrNotExist
		}
	}
}

//removes everything, including data, git, docker images, etc.
func (r *PluginManager) Remove(pluginId string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := validatePluginId(pluginId)
	if err != nil {
		return err
	}

	//remove all plugin data, like the git repo and app data
	pluginDir := filepath.Join(r.rootDir, pluginId)
	errDirRemove := os.RemoveAll(pluginDir)

	docker := r.getDocker(".")

	//shutdown any docker container instance with the appropriate label
	containers, err := docker.ListContainers()
	if err != nil {
		return err
	}
	var errDockerContainer error
	for _, con := range containers {
		value, _ := con.GetLabelValue(dockerLabelPlugin)
		if value == pluginId {
			err := docker.Stop(con.ID)
			if errDockerContainer != nil {
				errDockerContainer = err
			}
			if err != nil {
				log.Default.Error(log.New("failed to stop docker container").Put("id", con.ID))
			}
		}
	}

	//remove all docker images for the plugin (pluginId==repository)
	images, err := docker.ListImages()
	if err != nil {
		return err
	}
	var errDockerRemove error
	for _, img := range images {
		if img.Repository == pluginId {
			err := docker.RemoveImage(pluginId, img.Tag, true)
			if errDockerRemove == nil {
				errDockerRemove = err
			}
			if err != nil {
				log.Default.Error(log.New("failed to remove docker image").Put("repository", img.Repository).Put("tag", img.Tag))
			}

		}
	}

	//just GC everything
	docker.Prune()

	if errDirRemove != nil {
		return errDirRemove
	}

	if errDockerRemove != nil {
		return errDockerRemove
	}

	if errDockerContainer != nil {
		return errDockerContainer
	}
	return nil
}

func (r *PluginManager) Update(pluginId string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := validatePluginId(pluginId)
	if err != nil {
		return err
	}

	pluginAppDir := filepath.Join(r.rootDir, pluginId, pluginApp)
	git := r.getGit(pluginAppDir)
	err = git.Pull()
	if err != nil {
		return err
	}

	err = git.Checkout(pluginBranch)
	if err != nil {
		return err
	}

	return nil
}

//this is a security essential: avoid various filename attacks, like ../../etc/ because the id is used directly in the filesystem
func validatePluginId(id string) error {
	re := regexp.MustCompile("^[a-z0-9_.]+$")
	if !re.MatchString(id) {
		return fmt.Errorf("invalid format of id '" + id + "', use a format like 'com.mycompany.myplugin'")
	}
	return nil
}
