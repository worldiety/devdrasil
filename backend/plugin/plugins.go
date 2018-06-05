package plugin

import "github.com/worldiety/devdrasil/db"

const TABLE_PLUGIN = "plugin"

type Plugin struct {
	//Entity id, e.g. "0x373abc"
	Id db.PK

	//the unique developer id, e.g. "com.worldiety.buildserver"
	DevId string

	//the human readable name, e.g. "worldiety buildserver"
	Name string

	//the instances of this plugin
	Instances []*Instance
}

type Instance struct {
	//Entity id, e.g. "0xca31d"
	Id db.PK

	//the host, which this plugin runs on, e.g. "127.0.0.1"
	Host string

	//the port of the plugin, e.g. "9090"
	Port int
}

type Plugins struct {
	db   *db.Database
	crud *db.CRUD
}

func NewPlugins(d *db.Database) *Plugins {
	return &Plugins{d, db.NewCRUD(d)}
}

func (r *Plugins) List() ([]*Plugin, error) {
	tx := r.db.Partition(TABLE_PLUGIN).Begin(false)
	defer tx.Commit()
	return r.list(tx)
}

func (r *Plugins) list(tx db.Transaction) ([]*Plugin, error) {
	res := make([]*Plugin, 0)
	err := r.crud.ListTX(tx, "", &res)
	return res, err
}

func (r *Plugins) Create(plugin *Plugin) error {
	tx := r.db.Partition(TABLE_PLUGIN).Begin(true)
	defer tx.Commit()

	plugins, err := r.List()
	if err != nil {
		return err
	}
	//ensure unique DevId
	for _, group := range plugins {
		if group.DevId == group.DevId {
			return &db.NotUnique{group.DevId}
		}
	}

	plugin.Id = tx.NextKey()
	json := db.NewJSONDecorator(tx)
	return json.Put(plugin)
}

func (r *Plugins) Delete(id db.PK) error {
	return r.crud.Delete(TABLE_PLUGIN, id)
}

func (r *Plugins) Update(plugin *Plugin) error {
	return r.crud.Update(TABLE_PLUGIN, plugin)
}
