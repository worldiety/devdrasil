package tools

import (
	"sync"
	"encoding/json"
	"strconv"
	"strings"
	"os"
)

type Docker struct {
	Env   *Env
	mutex sync.Mutex
}

type ContainerId string

//A container is a spawned instance of an (ever) image identified by repository|tag.
// Once spawned, a container can be started and stopped as required.
type Container struct {
	Command      string
	CreatedAt    string
	ID           ContainerId
	Image        string
	Labels       string
	LocalVolumes string
	Mounts       string
	Names        string
	Networks     string
	Ports        string
	RunningFor   string
	Size         string
	Status       string
}

type StartOptions struct {
	//Docker image name == repository?
	Repository string

	//Docker image tag
	Tag string

	HostPort int

	ContainerPort int

	//key/value of labels may not contain spaces or ,
	Labels map[string]string

	//all mounts
	Mounts []*Mount

	//removes an exited "unuseful" stopped container automatically
	RemoveOnExit bool
}

type Mount struct {
	HostDir      string
	ContainerDir string
	ReadOnly     bool
}

//An Image is a readonly receipt to create instances of Container
type Image struct {
	Containers   string
	CreatedAt    string
	CreatedSince string
	Digest       string
	//the image id
	ID          string
	Repository  string
	SharedSize  string
	Size        string
	Tag         string
	UniqueSize  string
	VirtualSize string
}

//returns the value of the denoted label or os.ErrNotExist
func (c *Container) GetLabelValue(key string) (string, error) {
	labels := strings.Split(c.Labels, ",")
	for _, label := range labels {
		keyValue := strings.Split(label, "=")
		if keyValue[0] == key {
			return keyValue[1], nil
		}
	}
	return "", os.ErrNotExist
}

func NewDocker(env *Env) *Docker {
	return &Docker{Env: env}
}

//invokes a docker build -t in the current directory
func (d *Docker) Build(name string, tag string, removeIntermediateContainer bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	rm := "--rm=true"
	if !removeIntermediateContainer {
		rm = "--rm=false"
	}
	_, err := d.Env.ExecLines("docker", "build", rm, "-t", name+":"+tag, ".")
	if err != nil {
		return err
	}
	return nil
}

//stops a container
func (d *Docker) Stop(id ContainerId) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, err := d.Env.ExecLines("docker", "stop", string(id))
	if err != nil {
		return err
	}
	return nil
}

//starts a container, e.g. docker run -d -p 4000:80 blublub:xy --label x=y
func (d *Docker) Start(options *StartOptions) (ContainerId, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	args := make([]string, 0)
	args = append(args, "run", "-d")

	for key, value := range options.Labels {
		args = append(args, "--label", key+"="+value)
	}

	for _, mount := range options.Mounts {
		//type=bind,source="$(pwd)"/target,target=/app
		mountCmd := "type=bind,source=" + mount.HostDir + ",target=" + mount.ContainerDir + ""
		if mount.ReadOnly {
			mountCmd += ",readonly"
		}
		args = append(args, "--mount", mountCmd)
	}

	if options.RemoveOnExit {
		args = append(args, "--rm")
	}

	args = append(args, "-p", strconv.Itoa(options.HostPort)+":"+strconv.Itoa(options.ContainerPort), options.Repository+":"+options.Tag)

	lines, err := d.Env.ExecLines("docker", args...)
	if err != nil {
		return "", err
	}
	return ContainerId(lines[0]), nil
}

//lists all available images, e.g. docker images -a  --no-trunc --format="{{  json . }}"
func (d *Docker) ListImages() ([]*Image, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	lines, err := d.Env.ExecLines("docker", "images", "-a", "--no-trunc", "--format", "{{ json . }}")
	if err != nil {
		return nil, err
	}

	res := make([]*Image, 0)

	for _, line := range lines {
		c := &Image{}
		err := json.Unmarshal([]byte(line), c)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

//lists all running containers, e.g. docker container ls --no-trunc  --format="{{  json . }}"
func (d *Docker) ListContainers() ([]*Container, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	lines, err := d.Env.ExecLines("docker", "container", "ls", "-a", "--no-trunc", "--format", "{{  json . }}")
	if err != nil {
		return nil, err
	}

	res := make([]*Container, 0)

	for _, line := range lines {
		c := &Container{}
		err := json.Unmarshal([]byte(line), c)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

//removes an image, e.g. docker rmi blublub:xy -f
func (d *Docker) RemoveImage(repository string, tag string, force bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	args := []string{"rmi"}

	if force {
		args = append(args, "-f")
	}

	if len(tag) > 0 {
		args = append(args, repository+":"+tag)
	} else {
		args = append(args, repository)
	}

	_, err := d.Env.ExecLines("docker", args...)
	return err
}

//Removes all intermediate images and other temp stuff which is usually not required,  docker system prune -f
func (d *Docker) Prune() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, err := d.Env.ExecLines("docker", "system", "prune", "-f")
	return err
}
