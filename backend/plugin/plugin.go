package plugin

import (
	"path/filepath"
	"io/ioutil"
	"log"
	"github.com/worldiety/devdrasil/tools/exec"
)

const dockerImageVersionFilename = "docker-image.version"

type Plugin struct {
	pluginId  string
	rootDir   string
	baseDir   string
	dockerDir string
	dataDir   string
}

//creates a new plugin instance. BaseDir is something like ~/.devdrasil/plugins/de.worldiety.devdrasil.buildserver/
func NewPlugin(baseDir string) *Plugin {
	return &Plugin{pluginId: filepath.Base(baseDir), rootDir: filepath.Dir(baseDir), baseDir: baseDir, dockerDir: filepath.Join(baseDir, "docker"), dataDir: filepath.Join(baseDir, "data")}
}

//Returns the current docker image version, represented as a VCS version (e.g. git hash). If no version is available the empty string is returned
func (p *Plugin) GetDockerImageVersion() string {
	buf, err := ioutil.ReadFile(filepath.Join(p.baseDir, dockerImageVersionFilename))
	if err != nil {
		log.Printf("no %s is available: %s\n", dockerImageVersionFilename, err)
		return ""
	}
	return string(buf)
}

//checks if an update is available
func (p *Plugin) IsUpdateAvailable() (bool, error) {
	env := tools.NewEnv()

}

func (p *Plugin) Update() {

}
