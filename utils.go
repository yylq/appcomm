package appcomm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/yylq/log"
)

func HttpPost(obj interface{}, addr string) (int, []byte, error) {
	log.Debugf("%v", obj)
	buff, err := json.Marshal(obj)
	if err != nil {
		log.Error(err)
		return 500, nil, err
	}
	log.Debugf("%s", buff)
	bio := bytes.NewReader(buff)
	rsp, err := http.Post(addr, "text", bio)
	if err != nil {
		log.Error(err)
		return 500, nil, err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)

	return rsp.StatusCode, data, err
}
