package imredis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisClient struct {
	pool *redis.Pool
}

// protocol is unix or tcp
func NewRedis(protocol, server string) *redisClient {
	pool := &redis.Pool{
		IdleTimeout: 30 * time.Second,
		Wait:        false,
		Dial:        func() (redis.Conn, error) { return redis.Dial(protocol, server) },
	}

	return &redisClient{
		pool: pool,
	}
}

func (rc *redisClient) FetchConn() interface{} {
	return rc.pool.Get()
}

func (rc *redisClient) SingleSet(key string, value []byte) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func (rc *redisClient) SingleGet(key string) ([]byte, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("GET", key))
}

func (rc *redisClient) MultiSet(set map[string][]byte) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("MSET", redis.Args{}.AddFlat(set)...)
	return err
}

func (rc *redisClient) MultiGet(keys []string) ([][]byte, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	return redis.ByteSlices(conn.Do("MGET", redis.Args{}.AddFlat(keys)...))
}

func (rc *redisClient) SingleDelete(key string) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (rc *redisClient) MultiDelete(keys []string) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", redis.Args{}.AddFlat(keys)...)
	return err
}

func (rc *redisClient) SingleSetNX(key string, value []byte) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	ok, err := redis.Int(conn.Do("SETNX", key, value))
	if err != nil {
		return 0, err
	}
	return ok, nil
}

func (rc *redisClient) Increment(key string, delta uint64) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	length, _ := redis.Int(conn.Do("STRLEN", key))
	if length <= 0 {
		return 0, redis.ErrNil
	}

	ok, err := redis.Int(conn.Do("INCRBY", key, delta))
	if err != nil {
		return 0, err
	}
	return ok, nil
}

func (rc *redisClient) Flush() error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHALL")
	return err
}
