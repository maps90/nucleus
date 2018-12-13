package elastic

import (
	"context"
	"fmt"
	"sync"

	"github.com/maps90/nucleus/config"
	elasticDB "github.com/olivere/elastic"
)

var elasticConnections *Conn

// Conn struct
type Conn struct {
	es  map[string]*elasticDB.Client
	mux sync.Mutex
}

// Config struct
type Config struct {
	DNS      string
	Fallback string
}

// New : initialize new DB connections, if dbConnection is already exists, returns that instead
func New() *Conn {
	if elasticConnections == nil {
		elasticConnections = &Conn{es: make(map[string]*elasticDB.Client)}
	}

	return elasticConnections
}

// Get connection in thread-safe fashion
func (db *Conn) Get(id string) *elasticDB.Client {
	db.mux.Lock()
	defer db.mux.Unlock()
	if conn, ok := db.es[id]; ok {
		return conn
	}
	return nil
}

// Set connection in thread-safe fashion
func (db *Conn) Set(id string, es *elasticDB.Client) {
	db.mux.Lock()
	db.es[id] = es
	db.mux.Unlock()
}

// Connect to elastic based on config
func (db *Conn) Connect(id string) (*elasticDB.Client, error) {
	elasticConfig := config.GetStringMap("elastic")
	if _, ok := elasticConfig[id]; !ok {
		return nil, fmt.Errorf("elastic configuration for [%s] does not exists", id)
	}

	if con := db.Get(id); con != nil {
		return con, nil
	}

	con, err := db.openConnection(id)
	if err != nil {
		return nil, err
	}

	db.Set(id, con)
	return db.Get(id), nil
}

// Shutdown disconnecting all established elastic connections
func (db *Conn) Shutdown() (err error) {
	if elasticConnections == nil {
		return nil
	}
	for _, c := range db.es {
		c.Stop()
	}
	return err
}

func (db *Conn) openConnection(id string) (*elasticDB.Client, error) {
	opt := setupConfig(id)

	client, err := elasticDB.NewClient(
		elasticDB.SetSniff(false),
		elasticDB.SetURL(opt.DNS),
	)
	if err != nil {
		return nil, err
	}

	_, _, err = client.Ping(opt.DNS).Do(context.Background())
	if err != nil {
		if opt.Fallback == id { // prevent endless loop
			return client, err
		}
		return db.Connect(opt.Fallback)
	}

	db.Set(id, client)
	return db.Get(id), nil
}

func setupConfig(id string) *Config {
	option := &Config{
		DNS: config.GetString(getKey(id, "dns")),
	}

	return option
}

func getKey(id, types string) string {
	return fmt.Sprintf("elastic.%s.%s", id, types)
}
