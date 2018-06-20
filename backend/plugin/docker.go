package plugin

import (
	"path/filepath"
	"regexp"
	"fmt"
	"os"
	"io/ioutil"
	"log"
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
}

/*
Expects that the git is configured correctly and is able to just pull the repository without password or else.
It is a known problem for repositories which requires private/public ssh keys or user/password authentication. You
need to solve that for the devdrasil-user by hand on the devdrasil machine. We could use a go-only git implementation
but we would require to support all edge cases then - unclear if it is worth the effort.
 */
func (r *PluginManager) Install(pluginId string, gitUrl string) error {
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
			log.Printf("expected an empty or non-existing directory: %s\n", pluginDir)
			return fmt.Errorf("cannot install plugin '%s', directory is not empty", pluginId)
		}
	}

	os.MkdirAll(pluginDir, defaultFilePermission)
	stat, err := os.Stat(pluginDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		log.Printf("expected a directory: %s\n", pluginDir)
		return fmt.Errorf("cannot create plugin directory")
	}

	os.MkdirAll(pluginDir, defaultFilePermission)

	return nil
}

func (r *PluginManager) Remove(pluginId string) error {
	err := validatePluginId(pluginId)
	if err != nil {
		return err
	}
	return nil
}

func (r *PluginManager) Update(pluginId string) error {
	err := validatePluginId(pluginId)
	if err != nil {
		return err
	}
	return nil
}

func validatePluginId(id string) error {
	re := regexp.MustCompile("^[a-z0-9_.]+$")
	if !re.MatchString(id) {
		return fmt.Errorf("invalid format of id '" + id + "', use a format like 'com.mycompany.myplugin'")
	}
	return nil
}

