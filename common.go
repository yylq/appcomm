package common

import (
	"godsee/log"
	"os"

	"github.com/larspensjo/config"
)

func Load_conf(fn string, c interface{}) error {
	/*
		cf := config.Config{}
		cf.LoadFile(fn)

		if err := cf.Load(appconf); err != nil {

			return err
		}
	*/
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
