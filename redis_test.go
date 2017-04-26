package appcomm

import (
	"fmt"
	"reflect"
	"sort"
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
	err = c.Del("aaa")
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
	err = c.HMSet("box", "box_02", []byte(s))
	if err != nil {
		t.Fatal(err)
	}
	err = c.HDel("box", "box_02")
	if err != nil {
		t.Fatal(err)
	}

}
func TestHGet(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	host := "10.138.71.181"

	res, err := c.HMGetString("table_lost", host)
	if err != nil {
		t.Fatalf(" check lost %s err:%v", host, err)
	}
	t.Log(reflect.TypeOf(res))
	t.Log(res)
}
func TestHDel(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if err := c.HDel("table_lost", "10.138.71.181", "task_0"); err != nil {
		t.Fatal(err)
	}

}
func TestHMSet(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	host := "10.138.71.181"
	tid := "task_01"
	if err := c.HMSet("table_lost", host, tid, tid, host); err != nil {
		t.Fatal(err)
	}

}

func TestHMGetEx(t *testing.T) {

	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if err := c.HMSet("table_lost", "100", "domain", "dispacth", 500); err != nil {
		t.Fatal(err)
	}
	h, err := c.HMGetString("table_lost", "100")
	if err != nil {
		t.Fatal(err)
	}
	if h != "domain" {
		t.Fatal("check value fail")
	}
	i, err := c.HMGetInt64("table_lost", "dispacth")
	if err != nil {
		t.Fatal(err)
	}
	if i != 500 {
		t.Fatal("check value fail")
	}
	d, err := c.HMGetEx(int32(1), "table_lost", "dispacth")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#d", d)
	if d != int64(500) {
		t.Fatal("HMGetEx check value fail")
	}
}
func TestLPush(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := c.LPush("lang", "c", "c++", "java", "php"); err != nil {
		t.Fatal(err)
	}
}
func TestRPop(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.Del("lang")
	if err := c.LPush("lang", "c", "c++", "java", "php"); err != nil {
		t.Fatal(err)
	}
	langs := []string{"c", "c++", "java", "php"}
	for _, one := range langs {
		rone := c.RPop("lang")
		if rone != one {
			t.Logf("one %v rone %v", one, rone)
			t.Fatal("LPop fail")
		}
	}

}
func TestLRange(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.Del("lang")
	items, err := c.LRange("lang", -1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("lrange get wrong")
	}
	if err := c.LPush("lang", "c", "c++", "java", "php"); err != nil {
		t.Fatal(err)
	}
	langs := []string{"c", "c++", "java", "php"}

	for i := 0; i < len(langs); i++ {
		items, err = c.LRange("lang", -1, -1)
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 1 {
			t.Fatal("lrange get wrong")
		}
		if items[0] != langs[i] {
			t.Fatal("lrange get wrong ", i, items[0], langs[i])
		}
		c.RPop("lang")

	}

}
func TestLIndex(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.Del("lang")

	if err := c.LPush("lang", "c", "c++", "java", "php"); err != nil {
		t.Fatal(err)
	}
	langs := []string{"c", "c++", "java", "php"}
	for i := 1; i <= len(langs); i++ {
		ind := ^i + 1
		t.Logf("%d", ind)
		items, err := c.LIndex("lang", ind)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf(" %d %v", i, items)
		if items != langs[i-1] {
			t.Fatal("lrange get wrong ", i, items, langs[i-1])
		}

		t.Log(i, items, langs[i-1])
	}
	for i := 1; i < 10; i++ {
		ind := ^i + 1
		t.Logf("%d", ind)
		items, err := c.LIndex("lang", ind)
		if err != nil {
			t.Fatal(err)
		}
		if items == "" {
			break
		}
		t.Logf(" %d %v", i, items)
	}
}
func TestLILen(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.Del("lang")
	n := c.Llen("lang")

	t.Logf("n:%d", n)
	if err := c.LPush("lang", "c", "c++", "java", "php"); err != nil {
		t.Fatal(err)
	}
	langs := []string{"c", "c++", "java", "php"}
	n = c.Llen("lang")

	if n != len(langs) {
		t.Fatal("clen error", n, len(langs))
	}
	t.Logf("n:%d", n)
	/*
		ind := ^n + 1
		items, err := c.LRange("lang", ind, -1)
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != len(langs) {
			t.Fatal("LRange error", n, len(langs))
		}
		for i := 0; i < len(items); i++ {
			if items[i] != langs[i] {
				t.Fatal("LRange error", i, items[i], langs[i])
			}
		}
	*/
}

const (
	Last_File_Id  = "last_file_id"
	List_Refresh  = "list_refresh"
	Table_Refresh = "table_refresh"
)

func TestRm(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	n := c.Llen(List_Refresh)
	if n < 1 {
		return
	}
	t.Logf("n : %v", n)

	for i := n; i > 0; i-- {
		ind := i
		t.Logf("%v", ind)
		tid, err := c.LIndex(List_Refresh, ind)
		if err != nil {
			t.Fatal(err)
			break
		}
		t.Logf("ind:%v tid :%v", ind, tid)
		task, err := c.HMGetString(Table_Refresh, tid)
		if err != nil {
			t.Fatal(err)
			break
		}
		t.Logf("%v", task)

	}
}
func TestHKey(t *testing.T) {
	c, err := NewRedisConn(redisip, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.Del("table_keys")
	pos := 90996546
	for i := 0; i < 10; i++ {
		id := pos + i
		c.HMSet("table_keys", id, fmt.Sprintf("{id:%d}", id))
	}
	ids, err := c.HKeysInts("table_keys")

	if err != nil {
		t.Fatal(err)
	}
	sort.Ints(ids)
	for i := 0; i < 10; i++ {
		if ids[i] != pos+i {
			t.Fatal(fmt.Sprintf(" i:%v %v %v ", i, ids[i], pos+i))
		}
	}
	c.Del("table_keys")
	ids, err = c.HKeysInts("table_keys")

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("len(ids):%v", len(ids))
}
