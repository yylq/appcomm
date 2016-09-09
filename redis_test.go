package common

import (
	"fmt"
	"testing"
)

const (
	redisip string = "192.168.176.7:6379"
	redisdb int    = 3
)

func TestGet(t *testing.T) {
	c, err := NewRedisConn(redisip, redisdb)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	s := "dddd"
	err = c.Set("aaa", s)
	if err != nil {
		t.Fatal(err)
	}
	s1, err := c.Get("aaa")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%T %v", s1, s1)
	if s1 != s {
		t.Fatal("set and get fail")
	}

}
func TestSet(t *testing.T) {
	c, err := NewRedisConn(redisip, redisdb)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	s := []string{"dddd", "aaaa"}
	ss := fmt.Sprint(s)
	t.Logf("%v %s", s, ss)
	err = c.Set("aaa", s)
	if err != nil {
		t.Fatal(err)
	}
	s1, err := c.Get("aaa")
	if err != nil {
		t.Fatal(err)
	}

	if s1 != ss {
		t.Fatal("set and get fail")
	}
	t.Logf("%T %v", s1, s1)
	err = c.HDel("aaa")
	if err != nil {
		t.Fatal(err)
	}

}
func TestHSet(t *testing.T) {
	c, err := NewRedisConn(redisip, redisdb)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	s := `{"aaa":"1111","bbb":2222}`
	err = c.HSet("box", "box_02", []byte(s))
	if err != nil {
		t.Fatal(err)
	}
	err = c.HDel("box", "box_02")
	if err != nil {
		t.Fatal(err)
	}

}
