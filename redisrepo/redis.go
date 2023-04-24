package redisrepo

import (
	"github.com/gomodule/redigo/redis"
	"github.com/lightsaid/short-net/models"
)

const keyprefix = "shortnet:"

type RedisRepository interface {
	SetUserCache(user models.User) error
	GetUserCache(id uint) (models.User, error)
}

type redisRepo struct {
	pool *redis.Pool
}

func NewRedisRepository(pool *redis.Pool) RedisRepository {
	return &redisRepo{
		pool: pool,
	}
}
