package redis

import (
	"context"

	"github.com/go-redis/redis"
)

// Pipeline struct
type Pipeline struct {
	context context.Context
	redis.Pipeliner
}

// Exec : redis exec command
func (p *Pipeline) Exec() ([]redis.Cmder, error) {
	return p.Pipeliner.Exec()
}
