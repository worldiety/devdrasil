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

	git := r.getGit(pluginDir)
	err = git.Clone(gitUrl)
	if err != nil {
		return err
	}
	err = git.Checkout("master")
	if err != nil {
		return err
	}

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
	pluginDir := filepath.Join(r.rootDir, pluginId)

	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		} else {
			return "", err
		}
	} else {
		if len(files) > 0 {
			git := r.getGit(pluginDir)
			return git.GetHead()
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

	pluginDir := filepath.Join(r.rootDir, pluginId)
	return os.RemoveAll(pluginDir)
}

func (r *PluginManager) Update(pluginId string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := validatePluginId(pluginId)
	if err != nil {
		return err
	}

	pluginDir := filepath.Join(r.rootDir, pluginId)
	git := r.getGit(pluginDir)
	err = git.Pull()
	if err != nil {
		return err
	}

	err = git.Checkout("master")
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
