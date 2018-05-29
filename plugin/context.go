package plugin

type Context interface {
	//Returns the current user which is never nil. However if not authenticated the login and everything else is empty
	User() *User

	//your plugin id
	PluginId() PluginId

	//your plugin instance id
	PluginInstanceId() InstanceId

	//the current session id, if no session available is empty
	SessionId() SessionId

	//the actual method id
	MethodId() MethodId

	//parses the raw parameters into a destination, use map[string]interface{} for unspecified parsing
	Parameter(dst interface{}) error
}

type User struct {
	//Id of the user, which is the login name
	Id string

	//Each user has a bunch of properties, which is a key and multiple values. There is no additional nested configuration. Such things should be represented by namespacing like my.plugin.property0
	Properties Properties

	//the allowed plugin instance ids
	Plugins Plugins
}

//A map whose keys are the plugin-ids and the values are the available instance ids
type Plugins map[PluginId][]InstanceId

//A map with key and values
type Properties map[string][]string

type PluginId string
type InstanceId string
type SessionId string
type MethodId string

func (p Plugins) HasInstance(id InstanceId) bool {
	for _, values := range p {
		for _, v := range values {
			if v == id {
				return true
			}
		}
	}
	return false
}

func (p Properties) Has(key string) bool {
	if _, ok := p[key]; ok {
		return true
	}
	return false
}

func (p Properties) HasValue(key string, value string) bool {
	if values, ok := p[key]; ok {
		for _, v := range values {
			if v == value {
				return true
			}
		}
	}
	return false
}

type Description struct {
	id        string
	doc       string
	endpoints []*PluginEndpoint
}

func NewDescription(id string) *Description {
	return &Description{id: id}
}

func (p *Description) Doc(doc string) *Description {
	p.doc = doc
	return p
}

func (p *Description) AddEndpoint(id string, endpoint func(ctx Context) (interface{}, error)) *PluginEndpoint {
	ep := &PluginEndpoint{plugin: p, endpoint: endpoint}
	p.endpoints = append(p.endpoints, ep)
	ep.id = id
	return ep
}

type PluginEndpoint struct {
	id         string
	doc        string
	endpoint   func(ctx Context) (interface{}, error)
	parameters []interface{}
	out        []interface{}
	plugin     *Description
}

func (p *PluginEndpoint) DocParam(example interface{}) *PluginEndpoint {
	p.parameters = append(p.parameters, example)
	return p
}

func (p *PluginEndpoint) DocReturn(example interface{}) *PluginEndpoint {
	p.out = append(p.out, example)
	return p
}

func (p *PluginEndpoint) Doc(doc string) *PluginEndpoint {
	p.doc = doc
	return p
}

func (p *PluginEndpoint) And() *Description {
	return p.plugin
}

func (p *Description) CreateApiDoc() *ApiDoc {
	doc := &ApiDoc{}
	doc.Doc = p.doc
	doc.PluginId = p.id
	for _, ep := range p.endpoints {
		doc.Endpoints = append(doc.Endpoints, &ApiEndpointDoc{EndpointId: ep.id, Doc: ep.doc})
	}
	return doc
}

type Success struct {
}

type Void struct {
}

//manages a concrete plugin
type Manager interface {

	//Returns the plugin id, e.g. com.mydomain.myplugin
	Id() string

	//Creates and deploys a new instance and returns the id for it
	Deploy() (string, error)

	//Lists all deployed instances
	List() ([]ManagedInstance, error)

	//Tries to start an instance
	Start(id string) error

	//Tries to stop an instance
	Stop(id string) error

	//returns the api doc
	Doc() *ApiDoc
}

type ManagedInstance struct {
	//Instance id
	Id string

	//Current status of the instance
	Status Status

	//the host where this instance runs on
	Host string

	//the port where it runs
	Port int

	//a context path, if multiple plugins are on the same host
	ContextPath string

	//the secret which is used to sign the requests with an hmac. As long as an attacker does not know the secret, an attacker cannot tamper requests
	Secret string
}

type Status int

const STATUS_UNKOWN Status = 0
const STATUS_RUNNING Status = 1
const STATUS_STOPPED Status = 2

type ApiDoc struct {
	PluginId  string
	Doc       string
	Endpoints []*ApiEndpointDoc
}

type ApiEndpointDoc struct {
	EndpointId string
	Doc        string
}
