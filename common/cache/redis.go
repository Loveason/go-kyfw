package cache

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisCache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	password string
}

// actually do the redis cmds
func (rc *RedisCache) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

func (rc *RedisCache) RPop(key string) (v interface{}, err error) {
	return rc.Do("RPOP", key)
}

func (rc *RedisCache) LPush(key string, vals ...interface{}) error {
	_, err := rc.Do("LPUSH", key, vals)
	return err
}

func (rc *RedisCache) LLen(key string) int64 {
	if v, err := rc.Do("LLEN", key); err == nil {
		return v.(int64)
	}
	return 0
}

func (rc *RedisCache) SetEX(key string, val interface{}, timeout time.Duration) error {

	_, err := rc.Do("SETEX", key, int64(timeout/time.Second), val)

	return err
}

func (rc *RedisCache) Set(key string, val interface{}) error {

	_, err := rc.Do("SET", key, val)

	return err
}

func (rc *RedisCache) Get(key string) (ret interface{}, err error) {
	return rc.Do("GET", key)
}

func (rc *RedisCache) Del(key string) (ret interface{}, err error) {
	return rc.Do("DEL", key)
}

func (rc *RedisCache) Ltrim(key string, start, stop int64) error {
	_, err := rc.Do("LTRIM", key, start, stop)
	return err
}

func (rc *RedisCache) LRange(key string, start, stop int) (v interface{}, err error) {
	return rc.Do("LRANGE", key, start, stop)
}

func (rc *RedisCache) Incr(key string) (v interface{}, err error) {
	return rc.Do("INCR", key)
}

func (rc *RedisCache) HSet(key, field string, value interface{}) (v interface{}, err error) {
	return rc.Do("HSET", key, field, value)
}

func (rc *RedisCache) HKeys(key string) (v interface{}, err error) {
	return rc.Do("HKEYS", key)
}

func (rc *RedisCache) HVals(key string) (v interface{}, err error) {
	return rc.Do("HVALS", key)
}

func (rc *RedisCache) HGet(key, field string) (v interface{}, err error) {
	return rc.Do("HGET", key, field)
}

func (rc *RedisCache) HDel(key, field string) (v interface{}, err error) {
	return rc.Do("HDEL", key, field)
}

func (rc *RedisCache) HExists(key, field string) (v interface{}, err error) {
	return rc.Do("HEXISTS", key, field)
}

// connect to redis.
func (rc *RedisCache) Connect(config string) error {
	var cf map[string]string

	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()
	return c.Err()
}

func (rc *RedisCache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}

	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}
