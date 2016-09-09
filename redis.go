package common

import (
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
func (c *Conn) GetAll(key string) (map[string]string, error) {
	return redis.StringMap(c.conn.Do("HGETALL", key))
}
func (c *Conn) HMGet(v ...interface{}) (map[string]string, error) {
	return redis.StringMap(c.conn.Do("HMGET", v...))
}
func (c *Conn) HSet(v ...interface{}) error {
	_, err := redis.Int64(c.conn.Do("HSET", v...))
	return err

}
func (c *Conn) HDel(v ...interface{}) error {
	_, err := redis.Int64(c.conn.Do("Del", v...))
	return err

}
func (c *Conn) Set(key string, v interface{}) error {
	_, err := redis.String(c.conn.Do("SET", key, v))
	return err
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
