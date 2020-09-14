package pkg

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

// RedisConn -
type RedisConn struct {
	Conn redis.Conn
}

func createPool() {
	pool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", "redis:6379")
	}, 10)
}

// InitRedis -
func InitRedis() error {
	createPool()
	conn := GetConn()
	defer conn.Conn.Close()
	_, err := conn.Conn.Do("PING")
	return err
}

// GetConn -
func GetConn() RedisConn {
	return RedisConn{
		Conn: pool.Get(),
	}
}

// Set -
func (r RedisConn) Set(key, value string, duration time.Duration) error {
	_, err := r.Conn.Do("SET", key, value, "EX", duration.Seconds())
	return err
}

// Get -
func (r RedisConn) Get(key string) (string, error) {
	rep, err := r.Conn.Do("GET", key)
	return string(rep.([]byte)), err
}
