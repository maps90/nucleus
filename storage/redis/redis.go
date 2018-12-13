package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/maps90/nucleus/config"
)

var redisConnections = &Conn{DB: make(map[string]*Client)}

// Conn struct
type Conn struct {
	DB  map[string]*Client
	mux sync.Mutex
}

// Get connection in thread-safe fashion
func (s *Conn) Get(id string) *Client {
	s.mux.Lock()
	defer s.mux.Unlock()
	if conn, ok := s.DB[id]; ok {
		return conn
	}
	return nil
}

// Set connection in thread-safe fashion
func (s *Conn) Set(id string, client *Client) {
	s.mux.Lock()
	s.DB[id] = client
	s.mux.Unlock()
}

// Client struct
type Client struct {
	*redis.Client
	context context.Context
}

func setupConfig(id string) *redis.Options {
	option := &redis.Options{
		Addr:        config.GetString(getKey(id, "address")),
		Password:    config.GetString(getKey(id, "password")),
		DB:          config.GetInt(getKey(id, "db")),
		PoolSize:    config.GetInt(getKey(id, "pool_size")),
		PoolTimeout: time.Duration(config.GetInt(getKey(id, "pool_timeout"))) * time.Second,
	}
	if option.PoolTimeout == 0 {
		option.PoolTimeout = 30
	}

	if option.PoolSize == 0 {
		option.PoolSize = 10
	}

	return option
}

// Connect retrieve redis client connection
func Connect(id string) (*Client, error) {
	redisConfig := config.GetStringMap("redis")
	if _, ok := redisConfig[id]; !ok {
		return nil, fmt.Errorf("redis configuration for [%s] does not exists", id)
	}

	conn := redisConnections.Get(id)
	if conn != nil {
		if conn.Ping().Err() != nil {
			return newConnection(id)
		}
	}

	return newConnection(id)
}

// Shutdown close all established redis connections
func Shutdown() (err error) {
	if redisConnections == nil {
		return nil
	}
	for _, c := range redisConnections.DB {
		err = c.Close()
	}
	return err
}

func newConnection(id string) (*Client, error) {
	opt := setupConfig(id)
	r := redis.NewClient(opt)

	if err := r.Ping().Err(); err != nil {
		fallback := config.GetString(getKey(id, "fallback_to"))
		if fallback == id { // prevent endless loop
			return redisConnections.Get(id), err
		}
		return Connect(fallback)
	}

	redisConnections.Set(id, &Client{Client: r})
	return redisConnections.Get(id), nil
}

func getKey(id, types string) string {
	return fmt.Sprintf("redis.%s.%s", id, types)
}
