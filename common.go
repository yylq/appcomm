package appcomm

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/yylq/log"

	"strconv"

	"crypto/md5"

	"github.com/yylq/config"
)

func Load_conf(fn string, c interface{}) error {

	cfg, err := config.ReadDefault(fn)
	if err != nil {
		return err
	}

	err = cfg.ParseConf(c)
	if err != nil {
		return err
	}

	return nil
}

func GetHostId() string {
	return GetHostName()
}

func GetHostName() string {
	var id string
	var err error
	id = os.Getenv("HOSTNAME")
	if id == "" {
		id, err = os.Hostname()
		if err != nil {
			log.Errorf(" hostname err:%s", err)
			return ""
		}

		err = os.Setenv("HOSTNAME", id)
		if err != nil {
			log.Errorf(" Setenv err:%s", err)
			return ""
		}

	}
	return id
}
func PahtExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetPid(file string) (int, error) {
	buff, err := ioutil.ReadFile(file)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(buff))
}
func WritePid(file string, pid int) error {
	s := fmt.Sprintf("%d", pid)
	return ioutil.WriteFile(file, []byte(s), os.ModePerm)
}

func Md5(data string) string {
	var in []byte
	d := md5.New()
	d.Reset()
	d.Write([]byte(data))
	return fmt.Sprintf("%x", d.Sum(in))
}
func GetDateString() string {
	Y, M, D := time.Now().Date()
	return fmt.Sprintf("%4d%02d%02d", Y, M, D)
}
