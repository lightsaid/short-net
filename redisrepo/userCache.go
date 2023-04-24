package redisrepo

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
	"golang.org/x/exp/slog"
)

var (
	ErrNoFoundUser = errors.New("user no found")
)

// SetUserCache 设置用户缓存
func (r *redisRepo) SetUserCache(user models.User) error {
	conn := r.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s%s%d", keyprefix, "user:", user.ID)
	// NOTE: 这里没有给 models.User 设置 redis tag 因此，存的key根字段名一致, （不能改成小写开头），否则 redis.ScanStruct 扫描不到。
	_, err := conn.Do("HMSET", key,
		"ID", user.ID,
		"Name", user.Name,
		"Email", user.Email,
		"Avatar", user.Avatar,
		"Role", user.Role,
	)
	if err != nil {
		slog.Error("redis.Do HMSET failed", "err", err.Error(), "userid", user.ID)
		return err
	}

	// 设置一个随机过期时间
	expire := time.Minute * time.Duration(util.RandomInt(1, 2))
	_, err = conn.Do("EXPIRE", key, expire.Seconds())
	if err != nil {
		slog.Error("redis.Do EXPIRE failed", "err", err.Error(), "userid", user.ID)
		return err
	}

	return nil
}

// GetUserCache 获取用户缓存
func (r *redisRepo) GetUserCache(id uint) (models.User, error) {
	conn := r.pool.Get()
	defer conn.Close()

	var user models.User

	key := fmt.Sprintf("%s%s%d", keyprefix, "user:", id)

	// reply, err := conn.Do("HGETALL", key)
	// slog.Info("conn.Do(HGETALL, key)", "reply", reply, "err", err)
	values, err := redis.Values(conn.Do("HGETALL", key))
	// NOTE: 如果key，不存在或者key过期，nil 都是为 nil
	if err != nil {
		slog.Error("redis.Do HGETALL failed", "err", err.Error(), "userid", id, "is redis.ErrNil", redis.ErrNil == err)
		return user, err
	} else if len(values) == 0 {
		return user, ErrNoFoundUser
	}

	slog.Info("redis.Values - HGETALL - "+key, "values[0]", values[0], "err", err, "len", len(values))

	slog.Info("values all ", "values[0]", values[0], "values[1]", values[1], "values[2]", values[2], "values[3]", values[3])

	err = redis.ScanStruct(values, &user)
	if err != nil {
		slog.Error("redis.ScanStruct failed", "err", err.Error(), "userid", id, "values", values)
		return user, err
	}

	slog.Info("redis.ScanStruct - "+key, "user", user)

	return user, nil
}
