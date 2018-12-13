package redis

import (
	"context"

	"github.com/go-redis/redis"
)

// WithContext add context to redisClient
func (r *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		Client:  r.Client,
		context: ctx,
	}
}

// Get : redis GET
func (r *Client) Get(key string) *redis.StringCmd {
	return r.Client.Get(key)
}

// HSet : redis HSET
func (r *Client) HSet(key, field string, value interface{}) *redis.BoolCmd {
	return r.Client.HSet(key, field, value)
}

//HGet : redis HGET
func (r *Client) HGet(key, field string) *redis.StringCmd {
	return r.Client.HGet(key, field)
}

//Pipeline : redis pipeline
func (r *Client) Pipeline() redis.Pipeliner {
	return &Pipeline{Pipeliner: r.Client.Pipeline(), context: r.context}
}
