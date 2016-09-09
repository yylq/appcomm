package common

import (
	"godsee/log"
	"os"

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
