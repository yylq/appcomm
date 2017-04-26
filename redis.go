package appcomm

import (
	"errors"
	"fmt"
	//"reflect"

	"time"

	"github.com/garyburd/redigo/redis"
)

type Conn struct {
	conn    redis.Conn
	ch      string
	subconn redis.PubSubConn
}

func (c *Conn) Get(key string) (string, error) {
	return redis.String(c.conn.Do("GET", key))
}
func (c *Conn) GetInt64(key string) (int64, error) {
	n, err := redis.Int64(c.conn.Do("GET", key))
	if err == nil {
		return n, err
	}
	if err == redis.ErrNil {
		return 0, nil
	}
	return 0, err
}
func (c *Conn) GetAll(key string) (map[string]string, error) {
	return redis.StringMap(c.conn.Do("HGETALL", key))
}
func (c *Conn) HMGet(v ...interface{}) (map[string]string, error) {
	return redis.StringMap(c.conn.Do("HMGET", v...))
}

func (c *Conn) HMGetEx(rtype interface{}, v ...interface{}) (interface{}, error) {

	s, err := redis.Strings(c.conn.Do("HMGET", v...))
	if err != nil {
		return nil, err
	}
	switch rtype := rtype.(type) {
	case int64, int32:
		return redis.Int64([]byte(s[0]), err)
	case string:
		return s[0], err
	default:
		return nil, errors.New(fmt.Sprintf("unsurport rtyp %#v", rtype))
	}
}
func (c *Conn) HMGetInt64(v ...interface{}) (int64, error) {
	s, err := redis.Strings(c.conn.Do("HMGET", v...))
	if err != nil {
		return 0, err
	}
	return redis.Int64([]byte(s[0]), err)
}
func (c *Conn) HMGetString(v ...interface{}) (string, error) {
	s, err := redis.Strings(c.conn.Do("HMGET", v...))
	if err != nil {
		return "", err
	}
	return s[0], nil
}
func (c *Conn) HMSet(v ...interface{}) error {
	_, err := redis.String(c.conn.Do("HMSET", v...))
	return err

}

func (c *Conn) HDel(v ...interface{}) error {
	_, err := redis.Int64(c.conn.Do("HDEL", v...))
	return err

}
func (c *Conn) HExists(key, field string) (bool, error) {
	v, err := redis.Int64(c.conn.Do("HEXISTS", key, field))
	return v == 1, err

}
func (c *Conn) HKeys(key string) ([]string, error) {
	return redis.Strings(c.conn.Do("HKeys", key))
}
func (c *Conn) HKeysInts(key string) ([]int, error) {
	return redis.Ints(c.conn.Do("HKeys", key))
}
func (c *Conn) Del(v ...interface{}) error {
	_, err := redis.Int64(c.conn.Do("DEL", v...))
	return err

}
func (c *Conn) Set(key string, v interface{}) error {
	_, err := redis.String(c.conn.Do("SET", key, v))
	return err
}

func (c *Conn) LPush(v ...interface{}) error {
	_, err := redis.Int(c.conn.Do("LPush", v...))
	return err
}
func (c *Conn) LRange(key string, start, stop int) ([]string, error) {
	return redis.Strings(c.conn.Do("LRange", key, start, stop))
}
func (c *Conn) LIndex(key string, ind int) (string, error) {
	s, err := redis.String(c.conn.Do("LIndex", key, ind))
	if err == nil {
		return s, nil
	}
	if err == redis.ErrNil {
		return "", nil
	}
	return "", err
}
func (c *Conn) RPop(key string) string {
	s, err := redis.String(c.conn.Do("Rpop", key))
	if err != nil {
		return ""
	}
	return s
}
func (c *Conn) Llen(key string) int {
	n, err := redis.Int(c.conn.Do("Llen", key))
	if err != nil {
		return 0
	}
	return n
}
func (c *Conn) LRem(key, val string, count int) int {
	n, err := redis.Int(c.conn.Do("LRem", key, count, val))
	if err != nil {
		return 0
	}
	return n
}
func (c *Conn) LRemAll(key, val string) int {
	return c.LRem(key, val, 0)
}
func (c *Conn) Close() {
	c.conn.Close()
}
func (c *Conn) PSub() {

}
func NewRedisConn(host string, db int) (*Conn, error) {
	c, err := redis.Dial("tcp", host, redis.DialDatabase(db))
	if err != nil {
		return nil, err
	}
	return &Conn{conn: c}, nil
}
func NewRedisChannel(host string, channel string, timeout int) (*Conn, error) {

	c, err := redis.Dial("tcp", host, redis.DialConnectTimeout(time.Duration(timeout)*time.Second))
	if err != nil {

		return nil, err
	}
	psc := redis.PubSubConn{Conn: c}
	return &Conn{subconn: psc, ch: channel}, nil
}
